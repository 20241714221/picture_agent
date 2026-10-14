
package main

import (
	"context"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C9EFF")).
		MarginBottom(1)

	outputStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C9EFF")).
		Padding(1, 2)

	inputStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5ED6A0")).
		Padding(0, 1)

	hintStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))
)

type agentResultMsg struct {
	result string
}

func runAgent(agent *Agent, text string) tea.Cmd {
	return func() tea.Msg {
		result := agent.Run(context.Background(), text)
		return agentResultMsg{result: result}
	}
}

type TUIModel struct {
	agent    *Agent
	input    textinput.Model
	viewport viewport.Model

	output  string
	width   int
	height  int
	loading bool
}

func NewTUI(agent *Agent) TUIModel {
	input := textinput.New()
	input.Placeholder = "输入饮食描述或图片路径..."
	input.Focus()

	vp := viewport.New(
		viewport.WithWidth(60),
		viewport.WithHeight(15),
	)

	output := "你好！我是你的饮食管理助手。\n" +
		"请输入食物图片路径或者描述你的饮食。"

	vp.SetContent(output)

	return TUIModel{
		agent:    agent,
		input:    input,
		viewport: vp,
		output:   output,
	}
}

func (m TUIModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// 外层对话框宽度
		panelWidth := m.width - 4
		if panelWidth < 16 {
			panelWidth = 16
		}

		// 边框和内边距共占约 6 列
		contentWidth := panelWidth - 6
		if contentWidth < 10 {
			contentWidth = 10
		}

		// 为标题、输入框和提示文字预留高度
		contentHeight := m.height - 10
		if contentHeight < 3 {
			contentHeight = 3
		}

		m.viewport.SetWidth(contentWidth)
		m.viewport.SetHeight(contentHeight)
		m.input.SetWidth(contentWidth)

		m.viewport.SetContent(m.output)
		m.viewport.GotoBottom()

		return m, nil

	case agentResultMsg:
		m.loading = false

		// 删除之前显示的“正在思考...”
		m.output = strings.TrimSuffix(
			m.output,
			"\n\nAI:\n正在思考...",
		)

		// 只追加 AI 的回答，不再重复追加用户问题
		m.output += "\n\nAI:\n" + msg.result

		m.viewport.SetContent(m.output)
		m.viewport.GotoBottom()

		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			if m.loading {
				return m, nil
			}

			text := strings.TrimSpace(m.input.Value())
			if text == "" {
				return m, nil
			}

			m.input.SetValue("")
			m.loading = true

			// 用户的问题只追加一次
			m.output += "\n\n你:\n" + text
			m.output += "\n\nAI:\n正在思考..."

			m.viewport.SetContent(m.output)
			m.viewport.GotoBottom()

			return m, runAgent(m.agent, text)
		}
	}

	// 分别更新两个组件，避免命令互相覆盖
	var cmds []tea.Cmd

	newViewport, viewportCmd := m.viewport.Update(msg)
	m.viewport = newViewport
	if viewportCmd != nil {
		cmds = append(cmds, viewportCmd)
	}

	newInput, inputCmd := m.input.Update(msg)
	m.input = newInput
	if inputCmd != nil {
		cmds = append(cmds, inputCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m TUIModel) View() tea.View {
	title := titleStyle.Render("🥗 饮食管理 Agent")
	output := outputStyle.Render(m.viewport.View())
	input := inputStyle.Render(m.input.View())

	hintText := "Enter 发送 · ↑↓ 滚动 · PgUp/PgDn 翻页 · Esc 退出"
	if m.loading {
		hintText = "AI 正在思考... · Esc 退出"
	}

	hint := hintStyle.Render(hintText)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		output,
		input,
		hint,
	)

	return tea.View{
		Content:  content,
		AltScreen: true,
	}
}

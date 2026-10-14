package tools

import (
	"context"
	"encoding/json"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type ToolDef struct {
	Name        string
	Description string
	Params      openai.FunctionParameters
	Func        func(ctx context.Context, args map[string]any) (string, error)
}

func (t ToolDef) ToParam() responses.ToolUnionParam {
	return responses.ToolUnionParam{
		OfFunction : &responses.FunctionToolParam{
			Name: t.Name,
			Description: openai.String(t.Description), //转换成字符串指针
			Parameters: t.Params,
		},
	}
}

func (t ToolDef) Call(ctx context.Context, argsJSON string) (string, error) {
	var args map[string]any
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", err
	}
	return t.Func(ctx, args)

}
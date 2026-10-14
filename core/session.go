package core

import (
	"github.com/openai/openai-go/v3/responses"
)


// Session 用来保存Agent运行过程中的消息
type Session struct {

	// 保存所有输入输出消息
	Messages []responses.ResponseInputItemUnionParam

}



// 创建一个新的Session
func NewSession() *Session {

	return &Session{
		Messages: make(
			[]responses.ResponseInputItemUnionParam,
			0,
		),
	}

}



// 添加消息
func (s *Session) AddMessage(
	msg ...responses.ResponseInputItemUnionParam,
) {

	s.Messages = append(
		s.Messages,
		msg...,
	)

}



// 获取所有消息
func (s *Session) GetMessages() []responses.ResponseInputItemUnionParam {

	return s.Messages

}



// 清空消息
func (s *Session) Clear(){

	s.Messages = make(
		[]responses.ResponseInputItemUnionParam,
		0,
	)

}
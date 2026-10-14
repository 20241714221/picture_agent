package tools

import (
	"context"
	"time"
	"fmt"
	"github.com/openai/openai-go/v3"
)


func NewCurrentTimeTool() *ToolDef {

	return &ToolDef{

		Name: "get_current_time",

		Description:
		"获取当前日期和时间，用于判断用户饮食记录的日期和餐次",


		Params: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{},
		},


		Func: func(
			ctx context.Context,
			args map[string]any,
		)(string,error){


			now := time.Now()


      return fmt.Sprintf(
          `{
              "date":"%s",
              "time":"%s"
          }`,
          now.Format("2006-01-02"),
          now.Format("15:04:05"),
      ),nil
		},
	}
}
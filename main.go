package main

import (
	database "picture_agent/databases"
	"picture_agent/tools"

	"context"
	"database/sql"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)


func getLlmTools(
	client *openai.Client,
	db *sql.DB,
) map[string]*tools.ToolDef {


	llmTools := make(map[string]*tools.ToolDef)


	foodTool := tools.NewFoodRecognitionTool(client)

	timeTool := tools.NewCurrentTimeTool()

	dailyRecordTool := tools.NewDailyRecordTool(db)

	nutritionQueryTool := tools.NewNutritionQueryTool(db)



	llmTools[foodTool.Name] = foodTool

	llmTools[timeTool.Name] = timeTool

	llmTools[dailyRecordTool.Name] = dailyRecordTool

	llmTools[nutritionQueryTool.Name] = nutritionQueryTool


	return llmTools
}



func main() {


	ctx := context.Background()


	// 普通文本模型
	client := openai.NewClient(

		option.WithAPIKey(""),

		option.WithBaseURL(
			"https://api.deepseek.com",
		),

		option.WithRequestTimeout(
			30*time.Second,
		),
	)



	// 视觉模型
	visionClient := openai.NewClient(

		option.WithAPIKey(""),

		option.WithBaseURL(
			"https://maas.qianwenaiapi.com/compatible-mode/v1",
		),

		option.WithRequestTimeout(
			100*time.Second,
		),
	)



	// 数据库

	db,err:=database.NewMySQLDB()

	if err!=nil{

		panic(err)

	}


	defer db.Close()



	// 注册工具

	llmTools:=getLlmTools(
		&visionClient,
		db,
	)



	// 创建Agent

	agent:=NewAgent(

		ctx,

		&client,

		llmTools,

	)



	// 启动Bubble Tea


	p:=tea.NewProgram(

		NewTUI(agent),
		
	)



	if _,err:=p.Run();err!=nil{

		panic(err)

	}


}
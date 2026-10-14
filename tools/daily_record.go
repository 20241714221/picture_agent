package tools

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
)

func NewDailyRecordTool(db *sql.DB) *ToolDef {
	return &ToolDef{
		Name:        "daily_record",
		Description: "保存用户当天饮食记录，包括食物名称、重量、热量、蛋白质、脂肪和碳水",
		Params: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"meal_type": map[string]any{
					"type":        "string",
					"description": "餐次，例如早餐、午餐、晚餐",
				},
				"foods": map[string]any{
					"type":        "array",
					"description": "食物列表",
				},
				"date": map[string]any{
					"type":        "string",
					"description": "饮食日期，例如2026-10-12",
				},
			},
			"required": []string{
				"meal_type",
				"foods",
				"date",
			},
		},

		Func: func(
			ctx context.Context,
			args map[string]any,
		) (string, error) {

			mealType, _ := args["meal_type"].(string)
			foodsBytes, err := json.Marshal(
				args["foods"],
			)
			if err != nil {

				return "", err
			}
			var foods []map[string]any
			err = json.Unmarshal(
				foodsBytes,
				&foods,
			)

			if err != nil {
				return "", err
			}
			date, _ := args["date"].(string)
			for _, food := range foods {
				_, err = db.ExecContext(
					ctx,
					`
					INSERT INTO daily_food_records
					(
					user_id,
					eat_date,
					meal_type,
					food_name,
					weight,
					calorie,
					protein,
					fat,
					carbohydrate
					)
					VALUES
					(?,?,?,?,?,?,?,?,?)
					`,

					1,
					date,
					mealType,
					food["name"],
					food["weight"],
					food["calorie"],
					food["protein"],
					food["fat"],
					food["carbohydrate"],
				)
				if err != nil {
					return "",
						fmt.Errorf(
							"保存失败:%v",
							err,
						)
				}
			}
			return "饮食记录保存成功", nil
		},
	}
}

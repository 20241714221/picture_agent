package tools

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/openai/openai-go/v3"
)

func NewNutritionQueryTool(db *sql.DB) *ToolDef {
	return &ToolDef{
		Name:        "nutrition_query",
		Description: "查询指定日期的饮食营养摄入情况，包括总热量、蛋白质、脂肪和碳水",
		Params: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"date": map[string]any{
					"type":        "string",
					"description": "查询日期，例如2026-10-13",
				},
			},
			"required": []string{
				"date",
			},
		},
		Func: func(ctx context.Context, args map[string]any) (string, error) {
			date, _ := args["date"].(string)
			var name string
			var calorie float64
			var protein float64
			var fat float64
			var carbohydrate float64
			err := db.QueryRowContext(
				ctx,

				`
		  SELECT
            COALESCE(GROUP_CONCAT(food_name SEPARATOR '、'), ''),
            COALESCE(SUM(calorie), 0),
            COALESCE(SUM(protein), 0),
            COALESCE(SUM(fat), 0),
            COALESCE(SUM(carbohydrate), 0)
        FROM daily_food_records
        WHERE eat_date = ? AND user_id = ?
				`,
				date,
				1,
			).Scan(
				&name,
				&calorie,
				&protein,
				&fat,
				&carbohydrate,
			)
			if err != nil {
				return "", err
			}
			result := map[string]any{
				"date":         date,
				"name":         name,
				"calorie":      calorie,
				"protein":      protein,
				"fat":          fat,
				"carbohydrate": carbohydrate,
			}
			data, _ := json.Marshal(result)
			return string(data), nil
		},
	}
}

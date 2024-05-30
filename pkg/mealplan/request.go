package mealplan

import (
	"Food/pkg/recipe"
)

type AddRequest struct {
	WeekStartDate string                    `json:"week_start_date"` // ISO8601 date format
	Meals         map[string]map[string]int `json:"meals"`           // Nested map [dayOfWeek][mealType]recipeId
}

type GetResponse struct {
	UserID        int                                  `json:"user_id"`
	WeekStartDate string                               `json:"week_start_date"`
	Meals         map[string]map[string]recipe.Recipes `json:"meals"`
}

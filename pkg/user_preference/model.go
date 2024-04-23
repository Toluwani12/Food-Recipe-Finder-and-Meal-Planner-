package user_preference

import (
	"encoding/json"
)

type UserPreference struct {
	UserID                int             `json:"user_id"`
	Vegetarian            bool            `json:"vegetarian"`
	GlutenFree            bool            `json:"gluten_free"`
	CuisinePreference     string          `json:"cuisine_preference"`
	DislikedIngredients   []string        `json:"disliked_ingredients"`
	AdditionalPreferences json.RawMessage `json:"additional_preferences"` // RawMessage is used to delay JSON decoding
	DietaryGoals          string          `json:"dietary_goals"`
}

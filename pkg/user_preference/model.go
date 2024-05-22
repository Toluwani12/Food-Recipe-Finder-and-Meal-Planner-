package user_preference

import (
	"encoding/json"
	"github.com/lib/pq"
)

type UserPreference struct {
	UserID                string          `json:"user_id" db:"user_id"`
	Vegetarian            bool            `json:"vegetarian" db:"vegetarian"`
	GlutenFree            bool            `json:"gluten_free" db:"gluten_free"`
	CuisinePreference     string          `json:"cuisine_preference" db:"cuisine_preference"`
	DislikedIngredients   pq.StringArray  `json:"disliked_ingredients" db:"disliked_ingredients"`
	AdditionalPreferences json.RawMessage `json:"additional_preferences" db:"additional_preferences" ` // RawMessage is used to delay JSON decoding
	DietaryGoals          string          `json:"dietary_goals" db:"dietary_goals"`
}

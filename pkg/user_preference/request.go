package user_preference

import (
	"encoding/json"
	"github.com/lib/pq"
	"net/http"
)

type AddRequest struct {
	Vegetarian            bool            `json:"vegetarian"`
	GlutenFree            bool            `json:"gluten_free"`
	CuisinePreference     string          `json:"cuisine_preference"`
	DislikedIngredients   pq.StringArray  `json:"disliked_ingredients"`
	AdditionalPreferences json.RawMessage `json:"additional_preferences"` // JSON for flexibility
	DietaryGoals          string          `json:"dietary_goals"`
}

func (v *AddRequest) Bind(r *http.Request) error {
	return nil
}

type GetResponse struct {
	UserID                string          `json:"user_id"`
	Vegetarian            bool            `json:"vegetarian"`
	GlutenFree            bool            `json:"gluten_free"`
	CuisinePreference     string          `json:"cuisine_preference"`
	DislikedIngredients   pq.StringArray  `json:"disliked_ingredients"`
	AdditionalPreferences json.RawMessage `json:"additional_preferences"`
	DietaryGoals          string          `json:"dietary_goals"`
}

package user_preference

type AddRequest struct {
	Vegetarian            bool                   `json:"vegetarian"`
	GlutenFree            bool                   `json:"gluten_free"`
	CuisinePreference     string                 `json:"cuisine_preference"`
	DislikedIngredients   []string               `json:"disliked_ingredients"`
	AdditionalPreferences map[string]interface{} `json:"additional_preferences"` // JSON for flexibility
}

type GetResponse struct {
	UserID                int                    `json:"user_id"`
	Vegetarian            bool                   `json:"vegetarian"`
	GlutenFree            bool                   `json:"gluten_free"`
	CuisinePreference     string                 `json:"cuisine_preference"`
	DislikedIngredients   []string               `json:"disliked_ingredients"`
	AdditionalPreferences map[string]interface{} `json:"additional_preferences"`
}

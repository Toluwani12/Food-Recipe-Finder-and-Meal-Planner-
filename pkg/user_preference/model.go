package user_preference

type UserPreference struct {
	UserID       string   `json:"user_id" db:"user_id"`
	LikedRecipes []string `json:"liked_recipes" db:"liked_recipes"`
}

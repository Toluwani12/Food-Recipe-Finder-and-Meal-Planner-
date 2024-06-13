package user_preference

type UserPreference struct {
	UserID       string   `json:"user_id" db:"user_id"`
	LikedRecipes []Recipe `json:"liked_recipes" db:"liked_recipes"`
}

type Recipe struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

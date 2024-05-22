package recipe

import (
	"Food/pkg/ingredient"
	"time"
)

type Recipe struct {
	Id           string                  `json:"id" db:"id"`
	Name         string                  `json:"name" db:"name"`
	Description  string                  `json:"description" db:"description"`
	CookingTime  string                  `json:"cooking_time" db:"cooking_time"`
	Instructions string                  `json:"instructions" db:"instructions"`
	ImgUrl       string                  `json:"img_url" db:"img_url"`
	Ingredients  []ingredient.Ingredient `json:"ingredients"` // Convenient for handling full recipes
	CreatedAt    time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at" db:"updated_at"`
}

type Recipes = []Recipe

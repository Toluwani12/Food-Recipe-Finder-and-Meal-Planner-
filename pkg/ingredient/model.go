package ingredient

import (
	"time"
)

type Ingredient struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Alternatives []string  `json:"alternatives"`
	RecipeID     string    `json:"recipe_id"`
	Quantity     string    `json:"quantity"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Ingredients = []Ingredient

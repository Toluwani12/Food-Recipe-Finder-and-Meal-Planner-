package recipe

import (
	"time"
)

type Recipe struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	CookingTime  string    `json:"cooking_time"`
	Instructions string    `json:"instructions"`
	CreatedAt    time.Time `json:"created_at"`
	UpdateAt     time.Time `json:"updateAt"`
}

type Recipes = []Recipe

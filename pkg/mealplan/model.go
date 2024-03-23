package mealplan

import (
	"time"
)

type MealPlan struct {
	Id       string    `json:"id"`
	Date     time.Time `json:"date"`
	MealType time.Time `json:"meal_type"`
}

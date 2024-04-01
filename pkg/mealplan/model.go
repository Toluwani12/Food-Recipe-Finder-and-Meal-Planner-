package mealplan

import (
	"time"
)

type MealPlan struct {
	Id        string    `json:"id"`
	Date      time.Time `json:"date"`
	MealType  time.Time `json:"meal_type"`
	CreatedAt time.Time `json:"created_at"`
}

type MealPlans = []MealPlan

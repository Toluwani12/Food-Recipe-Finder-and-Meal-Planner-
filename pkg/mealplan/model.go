package mealplan

import (
	"time"
)

type MealPlan struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	DayOfWeek     DayOfWeek `json:"day_of_week"`
	MealType      MealType  `json:"meal_type"`
	RecipeID      string    `json:"recipe_id"`
	WeekStartDate time.Time `json:"week_start_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ImageURL      string    `json:"image_url"` // Placeholder image URL
}

type MealPlans = []MealPlan

type MealType string

const (
	Breakfast MealType = "breakfast"
	Lunch     MealType = "lunch"
	Dinner    MealType = "dinner"
)

type DayOfWeek string

const (
	Monday    DayOfWeek = "monday"
	Tuesday   DayOfWeek = "tuesday"
	Wednesday DayOfWeek = "wednesday"
	Thursday  DayOfWeek = "thursday"
	Friday    DayOfWeek = "friday"
	Saturday  DayOfWeek = "saturday"
	Sunday    DayOfWeek = "sunday"
)

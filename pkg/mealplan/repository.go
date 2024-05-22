package mealplan

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) save(mealPlans MealPlans) error {
	if len(mealPlans) == 0 {
		return nil
	}

	// Construct the base query
	query := `INSERT INTO meal_plans (user_id, day_of_week, meal_type, recipe_id, week_start_date) VALUES `
	values := []interface{}{}

	// Build the query and values slice dynamically
	for i, mealPlan := range mealPlans {
		// Add placeholders for each meal plan
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
		if i < len(mealPlans)-1 {
			query += ", "
		}
		values = append(values, mealPlan.UserID, mealPlan.DayOfWeek, mealPlan.MealType, mealPlan.RecipeID, mealPlan.WeekStartDate)
	}

	// Execute the query
	_, err := r.db.Exec(query, values...)
	return err
}

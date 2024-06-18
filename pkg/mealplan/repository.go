package mealplan

import (
	liberror "Food/internal/errors"
	"Food/pkg/recipe"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) save(ctx context.Context, mealPlans MealPlans) error {
	if len(mealPlans) == 0 {
		return nil
	}

	// Construct the base query
	query := `INSERT INTO meal_plans (user_id, day_of_week, meal_type, recipe_id, week_start_date, image_url)
			  VALUES `
	values := []interface{}{}

	// Build the query and values slice dynamically
	for i, mealPlan := range mealPlans {
		// Add placeholders for each meal plan
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
		if i < len(mealPlans)-1 {
			query += ", "
		}
		values = append(values, mealPlan.UserID, mealPlan.DayOfWeek, mealPlan.MealType, mealPlan.RecipeID, mealPlan.WeekStartDate, mealPlan.ImageURL)
	}

	// Add the ON CONFLICT clause to handle upsert
	query += ` ON CONFLICT (user_id, day_of_week, week_start_date, meal_type) DO UPDATE 
			   SET recipe_id = EXCLUDED.recipe_id, 
			       image_url = EXCLUDED.image_url`

	// Execute the query
	_, err := r.db.ExecContext(ctx, query, values...)
	return errors.Wrap(err, "ExecContext: failed to save meal plans")
}

type Ingredient struct {
	Name         string   `json:"name"`
	Quantity     string   `json:"quantity"`
	Alternatives []string `json:"alternatives"`
}

type Ingredients []Ingredient

type DetailedMealPlanDTO struct {
	ID           uuid.UUID      `db:"id"`
	Name         string         `json:"name" db:"name"`
	Description  string         `json:"description" db:"description"`
	CookingTime  string         `json:"cooking_time" db:"cooking_time"`
	Instructions pq.StringArray `json:"instructions" db:"instructions"`
	ImgUrl       string         `json:"img_url" db:"img_url"`
	Ingredients  Ingredients    `json:"ingredients" db:"-"`
}

type MealPlanPlaceholderDTO struct {
	DayOfWeek     DayOfWeek `json:"day_of_week" db:"day_of_week"`
	WeekStartDate time.Time `json:"week_start_date" db:"week_start_date"`
	ImageURL      string    `json:"image_url" db:"image_url"`
}

func (r *Repository) GetMealPlansForDay(userID string, dayOfWeek DayOfWeek, weekStartDate time.Time) ([]DetailedMealPlanDTO, error) {
	var recipes []DetailedMealPlanDTO
	query := `
		SELECT
			r.id,
			r.name,
			r.description,
			r.cooking_time,
			r.instructions,
			r.img_url
		FROM meal_plans mp
		JOIN recipes r ON mp.recipe_id = r.id
		WHERE mp.user_id = $1 AND mp.day_of_week = $2 AND mp.week_start_date = $3`
	err := r.db.Select(&recipes, query, userID, dayOfWeek, weekStartDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.New("No meal plans found for the specified day", http.StatusNotFound)
		}
		return nil, errors.Wrap(err, "Select: failed to get meal plans for day")
	}
	return recipes, nil
}

func (r *Repository) GetIngredientsForRecipes(recipeIDs []uuid.UUID) (map[uuid.UUID]Ingredients, error) {
	var ingredientsData []struct {
		RecipeID   uuid.UUID `db:"recipe_id"`
		Ingredient string    `db:"ingredient"`
		Quantity   string    `db:"quantity"`
	}

	query := `
		SELECT ri.recipe_id, i.name AS ingredient, ri.quantity
		FROM recipe_ingredients ri
		JOIN ingredients i ON ri.ingredient_id = i.id
		WHERE ri.recipe_id = ANY($1)`
	err := r.db.Select(&ingredientsData, query, pq.Array(recipeIDs))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.New("No ingredients found for the specified recipes", http.StatusNotFound)
		}
		return nil, errors.Wrap(err, "Select: failed to get ingredients for recipes")
	}

	// Map the ingredients by recipe ID
	ingredientsMap := make(map[uuid.UUID]Ingredients)
	for _, data := range ingredientsData {
		ingredient := Ingredient{
			Name:     data.Ingredient,
			Quantity: data.Quantity,
		}
		ingredientsMap[data.RecipeID] = append(ingredientsMap[data.RecipeID], ingredient)
	}
	return ingredientsMap, nil
}

func (r *Repository) GetMealPlanPlaceholders(userID string, weekStartDate time.Time) ([]MealPlanPlaceholderDTO, error) {
	var placeholders []MealPlanPlaceholderDTO
	query := `
        SELECT DISTINCT ON (day_of_week) day_of_week, week_start_date, image_url
        FROM meal_plans
        WHERE user_id = $1 AND week_start_date = $2
        ORDER BY day_of_week`
	err := r.db.Select(&placeholders, query, userID, weekStartDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.New("No meal plan placeholders found for the specified week", http.StatusNotFound)
		}
		return nil, errors.Wrap(err, "Select: failed to get meal plan placeholders")
	}
	return placeholders, nil
}

func (r *Repository) RecommendRecipes(ctx context.Context, userID uuid.UUID, limit int) ([]recipe.Recipe, error) {
	var recipes []recipe.Recipe

	query := `
    WITH recommended AS (
        SELECT recipe_id, similarity
        FROM recommend_recipes($1, $2)
    )
    SELECT
        r.id,
        r.name,
        r.description,
        r.cooking_time,
        r.instructions,
        r.img_url,
        rec.similarity
    FROM
        recommended rec
    JOIN recipes r ON rec.recipe_id = r.id
    ORDER BY
        rec.similarity DESC
    LIMIT $2;
    `

	err := r.db.SelectContext(ctx, &recipes, query, userID, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.New("No recommended recipes found", http.StatusNotFound)
		}
		return nil, errors.Wrap(err, "SelectContext: failed to recommend recipes")
	}

	return recipes, nil
}

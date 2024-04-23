package recipe

import (
	liberror "Food/internal/errors"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"net/http"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r Repository) get(ctx context.Context, id string) (*Recipe, error) {
	var recipe Recipe

	// Use Get to query and automatically scan the result into the struct
	err := r.db.GetContext(ctx, &recipe, "SELECT * FROM recipes WHERE id = $1", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.ErrNotFound
		}
	}

	return &recipe, errors.Wrap(err, "db.GetContext failed")
}

func (r Repository) getByName(ctx context.Context, name string) (*Recipe, error) {
	var recipe Recipe

	err := r.db.GetContext(ctx, &recipe, `SELECT r.name AS recipe_name, i.name AS ingredient_name, ri.quantity, alt.name AS alternative_name
								FROM recipes r
								JOIN recipe_ingredients ri ON r.id = ri.recipe_id
								JOIN ingredients i ON ri.ingredient_id = i.id
								LEFT JOIN ingredient_alternatives a ON i.id = a.ingredient_id
								LEFT JOIN ingredients alt ON a.alternative_id = alt.id
								WHERE r.name = 'Recipe Name';`, name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.ErrNotFound
		}
	}

	return &recipe, errors.Wrap(err, "db.GetContext failed")
}

func (r Repository) save(ctx context.Context, data Recipe) (*Recipe, error) {
	var recipe *Recipe
	err := r.db.GetContext(ctx, recipe, `BEGIN;
					-- Insert a recipe
					INSERT INTO recipes (name, description, cooking_time, instructions)
					VALUES ('Recipe Name', 'Recipe Description', 'Cooking Time', 'Cooking Instructions');
					
					-- Assume the recipe ID and ingredient ID are returned or known
					-- Insert an ingredient
					INSERT INTO ingredients (name)
					VALUES ('Ingredient Name');
					
					-- Link ingredient to the recipe
					INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity)
					VALUES ('Recipe ID', 'Ingredient ID', 'Quantity');
					
					-- Insert ingredient alternatives
					INSERT INTO ingredient_alternatives (ingredient_id, alternative_id)
					VALUES ('Ingredient ID', 'Alternative Ingredient ID');
					
					-- If all operations are correct
					COMMIT;
					
					-- If there is an error
					ROLLBACK;`, data.Name)

	if recipe != nil || !errors.Is(err, sql.ErrNoRows) {
		return nil, liberror.New("name already exist", http.StatusBadRequest)
	}
	if err != nil {
		return nil, errors.Wrap(err, "GetContext")
	}

	res, err := r.db.NamedExecContext(ctx, `INSERT INTO recipes (id, name, cooking_time, instructions, created_at, updated_at) VALUES (:id, :name, :cooking_time, :instructions, :created_at, :updated_at)`, data)
	if count, err := res.RowsAffected(); count != 1 {
		return nil, errors.Wrap(err, "RowsAffected")
	}

	return &data, errors.Wrap(err, "Db.NamedExecContext")
}

func (r Repository) delete(ctx context.Context, id string) (string, error) {
	_, err := r.db.ExecContext(ctx, `DELETE FROM recipes VALUES (:id)`, id)
	return id, errors.Wrap(err, "ExecContext")
}

func (r Repository) list(ctx context.Context) ([]Recipe, error) {
	var recipes []Recipe
	err := r.db.SelectContext(ctx, &recipes, "SELECT * FROM recipes")

	return recipes, errors.Wrap(err, "SelectContext")
}

func (r Repository) update(ctx context.Context, data Recipe) (*Recipe, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE recipes SET name = $1, cooking_time = $2, instructions = $3, updated_at = $4 WHERE id = $5", data.Name, data.CookingTime, data.Instructions, data.UpdatedAt, data.Id)
	if count, err := res.RowsAffected(); count != 1 {
		return nil, errors.Wrap(err, "RowsAffected")
	}

	return &data, errors.Wrap(err, "Db.NamedExecContext")
}

func (r Repository) findRecipes(ctx context.Context, ingredients []string) ([]Recipe, error) {
	var results []struct {
		Recipe          Recipe `db:"recipe"`
		IngredientsData string `db:"ingredients"`
	}

	query := buildQuery()

	// Execute the query
	err := r.db.SelectContext(ctx, &results, query, pq.Array(ingredients))
	if err != nil {
		return nil, errors.Wrap(err, "querying database")
	}

	// Convert results to []Recipe
	recipes := make([]Recipe, len(results))
	for idx, result := range results {
		recipes[idx] = result.Recipe
		// Parse the JSON ingredients
		if err := json.Unmarshal([]byte(result.IngredientsData), &recipes[idx].Ingredients); err != nil {
			return nil, errors.Wrap(err, "parsing ingredients")
		}
	}

	return recipes, nil
}

func buildQuery() string {
	return `WITH exact_matches AS (
    SELECT r.id, r.name, JSON_AGG(JSON_BUILD_OBJECT('id', i.id, 'name', i.name, 'alternatives', i.alternatives, 'quantity', i.quantity, 'created_at', i.created_at, 'updated_at', i.updated_at)) AS ingredients
    FROM recipes r
    JOIN ingredients i ON i.recipe_id = r.id
    GROUP BY r.id
    HAVING ARRAY_AGG(i.name ORDER BY i.name) = ARRAY[$1]::text[] -- Adjust to include any specific ordering or additional filtering
),
closest_matches AS (
    SELECT r.id, r.name, COUNT(*) AS ingredient_difference, JSON_AGG(JSON_BUILD_OBJECT('id', i.id, 'name', i.name, 'alternatives', i.alternatives, 'quantity', i.quantity, 'created_at', i.created_at, 'updated_at', i.updated_at)) AS ingredients
    FROM recipes r
    JOIN ingredients i ON i.recipe_id = r.id
    WHERE i.name <> ALL(ARRAY[$1]::text[]) -- This condition may need to be adjusted based on exact requirement
    GROUP BY r.id
    ORDER BY ingredient_difference ASC
    LIMIT 5
)
SELECT * FROM exact_matches
UNION ALL
SELECT * FROM closest_matches WHERE NOT EXISTS (SELECT 1 FROM exact_matches);`
}

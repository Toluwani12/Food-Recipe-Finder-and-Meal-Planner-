package recipe

import (
	liberror "Food/internal/errors"
	"Food/pkg"
	"Food/pkg/recipe/model"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) get(ctx context.Context, id string) (*Recipe, error) {
	var recipe Recipe

	err := r.db.GetContext(ctx, &recipe, "SELECT * FROM recipes WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.New("Recipe not found", http.StatusNotFound)
		}
		return nil, errors.Wrap(err, "db.GetContext failed")
	}

	return &recipe, nil
}

func (r *Repository) getByName(ctx context.Context, name string) (*Recipe, error) {
	var recipe Recipe

	err := r.db.GetContext(ctx, &recipe, `SELECT r.name AS recipe_name, i.name AS ingredient_name, ri.quantity, alt.name AS alternative_name
								FROM recipes r
								JOIN recipe_ingredients ri ON r.id = ri.recipe_id
								JOIN ingredients i ON ri.ingredient_id = i.id
								LEFT JOIN ingredient_alternatives a ON i.id = a.ingredient_id
								LEFT JOIN ingredients alt ON a.alternative_id = alt.id
								WHERE r.name = $1`, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.New("Recipe not found", http.StatusNotFound)
		}
		return nil, errors.Wrap(err, "db.GetContext failed")
	}

	return &recipe, nil
}

func (r *Repository) processRecipesAndIngredients(ctx context.Context, recipes model.Request) (map[string]bool, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "beginning transaction")
	}
	defer func() {
		if p := recover(); p != nil || err != nil {
			tx.Rollback()
			if p != nil {
				panic(p)
			}
		} else {
			err = tx.Commit()
			if err != nil {
				panic("commit failed")
			}
		}
	}()

	successMap, err := r.bulkUpsertRecipes(ctx, tx, recipes)
	if err != nil {
		return nil, err
	}

	newRecipes := make(model.Request, 0)
	for _, recipe := range recipes {
		if success, exists := successMap[recipe.Name]; exists && success {
			newRecipes = append(newRecipes, recipe)
		}
	}

	ingredientNames := make([]string, 0)
	for _, recipe := range newRecipes {
		for _, ingredient := range recipe.Ingredients {
			ingredientNames = append(ingredientNames, ingredient.Name)
		}
	}

	if len(ingredientNames) > 0 {
		ingredientIDs, err := r.bulkUpsertIngredients(ctx, tx, ingredientNames)
		if err != nil {
			return nil, err
		}

		if err = r.linkIngredients(ctx, tx, ingredientIDs, newRecipes); err != nil {
			return nil, err
		}
	}

	return successMap, nil
}

func (r *Repository) bulkUpsertRecipes(ctx context.Context, tx *sqlx.Tx, recipes model.Request) (map[string]bool, error) {
	existingRecipes, err := checkForDuplicateRecipes(tx, recipes)
	if err != nil {
		return nil, errors.Wrap(err, "checking for duplicate recipes")
	}

	var insertValues []string
	var insertArgs []interface{}
	index := 1
	for _, recipe := range recipes {
		if _, exists := existingRecipes[strings.TrimSpace(recipe.Name)]; !exists {
			insertValues = append(insertValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", index, index+1, index+2, index+3, index+4, index+5))
			insertArgs = append(insertArgs, recipe.ID, recipe.Name, recipe.Description, recipe.ImgUrl, recipe.CookingTime, pq.Array(recipe.Instructions))
			index += 6
		}
	}

	results := make(map[string]bool)

	if len(insertValues) > 0 {
		insertQuery := "INSERT INTO recipes (id, name, description, img_url, cooking_time, instructions) VALUES " +
			strings.Join(insertValues, ", ") + " RETURNING name"
		rows, err := tx.QueryContext(ctx, insertQuery, insertArgs...)
		if err != nil {
			return nil, errors.Wrap(err, "tx.QueryContext failed for insertQuery in bulkUpsertRecipes")
		}
		defer rows.Close()

		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				return nil, errors.Wrap(err, "row.scan failed for inserted recipes")
			}
			results[name] = true
		}

		if err := rows.Err(); err != nil {
			return nil, errors.Wrap(err, "rows.Err failed for inserted recipes")
		}
	}

	for name := range existingRecipes {
		results[name] = false
	}

	return results, nil
}

func checkForDuplicateRecipes(tx *sqlx.Tx, recipes model.Request) (map[string]string, error) {
	names := make([]interface{}, len(recipes))
	for i, recipe := range recipes {
		names[i] = recipe.Name
	}

	query, args, err := sqlx.In("SELECT name, id FROM recipes WHERE name IN (?)", names)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx.In failed for checking duplicate recipes")
	}
	query = tx.Rebind(query)

	rows, err := tx.Queryx(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "tx.Queryx failed for checking duplicate recipes")
	}
	defer rows.Close()

	existingRecipes := make(map[string]string)
	for rows.Next() {
		var name, id string
		if err := rows.Scan(&name, &id); err != nil {
			return nil, errors.Wrap(err, "rows.Scan failed for duplicate recipes")
		}
		existingRecipes[name] = id
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err failed for duplicate recipes")
	}

	return existingRecipes, nil
}

func (r *Repository) linkIngredients(ctx context.Context, tx *sqlx.Tx, ingredientIDs map[string]string, recipes model.Request) error {
	linkValues := []string{}
	linkArgs := []interface{}{}
	index := 1

	for _, recipe := range recipes {
		for _, ingredient := range recipe.Ingredients {
			if ingredientID, ok := ingredientIDs[ingredient.Name]; ok {
				linkValues = append(linkValues, fmt.Sprintf("($%d, $%d, $%d)", index, index+1, index+2))
				linkArgs = append(linkArgs, recipe.ID, ingredientID, ingredient.Quantity)
				index += 3
			}
		}
	}

	if len(linkValues) > 0 {
		linkQuery := "INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity) VALUES " +
			strings.Join(linkValues, ", ") + " ON CONFLICT (recipe_id, ingredient_id) DO NOTHING"
		if _, err := tx.ExecContext(ctx, linkQuery, linkArgs...); err != nil {
			return errors.Wrap(err, "tx.ExecContext failed for linkQuery in linkIngredients")
		}
	}

	return nil
}

func (r *Repository) bulkUpsertIngredients(ctx context.Context, tx *sqlx.Tx, ingredientNames []string) (map[string]string, error) {
	if len(ingredientNames) == 0 {
		return nil, nil
	}

	ingredientIDs := make(map[string]string)

	query, args, err := sqlx.In("SELECT id, name FROM ingredients WHERE name IN (?)", ingredientNames)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx.In failed for querying existing ingredients")
	}
	query = tx.Rebind(query)
	rows, err := tx.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "tx.QueryxContext failed for querying existing ingredients")
	}
	defer rows.Close()

	existingIngredients := make(map[string]bool)
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, errors.Wrap(err, "rows.Scan failed for existing ingredients")
		}
		ingredientIDs[name] = id
		existingIngredients[name] = true
	}

	newIngredients := []string{}
	for _, name := range ingredientNames {
		if !existingIngredients[name] {
			newIngredients = append(newIngredients, name)
		}
	}

	if len(newIngredients) > 0 {
		values := make([]string, 0, len(newIngredients))
		insertArgs := make([]interface{}, 0, len(newIngredients))
		for i, name := range newIngredients {
			values = append(values, fmt.Sprintf("($%d)", i+1))
			insertArgs = append(insertArgs, name)
		}
		insertQuery := "INSERT INTO ingredients (name) VALUES " + strings.Join(values, ", ") + " RETURNING id, name"
		insertStmt, err := tx.PrepareContext(ctx, insertQuery)
		if err != nil {
			return nil, errors.Wrap(err, "tx.PrepareContext failed for new ingredients")
		}
		defer insertStmt.Close()

		insertRows, err := insertStmt.QueryContext(ctx, insertArgs...)
		if err != nil {
			return nil, errors.Wrap(err, "insertStmt.QueryContext failed for new ingredients")
		}
		defer insertRows.Close()

		for insertRows.Next() {
			var id, name string
			if err := insertRows.Scan(&id, &name); err != nil {
				return nil, errors.Wrap(err, "insertRows.Scan failed for new ingredients")
			}
			ingredientIDs[name] = id
		}
	}

	return ingredientIDs, nil
}

func (r *Repository) delete(ctx context.Context, id string) (string, error) {
	_, err := r.db.ExecContext(ctx, `DELETE FROM recipes WHERE id = $1`, id)
	if err != nil {
		return id, errors.Wrap(err, "ExecContext")
	}
	return id, nil
}

func (r *Repository) list(ctx context.Context, userID string, recipeName string) ([]ListResponse, error) {
	var recipes []ListResponse
	var query string
	var args []interface{}

	if userID != "" {
		query = `
			SELECT r.id, r.name, COALESCE(l.liked, false) AS liked
			FROM recipes r
			LEFT JOIN (SELECT recipe_id, true AS liked FROM likes WHERE user_id = $1) l ON r.id = l.recipe_id
			WHERE r.name ILIKE $2
			ORDER BY r.name`
		args = append(args, userID, "%"+recipeName+"%")
	} else {
		query = `
			SELECT r.id, r.name, false AS liked
			FROM recipes r
			WHERE r.name ILIKE $1
			ORDER BY r.name`
		args = append(args, "%"+recipeName+"%")
	}

	err := r.db.SelectContext(ctx, &recipes, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "SelectContext")
	}
	return recipes, nil
}

func (r *Repository) update(ctx context.Context, data Recipe) (*Recipe, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE recipes SET name = $1, cooking_time = $2, instructions = $3, updated_at = $4 WHERE id = $5", data.Name, data.CookingTime, data.Instructions, data.UpdatedAt, data.Id)
	if err != nil {
		return nil, errors.Wrap(err, "ExecContext failed")
	}
	if count, err := res.RowsAffected(); err != nil {
		return nil, errors.Wrap(err, "RowsAffected")
	} else if count != 1 {
		return nil, errors.New("no rows affected")
	}

	return &data, nil
}

type Ingredient struct {
	Name         string   `json:"name"`
	Quantity     string   `json:"quantity"`
	Alternatives []string `json:"alternatives"`
}

type Ingredients []Ingredient

func (i *Ingredients) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, i)
}

type ResponseData struct {
	ID           uuid.UUID      `db:"id"`
	Name         string         `json:"name" db:"name"`
	Description  string         `json:"description" db:"description"`
	CookingTime  string         `json:"cooking_time" db:"cooking_time"`
	Instructions pq.StringArray `json:"instructions" db:"instructions"`
	ImgUrl       string         `json:"img_url" db:"img_url"`
	Ingredients  []Ingredient   `json:"ingredients" db:"ingredients"`
	Diff         int            `json:"diff" db:"diff"`
	FoodHealth   string         `json:"food_health,omitempty" db:"food_health"`
	FoodClass    string         `json:"food_class,omitempty" db:"food_class"`
	Region       string         `json:"region,omitempty" db:"region"`
	SpiceLevel   string         `json:"spice_level,omitempty" db:"spice_level"`
	Type         string         `json:"type,omitempty" db:"type"`
	Liked        bool           `json:"liked" db:"liked"`
}

func (r *Repository) search(ctx context.Context, ingredients []string, queryParams url.Values, userID string) ([]ResponseData, *pkg.Pagination, error) {
	page, pageSize, err := pkg.ParsePaginationParams(queryParams)
	if err != nil {
		return nil, nil, err
	}

	if (len(ingredients) == 0) || strings.TrimSpace(ingredients[0]) == "" {
		matches, pagination, err := r.findAllRecipes(ctx, queryParams, userID, page, pageSize)
		if err != nil {
			return nil, nil, errors.Wrap(err, "find all recipes failed")
		}
		return r.getIngredientsForRecipes(ctx, matches, pagination)
	}

	matches, pagination, err := r.getClosestRecipeWithDetails(ctx, ingredients, userID, page, pageSize)
	if err != nil {
		return nil, nil, errors.Wrap(err, "get closest recipe with details failed")
	}

	return r.getIngredientsForRecipes(ctx, matches, pagination)
}

func (r *Repository) findAllRecipes(ctx context.Context, queryParams url.Values, userID string, page, pageSize int) ([]ResponseData, *pkg.Pagination, error) {
	var recipes []ResponseData

	if queryParams == nil {
		query := `
			SELECT r.id, r.name, r.description, r.cooking_time, r.instructions, r.img_url, false AS liked
			FROM recipes r
			ORDER BY r.name`
		err := r.db.SelectContext(ctx, &recipes, query)
		if err != nil {
			return nil, nil, errors.Wrap(err, "db.SelectContext failed")
		}
		return recipes, nil, nil
	}

	var totalItems int
	countQuery := `SELECT COUNT(*) FROM recipes`
	err := r.db.GetContext(ctx, &totalItems, countQuery)
	if err != nil {
		return nil, nil, errors.Wrap(err, "db.GetContext failed")
	}

	var query string
	if userID == "" {
		query = `
			SELECT r.id, r.name, r.description, r.cooking_time, r.instructions, r.img_url, false AS liked
			FROM recipes r
			ORDER BY r.name`
		query = pkg.ApplyToQuery(query, page, pageSize)
		err = r.db.SelectContext(ctx, &recipes, query)
	} else {
		query = `
			SELECT r.id, r.name, r.description, r.cooking_time, r.instructions, r.img_url, COALESCE(l.liked, false) AS liked
			FROM recipes r
			LEFT JOIN (SELECT recipe_id, true AS liked FROM likes WHERE user_id = $1) l ON r.id = l.recipe_id
			ORDER BY r.name`
		query = pkg.ApplyToQuery(query, page, pageSize)
		err = r.db.SelectContext(ctx, &recipes, query, userID)
	}
	if err != nil {
		return nil, nil, errors.Wrap(err, "db.SelectContext failed")
	}

	pagination := pkg.NewPagination(page, pageSize, totalItems)
	return recipes, pagination, nil
}

func (r *Repository) getClosestRecipeWithDetails(ctx context.Context, ingredients []string, userID string, page, pageSize int) ([]ResponseData, *pkg.Pagination, error) {
	ingredientsArray := pq.Array(ingredients)

	var totalItems int
	countQuery := `
    SELECT COUNT(*)
    FROM find_recipes_by_jaccard_similarity($1::text[])
    `
	err := r.db.GetContext(ctx, &totalItems, countQuery, ingredientsArray)
	if err != nil {
		return nil, nil, errors.Wrap(err, "db.GetContext failed")
	}

	query := `
    WITH closest_recipes AS (
        SELECT recipe_id, similarity_score
        FROM find_recipes_by_jaccard_similarity($1::text[])
        ORDER BY similarity_score DESC
        LIMIT $2 OFFSET $3
    )
    SELECT 
        r.id, r.name, r.description, r.cooking_time, r.instructions, r.img_url, 
        %s AS liked
    FROM closest_recipes cr
    JOIN recipes r ON r.id = cr.recipe_id
    %s
    `

	likedQuery := "false"
	joinQuery := ""
	if userID != "" {
		likedQuery = "COALESCE(l.liked, false)"
		joinQuery = "LEFT JOIN (SELECT recipe_id, true AS liked FROM likes WHERE user_id = $4) l ON r.id = l.recipe_id"
	}

	finalQuery := fmt.Sprintf(query, likedQuery, joinQuery)

	var recipes []ResponseData
	if userID == "" {
		err = r.db.SelectContext(ctx, &recipes, finalQuery, ingredientsArray, pageSize, (page-1)*pageSize)
	} else {
		err = r.db.SelectContext(ctx, &recipes, finalQuery, ingredientsArray, pageSize, (page-1)*pageSize, userID)
	}
	if err != nil {
		return nil, nil, errors.Wrap(err, "db.SelectContext failed")
	}

	pagination := pkg.NewPagination(page, pageSize, totalItems)
	return recipes, pagination, nil
}

func (r *Repository) getIngredientsForRecipes(ctx context.Context, recipes []ResponseData, pagination *pkg.Pagination) ([]ResponseData, *pkg.Pagination, error) {
	if len(recipes) == 0 {
		return recipes, pagination, nil
	}

	recipeIDs := extractRecipeIDs(recipes)

	ingredientsData, err := r.fetchIngredients(ctx, recipeIDs)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fetching ingredients")
	}

	mapIngredientsToRecipes(recipes, ingredientsData)

	return recipes, pagination, nil
}

func extractRecipeIDs(recipes []ResponseData) []uuid.UUID {
	recipeIDs := make([]uuid.UUID, len(recipes))
	for i, recipe := range recipes {
		recipeIDs[i] = recipe.ID
	}
	return recipeIDs
}

func (r *Repository) fetchIngredients(ctx context.Context, recipeIDs []uuid.UUID) ([]struct {
	RecipeID   uuid.UUID `db:"recipe_id"`
	Ingredient string    `db:"ingredient"`
	Quantity   string    `db:"quantity"`
}, error) {
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
	err := r.db.SelectContext(ctx, &ingredientsData, query, pq.Array(recipeIDs))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.New("No ingredients found for the specified recipes", http.StatusNotFound)
		}
		return nil, errors.Wrap(err, "fetchIngredients: db.SelectContext failed")
	}

	return ingredientsData, nil
}

func mapIngredientsToRecipes(recipes []ResponseData, ingredientsData []struct {
	RecipeID   uuid.UUID `db:"recipe_id"`
	Ingredient string    `db:"ingredient"`
	Quantity   string    `db:"quantity"`
}) {
	ingredientsMap := make(map[uuid.UUID]Ingredients)
	for _, data := range ingredientsData {
		ingredientsMap[data.RecipeID] = append(ingredientsMap[data.RecipeID], Ingredient{
			Name:     data.Ingredient,
			Quantity: data.Quantity,
		})
	}

	for i, recipe := range recipes {
		if ingredients, ok := ingredientsMap[recipe.ID]; ok {
			recipes[i].Ingredients = ingredients
		}
	}
}

func (r *Repository) fetchLikes(ctx context.Context, recipeIDs []uuid.UUID) ([]struct {
	RecipeID uuid.UUID `db:"recipe_id"`
	Liked    bool      `db:"liked"`
}, error) {
	query, args, err := sqlx.In(`
        SELECT recipe_id, true AS liked
        FROM likes
        WHERE recipe_id IN (?)`, recipeIDs)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx.In failed for fetching likes")
	}

	query = r.db.Rebind(query)

	var results []struct {
		RecipeID uuid.UUID `db:"recipe_id"`
		Liked    bool      `db:"liked"`
	}

	err = r.db.SelectContext(ctx, &results, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, liberror.New("No likes found for the specified recipes", http.StatusNotFound)
		}
		return nil, errors.Wrap(err, "fetchLikes: db.SelectContext failed")
	}

	return results, nil
}

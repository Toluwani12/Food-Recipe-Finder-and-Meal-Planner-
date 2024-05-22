package recipe

import (
	liberror "Food/internal/errors"
	"Food/pkg"
	"Food/pkg/ingredient"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"net/url"
	"sort"
	"strconv"
	"strings"
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

func (r Repository) processRecipesAndIngredients(ctx context.Context, recipes Request) (map[string]bool, error) {
	// Start a transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "beginning transaction")
	}

	// Defer a rollback in case of an error
	defer func() {
		if p := recover(); p != nil || err != nil {
			tx.Rollback()
			if p != nil {
				panic(p) // re-throw panic after Rollback
			}
		} else {
			err = tx.Commit()
			if err != nil {
				panic("commit failed")
			}
		}
	}()

	// Upsert recipes and get which were successfully added
	successMap, err := r.bulkUpsertRecipes(ctx, tx, recipes)
	if err != nil {
		return nil, err
	}

	// Filter recipes to include only those successfully added
	newRecipes := make(Request, 0)
	for _, recipe := range recipes {
		if success, exists := successMap[recipe.Name]; exists && success {
			newRecipes = append(newRecipes, recipe)
		}
	}

	// Extract all ingredient names from the new recipes only
	ingredientNames := make([]string, 0)
	for _, recipe := range newRecipes {
		for _, ingredient := range recipe.Ingredients {
			ingredientNames = append(ingredientNames, ingredient.Name)
		}
	}

	if len(ingredientNames) > 0 {

		// Upsert ingredients and get their IDs
		ingredientIDs, err := r.bulkUpsertIngredients(ctx, tx, ingredientNames)
		if err != nil {
			return nil, err
		}

		// Link ingredients only for new recipes
		if err = r.linkIngredients(ctx, tx, ingredientIDs, newRecipes); err != nil {
			return nil, err
		}
	}

	return successMap, nil
}

// bulkUpsertRecipes processes a batch of recipes, inserts new ones in a batch, and reports duplicates with boolean flags.
func (r Repository) bulkUpsertRecipes(ctx context.Context, tx *sqlx.Tx, recipes Request) (map[string]bool, error) {
	// First, check for duplicates
	existingRecipes, err := checkForDuplicateRecipes(tx, recipes)
	if err != nil {
		return nil, err
	}

	// Prepare batch insert for new recipes
	var insertValues []string
	var insertArgs []interface{}
	index := 1
	for _, recipe := range recipes {
		if _, exists := existingRecipes[recipe.Name]; !exists {
			insertValues = append(insertValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", index, index+1, index+2, index+3, index+4, index+5))
			insertArgs = append(insertArgs, recipe.ID, recipe.Name, recipe.Description, recipe.ImgUrl, recipe.CookingTime, recipe.Instructions)
			index += 6
		}
	}

	// Initialize results map
	results := make(map[string]bool)

	// Execute batch insert if there are new recipes to add
	if len(insertValues) > 0 {
		insertQuery := "INSERT INTO recipes (id, name, description, img_url, cooking_time, instructions) VALUES " +
			strings.Join(insertValues, ", ") + " RETURNING name"
		rows, err := tx.QueryContext(ctx, insertQuery, insertArgs...)
		if err != nil {
			return nil, fmt.Errorf("error executing batch insert for new recipes: %v", err)
		}
		defer rows.Close()

		// Process inserted recipes
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				return nil, fmt.Errorf("error scanning inserted recipes: %v", err)
			}
			results[name] = true
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error reading results for inserted recipes: %v", err)
		}
	}

	// Mark duplicates as false in the results
	for name := range existingRecipes {
		results[name] = false
	}

	return results, nil
}

// checkForDuplicateRecipes checks the database for existing recipes with the same names.
func checkForDuplicateRecipes(tx *sqlx.Tx, recipes Request) (map[string]string, error) {
	names := make([]interface{}, len(recipes))
	for i, recipe := range recipes {
		names[i] = recipe.Name
	}

	query, args, err := sqlx.In("SELECT name, id FROM recipes WHERE name IN (?)", names)
	if err != nil {
		return nil, fmt.Errorf("error preparing query to check duplicates: %v", err)
	}
	query = tx.Rebind(query)

	rows, err := tx.Queryx(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying existing recipes: %v", err)
	}
	defer rows.Close()

	existingRecipes := make(map[string]string)
	for rows.Next() {
		var name, id string
		if err := rows.Scan(&name, &id); err != nil {
			return nil, fmt.Errorf("error scanning existing recipes: %v", err)
		}
		existingRecipes[name] = id
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading results for duplicate recipes: %v", err)
	}

	return existingRecipes, nil
}

func (r Repository) linkIngredients(ctx context.Context, tx *sqlx.Tx, ingredientIDs map[string]string, recipes Request) error {

	linkValues := []string{}
	linkArgs := []interface{}{}
	index := 1

	// Prepare insert values and arguments from the recipes and their ingredients
	for _, recipe := range recipes {
		for _, ingredient := range recipe.Ingredients {
			if ingredientID, ok := ingredientIDs[ingredient.Name]; ok {
				// Dynamically create placeholders for each ingredient link
				linkValues = append(linkValues, fmt.Sprintf("($%d, $%d, $%d)", index, index+1, index+2))
				linkArgs = append(linkArgs, recipe.ID, ingredientID, ingredient.Quantity)
				index += 3
			}
		}
	}

	// Perform the batch insert only if there are ingredients to link
	if len(linkValues) > 0 {
		linkQuery := "INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity) VALUES " +
			strings.Join(linkValues, ", ") + " ON CONFLICT (recipe_id, ingredient_id) DO NOTHING"

		// Execute the query with all arguments
		if _, err := tx.ExecContext(ctx, linkQuery, linkArgs...); err != nil {
			return errors.Wrap(err, "executing batch insert for recipe_ingredients with conflict handling")
		}
	}

	return nil
}

func (r Repository) bulkUpsertIngredients(ctx context.Context, tx *sqlx.Tx, ingredientNames []string) (map[string]string, error) {
	if len(ingredientNames) == 0 {
		return nil, nil
	}

	ingredientIDs := make(map[string]string)

	// Step 1: Query existing ingredients to avoid duplicates
	query, args, err := sqlx.In("SELECT id, name FROM ingredients WHERE name IN (?)", ingredientNames)
	if err != nil {
		return nil, errors.Wrap(err, "preparing query for existing ingredients")
	}
	query = tx.Rebind(query)
	rows, err := tx.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying existing ingredients")
	}
	defer rows.Close()

	existingIngredients := make(map[string]bool)
	for rows.Next() {
		var id string
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, errors.Wrap(err, "scanning ingredients")
		}
		ingredientIDs[name] = id
		existingIngredients[name] = true
	}

	// Step 2: Identify which ingredients need to be inserted
	newIngredients := []string{}
	for _, name := range ingredientNames {
		if !existingIngredients[name] {
			newIngredients = append(newIngredients, name)
		}
	}

	// Step 3: Perform bulk insert for new ingredients
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
			return nil, errors.Wrap(err, "preparing insert for new ingredients")
		}
		defer insertStmt.Close()

		insertRows, err := insertStmt.QueryContext(ctx, insertArgs...)
		if err != nil {
			return nil, errors.Wrap(err, "executing insert for new ingredients")
		}
		defer insertRows.Close()

		for insertRows.Next() {
			var id string
			var name string
			if err := insertRows.Scan(&id, &name); err != nil {
				return nil, errors.Wrap(err, "scanning inserted ingredients")
			}
			ingredientIDs[name] = id
		}
	}

	return ingredientIDs, nil
}

func (r Repository) delete(ctx context.Context, id string) (string, error) {
	_, err := r.db.ExecContext(ctx, `DELETE FROM recipes VALUES (:id)`, id)
	return id, errors.Wrap(err, "ExecContext")
}

func (r Repository) list(ctx context.Context) ([]ListResponse, error) {
	var recipes []ListResponse
	err := r.db.SelectContext(ctx, &recipes, "SELECT id, name, img_url FROM recipes")

	return recipes, errors.Wrap(err, "SelectContext")
}

func (r Repository) update(ctx context.Context, data Recipe) (*Recipe, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE recipes SET name = $1, cooking_time = $2, instructions = $3, updated_at = $4 WHERE id = $5", data.Name, data.CookingTime, data.Instructions, data.UpdatedAt, data.Id)
	if count, err := res.RowsAffected(); count != 1 {
		return nil, errors.Wrap(err, "RowsAffected")
	}

	return &data, errors.Wrap(err, "Db.NamedExecContext")
}

func (r Repository) search(ctx context.Context, ingredients []string, queryParams url.Values) (Request, *pkg.Pagination, error) {
	page, pageSize, err := pkg.ParsePaginationParams(queryParams)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing pagination parameters: %w", err)
	}

	if len(ingredients) == 0 {
		// No ingredients provided, return all recipes
		return r.findAllRecipes(ctx, page, pageSize)
	}

	recipes, pagination, err := r.findMatches(ctx, ingredients, page, pageSize)
	if err != nil {
		return nil, nil, fmt.Errorf("searching for matches: %w", err)
	}

	return recipes, pagination, nil
}

func (r Repository) findAllRecipes(ctx context.Context, page, pageSize int) (Request, *pkg.Pagination, error) {
	query := `
	   SELECT r.id, r.name, r.description, r.cooking_time, r.instructions, r.img_url,
	          ri.ingredient_id, i.name as ingredient_name, ri.quantity, a.name as alternative_name
	   FROM recipes r
	   JOIN recipe_ingredients ri ON r.id = ri.recipe_id
	   JOIN ingredients i ON i.id = ri.ingredient_id
	   LEFT JOIN ingredient_alternatives ia ON ia.ingredient_id = i.id
	   LEFT JOIN ingredients a ON ia.alternative_id = a.id
	   ORDER BY r.id
	   LIMIT $1 OFFSET $2;
	`

	offset := (page - 1) * pageSize
	rows, err := r.db.QueryxContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("executing find all recipes query: %w", err)
	}
	defer rows.Close()

	var recipes Request

	// Function to find a recipe by ID in the slice
	findRecipe := func(recipes Request, id string) (int, bool) {
		for i, r := range recipes {
			if r.ID.String() == id {
				return i, true
			}
		}
		return -1, false
	}

	// Process rows from the query
	for rows.Next() {
		var r RequestData
		var i ingredient.Request
		var alternativeName sql.NullString
		err = rows.Scan(&r.ID, &r.Name, &r.Description, &r.CookingTime, &r.Instructions, &r.ImgUrl, &i.ID, &i.Name, &i.Quantity, &alternativeName)
		if err != nil {
			return nil, nil, fmt.Errorf("scanning rows: %w", err)
		}

		index, found := findRecipe(recipes, r.ID.String())
		if !found {
			r.Ingredients = []ingredient.Request{}
			recipes = append(recipes, r)
			index = len(recipes) - 1
		}

		if alternativeName.Valid {
			i.Alternatives = append(i.Alternatives, alternativeName.String)
		}
		recipes[index].Ingredients = append(recipes[index].Ingredients, i)
	}

	totalItemsQuery := "SELECT COUNT(*) FROM recipes"
	var totalItems int
	if err := r.db.GetContext(ctx, &totalItems, totalItemsQuery); err != nil {
		return nil, nil, fmt.Errorf("counting total items: %w", err)
	}

	pagination := pkg.NewPagination(page, pageSize, totalItems)
	return recipes, pagination, nil
}

func (r Repository) findMatches(ctx context.Context, ingredients []string, page, pageSize int) (Request, *pkg.Pagination, error) {
	exactQuery := `
	   SELECT r.id, r.name, r.description, r.cooking_time, r.instructions, r.img_url,
	          ri.ingredient_id, i.name as ingredient_name, ri.quantity, a.name as alternative_name
	   FROM recipes r
	   JOIN recipe_ingredients ri ON r.id = ri.recipe_id
	   JOIN ingredients i ON i.id = ri.ingredient_id
	   LEFT JOIN ingredient_alternatives ia ON ia.ingredient_id = i.id
	   LEFT JOIN ingredients a ON ia.alternative_id = a.id
	   WHERE r.id IN (
	       SELECT recipe_id
	       FROM recipe_ingredients ri
	       JOIN ingredients i ON i.id = ri.ingredient_id
	       GROUP BY recipe_id
	       HAVING ARRAY_AGG(i.name ORDER BY i.name) @> $1::text[]
	          AND COUNT(i.name) = $2
	   )
	   ORDER BY r.id
	   LIMIT $3 OFFSET $4;
	`

	closestMatchQuery := `
        SELECT r.id, r.name, r.description, r.cooking_time, r.instructions, r.img_url,
               ri.ingredient_id, i.name as ingredient_name, ri.quantity, a.name as alternative_name,
               COUNT(i.name) FILTER (WHERE i.name = ANY($1)) AS matching_ingredients
        FROM recipes r
        JOIN recipe_ingredients ri ON r.id = ri.recipe_id
        JOIN ingredients i ON i.id = ri.ingredient_id
        LEFT JOIN ingredient_alternatives ia ON ia.ingredient_id = i.id
        LEFT JOIN ingredients a ON ia.alternative_id = a.id
        GROUP BY r.id, ri.ingredient_id, i.name, ri.quantity, a.name
        ORDER BY matching_ingredients DESC, r.id
        LIMIT $2 OFFSET $3;
    `

	offset := (page - 1) * pageSize

	var rows *sqlx.Rows
	var err error

	// Try exact match first
	rows, err = r.db.QueryxContext(ctx, exactQuery, pq.Array(ingredients), len(ingredients), pageSize, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("executing exact match query: %w", err)
	}
	defer rows.Close()

	var recipes Request

	// Function to find a recipe by ID in the slice
	findRecipe := func(recipes Request, id string) (int, bool) {
		for i, r := range recipes {
			if r.ID.String() == id {
				return i, true
			}
		}
		return -1, false
	}

	// Process rows from exact match query
	for rows.Next() {
		var r RequestData
		var i ingredient.Request
		var alternativeName sql.NullString
		err = rows.Scan(&r.ID, &r.Name, &r.Description, &r.CookingTime, &r.Instructions, &r.ImgUrl, &i.ID, &i.Name, &i.Quantity, &alternativeName)
		if err != nil {
			return nil, nil, fmt.Errorf("scanning rows: %w", err)
		}

		index, found := findRecipe(recipes, r.ID.String())
		if !found {
			r.Ingredients = []ingredient.Request{}
			recipes = append(recipes, r)
			index = len(recipes) - 1
		}

		if alternativeName.Valid {
			i.Alternatives = append(i.Alternatives, alternativeName.String)
		}
		recipes[index].Ingredients = append(recipes[index].Ingredients, i)
	}

	if len(recipes) > 0 {
		totalItemsQuery := `
		   SELECT COUNT(*) FROM (
		       SELECT r.id
		       FROM recipes r
		       JOIN recipe_ingredients ri ON r.id = ri.recipe_id
		       JOIN ingredients i ON i.id = ri.ingredient_id
		       WHERE r.id IN (
		           SELECT recipe_id
		           FROM recipe_ingredients ri
		           JOIN ingredients i ON i.id = ri.ingredient_id
		           GROUP BY recipe_id
		           HAVING ARRAY_AGG(i.name ORDER BY i.name) @> $1::text[]
		              AND COUNT(i.name) = $2
		       )
		   ) AS total
		`
		var totalItems int
		if err := r.db.GetContext(ctx, &totalItems, totalItemsQuery, pq.Array(ingredients), len(ingredients)); err != nil {
			return nil, nil, fmt.Errorf("counting total items: %w", err)
		}

		pagination := pkg.NewPagination(page, pageSize, totalItems)
		return recipes, pagination, nil
	}

	// If no exact matches are found, find the closest match
	rows, err = r.db.QueryxContext(ctx, closestMatchQuery, pq.Array(ingredients), pageSize, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("executing closest match query: %w", err)
	}
	defer rows.Close()

	// Process rows from closest match query
	for rows.Next() {
		var r RequestData
		var i ingredient.Request
		var alternativeName sql.NullString
		var matchingIngredients int
		err = rows.Scan(&r.ID, &r.Name, &r.Description, &r.CookingTime, &r.Instructions, &r.ImgUrl, &i.ID, &i.Name, &i.Quantity, &alternativeName, &matchingIngredients)
		if err != nil {
			return nil, nil, fmt.Errorf("scanning rows: %w", err)
		}

		index, found := findRecipe(recipes, r.ID.String())
		if !found {
			r.Ingredients = []ingredient.Request{}
			recipes = append(recipes, r)
			index = len(recipes) - 1
		}

		if alternativeName.Valid {
			i.Alternatives = append(i.Alternatives, alternativeName.String)
		}
		recipes[index].Ingredients = append(recipes[index].Ingredients, i)
		recipes[index].Diff = matchingIngredients
	}

	// Sort recipes based on the number of matching ingredients in descending order
	sort.SliceStable(recipes, func(i, j int) bool {
		return recipes[i].Diff > recipes[j].Diff
	})

	totalItemsQuery := `
	   SELECT COUNT(*) FROM (
	       SELECT r.id
	       FROM recipes r
	       JOIN recipe_ingredients ri ON r.id = ri.recipe_id
	       JOIN ingredients i ON i.id = ri.ingredient_id
	       WHERE r.id IN (
	           SELECT recipe_id
	           FROM recipe_ingredients ri
	           JOIN ingredients i ON i.id = ri.ingredient_id
	           GROUP BY recipe_id
	       )
	   ) AS total
	`
	var totalItems int
	if err := r.db.GetContext(ctx, &totalItems, totalItemsQuery, pq.Array(ingredients)); err != nil {
		return nil, nil, fmt.Errorf("counting total items: %w", err)
	}

	pagination := pkg.NewPagination(page, pageSize, totalItems)
	return recipes, pagination, nil
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

type RequestData3 struct {
	ID           uuid.UUID   `db:"id"`
	Name         string      `json:"name" db:"name"`
	Description  string      `json:"description" db:"description"`
	CookingTime  string      `json:"cooking_time" db:"cooking_time"`
	Instructions []string    `json:"instructions" db:"instructions"`
	ImgUrl       string      `json:"img_url" db:"img_url"`
	Ingredients  Ingredients `json:"ingredients" db:"-"`
}

func (r Repository) search2(ctx context.Context, ingredients []string, queryParams url.Values) ([]RequestData3, *pkg.Pagination, error) {
	if len(ingredients) == 0 {
		return r.findAllRecipes2(ctx, queryParams)
	}

	exactMatches, exactMatchPagination, err := r.findExactMatch(ctx, ingredients, queryParams)
	if err != nil {
		return nil, nil, err
	}

	if len(exactMatches) > 0 {
		return r.getIngredientsForRecipes(ctx, exactMatches, exactMatchPagination)
	}

	closestMatches, closestMatchPagination, err := r.findClosestMatches(ctx, ingredients, queryParams)
	if err != nil {
		return nil, nil, err
	}

	return r.getIngredientsForRecipes(ctx, closestMatches, closestMatchPagination)
}

func (r Repository) findAllRecipes2(ctx context.Context, queryParams url.Values) ([]RequestData3, *pkg.Pagination, error) {
	page, pageSize := getPagination(queryParams)

	var recipes []RequestData3
	query := `
	SELECT r.id, r.name, r.description, r.cooking_time, r.instructions, r.img_url
	FROM recipes r
	ORDER BY r.name
	LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &recipes, query, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, nil, err
	}

	return r.getIngredientsForRecipes(ctx, recipes, &pkg.Pagination{Page: page, PageSize: pageSize})
}

func (r Repository) findExactMatch(ctx context.Context, ingredients []string, queryParams url.Values) ([]RequestData3, *pkg.Pagination, error) {
	page, pageSize := getPagination(queryParams)

	var recipes []RequestData3
	ingredientList := strings.Join(ingredients, "|")
	query := `
	SELECT r.id, r.name, r.description, r.cooking_time, r.instructions, r.img_url
	FROM recipes r
	JOIN recipe_ingredients ri ON r.id = ri.recipe_id
	JOIN ingredients i ON ri.ingredient_id = i.id
	WHERE i.name ~* $1
	GROUP BY r.id
	HAVING COUNT(DISTINCT i.name) = $2
	ORDER BY r.name
	LIMIT $3 OFFSET $4`
	err := r.db.SelectContext(ctx, &recipes, query, ingredientList, len(ingredients), pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, nil, err
	}

	return recipes, &pkg.Pagination{Page: page, PageSize: pageSize}, nil
}

func (r Repository) findClosestMatches(ctx context.Context, ingredients []string, queryParams url.Values) ([]RequestData3, *pkg.Pagination, error) {
	page, pageSize := getPagination(queryParams)

	var recipes []RequestData3
	ingredientList := strings.Join(ingredients, "|")
	query := `
	SELECT r.id, r.name, r.description, r.cooking_time, r.instructions, r.img_url
	FROM recipes r
	JOIN recipe_ingredients ri ON r.id = ri.recipe_id
	JOIN ingredients i ON ri.ingredient_id = i.id
	WHERE i.name ~* $1
	GROUP BY r.id
	ORDER BY ABS(COUNT(DISTINCT i.name) - $2) ASC
	LIMIT $3 OFFSET $4`
	err := r.db.SelectContext(ctx, &recipes, query, ingredientList, len(ingredients), pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, nil, err
	}

	return recipes, &pkg.Pagination{Page: page, PageSize: pageSize}, nil
}

func (r Repository) getIngredientsForRecipes(ctx context.Context, recipes []RequestData3, pagination *pkg.Pagination) ([]RequestData3, *pkg.Pagination, error) {
	if len(recipes) == 0 {
		return recipes, pagination, nil
	}

	var recipeIDs []uuid.UUID
	for _, recipe := range recipes {
		recipeIDs = append(recipeIDs, recipe.ID)
	}

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
		return nil, nil, err
	}

	ingredientsMap := make(map[uuid.UUID]Ingredients)
	for _, ingredientData := range ingredientsData {
		ingredientsMap[ingredientData.RecipeID] = append(ingredientsMap[ingredientData.RecipeID], Ingredient{
			Name:     ingredientData.Ingredient,
			Quantity: ingredientData.Quantity,
		})
	}

	for i, recipe := range recipes {
		if ingredients, ok := ingredientsMap[recipe.ID]; ok {
			recipes[i].Ingredients = ingredients
		}
	}

	return recipes, pagination, nil
}

func getPagination(queryParams url.Values) (int, int) {
	page := 1
	pageSize := 10

	if p, ok := queryParams["page"]; ok && len(p) > 0 {
		page = atoi(p[0])
	}

	if ps, ok := queryParams["pageSize"]; ok && len(ps) > 0 {
		pageSize = atoi(ps[0])
	}

	return page, pageSize
}

func atoi(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

package mealplan

import (
	liberror "Food/internal/errors"
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/mat"
	"os"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) generateMealPlans(userID string, weekStartDate time.Time) ([]MealPlanPlaceholderDTO, error) {
	// Simulating a call to a recommendation engine
	// In a real scenario, this might involve an HTTP request to an external service
	placeholders, err := s.repo.GetMealPlanPlaceholders(userID, weekStartDate)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			log.WithFields(log.Fields{"service": "recipes/add",
				"repo": "recipes/add",
			}).WithError(err))
	}

	if len(placeholders) > 0 {
		return placeholders, nil
	}

	recommendedMealPlans, err := s.callRecommendationEngine(userID, weekStartDate)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			log.WithFields(log.Fields{"service": "recipes/add",
				"repo": "recipes/add",
			}).WithError(err))
	}

	for i := range recommendedMealPlans {
		recommendedMealPlans[i].WeekStartDate = weekStartDate
	}

	err = s.repo.save(recommendedMealPlans)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			log.WithFields(log.Fields{"service": "recipes/add",
				"repo": "recipes/add",
			}).WithError(err))
	}

	// Retrieve the placeholders using the repository method
	placeholders, err = s.repo.GetMealPlanPlaceholders(userID, weekStartDate)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			log.WithFields(log.Fields{"service": "recipes/add",
				"repo": "recipes/add",
			}).WithError(err))
	}

	return placeholders, nil
}
func (s *Service) GetMealPlansForDay(userID string, dayOfWeek DayOfWeek, weekStartDate time.Time) ([]DetailedMealPlanDTO, error) {
	recipes, err := s.repo.GetMealPlansForDay(userID, dayOfWeek, weekStartDate)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			log.WithFields(log.Fields{"service": "recipes/add",
				"repo": "recipes/add",
			}).WithError(err))
	}

	// Extract recipe IDs
	var recipeIDs []uuid.UUID
	for _, recipe := range recipes {
		recipeIDs = append(recipeIDs, recipe.ID)
	}

	// Get ingredients for the recipes
	ingredientsMap, err := s.repo.GetIngredientsForRecipes(recipeIDs)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			log.WithFields(log.Fields{"service": "recipes/add",
				"repo": "recipes/add",
			}).WithError(err))
	}

	// Attach ingredients to the recipes
	for i, recipe := range recipes {
		if ingredients, ok := ingredientsMap[recipe.ID]; ok {
			recipes[i].Ingredients = ingredients
		}
	}

	return recipes, nil
}

func (s *Service) callRecommendationEngine(userID string, weekStartDate time.Time) (MealPlans, error) {
	// Fetch random recipes
	randomRecipes, err := s.repo.GetRandomRecipes(21) // Fetch 21 recipes for 7 days, 3 meals per day
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			log.WithFields(log.Fields{"service": "recipes/add",
				"repo": "recipes/add",
			}).WithError(err))
	}

	// Generate meal plans using the random recipes
	mealPlans := MealPlans{}
	mealTypes := []MealType{Breakfast, Lunch, Dinner}
	daysOfWeek := []DayOfWeek{Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday}

	for i, day := range daysOfWeek {
		for j, mealType := range mealTypes {
			recipeIndex := i*3 + j
			mealPlans = append(mealPlans, MealPlan{
				UserID:        userID,
				DayOfWeek:     day,
				MealType:      mealType,
				RecipeID:      randomRecipes[recipeIndex].Id,
				WeekStartDate: weekStartDate,
				ImageURL:      randomRecipes[recipeIndex].ImgUrl,
			})
		}
	}

	return mealPlans, nil
}

func (s *Service) GetMealPlan() {
	s.saveEmbeddings()
	// List of selected food names
	selectedFoodNames := []string{"Efo Riro", "Jollof Rice", "Pounded Yam"}

	// Convert the list to a PostgreSQL array string
	foodNamesArray := "{" + string(selectedFoodNames[0])
	for _, name := range selectedFoodNames[1:] {
		foodNamesArray += "," + name
	}
	foodNamesArray += "}"

	// SQL query to find similar foods using the average vector
	// SQL query to fetch closest matches for each selected food
	query := `
    WITH selected_foods AS (
        SELECT $1::TEXT[] AS food_names
    ),
    average_vector_result AS (
        SELECT average_vector(food_names) AS avg_vector
        FROM selected_foods
    )
    SELECT
        fv.food_name,
        cosine_similarity(average_vector_result.avg_vector, fv.vector) AS similarity
    FROM
        food_vectors fv,
        average_vector_result
    ORDER BY
        similarity DESC
    LIMIT 20;
    `

	// Execute the query
	rows, err := s.repo.db.Query(query, foodNamesArray)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Print the results
	fmt.Println("Recommended foods:")
	for rows.Next() {
		var foodName string
		var similarity float64
		if err := rows.Scan(&foodName, &similarity); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s (similarity: %.2f)\n", foodName, similarity)
	}
}

func (s *Service) saveEmbeddings() {
	// Step 1: Read and clean the CSV file
	file, err := os.Open("./../Nigerian Foods.csv")
	if err != nil {
		log.Fatal(err)
	}
	//defer file.Close()

	df := dataframe.ReadCSV(file)

	// Fill missing 'Spice_Level' with 'Unknown'
	spiceLevel := df.Col("Spice_Level").Records()
	for i, v := range spiceLevel {
		if v == "" {
			spiceLevel[i] = "Unknown"
		}
	}
	df = df.Mutate(series.New(spiceLevel, series.String, "Spice_Level"))

	// Drop 'Price_Range' and 'Description' columns
	df = df.Drop([]string{"Price_Range", "Description"})

	// Step 2: Encode categorical variables
	categoricalColumns := []string{"Food_Name", "Main_Ingredients", "Food_Health", "Food_Class", "Region", "Spice_Level"}
	encodedData, _ := encodeCategorical(df, categoricalColumns)

	// Step 3: Store encoded data in PostgreSQL
	for i := 0; i < df.Nrow(); i++ {
		foodName := df.Elem(i, 0).String()
		vector := encodedData.RawRowView(i)

		// Convert vector to PostgreSQL array format
		vectorArray := fmt.Sprintf("{%v}", vector)
		for j := 0; j < len(vector); j++ {
			if j == 0 {
				vectorArray = fmt.Sprintf("{%v", vector[j])
			} else {
				vectorArray = fmt.Sprintf("%s,%v", vectorArray, vector[j])
			}
		}
		vectorArray = fmt.Sprintf("%s}", vectorArray)

		// Insert into PostgreSQL
		_, err := s.repo.db.Exec("INSERT INTO food_vectors (food_name, vector) VALUES ($1, $2) ON CONFLICT (food_name) DO NOTHING", foodName, vectorArray)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func encodeCategorical(df dataframe.DataFrame, cols []string) (*mat.Dense, map[string][]string) {
	uniqueValuesMap := make(map[string][]string)
	for _, col := range cols {
		uniqueValuesMap[col] = unique(df.Col(col).Records())
	}

	numRows := df.Nrow()
	numCols := 0
	for _, uniqueVals := range uniqueValuesMap {
		numCols += len(uniqueVals)
	}

	data := make([]float64, numRows*numCols)
	encodedData := mat.NewDense(numRows, numCols, data)

	colOffset := 0
	for _, col := range cols {
		uniqueVals := uniqueValuesMap[col]
		for rowIdx, val := range df.Col(col).Records() {
			for j, uniqueVal := range uniqueVals {
				if val == uniqueVal {
					encodedData.Set(rowIdx, colOffset+j, 1)
				} else {
					encodedData.Set(rowIdx, colOffset+j, 0)
				}
			}
		}
		colOffset += len(uniqueVals)
	}
	return encodedData, uniqueValuesMap
}

func unique(records []string) []string {
	uniqueMap := make(map[string]struct{})
	for _, record := range records {
		uniqueMap[record] = struct{}{}
	}
	uniqueSlice := make([]string, 0, len(uniqueMap))
	for k := range uniqueMap {
		uniqueSlice = append(uniqueSlice, k)
	}
	return uniqueSlice
}

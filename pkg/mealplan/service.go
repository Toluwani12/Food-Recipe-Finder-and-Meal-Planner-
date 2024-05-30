package mealplan

import (
	liberror "Food/internal/errors"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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

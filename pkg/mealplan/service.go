package mealplan

import (
	liberror "Food/internal/errors"
	"context"
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

func (s Service) getMealPLan(userID string, weekStartDate time.Time) ([]MealPlanPlaceholderDTO, error) {
	placeholders, err := s.repo.GetMealPlanPlaceholders(userID, weekStartDate)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			log.WithFields(log.Fields{"service": "recipes/add",
				"repo": "recipes/add",
			}).WithError(err))
	}

	return placeholders, nil
}

func (s *Service) generateMealPlans(ctx context.Context, userID uuid.UUID, weekStartDate time.Time) ([]MealPlanPlaceholderDTO, error) {

	recommendedMealPlans, err := s.callRecommendationEngine(ctx, userID, weekStartDate)
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
	placeholders, err := s.repo.GetMealPlanPlaceholders(userID.String(), weekStartDate)
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

func (s *Service) callRecommendationEngine(ctx context.Context, userID uuid.UUID, weekStartDate time.Time) (MealPlans, error) {
	// Fetch random recipes
	randomRecipes, err := s.repo.RecommendRecipes(ctx, userID, 21) // Fetch 21 recipes for 7 days, 3 meals per day
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
				UserID:        userID.String(),
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

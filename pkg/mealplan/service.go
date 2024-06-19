package mealplan

import (
	liberror "Food/internal/errors"
	"Food/pkg"
	"Food/pkg/recipe"
	"context"
	"fmt"
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

func (s *Service) getMealPlan(userID string, weekStartDate time.Time) ([]MealPlanPlaceholderDTO, error) {
	placeholders, err := s.repo.GetMealPlanPlaceholders(userID, weekStartDate)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			pkg.Log("mealplan.getMealPlan", "mealplan.GetMealPlanPlaceholders", userID, log.Fields{
				"week_start_date": weekStartDate,
			}).WithError(err))
	}

	return placeholders, nil
}

func (s *Service) generateMealPlans(ctx context.Context, userID uuid.UUID, weekStartDate time.Time) ([]MealPlanPlaceholderDTO, error) {
	recommendedMealPlans, err := s.callRecommendationEngine(ctx, userID, weekStartDate)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			pkg.Log("mealplan.generateMealPlans", "mealplan.callRecommendationEngine", userID.String(), log.Fields{
				"week_start_date": weekStartDate,
			}).WithError(err))
	}

	for i := range recommendedMealPlans {
		recommendedMealPlans[i].WeekStartDate = weekStartDate
	}

	err = s.repo.save(ctx, recommendedMealPlans)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			pkg.Log("mealplan.generateMealPlans", "mealplan.save", userID.String(), log.Fields{
				"week_start_date": weekStartDate,
			}).WithError(err))
	}

	placeholders, err := s.repo.GetMealPlanPlaceholders(userID.String(), weekStartDate)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			pkg.Log("mealplan.generateMealPlans", "mealplan.GetMealPlanPlaceholders", userID.String(), log.Fields{
				"week_start_date": weekStartDate,
			}).WithError(err))
	}

	return placeholders, nil
}

func (s *Service) GetMealPlansForDay(userID string, dayOfWeek DayOfWeek, weekStartDate time.Time) ([]DetailedMealPlanDTO, error) {
	recipes, err := s.repo.GetMealPlansForDay(userID, dayOfWeek, weekStartDate)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			pkg.Log("mealplan.GetMealPlansForDay", "mealplan.GetMealPlansForDay", userID, log.Fields{
				"day_of_week":     dayOfWeek,
				"week_start_date": weekStartDate,
			}).WithError(err))
	}

	if len(recipes) != 3 {
		return nil, liberror.CoverErr(fmt.Errorf("expected 3 recipes, got %d", len(recipes)),
			errors.New("service temporarily unavailable. Please try again later"),
			pkg.Log("mealplan.GetMealPlansForDay", "mealplan.GetMealPlansForDay", userID, log.Fields{
				"day_of_week":     dayOfWeek,
				"week_start_date": weekStartDate,
			}).WithError(err))
	}

	var recipeIDs []uuid.UUID
	for _, recipe := range recipes {
		recipeIDs = append(recipeIDs, recipe.ID)
	}

	ingredientsMap, err := s.repo.GetIngredientsForRecipes(recipeIDs)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			pkg.Log("mealplan.GetMealPlansForDay", "mealplan.GetIngredientsForRecipes", userID, log.Fields{
				"day_of_week":     dayOfWeek,
				"week_start_date": weekStartDate,
			}).WithError(err))
	}

	for i, recipe := range recipes {
		if ingredients, ok := ingredientsMap[recipe.ID]; ok {
			recipes[i].Ingredients = ingredients
		}
	}

	return recipes, nil
}

func (s *Service) callRecommendationEngine(ctx context.Context, userID uuid.UUID, weekStartDate time.Time) (MealPlans, error) {
	// Function to recommend recipes based on meal type
	recommendByMealType := func(mealType string, limit int) ([]recipe.Recipe, error) {
		return s.repo.RecommendRecipes(ctx, userID, limit, mealType)
	}

	breakfastRecipes, err := recommendByMealType("Breakfast", 7)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			pkg.Log("mealplan.callRecommendationEngine", "mealplan.RecommendRecipes", userID.String(), log.Fields{
				"week_start_date": weekStartDate,
			}).WithError(err))
	}

	lunchRecipes, err := recommendByMealType("Lunch", 7)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			pkg.Log("mealplan.callRecommendationEngine", "mealplan.RecommendRecipes", userID.String(), log.Fields{
				"week_start_date": weekStartDate,
			}).WithError(err))
	}

	dinnerRecipes, err := recommendByMealType("Dinner", 7)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			pkg.Log("mealplan.callRecommendationEngine", "mealplan.RecommendRecipes", userID.String(), log.Fields{
				"week_start_date": weekStartDate,
			}).WithError(err))
	}

	// Ensure we have exactly 7 recipes for each meal type
	if len(breakfastRecipes) != 7 || len(lunchRecipes) != 7 || len(dinnerRecipes) != 7 {
		return nil, liberror.CoverErr(fmt.Errorf("expected 7 recipes for each meal type, got %d breakfast, %d lunch, %d dinner", len(breakfastRecipes), len(lunchRecipes), len(dinnerRecipes)),
			errors.New("service temporarily unavailable. Please try again later"),
			pkg.Log("mealplan.callRecommendationEngine", "mealplan.RecommendRecipes", userID.String(), log.Fields{
				"week_start_date": weekStartDate,
			}).WithError(err))
	}

	mealPlans := MealPlans{}
	daysOfWeek := []DayOfWeek{Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday}

	for i, day := range daysOfWeek {
		mealPlans = append(mealPlans, MealPlan{
			UserID:        userID.String(),
			DayOfWeek:     day,
			MealType:      Breakfast,
			RecipeID:      breakfastRecipes[i].Id,
			WeekStartDate: weekStartDate,
			ImageURL:      breakfastRecipes[i].ImgUrl,
		})

		mealPlans = append(mealPlans, MealPlan{
			UserID:        userID.String(),
			DayOfWeek:     day,
			MealType:      Lunch,
			RecipeID:      lunchRecipes[i].Id,
			WeekStartDate: weekStartDate,
			ImageURL:      lunchRecipes[i].ImgUrl,
		})

		mealPlans = append(mealPlans, MealPlan{
			UserID:        userID.String(),
			DayOfWeek:     day,
			MealType:      Dinner,
			RecipeID:      dinnerRecipes[i].Id,
			WeekStartDate: weekStartDate,
			ImageURL:      dinnerRecipes[i].ImgUrl,
		})
	}

	return mealPlans, nil
}

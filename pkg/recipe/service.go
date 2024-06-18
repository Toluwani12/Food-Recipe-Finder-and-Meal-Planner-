package recipe

import (
	liberror "Food/internal/errors"
	"Food/pkg"
	"Food/pkg/recipe/crawler"
	"Food/pkg/recipe/model"
	"context"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/url"
)

type Service struct {
	repo        *Repository
	crawlerList []crawler.ICrawler
}

func NewService(repo *Repository, crawlerList []crawler.ICrawler) *Service {
	return &Service{
		repo:        repo,
		crawlerList: crawlerList,
	}
}

func (s Service) crawl(ctx context.Context) (map[string]bool, error) {
	recipeList := make([]model.RequestData, 0, len(s.crawlerList))
	for _, c := range s.crawlerList {
		data, err := c.CrawlRecipe()
		if err != nil {
			log.WithFields(log.Fields{"service": "recipes/fetchAndSaveRecipes", "crawler": c}).WithError(err).Error("Failed to crawl recipe")
			continue
		}
		recipeList = append(recipeList, *data...)
	}

	req := model.Request(recipeList)

	return s.save(ctx, req)
}

func (s Service) save(ctx context.Context, recipes model.Request) (map[string]bool, error) {
	dupl := make(map[string]bool)
	var uniqueRecipes []model.RequestData
	for _, recipe := range recipes {
		if _, ok := dupl[recipe.Name]; !ok {
			dupl[recipe.Name] = true
			uniqueRecipes = append(uniqueRecipes, recipe)
		}
	}

	recipes = uniqueRecipes
	for i := range recipes {
		if recipes[i].ID == uuid.Nil {
			recipes[i].ID = uuid.New()
		}
	}
	successMap, err := s.repo.processRecipesAndIngredients(ctx, recipes)
	return successMap, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		pkg.Log("recipes.save", "recipe.processRecipesAndIngredients", "").WithError(err))
}

func (s Service) delete(ctx context.Context, id string) (string, error) {
	resp, err := s.repo.delete(ctx, id)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		pkg.Log("recipes.delete", "recipe.delete", id).WithError(err))
}

func (s Service) get(ctx context.Context, id string) (*Recipe, error) {
	resp, err := s.repo.get(ctx, id)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		pkg.Log("recipes.get", "recipe.get", id).WithError(err))
}

func (s Service) list(ctx context.Context, userID string, recipeName string) ([]ListResponse, error) {
	resp, err := s.repo.list(ctx, userID, recipeName)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		pkg.Log("recipes.list", "recipe.list", "").WithError(err))
}

func (s Service) search(ctx context.Context, ingredients []string, queryParams url.Values, userID string) ([]ResponseData, *pkg.Pagination, error) {
	recipes, pg, err := s.repo.search(ctx, ingredients, queryParams, userID)
	return recipes, pg, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		pkg.Log("recipes.search", "recipe.search", "").WithError(err))
}

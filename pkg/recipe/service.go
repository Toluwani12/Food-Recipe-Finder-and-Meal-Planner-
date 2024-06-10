package recipe

import (
	liberror "Food/internal/errors"
	"Food/pkg"
	"context"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/url"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// AddRecipes adds a list of new recipes along with their ingredients to the database.
func (s Service) save(ctx context.Context, recipes Request) (map[string]bool, error) {

	// Ensure all IDs are set for recipes, ingredients, and their alternatives
	for i := range recipes {
		if recipes[i].ID == uuid.Nil {
			recipes[i].ID = uuid.New() // Ensure the recipe has a UUID
		}
	}
	successMap, err := s.repo.processRecipesAndIngredients(ctx, recipes)
	return successMap, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "recipes/add",
			"repo": "recipes/add",
		}).WithError(err))
}

//func (s Service) update(ctx context.Context, id string, data AddRequest) (*Recipe, error) {
//	recipe := Recipe{
//		Id:           id,
//		Name:         data.Name,
//		CookingTime:  data.CookingTime,
//		Instructions: data.Instructions,
//		UpdatedAt:    time.Now(),
//	}
//	resp, err := s.repo.update(ctx, recipe)
//	return resp, liberror.CoverErr(err,
//		errors.New("service temporarily unavailable. Please try again later"),
//		log.WithFields(log.Fields{"service": "recipes/update", "repo": "recipes/update"}).WithError(err))
//}

func (s Service) delete(ctx context.Context, id string) (string, error) {
	resp, err := s.repo.delete(ctx, id)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "recipes/delete", "repo": "recipes/delete"}).WithError(err))
}

func (s Service) get(ctx context.Context, id string) (*Recipe, error) {
	resp, err := s.repo.get(ctx, id)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "recipes/get", "repo": "recipes/get"}).WithError(err))
}

func (s Service) list(ctx context.Context, userID string, recipeName string) ([]ListResponse, error) {
	resp, err := s.repo.list(ctx, userID, recipeName)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "recipes/list", "repo": "recipes/list"}).WithError(err))
}

func (s Service) search(ctx context.Context, ingredients []string, queryParams url.Values, userId string) (interface{}, *pkg.Pagination, error) {
	recipes, pg, err := s.repo.search(ctx, ingredients, queryParams, userId)
	return recipes, pg, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "recipes/findRecipes", "repo": "recipes/findRecipes"}).WithError(err))
}

//func (s Service) search2(ctx context.Context, ingredients []string, queryParams url.Values) (interface{}, *pkg.Pagination, error) {
//	recipes, pg, err := s.repo.Search(ctx, ingredients, queryParams)
//	return recipes, pg, liberror.CoverErr(err,
//		errors.New("service temporarily unavailable. Please try again later"),
//		log.WithFields(log.Fields{"service": "recipes/findRecipes", "repo": "recipes/findRecipes"}).WithError(err))
//}

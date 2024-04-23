package recipe

import (
	liberror "Food/internal/errors"
	"context"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) save(ctx context.Context, data AddRequest) (*Recipe, error) {
	recipe := Recipe{
		Id:           uuid.NewString(),
		Name:         data.Name,
		CookingTime:  data.CookingTime,
		Instructions: data.Instructions,
		CreatedAt:    time.Now(),
	}
	resp, err := s.repo.save(ctx, recipe)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "recipes/save", "repo": "recipes/save"}).WithError(err))
}

func (s Service) update(ctx context.Context, id string, data AddRequest) (*Recipe, error) {
	recipe := Recipe{
		Id:           id,
		Name:         data.Name,
		CookingTime:  data.CookingTime,
		Instructions: data.Instructions,
		UpdatedAt:    time.Now(),
	}
	resp, err := s.repo.update(ctx, recipe)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "recipes/update", "repo": "recipes/update"}).WithError(err))
}

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

func (s Service) list(ctx context.Context) (Recipes, error) {
	resp, err := s.repo.list(ctx)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "recipes/list", "repo": "recipes/list"}).WithError(err))
}

func (s Service) findRecipes(ctx context.Context, ingredients []string) ([]Recipe, error) {
	if len(ingredients) == 0 {
		return nil, liberror.New("no ingredients provided", http.StatusBadRequest)
	}
	recipes, err := s.repo.findRecipes(ctx, ingredients)
	return recipes, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "recipes/findRecipes", "repo": "recipes/findRecipes"}).WithError(err))
}

package ingredient

import (
	liberror "Food/internal/errors"
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) save(ctx context.Context, data AddRequest) (*Ingredient, error) {
	ingredient := Ingredient{
		Name:         data.Name,
		Alternatives: data.Alternatives,
		Quantity:     data.Quantity,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	resp, err := s.repo.save(ctx, ingredient)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "ingredients/save", "repo": "ingredients/save"}).WithError(err))
}

func (s Service) delete(ctx context.Context, id string) (string, error) {
	resp, err := s.repo.delete(ctx, id)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "ingredients/delete", "repo": "ingredients/delete"}).WithError(err))
}

func (s Service) update(ctx context.Context, id string, data AddRequest) (*Ingredient, error) {
	// Check if an ingredient with the same name already exists
	existingIngredient, err := s.repo.get(ctx, id)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			log.WithFields(log.Fields{"service": "ingredients/save", "repo": "ingredients/save"}).WithError(err))
	}
	ingredient := Ingredient{
		ID:           id,
		Name:         data.Name,
		Alternatives: data.Alternatives,
		Quantity:     data.Quantity,
		CreatedAt:    existingIngredient.CreatedAt,
		UpdatedAt:    time.Now(),
	}

	// Update the ingredient
	updatedIngredient, err := s.repo.update(ctx, ingredient)
	return updatedIngredient, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "ingredients/update", "repo": "ingredients/update"}).WithError(err))
}

func (s Service) get(ctx context.Context, id string) (*Ingredient, error) {
	resp, err := s.repo.get(ctx, id)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "ingredients/get", "repo": "ingredients/get"}).WithError(err))
}

func (s Service) list(ctx context.Context) (Ingredients, error) {
	resp, err := s.repo.list(ctx)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "ingredients/list", "repo": "ingredients/list"}).WithError(err))
}

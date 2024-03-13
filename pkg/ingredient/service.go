package ingredient

import (
	"context"
	"errors"
	"strconv"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) add(ctx context.Context, data AddRequest) (*Ingredient, error) {
	ingredient, err := s.repo.getByReference(data.Name)
	if ingredient != nil {
		return nil, errors.New("ingredient with this reference already exist")
	}

	if err != nil {
		return nil, err
	}

	resp, err := s.repo.save(ctx, data)

	if err != nil {
		return nil, err
	}

	return resp, nil

}

func (s Service) delete(ctx context.Context, id string) (*Ingredient, error) {
	resp, err := s.repo.delete(ctx, id)

	if err != nil {
		return nil, err
	}

	return resp, nil

}

func (s Service) update(ctx context.Context, id string, data AddRequest) (*Ingredient, error) {
	// Check if an ingredient with the same name already exists
	existingIngredient, err := s.repo.getByReference(data.Name)
	if err != nil {
		return nil, err
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err // Handle the error if the conversion fails
	}

	// Then use idInt in the comparison
	if existingIngredient != nil && existingIngredient.ID != idInt {
		return nil, errors.New("ingredient with this reference already exists")
	}

	// Update the ingredient
	updatedIngredient, err := s.repo.update(id, data)
	if err != nil {
		return nil, err
	}

	return updatedIngredient, nil
}

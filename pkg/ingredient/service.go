package ingredient

import (
	"context"
	"errors"
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

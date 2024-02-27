package ingredient

import "errors"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) add(data Ingredient) (*Ingredient, error) {
	ingredient, err := s.repo.getByReference(data.Name)
	if ingredient != nil {
		return nil, errors.New("ingredient with this reference already exist")
	}

	if err != nil {
		return nil, err
	}

	resp, err := s.repo.save(data)

	if err != nil {
		return nil, err
	}

	return resp, nil

}

func (s Service) delete(data Ingredient) (*Ingredient, error) {
	ingredient, err := s.repo.getByReference(data.Name)
	if ingredient != nil {
		return nil, errors.New("ingredient with this reference doesn't exist")
	}

	if err != nil {
		return nil, err
	}

	resp, err := s.repo.delete(data)

	if err != nil {
		return nil, err
	}

	return resp, nil

}

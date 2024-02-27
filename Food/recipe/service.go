package recipe

import "errors"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) add(data Recipe) (*Recipe, error) {
	recipe, err := s.repo.getByReference(data.Name)
	if recipe != nil {
		return nil, errors.New("recipe with this reference already exist")
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

func (s Service) delete(data Recipe) (*Recipe, error) {
	recipe, err := s.repo.getByReference(data.Name)
	if recipe != nil {
		return nil, errors.New("recipe with this reference doesn't exist")
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

package user_preference

import (
	liberror "Food/internal/errors"
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) save(ctx context.Context, userID string, req AddRequest) error {
	//first retrieve user preference by user id
	userPreference, err := s.repo.get(ctx, userID)
	if err != nil && !errors.Is(err, liberror.ErrNotFound) {
		return liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			log.WithFields(log.Fields{"service": "users/findByEmail", "repo": "users/findByEmail"}).WithError(err))
	}

	if userPreference != nil {
		return liberror.New("user preference already exists for this user", http.StatusBadRequest)
	}

	return s.repo.add(ctx, userID, req)
}

func (s *Service) delete(ctx context.Context, id string) error {
	return s.repo.delete(ctx, id)
}

func (s *Service) update(ctx context.Context, id string, data AddRequest) error {
	err := s.repo.update(ctx, id, data)
	return liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "users/findByEmail", "repo": "users/findByEmail"}).WithError(err))
}

func (s *Service) get(ctx context.Context, id string) (*UserPreference, error) {
	u, err := s.repo.get(ctx, id)
	return u, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "users/findByEmail", "repo": "users/findByEmail"}).WithError(err))
}

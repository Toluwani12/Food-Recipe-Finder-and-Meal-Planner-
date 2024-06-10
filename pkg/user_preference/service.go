package user_preference

import (
	liberror "Food/internal/errors"
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

//func (s *Service) save(ctx context.Context, userID string, req AddRequest) error {
//
//	return s.repo.add(ctx, s.repo.db.ExecContext, userID, req.RecipeIds)
//}
//
//func (s *Service) delete(ctx context.Context, userID string, req AddRequest) error {
//	return s.repo.remove(ctx, s.repo.db.ExecContext, userID, req.RecipeIds)
//}

func (s *Service) Save(ctx context.Context, userID string, req AddRequest, liked bool) error {
	return s.repo.setLikeStatus(ctx, userID, req.RecipeIds, liked)
}

//
//func (s *Service) delete(ctx context.Context, userID string, req AddRequest) error {
//	return s.repo.remove(ctx, s.repo.db.ExecContext, userID, req.RecipeIds)
//}

func (s *Service) get(ctx context.Context, id string) (*UserPreference, error) {
	u, err := s.repo.get(ctx, id)
	return &UserPreference{
			UserID:       id,
			LikedRecipes: u,
		}, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			log.WithFields(log.Fields{"service": "users/findByEmail", "repo": "users/findByEmail"}).WithError(err))
}

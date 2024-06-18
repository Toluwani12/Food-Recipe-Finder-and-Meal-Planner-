package user_preference

import (
	liberror "Food/internal/errors"
	"Food/pkg"
	"context"
	"github.com/pkg/errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Save(ctx context.Context, userID string, req AddRequest, liked bool) error {
	err := s.repo.setLikeStatus(ctx, userID, req.RecipeIds, liked)
	return liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		pkg.Log("user_preference.Save", "user_preference.setLikeStatus", userID).WithError(err))
}

func (s *Service) Get(ctx context.Context, userID string) (*UserPreference, error) {
	recipes, err := s.repo.get(ctx, userID)
	if err != nil {
		return nil, liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			pkg.Log("user_preference.Get", "user_preference.get", userID).WithError(err))
	}

	return &UserPreference{
		UserID:       userID,
		LikedRecipes: recipes,
	}, nil
}

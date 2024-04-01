package mealplan

import (
	liberror "Food/internal/errors"
	"context"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) save(ctx context.Context, data AddRequest) (*MealPlan, error) {
	mealPlan := MealPlan{
		Id:       uuid.NewString(),
		Date:     time.Now(),
		MealType: time.Time{},
	}
	resp, err := s.repo.save(ctx, mealPlan)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "mealPlans/save", "repo": "mealPlans/save"}).WithError(err))
}

func (s Service) delete(ctx context.Context, id string) (string, error) {
	resp, err := s.repo.delete(ctx, id)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "mealPlans/delete", "repo": "mealPlans/delete"}).WithError(err))
}

func (s Service) get(ctx context.Context, id string) (*MealPlan, error) {
	resp, err := s.repo.get(ctx, id)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "mealPlans/get", "repo": "mealPlans/get"}).WithError(err))
}

func (s Service) list(ctx context.Context) (MealPlans, error) {
	resp, err := s.repo.list(ctx)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "mealPlans/list", "repo": "mealPlans/list"}).WithError(err))
}

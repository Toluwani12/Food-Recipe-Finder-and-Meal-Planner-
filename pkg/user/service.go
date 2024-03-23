package user

import (
	liberror "Food/internal/errors"
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) delete(ctx context.Context, id string) (interface{}, interface{}) {
	resp, err := s.repo.delete(ctx, id)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s Service) save(ctx context.Context, request AddRequest) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, liberror.CoverErr(err, errors.New("service temporoaily unavailable"),
			log.WithFields(log.Fields{"service": "users/save", "repo": "users/save"}).WithError(err))
	}
	user := User{
		ID:       uuid.NewString(),
		Email:    request.Email,
		Password: string(hashedPassword),
	}
	save, err := s.repo.save(ctx, user)
	if err != nil {
		return Users{}, err
	}
	return save, nil
}

func (s Service) get(ctx context.Context, id string) (*User, error) {
	resp, err := s.repo.get(ctx, id)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "users/get", "repo": "users/get"}).WithError(err))
}

func (s Service) list(ctx context.Context) (Users, error) {
	resp, err := s.repo.list(ctx)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		log.WithFields(log.Fields{"service": "users/list", "repo": "users/list"}).WithError(err))
}

func (s Service) login(ctx context.Context, request LoginRequest) {

}

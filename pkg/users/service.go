package users

import (
	liberror "Food/internal/errors"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) delete(ctx context.Context, id string) (string, error) {
	resp, err := s.repo.delete(ctx, id)

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (s Service) update(ctx context.Context, id string, request UpdateRequest) (string, error) {

	save, err := s.repo.update(ctx, id, request)
	if err != nil {
		return "", liberror.CoverErr(err, errors.New("service temporally unavailable"),
			log.WithFields(log.Fields{"service": "users/update", "repo": "users/update"}).WithError(err))
	}
	return save, nil

}

func (s Service) save(ctx context.Context, request AddRequest) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, liberror.CoverErr(err, errors.New("service temporally unavailable"),
			log.WithFields(log.Fields{"service": "users/bcrypt.GenerateFromPassword", "repo": "users/bcrypt.GenerateFromPassword"}).WithError(err))
	}

	request.Password = string(hashedPassword)

	save, err := s.repo.save(ctx, request)
	if err != nil {
		return nil, liberror.CoverErr(err, errors.New("service temporally unavailable"),
			log.WithFields(log.Fields{"service": "users/save", "repo": "users/save"}).WithError(err))
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

func (s Service) login(ctx context.Context, request LoginRequest) (*User, string, error) {
	user, err := s.repo.findByEmail(ctx, request.Email)
	if err != nil {
		return nil, "", liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			log.WithFields(log.Fields{"service": "users/findByEmail", "repo": "users/findByEmail"}).WithError(err))
	}

	if user == nil {
		return nil, "", errors.New("users not exit, please sign up!")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, "", errors.New("invalid password")
	}

	// Create a new JWT token for the authenticated users
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte("your_secret_key")) // Use your secret key here
	if err != nil {
		log.WithFields(log.Fields{"service": "users/login", "repo": "users/login"}).WithError(err)
		return nil, "",
			errors.New("service temporarily unavailable. Please try again later")
	}

	return user, tokenString, nil
}

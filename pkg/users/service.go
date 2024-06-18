package users

import (
	liberror "Food/internal/errors"
	"Food/pkg"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		pkg.Log("users.delete", "users.delete", id).WithError(err))
}

func (s Service) update(ctx context.Context, id string, request UpdateRequest) (string, error) {
	save, err := s.repo.update(ctx, id, request)
	return save, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		pkg.Log("users.update", "users.update", id).WithError(err))
}

func (s Service) save(ctx context.Context, request AddRequest) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, liberror.CoverErr(err, errors.New("service temporarily unavailable. Please try again later"),
			pkg.Log("users.save", "bcrypt.GenerateFromPassword", "").WithError(err))
	}

	request.PasswordHash = string(hashedPassword)

	save, err := s.repo.save(ctx, request)
	return save, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		pkg.Log("users.save", "users.save", "").WithError(err))
}

func (s Service) get(ctx context.Context, id string) (*User, error) {
	resp, err := s.repo.get(ctx, id)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		pkg.Log("users.get", "users.get", id).WithError(err))
}

func (s Service) list(ctx context.Context) (Users, error) {
	resp, err := s.repo.list(ctx)
	return resp, liberror.CoverErr(err,
		errors.New("service temporarily unavailable. Please try again later"),
		pkg.Log("users.list", "users.list", "").WithError(err))
}

func (s Service) login(ctx context.Context, request LoginRequest) (*User, string, error) {
	user, err := s.repo.findByEmail(ctx, request.Email)
	if err != nil {
		return nil, "", liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			pkg.Log("users.login", "users.findByEmail", request.Email).WithError(err))
	}

	if user == nil {
		return nil, "", liberror.New("User does not exist, please sign up!", http.StatusBadRequest)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		return nil, "", liberror.New("Invalid password", http.StatusBadRequest)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		return nil, "", liberror.CoverErr(err,
			errors.New("service temporarily unavailable. Please try again later"),
			pkg.Log("users.login", "jwt.SignedString", request.Email).WithError(err))
	}

	return user, tokenString, nil
}

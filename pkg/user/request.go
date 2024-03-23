package user

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"net/http"
	"strings"
)

type AddRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password"`
}

func (v *AddRequest) Bind(r *http.Request) error {
	if err := render.Bind(r, v); err != nil {
		return err
	}

	err1 := validate.Validate(
		&validators.StringIsPresent{Name: "name", Field: v.Name, Message: fmt.Sprintf("%s is missing", "name")},
		&validators.EmailIsPresent{Name: "email", Field: v.Email, Message: fmt.Sprintf("%s is invalid", "email")},
		&validators.StringIsPresent{Name: "password", Field: v.Password, Message: fmt.Sprintf("%s is missing", "password")},
		&validators.StringIsPresent{Name: "confirm_password", Field: v.ConfirmPassword, Message: fmt.Sprintf("%s is missing", "confirm_password")},
		&validators.FuncValidator{
			Fn: func() bool {
				return v.Password == v.ConfirmPassword
			},
			Message: "Password doesn't match",
		},
	)

	v.Name = strings.TrimSpace(strings.ToLower(v.Name))
	if err1.HasAny() {
		return err1
	}

	return nil
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (v *LoginRequest) Bind(r *http.Request) error {
	if err := render.Bind(r, v); err != nil {
		return err
	}

	err1 := validate.Validate(
		&validators.EmailIsPresent{Name: "email", Field: v.Email, Message: fmt.Sprintf("%s is invalid", "email")},
		&validators.StringIsPresent{Name: "password", Field: v.Password, Message: fmt.Sprintf("%s is missing", "password")},
	)

	v.Email = strings.TrimSpace(strings.ToLower(v.Email))
	if err1.HasAny() {
		return err1
	}

	return nil
}

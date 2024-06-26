package users

import (
	"fmt"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"net/http"
	"strings"
)

type AddRequest struct {
	Username        string `json:"username" db:"username"`
	Email           string `json:"email" db:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	PasswordHash    string `db:"password_hash"`
}

func (v *AddRequest) Bind(r *http.Request) error {
	err1 := validate.Validate(
		&validators.StringIsPresent{Name: "username", Field: v.Username, Message: fmt.Sprintf("%v is missing", "username")},
		&validators.EmailIsPresent{Name: "email", Field: v.Email, Message: fmt.Sprintf("%v is invalid", "email")},
		&validators.StringIsPresent{Name: "password", Field: v.Password, Message: fmt.Sprintf("%v is missing", "password")},
		&validators.StringIsPresent{Name: "confirm_password", Field: v.ConfirmPassword},
		&validators.FuncValidator{
			Fn: func() bool {
				return v.Password == v.ConfirmPassword
			},
			Message: "Password doesn't match",
		},
	)

	v.Username = strings.TrimSpace(strings.ToLower(v.Username))
	if err1.HasAny() {
		return err1
	}

	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password" `
}

func (v *LoginRequest) Bind(r *http.Request) error {
	err1 := validate.Validate(
		&validators.StringIsPresent{Name: "email", Field: v.Email, Message: fmt.Sprintf("%s is invalid", "email")},
		&validators.StringIsPresent{Name: "password", Field: v.Password, Message: fmt.Sprintf("%s is missing", "password")},
	)

	v.Email = strings.TrimSpace(strings.ToLower(v.Email))
	if err1.HasAny() {
		return err1
	}

	return nil
}

type UpdateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u UpdateRequest) Bind(r *http.Request) error {
	return nil
}

package user

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var validate = validator.New()

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (u UserLogin) Bind(r *http.Request) error {
	//TODO implement me
	panic("implement me")
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	if _, err = db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, string(hashedPassword)); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var userLogin UserLogin
	if err := render.Bind(r, &userLogin); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := validate.Struct(userLogin); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

}

func ErrInvalidRequest(err error) render.Renderer {
	return &render.Text{StatusCode: 400, Format: "Invalid request: %s", Data: err.Error()}
}

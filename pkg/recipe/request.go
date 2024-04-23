package recipe

import (
	"Food/pkg/ingredient"
	"fmt"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"net/http"
	"strings"
)

type AddRequest struct {
	Name         string                  `json:"name"`
	Description  string                  `json:"description"`
	CookingTime  int                     `json:"cooking_time"` // minutes
	Instructions string                  `json:"instructions"`
	Ingredients  []ingredient.AddRequest `json:"ingredients"`
}

func (v *AddRequest) Bind(r *http.Request) error {

	err1 := validate.Validate(
		&validators.StringIsPresent{Name: "name", Field: v.Name, Message: fmt.Sprintf("%s is missing", "name")},
	)

	v.Name = strings.TrimSpace(strings.ToLower(v.Name))
	if err1.HasAny() {
		return err1
	}

	return nil
}

type GetResponse struct {
	ID           int                     `json:"id"`
	Name         string                  `json:"name"`
	Description  string                  `json:"description"`
	CookingTime  int                     `json:"cooking_time"`
	Instructions string                  `json:"instructions"`
	Ingredients  []ingredient.AddRequest `json:"ingredients"`
}

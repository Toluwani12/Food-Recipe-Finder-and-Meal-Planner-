package ingredient

import (
	"fmt"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type Request struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `json:"name" db:"name"`
	Quantity     string    `json:"quantity" db:"quantity"` // Quantity for the recipe_ingredients table
	Alternatives []string  `json:"alternative,omitempty" db:"alternative"`
}

func (v *Request) Bind(r *http.Request) error {

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
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Alternatives []string `json:"alternatives"`
}

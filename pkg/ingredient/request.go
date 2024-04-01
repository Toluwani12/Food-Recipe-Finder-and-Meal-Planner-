package ingredient

import (
	"fmt"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"net/http"
	"strings"
)

type AddRequest struct {
	Name        string `json:"name"`
	Alternative string `json:"alternative,omitempty"`
	Quantity    string `json:"quantity,omitempty"`
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

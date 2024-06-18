package model

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type IngredientRequest struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `json:"name" db:"name"`
	Quantity     string    `json:"quantity" db:"quantity"` // Quantity for the recipe_ingredients table
	Alternatives []string  `json:"alternative,omitempty" db:"alternative"`
}

type RequestData struct {
	ID           uuid.UUID           `db:"id"`
	Name         string              `json:"name" db:"name"`
	Description  string              `json:"description" db:"description"`
	CookingTime  string              `json:"cooking_time" db:"cooking_time"`
	Instructions pq.StringArray      `json:"instructions" db:"instructions"`
	ImgUrl       string              `json:"img_url" db:"img_url"`
	Ingredients  []IngredientRequest `json:"ingredients" db:"ingredients"`
	Diff         int                 `json:"diff" db:"diff"`
}

type Request []RequestData

func (v *Request) Bind(r *http.Request) error {

	log.Printf("Binding recipe request: %v", v)

	//err1 := validate.Validate(
	//	&validators.StringIsPresent{Name: "name", Field: v.Name, Message: fmt.Sprintf("%s is missing", "name")},
	//)
	//
	//v.Name = strings.TrimSpace(strings.ToLower(v.Name))
	//if err1.HasAny() {
	//	return err1
	//}

	return nil
}

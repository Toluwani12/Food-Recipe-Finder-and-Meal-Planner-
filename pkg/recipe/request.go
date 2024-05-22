package recipe

import (
	"Food/pkg/ingredient"
	"github.com/google/uuid"
	"net/http"
)

type RequestData struct {
	ID           uuid.UUID            `db:"id"`
	Name         string               `json:"name" db:"name"`
	Description  string               `json:"description" db:"description"`
	CookingTime  string               `json:"cooking_time" db:"cooking_time"`
	Instructions string               `json:"instructions" db:"instructions"`
	ImgUrl       string               `json:"img_url" db:"img_url"`
	Ingredients  []ingredient.Request `json:"ingredients" db:"ingredients"`
	Diff         int                  `json:"diff" db:"diff"`
}

type Request []RequestData

func (v *Request) Bind(r *http.Request) error {

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

type ListResponse struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `json:"name" db:"name"`
	ImgUrl      string    `json:"img_url" db:"img_url"`
	Description string    `json:"description" db:"description"`
}

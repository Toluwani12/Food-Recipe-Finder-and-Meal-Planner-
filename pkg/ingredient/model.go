package ingredient

import (
	"Food/internal/errors"
	"encoding/json"
	"time"
)

type Ingredient struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Alternatives []string  `json:"alternatives"`
	RecipeID     string    `json:"recipe_id"`
	Quantity     string    `json:"quantity"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Ingredients []Ingredient

func (i *Ingredients) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, i)
}

package recipe

import (
	"encoding/json"
	"net/http"
)

type Recipe struct {
	Id           string `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Cuisine      string `json:"cuisine"`
	MealType     string `json:"mealType"`
	CookingTime  string `json:"cookingTime"`
	Instructions string `json:"instructions"`
	Servings     string `json:"servings"`
}

func (data *Recipe) bind(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return err
	}

	return nil
}

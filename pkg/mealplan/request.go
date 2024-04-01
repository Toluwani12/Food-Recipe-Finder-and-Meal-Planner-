package mealplan

type AddRequest struct {
	MealType string `json:"mealType"`
	Date     string `json:"date"`
	Quantity string `json:"quantity,omitempty"`
}

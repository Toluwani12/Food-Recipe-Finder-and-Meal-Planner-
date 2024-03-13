package ingredient

import (
	"time"
)

type Ingredient struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Alternative string    `json:"alternative"`
	Quantity    string    `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

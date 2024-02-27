package ingredient

import (
	"encoding/json"
	"net/http"
	"time"
)

type ApiResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Error   string      `json:"error,omitempty"`
}

func NewApiResponse(data interface{}, message string, code int, err string) *ApiResponse {
	return &ApiResponse{
		Data:    data,
		Message: message,
		Code:    code,
		Error:   err,
	}
}

type Ingredient struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Alternative string    `json:"alternative"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (data *Ingredient) bind(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return err
	}

	return nil
}

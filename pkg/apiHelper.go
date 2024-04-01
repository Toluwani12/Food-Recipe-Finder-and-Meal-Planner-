package pkg

import (
	"github.com/go-chi/render"
	"net/http"
)

type ApiResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Error   string      `json:"errors,omitempty"`
}

type Response struct {
	Message interface{} `json:"message,omitempty"`
	Err     interface{} `json:"errors,omitempty"`
}

func Render(w http.ResponseWriter, r *http.Request, res interface{}) {
	w.Header().Set("Content-Type", "application/json")
	switch res.(type) {
	case render.Renderer:
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, Response{Message: res.(render.Renderer), Err: ""})
	case error:
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, Response{Message: "", Err: res.(error).Error()})
	default:
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, Response{Message: res, Err: ""})
	}
}

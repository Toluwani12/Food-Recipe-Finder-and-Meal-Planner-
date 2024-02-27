package recipe

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h Handler) add(w http.ResponseWriter, r *http.Request) {

	var data Recipe

	w.Header().Set("Content-Type", "application/json")

	if err := data.bind(r); err != nil {
		resp := NewApiResponse(nil, "could not validate request", http.StatusBadRequest, err.Error())
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	recipe, err := h.svc.add(data)

	if err != nil {
		resp := NewApiResponse(nil, "could not add recipe", http.StatusInternalServerError, err.Error())
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	resp := NewApiResponse(recipe, "recipe added successfully", http.StatusCreated, "")

	_ = json.NewEncoder(w).Encode(resp)
}

func (h Handler) delete(w http.ResponseWriter, r *http.Request) {

	var data Recipe

	w.Header().Set("Content-Type", "application/json")

	if err := data.bind(r); err != nil {
		resp := NewApiResponse(nil, "could not validate request", http.StatusBadRequest, err.Error())
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	recipe, err := h.svc.delete(data)

	if err != nil {
		resp := NewApiResponse(nil, "could not add recipe", http.StatusInternalServerError, err.Error())
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	resp := NewApiResponse(recipe, "recipe added successfully", http.StatusCreated, "")

	_ = json.NewEncoder(w).Encode(resp)
}

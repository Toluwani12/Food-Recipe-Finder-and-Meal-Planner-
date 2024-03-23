package mealplan

import (
	"Food/pkg"
	"github.com/go-chi/chi/v5"
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

	// binding or extracting the data from the request
	var data AddRequest
	if err := data.Bind(r); err != nil {
		pkg.Render(w, r, err)
		return
	}

	// add the data through the add service
	ingredient, err := h.svc.save(r.Context(), data)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	// returning the response to the user
	pkg.Render(w, r, pkg.ApiResponse{
		Data:    ingredient,
		Message: "ingredient added successfully",
		Code:    201,
	})
}

func (h Handler) delete(w http.ResponseWriter, r *http.Request) {

	// extract id from url param
	id := chi.URLParam(r, "id")

	// call the delete service to delete by id
	ingredient, err := h.svc.delete(r.Context(), id)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	// return a success response to user
	pkg.Render(w, r, pkg.ApiResponse{
		Data:    ingredient,
		Message: "ingredient successfully deleted",
		Code:    200,
	})
}

func (h Handler) update(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL
	id := chi.URLParam(r, "id")

	// Binding or extracting the data from the request
	var data AddRequest
	if err := data.Bind(r); err != nil {
		pkg.Render(w, r, err)
		return
	}

	// Update the data through the update service
	ingredient, err := h.svc.update(r.Context(), id, data)
	if err != nil {
		// Handle errors
		pkg.Render(w, r, err)
		return
	}

	// Returning the response to the user
	pkg.Render(w, r, pkg.ApiResponse{
		Data:    ingredient,
		Message: "ingredient updated successfully",
		Code:    200,
	})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Get the ingredient from the service layer
	ingredient, err := h.svc.get(r.Context(), id)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    ingredient,
		Message: "Ingredient retrieved successfully",
		Code:    http.StatusOK,
	})
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	ingredients, err := h.svc.list(r.Context())
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    ingredients,
		Message: "Ingredient retrieved successfully",
		Code:    http.StatusOK,
	})
}

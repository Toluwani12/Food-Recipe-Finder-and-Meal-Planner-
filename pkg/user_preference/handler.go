package user_preference

import (
	"Food/internal/errors"
	"Food/pkg"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"strings"
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
	if err := render.Bind(r, &data); err != nil {
		pkg.Render(w, r, err)
		return
	}

	// add the data through the add service
	recipe, err := h.svc.save(r.Context(), data)
	if err != nil {
		pkg.Render(w, r, nil)
		return
	}

	// returning the response to the users
	pkg.Render(w, r, pkg.ApiResponse{
		Data:    recipe,
		Message: "recipe added successfully",
		Code:    201,
	})
}

func (h Handler) delete(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	recipe, err := h.svc.delete(r.Context(), id)
	if err != nil {
		pkg.Render(w, r, nil)
		return
	}

	// return a success response to users
	pkg.Render(w, r, pkg.ApiResponse{
		Data:    recipe,
		Message: "recipe successfully deleted",
		Code:    200,
	})
}

func (h Handler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var data AddRequest
	if err := render.Bind(r, &data); err != nil {
		pkg.Render(w, r, err)
		return
	}

	recipe, err := h.svc.update(r.Context(), id, data)
	if err != nil {
		pkg.Render(w, r, nil)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    recipe,
		Message: "recipe updated successfully",
		Code:    200,
	})
}

func (h Handler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Get the recipe from the service layer
	recipe, err := h.svc.get(r.Context(), id)
	if err != nil {
		// If an errors occurs, send an appropriate HTTP response
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// If the recipe is found, marshal it to JSON and send it in the response
	pkg.Render(w, r, pkg.ApiResponse{
		Data:    recipe,
		Message: "Recipe retrieved successfully",
		Code:    http.StatusOK,
	})
}

func (h Handler) list(w http.ResponseWriter, r *http.Request) {
	recipe, err := h.svc.list(r.Context())
	if err != nil {
		// If an errors occurs, send an appropriate HTTP response
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    recipe,
		Message: "Recipe retrieved successfully",
		Code:    http.StatusOK,
	})
}

func (h Handler) findRecipesHandler(w http.ResponseWriter, r *http.Request) {
	// Parse ingredients from query parameter, expecting comma-separated values
	queryValues := r.URL.Query()
	ingredientsParam := queryValues.Get("ingredients")
	if ingredientsParam == "" {
		pkg.Render(w, r, errors.New("ingredients parameter is required"))
		return
	}
	ingredients := strings.Split(ingredientsParam, ",")

	// Call the service layer to find recipes
	recipes, err := h.svc.findRecipes(r.Context(), ingredients)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}
	if len(recipes) == 0 {
		pkg.Render(w, r, errors.ErrWithCustomText(nil, "No recipes found", http.StatusNotFound))
		return
	}

	response, err := json.Marshal(recipes)
	if err != nil {
		pkg.Render(w, r, errors.ErrWithCustomText(err, "Error marshalling response", http.StatusInternalServerError))
		return
	}
	pkg.Render(w, r, pkg.ApiResponse{
		Data:    response,
		Message: "Recipes found successfully",
		Code:    http.StatusOK,
	})
}

package recipe

import (
	"Food/pkg"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
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

// addRecipe is a HTTP handler for adding a new recipe.
func (h Handler) save(w http.ResponseWriter, r *http.Request) {
	var recipes Request
	// Decode the JSON body into the recipe DTO
	userID := r.Context().Value("user_id").(string)
	log.Infoln("user id", userID)
	// binding or extracting the data from the request

	if err := render.Bind(r, &recipes); err != nil {
		pkg.Render(w, r, err)
		return
	}

	// Call the service to add the recipe
	successMap, err := h.svc.save(r.Context(), recipes)
	if err != nil {
		pkg.Render(w, r, pkg.ApiResponse{
			Data:    nil,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Return the newly added recipe and a success message
	pkg.Render(w, r, pkg.ApiResponse{
		Data:    successMap,
		Message: "Recipe added successfully",
		Code:    http.StatusCreated,
	})
}

func (h Handler) delete(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	recipe, err := h.svc.delete(r.Context(), id)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	// return a success response to users
	pkg.Render(w, r, pkg.ApiResponse{
		Data:    recipe,
		Message: "recipe successfully deleted",
		Code:    200,
	})
}

//func (h Handler) update(w http.ResponseWriter, r *http.Request) {
//	id := chi.URLParam(r, "id")
//
//	var data AddRequest
//	if err := render.Bind(r, &data); err != nil {
//		pkg.Render(w, r, err)
//		return
//	}
//
//	recipe, err := h.svc.update(r.Context(), id, data)
//	if err != nil {
//		pkg.Render(w, r, nil)
//		return
//	}
//
//	pkg.Render(w, r, pkg.ApiResponse{
//		Data:    recipe,
//		Message: "recipe updated successfully",
//		Code:    200,
//	})
//}

func (h Handler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Get the recipe from the service layer
	recipe, err := h.svc.get(r.Context(), id)
	if err != nil {
		pkg.Render(w, r, err)
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

type SearchRequest struct {
	Ingredients []string `json:"ingredients"`
}

func (h Handler) search(w http.ResponseWriter, r *http.Request) {
	var req SearchRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Infoln("Searching for recipes with ingredients: ", req.Ingredients)

	queryParams := r.URL.Query()
	recipes, pagination, err := h.svc.search(r.Context(), req.Ingredients, queryParams)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:       recipes,
		Message:    "Recipes retrieved successfully",
		Code:       http.StatusOK,
		Pagination: pagination,
	})
}

//func (h Handler) search2(w http.ResponseWriter, r *http.Request) {
//	log.Infoln("Searching for recipes with ingredientsssssssssssssssssss: ", r.URL.Query().Get("ingredients"))
//	var req SearchRequest
//
//	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//		http.Error(w, "Invalid request body", http.StatusBadRequest)
//		return
//	}
//
//	log.Infoln("Searching for recipes with ingredients: ", req.Ingredients)
//
//	queryParams := r.URL.Query()
//	recipes, pagination, err := h.svc.search2(r.Context(), req.Ingredients, queryParams)
//	if err != nil {
//		pkg.Render(w, r, err)
//		return
//	}
//
//	pkg.Render(w, r, pkg.ApiResponse{
//		Data:       recipes,
//		Message:    "Recipes retrieved successfully",
//		Code:       http.StatusOK,
//		Pagination: pagination,
//	})
//}

package recipe

import (
	"Food/pkg"
	"Food/pkg/recipe/model"
	"Food/pkg/user_preference"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	svc        *Service
	usrPrefSvc *user_preference.Service
}

func NewHandler(svc *Service, usrPrefSvc *user_preference.Service) *Handler {
	return &Handler{
		svc:        svc,
		usrPrefSvc: usrPrefSvc,
	}
}

func (h Handler) crawl(w http.ResponseWriter, r *http.Request) {
	recipeList, err := h.svc.crawl(r.Context())
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    recipeList,
		Message: "Recipes retrieved successfully",
		Code:    http.StatusOK,
	})
}

func (h Handler) like(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	queryParams := r.URL.Query()
	like := queryParams.Get("liked") == "true"

	userID := r.Context().Value("user_id").(string)
	req := user_preference.AddRequest{RecipeIds: []string{id}}

	err := h.usrPrefSvc.Save(r.Context(), userID, req, like)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	// return a success response to users
	pkg.Render(w, r, pkg.ApiResponse{
		Message: "successful",
		Code:    200,
	})
}

// addRecipe is a HTTP handler for adding a new recipe.
func (h Handler) save(w http.ResponseWriter, r *http.Request) {
	var recipes model.Request
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
	// Extract the userID from the request, for example, from a query parameter
	userID := r.Context().Value("user_id").(string)

	// Extract the recipeName from the request query parameters
	recipeName := r.URL.Query().Get("recipe_name")

	recipes, err := h.svc.list(r.Context(), userID, recipeName)
	if err != nil {
		// If an error occurs, send an appropriate HTTP response
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    recipes,
		Message: "Recipes retrieved successfully",
		Code:    http.StatusOK,
	})
}

type SearchRequest struct {
	Ingredients []string `json:"ingredients"`
}

func (h Handler) search(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("user_id")

	id := ""
	if userID != nil {
		id = userID.(string)
	}

	var req SearchRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Infoln("Searching for recipes with ingredients: ", req.Ingredients)

	queryParams := r.URL.Query()
	recipes, pagination, err := h.svc.search(r.Context(), req.Ingredients, queryParams, id)
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

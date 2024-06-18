package recipe

import (
	"Food/pkg"
	"Food/pkg/recipe/model"
	"Food/pkg/user_preference"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
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
		Message: "Recipes crawled successfully",
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

	pkg.Render(w, r, pkg.ApiResponse{
		Message: "Recipe like status updated successfully",
		Code:    http.StatusOK,
	})
}

func (h Handler) save(w http.ResponseWriter, r *http.Request) {
	var recipes model.Request

	if err := render.Bind(r, &recipes); err != nil {
		pkg.Render(w, r, err)
		return
	}

	successMap, err := h.svc.save(r.Context(), recipes)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

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

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    recipe,
		Message: "Recipe successfully deleted",
		Code:    http.StatusOK,
	})
}

func (h Handler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	recipe, err := h.svc.get(r.Context(), id)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    recipe,
		Message: "Recipe retrieved successfully",
		Code:    http.StatusOK,
	})
}

func (h Handler) list(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	recipeName := r.URL.Query().Get("recipe_name")

	recipes, err := h.svc.list(r.Context(), userID, recipeName)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    recipes,
		Message: "Recipes retrieved successfully",
		Code:    http.StatusOK,
	})
}

func (h Handler) search(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")

	id := ""
	if userID != nil {
		id = userID.(string)
	}

	var req SearchRequest

	if err := render.Bind(r, &req); err != nil {
		pkg.Render(w, r, err)
		return
	}

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

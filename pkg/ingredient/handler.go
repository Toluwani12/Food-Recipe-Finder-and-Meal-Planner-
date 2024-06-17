package ingredient

import (
	"Food/pkg"
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

//func (h Handler) add(w http.ResponseWriter, r *http.Request) {
//
//	var data AddRequest
//	if err := render.Bind(r, &data); err != nil {
//		pkg.Render(w, r, err)
//		return
//	}
//
//	ingredient, err := h.svc.save(r.Context(), data)
//	if err != nil {
//		pkg.Render(w, r, err)
//		return
//	}
//
//	pkg.Render(w, r, pkg.ApiResponse{
//		Data:    ingredient,
//		Message: "ingredient added successfully",
//		Code:    201,
//	})
//}

//func (h Handler) delete(w http.ResponseWriter, r *http.Request) {
//
//	id := chi.URLParam(r, "id")
//
//	ingredient, err := h.svc.delete(r.Context(), id)
//	if err != nil {
//		pkg.Render(w, r, err)
//		return
//	}
//
//	pkg.Render(w, r, pkg.ApiResponse{
//		Data:    ingredient,
//		Message: "ingredient successfully deleted",
//		Code:    200,
//	})
//}
//
//func (h Handler) update(w http.ResponseWriter, r *http.Request) {
//	id := chi.URLParam(r, "id")
//
//	var data AddRequest
//
//	if err := render.Bind(r, &data); err != nil {
//		pkg.Render(w, r, err)
//		return
//	}
//
//	ingredient, err := h.svc.update(r.Context(), id, data)
//	if err != nil {
//		pkg.Render(w, r, err)
//		return
//	}
//
//	// Returning the response to the users
//	pkg.Render(w, r, pkg.ApiResponse{
//		Data:    ingredient,
//		Message: "ingredient updated successfully",
//		Code:    200,
//	})
//}
//
//func (h Handler) get(w http.ResponseWriter, r *http.Request) {
//	id := chi.URLParam(r, "id")
//
//	ingredient, err := h.svc.get(r.Context(), id)
//	if err != nil {
//		pkg.Render(w, r, err)
//		return
//	}
//
//	pkg.Render(w, r, pkg.ApiResponse{
//		Data:    ingredient,
//		Message: "Ingredient retrieved successfully",
//		Code:    http.StatusOK,
//	})
//}

func (h Handler) list(w http.ResponseWriter, r *http.Request) {
	ingredient := r.URL.Query().Get("ingredient-name")

	ingredients, err := h.svc.list(r.Context(), ingredient)
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

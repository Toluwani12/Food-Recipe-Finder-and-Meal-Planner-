package user_preference

import (
	"Food/pkg"
	"github.com/go-chi/render"
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
	userID := r.Context().Value("user_id").(string)

	var data AddRequest
	if err := render.Bind(r, &data); err != nil {
		pkg.Render(w, r, err)
		return
	}

	err := h.svc.Save(r.Context(), userID, data, true)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Message: "User preference added successfully",
		Code:    http.StatusCreated,
	})
}

func (h Handler) delete(w http.ResponseWriter, r *http.Request) {
	var data AddRequest
	if err := render.Bind(r, &data); err != nil {
		pkg.Render(w, r, err)
		return
	}

	userID := r.Context().Value("user_id").(string)

	err := h.svc.Save(r.Context(), userID, data, false)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Message: "User preference removed successfully",
		Code:    http.StatusOK,
	})
}

func (h Handler) get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	preference, err := h.svc.Get(r.Context(), userID)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    preference,
		Message: "User preference retrieved successfully",
		Code:    http.StatusOK,
	})
}

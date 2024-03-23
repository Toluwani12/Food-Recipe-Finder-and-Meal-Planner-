package user

import (
	liberror "Food/internal/errors"
	"Food/pkg"
	"github.com/go-chi/chi/v5"
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

func (h Handler) register(w http.ResponseWriter, r *http.Request) {
	var addRequest AddRequest
	if err := addRequest.Bind(r); err != nil {
		pkg.Render(w, r, err)
		return
	}

	user, err := h.svc.save(r.Context(), addRequest)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    user,
		Message: "user registered successfully",
		Code:    201,
	})
}

func (h Handler) login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	if err := render.Bind(r, &request); err != nil {
		err := render.Render(w, r, liberror.ErrInvalidLogin)
		if err != nil {
			return
		}
		return
	}
	h.svc.login(r.Context(), request)
}

func (h Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var userLogin User
	if err := render.Bind(r, &userLogin); err != nil {
		render.Render(w, r, liberror.ErrInvalidLogin)
		return
	}

	// call the delete service to delete by id
	user, err := h.svc.delete(r.Context(), id)
	if err != nil {
		pkg.Render(w, r, nil)
		return
	}

	// return a success response to user
	pkg.Render(w, r, pkg.ApiResponse{
		Data:    user,
		Message: "user successfully deleted",
		Code:    200,
	})

}

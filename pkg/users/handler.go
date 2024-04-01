package users

import (
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

func (h Handler) save(w http.ResponseWriter, r *http.Request) {
	var addRequest AddRequest
	if err := render.Bind(r, &addRequest); err != nil {
		pkg.Render(w, r, err)
		return
	}

	user, err := h.svc.save(r.Context(), addRequest)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    user.Response(),
		Message: "users registered successfully",
		Code:    201,
	})
}

func (h Handler) login(w http.ResponseWriter, r *http.Request) {
	loginReq := LoginRequest{}

	if err := render.Bind(r, &loginReq); err != nil {
		pkg.Render(w, r, err)
		return
	}
	user, s, err := h.svc.login(r.Context(), loginReq)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data: struct {
			User  UserResponse `json:"user"`
			Token string       `json:"token"`
		}{User: user.Response(), Token: s},
		Message: "users registered successfully",
		Code:    200,
	})

}

func (h Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	delId, err := h.svc.delete(r.Context(), id)
	if err != nil {
		pkg.Render(w, r, nil)
		return
	}

	// return a success response to users
	pkg.Render(w, r, pkg.ApiResponse{
		Data:    delId,
		Message: "users successfully deleted",
		Code:    200,
	})

}

func (h Handler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := h.svc.get(r.Context(), id)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    user.Response(),
		Message: "users retrieved successfully",
		Code:    200,
	})
}

func (h Handler) list(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.list(r.Context())
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    users.Response(),
		Message: "users retrieved successfully",
		Code:    200,
	})
}

func (h Handler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var updateRequest UpdateRequest
	if err := render.Bind(r, &updateRequest); err != nil {
		pkg.Render(w, r, err)
		return
	}

	uId, err := h.svc.update(r.Context(), id, updateRequest)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    uId,
		Message: "users updated successfully",
		Code:    200,
	})
}

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
		Message: "User registered successfully",
		Code:    http.StatusCreated,
	})
}

func (h Handler) login(w http.ResponseWriter, r *http.Request) {
	loginReq := LoginRequest{}

	if err := render.Bind(r, &loginReq); err != nil {
		pkg.Render(w, r, err)
		return
	}
	user, token, err := h.svc.login(r.Context(), loginReq)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data: struct {
			User  UserResponse `json:"user"`
			Token string       `json:"token"`
		}{User: user.Response(), Token: token},
		Message: "User logged in successfully",
		Code:    http.StatusOK,
	})
}

func (h Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	delId, err := h.svc.delete(r.Context(), id)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    delId,
		Message: "User successfully deleted",
		Code:    http.StatusOK,
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
		Message: "User retrieved successfully",
		Code:    http.StatusOK,
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
		Message: "Users retrieved successfully",
		Code:    http.StatusOK,
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
		Message: "User updated successfully",
		Code:    http.StatusOK,
	})
}

package user_preference

import (
	"Food/pkg"
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

func (h Handler) add(w http.ResponseWriter, r *http.Request) {
	//extract the user id from the request
	userID := r.Context().Value("user_id").(string)
	log.Infoln("user id", userID)
	// binding or extracting the data from the request
	var data AddRequest
	if err := render.Bind(r, &data); err != nil {
		pkg.Render(w, r, err)
		return
	}

	// add the data through the add service
	err := h.svc.save(r.Context(), userID, data)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	// returning the response to the users
	pkg.Render(w, r, pkg.ApiResponse{
		Message: "user preference added successfully",
		Code:    201,
	})
}

func (h Handler) delete(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	err := h.svc.delete(r.Context(), id)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	// return a success response to users
	pkg.Render(w, r, pkg.ApiResponse{
		Message: "user preference deleted successfully",
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

	err := h.svc.update(r.Context(), id, data)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Message: "recipe updated successfully",
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

//func (h Handler) list(w http.ResponseWriter, r *http.Request) {
//	recipe, err := h.svc.list(r.Context())
//	if err != nil {
//		// If an errors occurs, send an appropriate HTTP response
//		http.Error(w, err.Error(), http.StatusNotFound)
//		return
//	}
//
//	pkg.Render(w, r, pkg.ApiResponse{
//		Data:    recipe,
//		Message: "Recipe retrieved successfully",
//		Code:    http.StatusOK,
//	})
//}

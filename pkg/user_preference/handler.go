package user_preference

import (
	"Food/pkg"
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
	err := h.svc.Save(r.Context(), userID, data, true)
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

	// return a success response to users
	pkg.Render(w, r, pkg.ApiResponse{
		Message: "user preference removed successfully",
		Code:    200,
	})
}

func (h Handler) get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	// Get the recipe from the service layer
	recipe, err := h.svc.get(r.Context(), userID)
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

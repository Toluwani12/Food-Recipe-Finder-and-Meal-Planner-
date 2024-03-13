package ingredient

import (
	"Food/pkg"
	"github.com/go-chi/chi/v5"
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

	// binding or extracting the data from the request
	var data AddRequest
	if err := data.Bind(r); err != nil {
		pkg.Render(w, r, nil)
		return
	}

	// add the data through the add service
	ingredient, err := h.svc.add(r.Context(), data)
	if err != nil {
		pkg.Render(w, r, nil)
		return
	}

	// returning the response to the user
	pkg.Render(w, r, pkg.ApiResponse{
		Data:    ingredient,
		Message: "ingredient added successfully",
		Code:    201,
	})
}

func (h Handler) delete(w http.ResponseWriter, r *http.Request) {

	// extract id from url param
	id := chi.URLParam(r, "id")

	// call the delete service to delete by id
	ingredient, err := h.svc.delete(r.Context(), id)
	if err != nil {
		pkg.Render(w, r, nil)
		return
	}

	// return a success response to user
	pkg.Render(w, r, pkg.ApiResponse{
		Data:    ingredient,
		Message: "ingredient successfully deleted",
		Code:    200,
	})
}

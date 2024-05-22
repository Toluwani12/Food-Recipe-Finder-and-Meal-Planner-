package mealplan

import (
	"Food/pkg"
	"net/http"
	"time"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) generate(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("user_id").(string)
	weekStartDate := getStartOfWeek()

	mealPlans, err := h.svc.generateMealPlans(userID, weekStartDate)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    mealPlans,
		Message: "Meal plans generated successfully",
		Code:    http.StatusOK,
	})
}

func getStartOfWeek() time.Time {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}
	startOfWeek := now.AddDate(0, 0, -weekday+1)
	return time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
}

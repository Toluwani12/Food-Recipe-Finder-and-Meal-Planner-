package mealplan

import (
	"Food/internal/errors"
	"Food/pkg"
	"encoding/json"
	"github.com/google/uuid"
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

	placeholders, err := h.svc.generateMealPlans(r.Context(), uuid.MustParse(userID), weekStartDate)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    placeholders,
		Message: "Meal plans generated successfully",
		Code:    http.StatusOK,
	})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("user_id").(string)
	weekStartDate := getStartOfWeek()

	placeholders, err := h.svc.getMealPLan(userID, weekStartDate)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	pkg.Render(w, r, pkg.ApiResponse{
		Data:    placeholders,
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

func (h *Handler) GetMealPlansForDay(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	dayOfWeekStr := r.URL.Query().Get("day_of_week")
	if dayOfWeekStr == "" {
		pkg.Render(w, r, errors.New("day_of_week query parameter is required", http.StatusBadRequest))
		return
	}
	dayOfWeek := DayOfWeek(dayOfWeekStr)
	weekStartDate := getStartOfWeek()

	mealPlans, err := h.svc.GetMealPlansForDay(userID, dayOfWeek, weekStartDate)
	if err != nil {
		pkg.Render(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mealPlans)
}

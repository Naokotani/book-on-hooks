package metrics

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/logger"
	"booksonhooks.ca/internal/requestx"
)

type Repo interface {
	GetMetricsDashboard(ctx context.Context, month string, startDate, endDate time.Time, qr *bool) (*domain.MetricsDashboard, error)
}

type Handler struct {
	log  *logger.Logger
	repo Repo
}

func New(log *logger.Logger, repo Repo) *Handler {
	return &Handler{
		log:  log,
		repo: repo,
	}
}

func (h *Handler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	month, startDate, endDate, err := parseMonth(r.URL.Query().Get("month"))
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "month must use YYYY-MM format",
		})
		return
	}

	qr, err := parseOptionalBoolPtr(r.URL.Query().Get("qr"))
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "qr must be true or false",
		})
		return
	}

	dashboard, err := h.repo.GetMetricsDashboard(r.Context(), month, startDate, endDate, qr)
	if err != nil {
		h.log.Error("failed to get metrics dashboard: %v", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to get metrics dashboard",
		})
		return
	}

	h.writeJSON(w, http.StatusOK, dashboard)
}

func parseMonth(value string) (string, time.Time, time.Time, error) {
	if value == "" {
		now := time.Now()
		value = now.Format("2006-01")
	}

	startDate, err := time.Parse("2006-01", value)
	if err != nil {
		return "", time.Time{}, time.Time{}, err
	}

	return value, startDate, startDate.AddDate(0, 1, 0), nil
}

func parseOptionalBoolPtr(value string) (*bool, error) {
	parsed, err := requestx.ParseOptionalBool(value)
	if err != nil {
		return nil, err
	}
	if value == "" {
		return nil, nil
	}
	return &parsed, nil
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

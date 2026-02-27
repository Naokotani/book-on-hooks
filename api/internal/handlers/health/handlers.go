package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"booksonhooks.ca/internal/logger"
)

type HealthRepo interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	log  *logger.Logger
	repo HealthRepo
}

func New(log *logger.Logger, repo HealthRepo) *HealthHandler {
	return &HealthHandler{
		log:  log,
		repo: repo,
	}
}

func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	if err := h.repo.Ping(ctx); err != nil {
		h.log.Warn("health check db ping failed: %v", err)
		h.writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"status":  "degraded",
			"service": "api",
			"error":   "db unavailable",
		})
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"service": "api",
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *HealthHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

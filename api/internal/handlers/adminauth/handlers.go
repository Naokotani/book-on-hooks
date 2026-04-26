package adminauth

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"booksonhooks.ca/internal/logger"
	"github.com/alexedwards/scs/v2"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	log      *logger.Logger
	sessions *scs.SessionManager
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func New(log *logger.Logger, sessions *scs.SessionManager) *Handler {
	return &Handler{
		log:      log,
		sessions: sessions,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var payload loginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(payload.Username)
	password := payload.Password
	expectedUsername := strings.TrimSpace(os.Getenv("ADMIN_USERNAME"))
	passwordHash := os.Getenv("ADMIN_PASSWORD_HASH")

	if username == "" || password == "" || expectedUsername == "" || passwordHash == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if username != expectedUsername || bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if err := h.sessions.RenewToken(r.Context()); err != nil {
		h.log.Error("failed to renew session token: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	h.sessions.Put(r.Context(), "is_admin", true)
	h.sessions.Put(r.Context(), "admin_username", expectedUsername)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.sessions.Destroy(r.Context()); err != nil {
		h.log.Error("failed to destroy session: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	authenticated := h.sessions.GetBool(r.Context(), "is_admin")
	username := ""
	if authenticated {
		username = h.sessions.GetString(r.Context(), "admin_username")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"authenticated": authenticated,
		"username":      username,
	})
}

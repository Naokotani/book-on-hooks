package app

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func newSessionManager(pool *pgxpool.Pool) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(stdlib.OpenDBFromPool(pool))
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Path = "/"
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = useSecureCookies()

	return sessionManager
}

func useSecureCookies() bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv("COOKIE_SECURE")))
	return value == "1" || value == "true" || value == "yes"
}

package app

import (
	"net/http"

	"booksonhooks.ca/internal/httpErrors"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *Application) Routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpErrors.NotFound(w)
	})

	// Health routes
	router.HandlerFunc(http.MethodGet, "/api/healthz", app.healthHandlers.Healthz)

	// Admin auth routes
	router.HandlerFunc(http.MethodPost, "/api/admin/login", app.authHandlers.Login)
	router.HandlerFunc(http.MethodPost, "/api/admin/logout", app.authHandlers.Logout)
	router.HandlerFunc(http.MethodGet, "/api/admin/me", app.authHandlers.Me)
	router.Handler(http.MethodGet, "/api/admin/metrics", app.requireAdmin(http.HandlerFunc(app.metricsHandlers.GetDashboard)))

	// Book routes
	router.HandlerFunc(http.MethodGet, "/api/books", app.bookHandlers.GetBooks)
	router.HandlerFunc(http.MethodGet, "/api/books/book/:id", app.bookHandlers.GetBook)
	router.HandlerFunc(http.MethodGet, "/api/books/summary/:id", app.bookHandlers.GetBookSummary)
	router.HandlerFunc(http.MethodGet, "/api/books/images/:image", app.bookHandlers.GetBookImage)

	// Machine routes
	router.HandlerFunc(http.MethodGet, "/api/machines", app.machineHandlers.GetMachines)
	router.HandlerFunc(http.MethodGet, "/api/machines/:id", app.machineHandlers.GetMachine)
	router.HandlerFunc(http.MethodGet, "/api/machines/:id/books", app.machineHandlers.GetMachineWithBooks)

	// Admin book routes
	router.Handler(http.MethodPost, "/api/admin/books", app.requireAdmin(http.HandlerFunc(app.bookHandlers.CreateBook)))
	router.Handler(http.MethodPatch, "/api/admin/books/:id", app.requireAdmin(http.HandlerFunc(app.bookHandlers.UpdateBook)))
	router.Handler(http.MethodPatch, "/api/admin/books/:id/image", app.requireAdmin(http.HandlerFunc(app.bookHandlers.UpdateBookImage)))
	router.Handler(http.MethodDelete, "/api/admin/books/:id", app.requireAdmin(http.HandlerFunc(app.bookHandlers.DeleteBook)))

	// Admin machine routes
	router.Handler(http.MethodPost, "/api/admin/machines", app.requireAdmin(http.HandlerFunc(app.machineHandlers.CreateMachine)))
	router.Handler(http.MethodPatch, "/api/admin/machines/:id", app.requireAdmin(http.HandlerFunc(app.machineHandlers.UpdateMachine)))
	router.Handler(http.MethodPatch, "/api/admin/machines/:id/grid", app.requireAdmin(http.HandlerFunc(app.machineHandlers.UpdateMachineRowsCols)))
	router.Handler(http.MethodDelete, "/api/admin/machines/:id", app.requireAdmin(http.HandlerFunc(app.machineHandlers.DeleteMachine)))
	router.Handler(http.MethodDelete, "/api/admin/machines/:id/books", app.requireAdmin(http.HandlerFunc(app.machineHandlers.ClearMachineBooks)))
	router.Handler(http.MethodPut, "/api/admin/machines/:id/books", app.requireAdmin(http.HandlerFunc(app.machineHandlers.LoadMachine)))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return app.sessions.LoadAndSave(standard.Then(router))
}

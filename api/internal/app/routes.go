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

	// Book routes
	router.HandlerFunc(http.MethodGet, "/api/books", app.bookHandlers.GetBooks)
	router.HandlerFunc(http.MethodPost, "/api/books", app.bookHandlers.CreateBook)
	router.HandlerFunc(http.MethodGet, "/api/books/:id", app.bookHandlers.GetBook)
	router.HandlerFunc(http.MethodGet, "/api/books/:id/locations", app.bookHandlers.GetBookLocations)
	router.HandlerFunc(http.MethodGet, "/api/images/:image", app.bookHandlers.GetImage)
	router.HandlerFunc(http.MethodPatch, "/api/books/:id", app.bookHandlers.UpdateBook)
	router.HandlerFunc(http.MethodDelete, "/api/books/:id", app.bookHandlers.DeleteBook)

	// Machine routes
	router.HandlerFunc(http.MethodGet, "/api/machines", app.machineHandlers.MachinesView)
	router.HandlerFunc(http.MethodPost, "/api/machines", app.machineHandlers.MachineCreate)
	router.HandlerFunc(http.MethodGet, "/api/machines/:id", app.machineHandlers.GetMachine)
	router.HandlerFunc(http.MethodPatch, "/api/machines/:id", app.machineHandlers.UpdateMachine)
	router.HandlerFunc(http.MethodDelete, "/api/machines/:id", app.machineHandlers.DeleteMachine)
	router.HandlerFunc(http.MethodDelete, "/api/machines/:id/books", app.machineHandlers.ClearMachineBooks)
	router.HandlerFunc(http.MethodGet, "/api/machines/:id/books", app.machineHandlers.GetMachineWithBooks)
	router.HandlerFunc(http.MethodPut, "/api/machines/:id/books", app.machineHandlers.LoadMachine)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}

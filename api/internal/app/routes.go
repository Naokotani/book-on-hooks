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

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Home routes
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	router.HandlerFunc(http.MethodGet, "/", app.homeHandlers.Home)
	router.HandlerFunc(http.MethodGet, "/covers/:row/:col", app.bookHandlers.Cover)
	router.HandlerFunc(http.MethodGet, "/book/:row/:col", app.bookHandlers.Book)

	// Book routes
	router.HandlerFunc(http.MethodGet, "/api/books", app.bookHandlers.GetBooks)
	router.HandlerFunc(http.MethodPost, "/api/admin/book/create", app.bookHandlers.CreateBook)
	router.HandlerFunc(http.MethodGet, "/api/book/find/:id", app.bookHandlers.GetBook)
	router.HandlerFunc(http.MethodGet, "/api/books/:id", app.bookHandlers.GetBookByIDJSON)
	router.HandlerFunc(http.MethodPatch, "/api/books/:id", app.bookHandlers.UpdateBook)
	router.HandlerFunc(http.MethodDelete, "/api/books/:id", app.bookHandlers.DeleteBook)

	// Machine routes
	router.HandlerFunc(http.MethodGet, "/api/machines", app.machineHandlers.MachinesView)
	router.HandlerFunc(http.MethodPost, "/api/admin/machine/create", app.machineHandlers.MachineCreate)
	router.HandlerFunc(http.MethodPost, "/api/admin/machine/load", app.machineHandlers.LoadMachine)
	router.HandlerFunc(http.MethodGet, "/api/machines/:id", app.machineHandlers.GetMachine)
	router.HandlerFunc(http.MethodPatch, "/api/machines/:id", app.machineHandlers.UpdateMachine)
	router.HandlerFunc(http.MethodDelete, "/api/machines/:id", app.machineHandlers.DeleteMachine)
	router.HandlerFunc(http.MethodDelete, "/api/machines/:id/books", app.machineHandlers.ClearMachineBooks)
	router.HandlerFunc(http.MethodGet, "/api/machines/:id/books", app.machineHandlers.GetMachineWithBooks)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}

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

	// Home routes
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	router.HandlerFunc(http.MethodGet, "/", app.homeHandlers.Home)
	router.HandlerFunc(http.MethodGet, "/covers/:row/:col", app.bookHandlers.Cover)
	router.HandlerFunc(http.MethodGet, "/book/:row/:col", app.bookHandlers.Book)

	// Book routes
	router.HandlerFunc(http.MethodGet, "/api/books", app.bookHandlers.GetBooks)
	router.HandlerFunc(http.MethodPost, "/api/book/create", app.bookHandlers.CreateBook)
	router.HandlerFunc(http.MethodGet, "/api/book/create", app.bookHandlers.BookCreateView)
	router.HandlerFunc(http.MethodGet, "/api/book/find/:id", app.bookHandlers.GetBook)

	// Machine routes
	router.HandlerFunc(http.MethodGet, "/api/machines", app.machineHandlers.MachinesView)
	router.HandlerFunc(http.MethodPost, "/api/machine/create", app.machineHandlers.MachineCreate)
	router.HandlerFunc(http.MethodGet, "/api/machine/create", app.machineHandlers.MachineCreateView)
	router.HandlerFunc(http.MethodGet, "/api/machine/load/:id", app.machineHandlers.MachineLoadView)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}

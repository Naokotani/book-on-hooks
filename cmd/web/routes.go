package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/covers/:row/:col", app.cover)
	router.HandlerFunc(http.MethodGet, "/book/:row/:col", app.book)

	// Book creation

	router.HandlerFunc(http.MethodPost, "/inventory/book/create", app.createBook)
	router.HandlerFunc(http.MethodGet, "/inventory/book/create", app.bookCreateView)

	// Machine creation and management

	router.HandlerFunc(http.MethodGet, "/inventory/machines", app.machinesView)
	router.HandlerFunc(http.MethodPost, "/inventory/machine/create", app.machineCreate)
	router.HandlerFunc(http.MethodGet, "/inventory/machine/create", app.machineCreateView)
	router.HandlerFunc(http.MethodPost, "/inventory/machine/:id", app.machineUpdate)
	router.HandlerFunc(http.MethodGet, "/inventory/machine/:id", app.machineCreateView)

	router.HandlerFunc(http.MethodGet, "/invenrtory/book/:id", app.getBook)
	router.HandlerFunc(http.MethodDelete, "/inventory/boo/:id/", app.deleteBook)
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}

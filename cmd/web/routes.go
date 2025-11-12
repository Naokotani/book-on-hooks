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
	router.HandlerFunc(http.MethodGet, "/covers/:row/:col", app.covers)
	router.HandlerFunc(http.MethodGet, "/book/:row/:col", app.book)
	router.HandlerFunc(http.MethodPost, "/inventory/book/create", app.createBook)
	router.HandlerFunc(http.MethodGet, "/invenrtory/book/:id", app.getBook)
	router.HandlerFunc(http.MethodDelete, "/inventory/boo/:id/", app.deleteBook)
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}

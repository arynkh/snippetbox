package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	//think of ServeMux as a special kind of handler, that passes the request to a second handler
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	//middleware chain used for every request our app receives
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux) //return the standard middleware followed by the ServeMux
}

package main

import (
	"net/http"

	"github.com/arynkh/snippetbox/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	//think of ServeMux as a special kind of handler, that passes the request to a second handler
	mux := http.NewServeMux()

	//create a HTTP handler which serves the embedd files in the ui.Files.
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.HandleFunc("GET /ping", ping)

	//new middleware chain containing the middleware specific to our dynamic application routes.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	//protected (authenticated-only) routes
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	//middleware chain used for every request our app receives
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux) //return the standard middleware followed by the ServeMux
}

package main

import (
	"net/http"
	"runtime/debug"
)

// helps write a log entry at Error level (including the reques method and URI as attributes), then sends a generic 500 Internal Server Error response to the client/user
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// helps send a specific status code & corresponding description to the user. used to send responses like 404 Bad Request when theres a problem with the request that the user sent
// func (app *application) clientError(w http.ResponseWriter, r *http.Request, status int) {
// 	http.Error(w, http.StatusText(status), status)
// }

package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
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
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	//retrieve the appropriate template set from the cache based on the page name.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status) //write out the provided HTTP status code ('200 OK', '400 Bad Request', etc).

	buf.WriteTo(w) //write the contents of the buffer to the http.ResponseWriter
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"), //retrieve the flash message from the session and remove it
		IsAuthenticated: app.isAuthenticated(r),                             //check if the user is authenticated
		CSRFToken:       nosurf.Token(r),                                    //get the CSRF token from the request
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	//call Decode on the formDecoder field of the app struct, passing in the destination as the first parameter
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderErr *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderErr) {
			panic(err)
		}
		return err
	}

	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

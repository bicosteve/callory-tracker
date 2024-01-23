package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError() writes error messages and trace to the error line
// Also send generic 500 internal server error
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace) // returns which line error is.
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError() -> sends specific status code and description to user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound() returns 404 not found response to user
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// renderATemplate() -> renders specific template in the specific handler
func (app *application) renderATemplate(
	w http.ResponseWriter, r *http.Request, templateName string, data *templateData,
) {
	/*
		Get the appropriate template set from the cache base on page name eg
		home.page.html. If it does not exist in the cache with provided name,
		call the serverError helper method
	*/

	ts, ok := app.templateCache[templateName]
	if !ok {
		app.serverError(w, fmt.Errorf("template %s does not exist", templateName))
		return
	}

	// Executes template set, passing in any dynamic data.
	err := ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}

}

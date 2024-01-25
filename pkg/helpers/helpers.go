package helpers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime/debug"

	"github.com/bicosteve/callory-tracker/pkg/configs"
)

type Application configs.Application

func LoadEnv(name string) (string, error) {
	file, err := filepath.Abs(name)
	if err != nil {
		return "", err
	}

	return file, nil
}

func (app *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace) // returns which line error is.
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func CheckFormInput(input int) bool {
	if input < 1 {
		return false
	}

	return true
}

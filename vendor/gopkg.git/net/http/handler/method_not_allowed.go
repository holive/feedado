package handler

import (
	"net/http"

	"github.com/pkg/errors"
)

// MethodNotAllowed returns a http.Handler that follows our corporative guideline.
func MethodNotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, errors.New("Método não permitido."), http.StatusMethodNotAllowed)
	}
}

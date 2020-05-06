package handler

import (
	"net/http"

	"github.com/pkg/errors"
)

// NotFound returns a http.Handler that follows our corporative guideline.
func NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, errors.New("Recurso não encontrado."), http.StatusNotFound)
	}
}

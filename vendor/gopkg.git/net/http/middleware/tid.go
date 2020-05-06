package middleware

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// NewTID is a middleware that looks for a XTID value inside the http.Request
// and generate one if it does not exists.
func NewTID(appname string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tid := r.Header.Get("x-tid")
			if tid == "" {
				tid = appname + "-" + uuid.NewV4().String()
				r.Header.Set("x-tid", tid)
			}
			w.Header().Set("x-tid", tid)
			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

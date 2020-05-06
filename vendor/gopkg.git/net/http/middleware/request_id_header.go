package middleware

import (
	"context"
	"net/http"
)

// RequestIdHeader is a middleware that creates a unique request id
func RequestIdHeader(
	getRequestId func(context.Context) string,
	headerName string,
) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			requestId := getRequestId(ctx)
			if requestId != "" {
				w.Header().Set("x-"+headerName, requestId)
			}
			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

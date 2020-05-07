package http

import (
	"context"
	"net/http"
	"strings"
)

// AtlasRetryPolicy provider a callback for Client.CheckRetry, which
// will retry on connection errors and server errors.
// (including workarounds for atlas backends)
func AtlasRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if resp != nil && resp.StatusCode == http.StatusNotFound &&
		!strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		return true, nil
	}

	return DefaultRetryPolicy(ctx, resp, err)
}

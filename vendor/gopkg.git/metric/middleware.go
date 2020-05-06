package metric

import (
	"fmt"
	"net/http"
	"strings"
)

// Middleware is the http handling middleware for metric.
func (c *Client) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			url := r.URL.String()
			path := r.URL.Path

			txn := c.StartTransaction(path, w, r)
			defer txn.End()

			txn.AddAttribute("request.url", fmt.Sprintf("http://%s%s", r.Host, url))

			qs := r.URL.Query()
			for key, value := range qs {
				txn.AddAttribute(key, strings.Join(value, "|"))
			}

			next.ServeHTTP(w, r)
		})
}

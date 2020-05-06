package endpoint

import (
	"net/http"

	"github.com/pkg/errors"

	"gopkg.git/health"
	infraHTTP "gopkg.git/net/http"
)

// Health exposes a handler to monitor application health status.
type Health struct {
	writer  *infraHTTP.Writer
	checker *health.Health
}

// Handler accepts http requests about the application health status.
func (h *Health) Handler(w http.ResponseWriter, r *http.Request) {
	healthResult := h.checker.Check(r.Context())

	var statusCode = http.StatusServiceUnavailable
	if healthResult.Status == health.StatusOK || healthResult.Status == health.StatusPartial {
		statusCode = http.StatusOK
	}

	h.writer.Response(w, healthResult, statusCode, nil)
}

// NewHealth returns a initialized health handler.
func NewHealth(options ...func(*Health)) (*Health, error) {
	h := &Health{}

	for _, option := range options {
		option(h)
	}

	if h.writer == nil {
		return nil, errors.New("invalid writer")
	}

	if h.checker == nil {
		return nil, errors.New("invalid checker")
	}

	return h, nil
}

// HealthWriter defines the http.Writer that will be used during response.
func HealthWriter(writer *infraHTTP.Writer) func(*Health) {
	return func(h *Health) {
		h.writer = writer
	}
}

// HealthChecker defines the checker that will be asked during health check.
func HealthChecker(checker *health.Health) func(*Health) {
	return func(h *Health) {
		h.checker = checker
	}
}

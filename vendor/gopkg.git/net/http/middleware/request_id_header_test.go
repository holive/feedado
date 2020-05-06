package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/middleware"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRequestId(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	requestNewIdMiddleware := RequestIdHeader(middleware.GetReqID, "request-id")(next)
	requestIdMiddleware := middleware.RequestID(requestNewIdMiddleware)

	Convey("Given a request, generate a x-request-id header", t, func() {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/search", nil)
		requestIdMiddleware.ServeHTTP(w, r)
		requestIdHeader := w.Header().Get("x-request-id")
		So(requestIdHeader, ShouldNotBeEmpty)
	})
}

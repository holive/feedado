package endpoint

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-kit/kit/log"

	. "github.com/smartystreets/goconvey/convey"

	"gopkg.git/health"
	infraHTTP "gopkg.git/net/http"
	infraHTTPTest "gopkg.git/net/http/test"
	infraTest "gopkg.git/test"
)

type TestChecker struct {
	Result health.DependencyResult
}

func (c TestChecker) Check(ctx context.Context) health.DependencyResult {
	time.Sleep(40 * time.Millisecond)
	return c.Result
}

func TestNewHealth(t *testing.T) {
	Convey("Given a list of options", t, func() {
		tests := []struct {
			name    string
			options []func(*Health)
			err     error
		}{
			{
				"It should not have errors",
				[]func(*Health){
					HealthChecker(&health.Health{}),
					HealthWriter(&infraHTTP.Writer{}),
				},
				nil,
			},
			{
				"It should error about invalid writer",
				[]func(*Health){
					HealthChecker(&health.Health{}),
				},
				errors.New("invalid writer"),
			},
			{
				"It should error about invalid Health",
				[]func(*Health){
					HealthWriter(&infraHTTP.Writer{}),
				},
				errors.New("invalid checker"),
			},
		}

		for _, tt := range tests {
			Convey(tt.name, func() {
				h, err := NewHealth(tt.options...)
				if tt.err == nil {
					So(err, ShouldBeNil)
					So(h, ShouldNotBeNil)
					return
				}

				So(h, ShouldBeNil)
				So(tt.err.Error(), ShouldEqual, err.Error())
			})
		}
	})
}

func TestHandlerHealth(t *testing.T) {
	Convey("Given a health request", t, func() {
		tests := []struct {
			title   string
			checker []health.Checker
			req     *http.Request
			status  int
			header  http.Header
			body    []byte
		}{
			{
				"In case of no dependencies it should reply 200",
				[]health.Checker{},
				httptest.NewRequest("GET", "http://application/health", nil),
				http.StatusOK,
				http.Header{"Content-Type": []string{"application/json"}},
				infraTest.Load("healthHandler.empty.json"),
			},
			{
				"In case of success it should reply 200",
				[]health.Checker{TestChecker{health.DependencyResult{
					Status: health.DependencyOK,
				}}},
				httptest.NewRequest("GET", "http://application/health", nil),
				http.StatusOK,
				http.Header{"Content-Type": []string{"application/json"}},
				infraTest.Load("healthHandler.ok.json"),
			},
			{
				"In case of partial success it should reply 200",
				[]health.Checker{TestChecker{health.DependencyResult{
					Status: health.DependencyFail,
				}}},
				httptest.NewRequest("GET", "http://application/health", nil),
				http.StatusOK,
				http.Header{"Content-Type": []string{"application/json"}},
				infraTest.Load("healthHandler.partial.json"),
			},
			{
				"In case of permanent failure it should reply 503",
				[]health.Checker{TestChecker{health.DependencyResult{
					Status:   health.DependencyFail,
					Critical: true,
				}}},
				httptest.NewRequest("GET", "http://application/health", nil),
				http.StatusServiceUnavailable,
				http.Header{"Content-Type": []string{"application/json"}},
				infraTest.Load("healthHandler.fail.json"),
			},
		}

		for _, tt := range tests {
			Convey(tt.title, func() {
				logger := log.NewNopLogger()
				writer, err := infraHTTP.NewWriter(logger)
				So(err, ShouldBeNil)

				h, err := NewHealth(
					HealthChecker(health.NewHealth(tt.checker)),
					HealthWriter(writer),
				)
				So(err, ShouldBeNil)

				infraHTTPTest.Runner(tt.status, tt.header, h.Handler, tt.req, tt.body)
			})
		}
	})
}

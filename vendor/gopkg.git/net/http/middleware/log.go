package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
type Logger struct {
	logger log.Logger
}

// Handler to log each request.
func (l *Logger) Handler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		reqId := middleware.GetReqID(r.Context())
		preReqContent := []interface{}{
			"message", "request started",
			"time", t1,
			"requestId", reqId,
			"method", r.Method,
			"endpoint", r.RequestURI,
			"protocol", r.Proto,
		}

		if r.RemoteAddr != "" {
			preReqContent = append(preReqContent, "ip", r.RemoteAddr)
		}

		tid := r.Header.Get("X-TID")
		if tid != "" {
			preReqContent = append(preReqContent, "tid", tid)
		}

		l.logger.Log(preReqContent...)

		defer func() {
			if err := recover(); err != nil {
				level.Error(l.logger).Log(
					"requestId", reqId,
					"duration", time.Since(t1),
					"status", 500,
					"message", "request finished with panic",
					"stacktrace", string(debug.Stack()),
				)

				panic(err)
			}
		}()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		status := ww.Status()
		postReqContent := []interface{}{
			"requestId", reqId,
			"duration", time.Since(t1),
			"contentLength", ww.BytesWritten(),
			"status", status,
			"message", "request finished",
		}

		if cache := ww.Header().Get("x-cache"); cache != "" {
			postReqContent = append(postReqContent, "cache", cache)
		}

		logger := log.WithPrefix(l.logger, postReqContent...)
		if status >= 100 && status < 400 {
			logger.Log("message", "request finished")
		} else if status == 500 {
			level.Error(logger).Log(
				"stacktrace", string(debug.Stack()),
				"message", "internal error during request",
			)
		} else {
			message := "request finished"

			// FIX: For some reason, the 'context.deadlineExceededError{}' isn't getting into here, we
			// did a quick fix checking the status code and returing the same message as the error., but
			// something is wrong and we need fix it.
			if status == 504 {
				message += ": context deadline exceeded"
			} else {
				if err := r.Context().Err(); err != nil {
					message += fmt.Sprintf(": %s", err.Error())
				}
			}
			level.Error(logger).Log("message", message)
		}
	}

	return http.HandlerFunc(fn)
}

// NewLogger returns a initialized logger middleware.
func NewLogger(logger log.Logger) func(http.Handler) http.Handler {
	l := &Logger{logger}
	return l.Handler
}

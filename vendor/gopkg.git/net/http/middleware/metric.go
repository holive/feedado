package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log/level"

	"github.com/go-kit/kit/log"
	"gopkg.git/metric"
)

// NewMetricer is used to return or generate the metric handler.
func NewMetricer(
	metricer metric.Metricer,
	getRequestId func(context.Context) string,
	log log.Logger,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			txn := metricer.StartTransaction(r.URL.Path, w, r)
			defer txn.End()

			defer func() {
				if err := recover(); err != nil {
					txn.NoticeError(errors.New(fmt.Sprint("recovered from panic: ", err)))
					panic(err)
				}
			}()

			ctx := r.Context()
			metricEnhanceTxnPre(ctx, r, w, txn, log, getRequestId)
			r = r.WithContext(context.WithValue(ctx, metric.ContextTransaction, txn))
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			metricEnhanceTxnPos(txn, log, ww)
		}

		return http.HandlerFunc(fn)
	}
}

func metricSetAttr(txn metric.Transaction, log log.Logger, txnKey string, value interface{}) {
	if err := txn.AddAttribute(txnKey, value); err != nil {
		level.Error(log).Log(
			"error", err.Error(), "message", fmt.Sprintf("failed to add txn attribute '%s'", txnKey),
		)
	}
}

func metricEnhanceTxnPre(
	ctx context.Context,
	r *http.Request,
	w http.ResponseWriter,
	txn metric.Transaction,
	log log.Logger,
	getRequestId func(context.Context) string,
) {
	for key, values := range r.URL.Query() {
		for i, value := range values {
			txnKey := fmt.Sprintf("request.query.%s.%d", key, i)
			metricSetAttr(txn, log, txnKey, value)
		}
	}

	for key, values := range r.Header {
		for i, value := range values {
			txnKey := fmt.Sprintf("request.headers.%s.%d", key, i)
			metricSetAttr(txn, log, txnKey, value)
		}
	}

	if r.RemoteAddr != "" {
		txnKey := fmt.Sprintf("request.header.%s.0", "X-Real-IP")
		metricSetAttr(txn, log, txnKey, r.RemoteAddr)
	}

	if reqId := getRequestId(ctx); reqId != "" {
		txnKey := fmt.Sprintf("httpRequestId")
		metricSetAttr(txn, log, txnKey, reqId)
	}

	metricSetAttr(txn, log, "request.method", r.Method)
}

func metricEnhanceTxnPos(txn metric.Transaction, log log.Logger, w middleware.WrapResponseWriter) {
	metricSetAttr(txn, log, "httpResponseCode", w.Status())
}

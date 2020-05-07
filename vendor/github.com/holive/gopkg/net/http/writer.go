package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

// Writer is used to send the content to the client.
type Writer struct {
	logger log.Logger
}

// Response is used write response on http.ResponseWriter.
func (wrt *Writer) Response(w http.ResponseWriter, r interface{}, status int, headers http.Header) {
	if headers != nil {
		for key, values := range headers {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}

	if r == nil {
		w.WriteHeader(status)
		return
	}

	content, err := json.Marshal(r)
	if err != nil {
		wrt.logger.Log("error", err.Error(), "message", "error during json.Marshal")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	writed, err := w.Write(content)
	if err != nil {
		wrt.logger.Log("error", err.Error(), "message", "error during write at http.ResponseWriter")
	}
	if writed != len(content) {
		wrt.logger.Log(
			"message",
			fmt.Sprintf("invalid quantity of writed bytes, expected %d and got %d", len(content), writed),
		)
	}
}

// Error is used to generate a proper error content to be sent to the client.
func (wrt *Writer) Error(w http.ResponseWriter, err error, status int) {
	resp := struct {
		HTTPStatusCode string `json:"httpStatusCode"`
		ErrorCode      int    `json:"errorCode"`
		Message        string `json:"message"`
	}{
		HTTPStatusCode: strconv.Itoa(status),
		ErrorCode:      status,
	}

	if err != nil {
		resp.Message = err.Error()
	}

	wrt.Response(w, &resp, status, nil)
}

// NewWriter returns a configured writer.
func NewWriter(logger log.Logger) (*Writer, error) {
	if logger == nil {
		return nil, errors.New("logger not found")
	}
	return &Writer{logger}, nil
}

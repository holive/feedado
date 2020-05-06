package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

// stateless functions that should not be used from outside this package.
// borrowed from the parent http package because we don't have logger here.

// writeResponse is used write response on http.ResponseWriter.
func writeResponse(w http.ResponseWriter, r interface{}, status int, headers http.Header) {
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errors.Wrap(err, "error during json.Marshal").Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	w.Write(content)
}

// writeError is used to generate a proper error content to be sent to the client.
func writeError(w http.ResponseWriter, err error, status int) {
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

	writeResponse(w, &resp, status, nil)
}

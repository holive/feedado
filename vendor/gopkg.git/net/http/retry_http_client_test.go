package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	count := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		if count == 5 {
			w.Write([]byte("OK"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("FAIL"))
	}))
	defer ts.Close()

	runner := NewRetryHTTP(http.DefaultClient)
	runner.RetryWaitMax = 1 * time.Nanosecond

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := runner.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	btbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	got := string(btbody)
	if got != "OK" {
		t.Errorf("got %v, want %v", got, "OK")
	}
}

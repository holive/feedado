package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPRetry(t *testing.T) {
	t.Parallel()

	retries := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retries++

		if retries == 2 {
			w.Write([]byte("OK"))
			return
		}

		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			panic(err)
		}
		conn.Close()
	}))
	defer ts.Close()

	runner := NewRetryHTTP(http.DefaultClient)
	runner.CheckRetry = AtlasRetryPolicy

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := runner.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
}

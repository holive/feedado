package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClientHeaders(t *testing.T) {
	Convey("Given some Headers", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bt, err := json.Marshal(r.Header)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(bt)
		}))
		defer ts.Close()

		tests := []struct {
			name           string
			clientHeader   http.Header
			requestHeader  http.Header
			expectedHeader http.Header
		}{
			{
				name:           "It should add the client header into the final request",
				clientHeader:   http.Header{"User-Agent": {"alpha"}},
				expectedHeader: http.Header{"Accept-Encoding": {"gzip"}, "User-Agent": {"alpha"}},
			},
			{
				name:          "It should merge the client request header with the request header",
				clientHeader:  http.Header{"User-Agent": {"alpha"}},
				requestHeader: http.Header{"X-TID": {"alpha-123"}},
				expectedHeader: http.Header{
					"Accept-Encoding": {"gzip"}, "User-Agent": {"alpha"}, "X-Tid": {"alpha-123"},
				},
			},
			{
				name:          "It should not override the request header",
				clientHeader:  http.Header{"User-Agent": {"chrome"}},
				requestHeader: http.Header{"User-Agent": {"firefox"}},
				expectedHeader: http.Header{
					"Accept-Encoding": []string{"gzip"}, "User-Agent": []string{"firefox"},
				},
			},
		}

		for _, tt := range tests {
			Convey(tt.name, func() {
				client, err := NewClient(ClientHeaders(tt.clientHeader), ClientTimeout(1*time.Second))
				So(err, ShouldBeNil)

				req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
				So(err, ShouldBeNil)
				for key, values := range tt.requestHeader {
					for _, value := range values {
						req.Header.Add(key, value)
					}
				}

				resp, err := client.Do(req)
				So(err, ShouldBeNil)

				bt, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)

				requestHeader := make(http.Header)
				err = json.Unmarshal(bt, &requestHeader)
				So(err, ShouldBeNil)

				So(requestHeader, ShouldResemble, tt.expectedHeader)
			})
		}
	})
}

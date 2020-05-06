package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/smartystreets/goconvey/convey"
)

// Runner is used to execute http requests and check if it matchs.
func Runner(
	status int,
	header http.Header,
	handler func(w http.ResponseWriter, r *http.Request),
	req *http.Request,
	expectedBody []byte,
) {
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	convey.So(err, convey.ShouldBeNil)
	convey.So(resp.StatusCode, convey.ShouldEqual, status)
	convey.So(resp.Header, convey.ShouldResemble, header)

	if len(body) == 0 && expectedBody == nil {
		return
	}

	b1, b2 := make(map[string]interface{}), make(map[string]interface{})
	err = json.Unmarshal(body, &b1)
	convey.So(err, convey.ShouldBeNil)

	err = json.Unmarshal(expectedBody, &b2)
	convey.So(err, convey.ShouldBeNil)

	convey.So(b1, convey.ShouldResemble, b2)
}

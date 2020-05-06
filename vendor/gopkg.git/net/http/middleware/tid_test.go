package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/satori/go.uuid"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTID(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	appname := "app-test-app"
	tidmiddleware := NewTID(appname)(next)

	Convey("Given a request without x-tid header", t, func() {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/search", nil)
		tidmiddleware.ServeHTTP(w, r)
		So(r.Header.Get("x-tid"), ShouldStartWith, appname)
	})

	Convey("Given a request with x-tid header", t, func() {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/search", nil)
		xtid := "otherapp-0130-ab"
		r.Header.Add("x-tid", xtid)
		tidmiddleware.ServeHTTP(w, r)
		So(r.Header.Get("x-tid"), ShouldResemble, xtid)
	})
}

func BenchmarkTIDPresent(b *testing.B) {
	reqs := make([]*http.Request, b.N)
	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest(http.MethodGet, "/search", nil)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
		req.Header.Add("x-tid", uuid.NewV4().String())
		reqs[i] = req
	}

	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Use(NewTID("infra"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, reqs[i])
	}
}

func BenchmarkTIDNotPresent(b *testing.B) {
	reqs := make([]*http.Request, b.N)
	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest(http.MethodGet, "/search", nil)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
		reqs[i] = req
	}

	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Use(NewTID("infra"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, reqs[i])
	}
}

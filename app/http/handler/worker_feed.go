package handler

import (
	"context"
	"net/http"
	url2 "net/url"

	"github.com/go-chi/chi"
)

func (wh *WorkerHandler) ScrollFeeds(w http.ResponseWriter, r *http.Request) {
	err := wh.Services.RSS.FindAllFeeds(context.Background())
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
	return
}

func (wh *WorkerHandler) FindFeed(w http.ResponseWriter, r *http.Request) {
	source := chi.URLParam(r, "source")

	url, err := url2.QueryUnescape(source)
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	err = wh.Services.RSS.FindFeedBySource(context.Background(), url)
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
	return
}

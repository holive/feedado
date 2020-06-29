package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
)

func (wh *WorkerHandler) ReindexFeeds(w http.ResponseWriter, r *http.Request) {
	err := wh.Services.RSS.FindAllFeeds(context.Background())
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
	return
}

func (wh *WorkerHandler) ReindexFeedsByCategory(w http.ResponseWriter, r *http.Request) {
	category := chi.URLParam(r, "category")

	err := wh.Services.RSS.FindFeedByCategory(context.Background(), category)
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
	return
}

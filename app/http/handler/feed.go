package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/holive/feedado/app/feed"
)

func (h *Handler) CreateFeed(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	var s feed.Feed
	if err := json.Unmarshal(payload, &s); err != nil {
		respondWithJSONError(w, http.StatusBadRequest, err)
		return
	}

	newFeed, err := h.Services.Feed.Create(r.Context(), &s)
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, *newFeed)
	return
}

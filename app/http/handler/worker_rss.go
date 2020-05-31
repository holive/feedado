package handler

import (
	"net/http"
)

func (h *WorkerHandler) RSS(w http.ResponseWriter, r *http.Request) {

	newFeed, err := h.Services.Feed.Create(r.Context(), &s)
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, *newFeed)
	return
}

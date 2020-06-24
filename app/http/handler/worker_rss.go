package handler

import (
	"net/http"
)

// RSS - triggers the scroll and fetch of all schemas
func (wh *WorkerHandler) RSS(w http.ResponseWriter, r *http.Request) {
	err := wh.Services.RSS.FindAll(r.Context())
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
	return
}

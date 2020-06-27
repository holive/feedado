package handler

import (
	"net/http"
	url2 "net/url"

	"github.com/go-chi/chi"
)

func (h *Handler) GetAllRSS(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	results, err := h.Services.RSS.FindAll(r.Context(), limit, offset)
	if err != nil {
		respondWithJSONError(w, http.StatusNotFound, err)
		return
	}

	respondWithJSON(w, http.StatusOK, results)
}

func (h *Handler) DeleteRSS(w http.ResponseWriter, r *http.Request) {
	url := chi.URLParam(r, "url")

	url, err := url2.QueryUnescape(url)
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.Services.RSS.Delete(r.Context(), url); err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{})
}

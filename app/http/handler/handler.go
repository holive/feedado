package handler

import (
	"net/http"

	"github.com/holive/feedado/app/feedado"
)

type Handler struct {
	Services *feedado.Services
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

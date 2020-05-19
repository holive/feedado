package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/holive/feedado/app/feedado"
)

type Handler struct {
	Services *feedado.Services
}

type Message struct {
	Name string
	Body string
	Time int64
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(map[string]string{"status": "ok"})
	if err != nil {
		fmt.Println("could not Marshal health json response")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

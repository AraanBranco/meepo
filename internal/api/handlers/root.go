package handlers

import (
	"encoding/json"
	"net/http"
)

type DefaultResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func Default(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(DefaultResponse{
		Message: "Welcome to Meepo",
		Status:  http.StatusOK,
	})
}

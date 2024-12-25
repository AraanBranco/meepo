package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/AraanBranco/meepo/internal/config"
	"github.com/AraanBranco/meepo/internal/core/interfaces"
	"github.com/AraanBranco/meepo/internal/service"
)

func NewLobby(w http.ResponseWriter, r *http.Request, configs config.Config) {
	var req interfaces.PostLobbyRequest
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &req)

	hasError, errors := req.Validate(configs)
	if hasError {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"errors": errors})
		return
	}

	lobbyManager, err := service.NewLobbyManager(configs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	lobbyStatus := lobbyManager.CreateLobby(req)

	res := interfaces.PostLobbyResponse{
		ReferenceID: req.ReferenceID,
		LobbyStatus: lobbyStatus,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func StatusLobby(w http.ResponseWriter, r *http.Request, configs config.Config, referenceID string) {
	lobbyManager, err := service.NewLobbyManager(configs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	lobbyStatus := lobbyManager.StatusLobby()
	res := interfaces.GetLobbyResponse{
		ReferenceID: referenceID,
		LobbyStatus: lobbyStatus,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

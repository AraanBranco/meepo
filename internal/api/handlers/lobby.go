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

	hasErrors, errors := req.Validate(configs)
	if hasErrors {
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

	// Send message to MQTT
	lobbyStatus := lobbyManager.CreateLobby(req)

	// Create status in Redis
	err = lobbyManager.EntityInRedis(req.ReferenceID, lobbyStatus)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

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

	lobbyStatus := lobbyManager.StatusLobby(referenceID)
	res := interfaces.GetLobbyResponse{
		ReferenceID: referenceID,
		LobbyStatus: lobbyStatus,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

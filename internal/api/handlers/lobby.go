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

	lobbyManager := service.NewLobbyManager(configs)

	lobbyStatus := lobbyManager.CreateLobby(req)

	res := interfaces.PostLobbyResponse{
		ReferenceID: req.ReferenceID,
		LobbyStatus: lobbyStatus,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func StatusLobby(w http.ResponseWriter, r *http.Request, configs config.Config, referenceID string) {
	lobbyManager := service.NewLobbyManager(configs)

	lobbyStatus, lobbyData := lobbyManager.StatusLobby(referenceID)
	if lobbyStatus == "error" {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": "lobby not found"})
		return
	}

	var response interfaces.GetLobbyResponse
	_ = json.Unmarshal([]byte(lobbyData), &response)

	response.LobbyStatus = lobbyStatus

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

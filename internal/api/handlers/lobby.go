package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/AraanBranco/meepow/internal/config"
	"github.com/AraanBranco/meepow/internal/core/interfaces"
	"github.com/AraanBranco/meepow/internal/core/services/lobby"
)

func NewLobby(w http.ResponseWriter, r *http.Request, configs config.Config, lobbyManager *lobby.LobbyManager) {
	var req interfaces.PostLobbyRequest
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &req)

	hasErrors, errors := req.Validate(configs)
	if hasErrors {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"errors": errors})
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

func StatusLobby(w http.ResponseWriter, r *http.Request, configs config.Config, referenceID string, lobbyManager *lobby.LobbyManager) {
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

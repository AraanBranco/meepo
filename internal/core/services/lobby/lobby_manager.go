package lobby

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AraanBranco/meepo/internal/config"
	"github.com/AraanBranco/meepo/internal/core/interfaces"
	"github.com/go-stomp/stomp/v3"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type LobbyManager struct {
	Stomp  *stomp.Conn
	Redis  *redis.Client
	Config config.Config
	Logger *zap.Logger
}

func New(conf config.Config, st *stomp.Conn, rs *redis.Client) *LobbyManager {
	return &LobbyManager{
		Stomp:  st,
		Redis:  rs,
		Config: conf,
		Logger: zap.L(),
	}
}

func (l *LobbyManager) PublishNewLobby(lobbyData []byte) error {
	return l.Stomp.Send("/queue/new-lobby", "text/plain", lobbyData, nil)
}

func (l *LobbyManager) EntityInRedis(referenceID string, status string) error {
	return l.Redis.Set(context.Background(), fmt.Sprintf("lobby:%s:status", referenceID), status, 0).Err()
}

func (l *LobbyManager) GetEntityInRedis(referenceID string) (string, error) {
	result, err := l.Redis.Get(context.Background(), fmt.Sprintf("lobby:%s:status", referenceID)).Result()
	if err != nil {
		if err == redis.Nil {
			return "not_found", nil
		}
		l.Logger.Error("Error getting lobby status from Redis", zap.Error(err))
		return "", err
	}

	return result, nil
}

func (l *LobbyManager) CreateLobby(params interfaces.PostLobbyRequest) string {
	l.Logger.Info("Creating lobby", zap.String("reference_id", params.ReferenceID), zap.String("lobby_name", params.LobbyName))
	data, err := json.Marshal(params)
	if err != nil {
		l.Logger.Error("Error marshalling lobby data", zap.Error(err))
		return "error"
	}

	err = l.PublishNewLobby(data)
	if err != nil {
		l.Logger.Error("Error publishing lobby data", zap.Error(err))
		return "error"
	}

	return "created"
}

func (l *LobbyManager) StatusLobby(referenceID string) string {
	lobbyStatus, err := l.GetEntityInRedis(referenceID)
	if err != nil {
		l.Logger.Error("Error getting lobby status from Redis", zap.Error(err))
		return "error"
	}

	return lobbyStatus
}

package lobby

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AraanBranco/meepo/internal/config"
	"github.com/AraanBranco/meepo/internal/core/interfaces"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type LobbyManager struct {
	Redis  *redis.Client
	Config config.Config
	Logger *zap.Logger
}

func New(conf config.Config, rs *redis.Client) *LobbyManager {
	return &LobbyManager{
		Redis:  rs,
		Config: conf,
		Logger: zap.L().With(zap.String("service", "lobby")),
	}
}

func (l *LobbyManager) CreateEntityInRedis(referenceID string, lobbyData string) error {
	_ = l.Redis.Set(context.Background(), fmt.Sprintf("lobby:%s", referenceID), lobbyData, 0).Err()
	return l.Redis.Set(context.Background(), fmt.Sprintf("lobby:%s:status", referenceID), interfaces.LOBBY_CREATING, 0).Err()
}

func (l *LobbyManager) GetEntityInRedis(referenceID string) (string, string, error) {
	status, err := l.Redis.Get(context.Background(), fmt.Sprintf("lobby:%s:status", referenceID)).Result()
	if err != nil {
		if err == redis.Nil {
			return "not_found", "", nil
		}
		l.Logger.Error("Error getting lobby status from Redis", zap.Error(err))
		return "", "", err
	}

	lobbyData, _ := l.Redis.Get(context.Background(), fmt.Sprintf("lobby:%s", referenceID)).Result()

	return status, lobbyData, nil
}

func (l *LobbyManager) GetLobbyData(referenceID string) (string, error) {
	lobbyData, err := l.Redis.Get(context.Background(), fmt.Sprintf("lobby:%s", referenceID)).Result()
	if err != nil {
		if err == redis.Nil {
			return "not_found", nil
		}
		l.Logger.Error("Error getting lobby status from Redis", zap.Error(err))
		return "", err
	}

	return lobbyData, nil
}

func (l *LobbyManager) CreateLobby(params interfaces.PostLobbyRequest) string {
	l.Logger.Info("Creating lobby", zap.String("reference_id", params.ReferenceID), zap.String("lobby_name", params.LobbyName))
	data, err := json.Marshal(params)
	if err != nil {
		l.Logger.Error("Error marshalling lobby data", zap.Error(err))
		return "error"
	}

	// Save lobby data in Redis
	err = l.CreateEntityInRedis(params.ReferenceID, string(data))
	if err != nil {
		l.Logger.Error("Error saving lobby data in Redis", zap.Error(err))
		return "error"
	}

	// Start ECS task asynchronously

	return "created"
}

func (l *LobbyManager) StatusLobby(referenceID string) (string, string) {
	lobbyStatus, lobbyData, err := l.GetEntityInRedis(referenceID)
	if err != nil {
		l.Logger.Error("Error getting lobby status from Redis", zap.Error(err))
		return "error", ""
	}

	return lobbyStatus, lobbyData
}

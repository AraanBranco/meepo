package lobby

import (
	"github.com/AraanBranco/meepo/internal/config"
	"github.com/AraanBranco/meepo/internal/core/interfaces"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type LobbyManager struct {
	Config config.Config
	Logger *zap.Logger
}

func New(c config.Config, redis *redis.Client) *LobbyManager {
	return &LobbyManager{
		Config: c,
		Logger: zap.L(),
	}
}

func (l *LobbyManager) CreateLobby(params interfaces.PostLobbyRequest) string {
	l.Logger.Info("Creating lobby")
	l.Logger.Info("Params", zap.Any("params", params))

	return "created"
}

func (l *LobbyManager) StatusLobby() string {
	l.Logger.Info("Status lobby")

	return "finish"
}

package service

import (
	"context"

	"github.com/AraanBranco/meepow/internal/config"
	"github.com/AraanBranco/meepow/internal/core/services/bot"
	"github.com/AraanBranco/meepow/internal/core/services/lobby"
	"github.com/paralin/go-steam"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	cfgAws "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

// Configs Paths for adapters
const (
	redisURIPath = "adapters.redis.uri"
)

func NewLobbyManager(c config.Config) *lobby.LobbyManager {
	redisClient := createRedisClient(c)

	cfg, err := cfgAws.LoadDefaultConfig(context.TODO(), cfgAws.WithRegion(c.GetString("providers.aws.region")))
	if err != nil {
		zap.L().Error("Erro ao carregar a configuração da AWS", zap.Error(err))
	}

	client := ecs.NewFromConfig(cfg)

	return lobby.New(c, redisClient, client)
}

func NewBotManager(c config.Config) *bot.BotManager {
	redisClient := createRedisClient(c)
	steamClient := steam.NewClient()
	steam.InitializeSteamDirectory()

	return bot.New(c, redisClient, steamClient)
}

func createRedisClient(c config.Config) *redis.Client {
	logger := zap.L().With(zap.String("adapter", "redis"))

	logger.Info("Creating Redis client")

	conf, err := redis.ParseURL(c.GetString(redisURIPath))
	if err != nil {
		logger.Error("Error parsing Redis URI", zap.Error(err))
	}

	rdb := redis.NewClient(conf)

	redisStatus, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logger.Error("Error connecting to Redis", zap.Error(err))
	}

	logger.Info("Redis status", zap.String("status", redisStatus))

	return rdb
}

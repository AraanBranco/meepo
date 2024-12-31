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
	redisPoolSizePath = "adapters.redis.poolSize"
	redisURIPath      = "adapters.redis.uri"
	redisUserPath     = "adapters.redis.user"
	redisPassPath     = "adapters.redis.password"
	redisDBPath       = "adapters.redis.db"

	stompURLPath = "adapters.stomp.url"
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
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.GetString(redisURIPath),
		Username: c.GetString(redisUserPath),
		Password: c.GetString(redisPassPath),
		DB:       c.GetInt(redisDBPath),
		PoolSize: c.GetInt(redisPoolSizePath),
	})

	return rdb
}

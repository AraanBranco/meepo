package service

import (
	"github.com/AraanBranco/meepo/internal/config"
	"github.com/AraanBranco/meepo/internal/core/services/bot"
	"github.com/AraanBranco/meepo/internal/core/services/lobby"
	"github.com/paralin/go-steam"
	"github.com/redis/go-redis/v9"
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

	return lobby.New(c, redisClient)
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

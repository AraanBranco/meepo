package service

import (
	"fmt"

	"github.com/AraanBranco/meepo/internal/config"
	"github.com/AraanBranco/meepo/internal/core/services/lobby"
	"github.com/redis/go-redis/v9"
)

// Configs Paths for adapters
const (
	redisPoolSizePath = "adapters.redis.poolSize"
	redisURLPath      = "adapters.redis.url"
)

func NewLobbyManager(c config.Config) (*lobby.LobbyManager, error) {
	redisURL := c.GetString(redisURLPath)
	redisClient, err := createRedisClient(c, redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create redis client: %w", err)
	}

	return lobby.New(c, redisClient), nil
}

func createRedisClient(c config.Config, url string) (*redis.Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("invalid redis URL: %w", err)
	}
	opts.PoolSize = c.GetInt(redisPoolSizePath)
	if opts.TLSConfig != nil {
		opts.TLSConfig.InsecureSkipVerify = true
	}

	client := redis.NewClient(opts)

	return client, nil
}

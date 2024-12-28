package service

import (
	"fmt"

	"github.com/AraanBranco/meepo/internal/config"
	"github.com/AraanBranco/meepo/internal/core/services/lobby"
	"github.com/go-stomp/stomp/v3"
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

func NewLobbyManager(c config.Config) (*lobby.LobbyManager, error) {
	connStomp, err := createStompConn(c.GetString(stompURLPath))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to stomp: %w", err)
	}

	redisClient := createRedisClient(c)

	return lobby.New(c, connStomp, redisClient), nil
}

func createStompConn(url string) (*stomp.Conn, error) {
	conn, err := stomp.Dial("tcp", url)
	if err != nil {
		return nil, err
	}

	return conn, nil
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

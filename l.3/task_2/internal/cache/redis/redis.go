package redis

import (
	"context"
	"os"
	"time"

	"github.com/wb-go/wbf/redis"
	"task-2/internal/config"
)

type Redis struct {
	client redis.Client
}

func New() *Redis {
	password := os.Getenv("REDIS_PASSWORD")

	client := redis.New(
		config.Cfg.Redis.Host+":"+config.Cfg.Redis.Port,
		password,
		0,
	)

	return &Redis{
		client: *client,
	}
}

func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key)
}

func (r *Redis) Set(key string, value interface{}) error {
	return r.client.SetEX(context.Background(), key, value, time.Hour*24).Err()
}

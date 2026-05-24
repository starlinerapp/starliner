package redis

import (
	"github.com/redis/go-redis/v9"
	"starliner.app/internal/core/conf"
)

func Connect(cfg conf.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.GetRedisPassword(),
		DB:       0,
		Protocol: 2,
	})
}

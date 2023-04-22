package model

import (
	"appeal-gateway/rpc/appeal/internal/config"

	"github.com/go-redis/redis/v8"
)

func InitRedis(c *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

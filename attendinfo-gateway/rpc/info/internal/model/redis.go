package model

import (
	"info/internal/config"

	"github.com/go-redis/redis/v8"
)

var RDB5 *redis.Client
var RDB *redis.Client
var RDB6 *redis.Client
var RDB7 *redis.Client

func InitRedis(c *config.Config) *redis.Client {
	RDB = redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: "147258", // no password set
		DB:       1,        // use default DB
	})
	return RDB
}
func InitRedis4(c *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: "147258", // no password set
		DB:       4,        // use default DB
	})
}
func InitRedis3(c *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: "147258", // no password set
		DB:       3,        // use default DB
	})
}
func InitRedis5(c *config.Config) *redis.Client {
	RDB5 = redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: "147258", // no password set
		DB:       5,        // use default DB
	})
	return RDB5
}
func InitRedis6(c *config.Config) *redis.Client {
	RDB6 = redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: "147258", // no password set
		DB:       6,        // use default DB
	})
	return RDB6
}
func InitRedis7(c *config.Config) *redis.Client {
	RDB7 = redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: "147258", // no password set
		DB:       7,        // use default DB
	})
	return RDB7
}

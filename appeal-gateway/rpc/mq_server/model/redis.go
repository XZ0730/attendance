package model

import (
	"mq_server/internal/config"

	"github.com/go-redis/redis/v8"
)

var RDB1 *redis.Client
var RDB3 *redis.Client
var RDB5 *redis.Client
var RDB6 *redis.Client
var RDB7 *redis.Client

func InitRedis(c *config.Config) *redis.Client {
	RDB1 = redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: "147258", // no password set
		DB:       1,        // use default DB
	})
	return RDB1
}
func InitRedis3(c *config.Config) *redis.Client {
	RDB3 = redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: "147258", // no password set
		DB:       3,        // use default DB
	})
	return RDB3
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

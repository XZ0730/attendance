package svc

import (
	"info/internal/config"
	"info/internal/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	RDB    *redis.Client
	RDB6   *redis.Client
	RDB7   *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		DB:     model.Init(&c),
		RDB:    model.InitRedis(&c),
		RDB6:   model.InitRedis6(&c),
		RDB7:   model.InitRedis7(&c),
	}
}

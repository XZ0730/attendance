package svc

import (
	"appeal-gateway/rpc/appeal/internal/config"
	"appeal-gateway/rpc/appeal/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	MysqlDB *gorm.DB
	RDB     *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		MysqlDB: model.Init(&c),
		RDB:     model.InitRedis(&c),
	}
}

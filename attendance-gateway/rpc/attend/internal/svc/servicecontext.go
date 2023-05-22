package svc

import (
	"attend/internal/config"
	"attend/model"
	"mq_server/mqclient"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	RDB    *redis.Client
	RDB4   *redis.Client
	RDB3   *redis.Client
	RDB5   *redis.Client
	RDB6   *redis.Client
	DB     *gorm.DB
	MQ     mqclient.Mq
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		RDB:    model.InitRedis(&c),
		RDB3:   model.InitRedis3(&c),
		RDB4:   model.InitRedis4(&c),
		RDB5:   model.InitRedis5(&c),
		RDB6:   model.InitRedis6(&c),
		DB:     model.Init(&c),
		MQ:     mqclient.NewMq(zrpc.MustNewClient(c.MqCli)),
	}
}

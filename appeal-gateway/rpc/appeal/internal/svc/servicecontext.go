package svc

import (
	"appeal/internal/config"
	"appeal/model"
	"mq_server/mqclient"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	MysqlDB *gorm.DB
	RDB     *redis.Client
	MQ      mqclient.Mq
	RDB2    *redis.Client
	RDB3    *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		MysqlDB: model.Init(&c),
		RDB:     model.InitRedis(&c),
		MQ:      mqclient.NewMq(zrpc.MustNewClient(c.MqCli)),
		RDB2:    model.InitRedis2(&c),
		RDB3:    model.InitRedis3(&c),
	}
}

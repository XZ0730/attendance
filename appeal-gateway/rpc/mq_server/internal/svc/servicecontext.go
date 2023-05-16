package svc

import (
	"mq_server/internal/config"
	model "mq_server/model"
	"mq_server/rabbitmq"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	Mq_server *rabbitmq.RabbitMQ
	MysqlDB   *gorm.DB
	RDB1      *redis.Client
	RDB5      *redis.Client
	RDB6      *redis.Client
	RDB7      *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		Mq_server: rabbitmq.InitRabbitMQ(&c),
		MysqlDB:   model.Init(&c),
		RDB1:      model.InitRedis(&c),
		RDB5:      model.InitRedis5(&c),
		RDB6:      model.InitRedis6(&c),
		RDB7:      model.InitRedis7(&c),
	}
}

package svc

import (
	"mq_server/internal/config"
	model "mq_server/model"
	"mq_server/rabbitmq"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	Mq_server *rabbitmq.RabbitMQ
	MysqlDB   *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		Mq_server: rabbitmq.InitRabbitMQ(&c),
		MysqlDB:   model.Init(&c),
	}
}

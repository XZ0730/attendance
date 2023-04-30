package rabbitmq

import (
	"mq_server/internal/config"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn  *amqp.Connection
	mqurl string
}

var Rmq *RabbitMQ

func InitRabbitMQ(c *config.Config) *RabbitMQ {

	Rmq = &RabbitMQ{
		mqurl: c.MQURL,
	}
	dial, err := amqp.Dial(Rmq.mqurl)
	if err != nil {
		return nil
	}
	Rmq.conn = dial
	return Rmq
}
func (r *RabbitMQ) destroy() {
	r.conn.Close()
}

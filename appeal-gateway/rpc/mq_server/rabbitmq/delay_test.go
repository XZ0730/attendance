package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/streadway/amqp"
)

func TestRabbit(t *testing.T) {
	fmt.Println("--------------------")
	InitRabbitMQ1()
	NewDelayMQ("delayto").Publish("hellohello")

}
func TestCosume(t *testing.T) {
	fmt.Println("--------------------")
	InitRabbitMQ1()
	fmt.Println("-------------2-------")
	NewDelayMQ("delayto").Consumer()
}
func InitRabbitMQ1() *RabbitMQ {

	Rmq = &RabbitMQ{
		mqurl: "amqp://admin:admin123@43.136.122.18:5672/",
	}
	dial, err := amqp.Dial(Rmq.mqurl)
	if err != nil {
		fmt.Println("err:", err)
		return nil
	}
	Rmq.conn = dial
	return Rmq
}

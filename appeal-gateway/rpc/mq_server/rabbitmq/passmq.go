package rabbitmq

import (
	"encoding/json"
	"fmt"
	"mq_server/model"
	"mq_server/mq"

	"github.com/streadway/amqp"
)

type PassMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	exchange  string
	queueName string
}

func NewPassMQ(queueName string) *PassMQ {
	passMQ := &PassMQ{
		RabbitMQ:  *Rmq,
		queueName: queueName, //friendQue groupQue
	}

	ch, err := passMQ.conn.Channel()
	if err != nil {

		return nil
	}
	passMQ.channel = ch
	return passMQ
}
func InitPassMQ() {
	passmq := NewPassMQ("passmq")
	go passmq.Consumer()
}
func (c *PassMQ) Publish(message string) error {
	fmt.Println(":", c.queueName)
	_, err := c.channel.QueueDeclare(
		c.queueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil,
	)
	if err != nil {
		return err
	}
	//json.marshal 可序列化结构体为二进制byte类型
	//然后就可以通过消息队列进行传参，
	//在消费者方面只需要通过unmarshal进行反序列化就可以得到结构体

	err1 := c.channel.Publish(
		c.exchange,
		c.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err1 != nil {
		return err
	}
	return nil
}
func (r *PassMQ) Consumer() {
	defer r.destroy()
	_, err := r.channel.QueueDeclare(r.queueName, false, false, false, false, nil)

	if err != nil {
		return
	}

	//2、接收消息
	msg, err := r.channel.Consume(
		r.queueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//消息队列是否阻塞
		false,
		nil,
	)
	if err != nil {
		return
	}
	go r.PassTo(msg)
	//log.Printf("[*] Waiting for messages,To exit press CTRL+C")
	forever := make(chan bool)
	<-forever
}
func (r *PassMQ) PassTo(msg <-chan amqp.Delivery) {
	for req := range msg {
		fmt.Println("req:", string(req.Body))
		rq := &mq.Request{}

		err := json.Unmarshal(req.Body, rq)
		if err != nil {
			fmt.Println("err:", err)
			continue
		}
		fmt.Println("rq:", rq)
		err1 := model.DB.
			Where("id=? AND counsellor_id=?", rq.Cid, rq.CounsellorId).
			Delete(&model.ComplainTable{}).Error
		if err1 != nil {
			fmt.Println("err1:", err1)
			continue
		}
	}
}

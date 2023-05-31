package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mq_server/model"
	"mq_server/pkg"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

type DelayMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	exchange  string
	queueName string
}

func InitDelayMQ() {
	delay := NewDelayMQ("delay1")
	go delay.Consumer()
}
func NewDelayMQ(queueName string) *DelayMQ {
	delayMQ := &DelayMQ{
		RabbitMQ:  *Rmq,
		queueName: queueName, //friendQue groupQue
	}

	ch, err := delayMQ.conn.Channel()
	if err != nil {

		return nil
	}
	delayMQ.channel = ch
	return delayMQ
}

func (d *DelayMQ) Publish(message string) {

	// 将消息发送到延时队列上
	_, err := d.channel.QueueDeclare(
		d.queueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		amqp.Table{
			"x-dead-letter-exchange": "dealWithdelay",
		},
	)
	_ = d.channel.Publish(
		"",      // exchange 这里为空则不选择 exchange
		"delay", // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			Expiration:  "180000", // 设置两分钟的过期时间
		})

	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", message)

}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (d *DelayMQ) Consumer() {
	err := d.channel.ExchangeDeclare(
		"dealWithdelay",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")
	_, err = d.channel.QueueDeclare( //正常队列
		d.queueName,
		false,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange": "dealWithdelay",
		},
	)
	failOnError(err, "Failed to declare an queue")
	//正常队列超时---放入延时队列
	//延时队列

	err = d.channel.QueueBind(
		d.queueName,
		"",
		"dealWithdelay",
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := d.channel.Consume(
		d.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")
	go d.unRegisterAtt(msgs)
	forever := make(chan bool)
	<-forever
}

func (*DelayMQ) unRegisterAtt(msg <-chan amqp.Delivery) {
	for req := range msg {
		fmt.Println("----------------------------------------------------------------------------")
		fmt.Println(time.Now())
		fmt.Println("req:", string(req.Body))
		mp := &model.MarshalPull{}
		err := json.Unmarshal(req.Body, mp)
		if err != nil {
			fmt.Println("err:", err)
			req.Ack(false)
			continue
		}
		cr := &model.Course{}
		week := pkg.GetWeek(model.RDB5)
		if week == -1 {
			fmt.Println("err:", err)
			req.Ack(false)
			continue
		}
		model.DB.Table("course").Where("course_id=? AND university=?", mp.CourseID, mp.University).First(&cr)
		//计算周数--当前时间减去开学时间

		// week, err2 := model.RDB1.HGet(model.RDB1.Context(), cr.University, strconv.Itoa(int(cr.Id))).Result()
		// if err2 != nil {
		// 	fmt.Println("week:", week)
		// 	fmt.Println("err2:", err2)
		// 	req.Ack(false)
		// 	continue
		// }
		// weekto, err := strconv.Atoi(week)
		// if err != nil {
		// 	fmt.Println("err:", err)
		// 	req.Ack(false)
		// 	continue
		// }
		stu := make([]*model.LeaveTable, 0)
		model.DB.Table("leave_table").
			Where("course_id=? AND is_audit=3 AND tag_as=1 AND leave_course_from<=? AND leave_course_to >=? AND school_name=?", mp.CourseID, week, week, mp.University).
			Find(&stu)
		for _, s := range stu { //将请假的同学考勤状态设置为1
			model.RDB5.ZAdd(model.RDB5.Context(), strconv.Itoa(int(cr.Id)), &redis.Z{
				Score:  1,
				Member: s.StudentID,
			})
		}
		stuid, _ := model.RDB5.ZRangeByScore(model.RDB5.Context(), strconv.Itoa(int(cr.Id)), &redis.ZRangeBy{
			Max: "0",
			Min: "0",
		}).Result()
		fmt.Println("stuid:", stuid)
		fmt.Println("测试1")
		results := make([]*model.Result, 0) //存放最终考勤结果
		model.DB.Table("course_group").
			Select("character_msg.code as code,character_msg.name as student_name,course_group.course_id as course_id").Joins("left join character_msg on course_group.student_id=character_msg.code where course_id=? AND course_group.university=?", mp.CourseID, cr.University).
			Scan(&results)
		fmt.Println("测试2")
		stus := make(map[string]*model.Result, len(results))
		// weekto = weekto + 1
		// weeknow := strconv.Itoa(weekto)
		for _, v := range results {
			// v.Week = uint(weekto)
			v.Week = uint(week)
			stus[v.Code] = v
		}
		fmt.Println("测试3")
		name := ""
		id := ""
		cnt := 0
		for _, v := range stuid {
			stus[v].MissAttend = 1
			if cnt < len(stuid)-1 {
				temp := stus[v].StudentName + ","
				name += temp
				temp = v + ","
				id += temp
			} else {
				name += stus[v].StudentName
				id += v
			}
			cnt++
		}
		// fmt.Println("测试4")

		// fmt.Println("测试5")
		// fmt.Println("mpL:", mp)
		// fmt.Println("week:", week)
		//计算是否存在记录，如果缺勤人数为0，
		var cntf int64
		att := &model.AttendTable{}
		model.DB.Table("attend_table").
			Where("week=? AND course_id=? AND university=?", week, mp.CourseID, cr.University).First(&att).Count(&cntf)
		if cntf > 0 {
			if att.Unpresent != 0 {
				if att.Unpresent > 1 {
					sid := strings.Split(att.UnpresenterID, ",")
					for _, v := range sid {
						//两次考勤进行合并
						if stus[v].MissAttend == 2 {
							stus[v].MissAttend = 1
							name = name + "," + stus[v].StudentName
							id = id + "," + stus[v].Code
						}
					}
				}
			}
			err4 := model.DB.Unscoped().Delete(&att).Error
			if err4 != nil {
				fmt.Println("err4:", err4)
			}
		}
		at := &model.AttendTable{
			CourseID:      mp.CourseID,
			CourseName:    cr.Name,
			University:    cr.University,
			Week:          uint(week),
			Teacher:       cr.TeacherName,
			Unpresent:     uint(len(stuid)),
			Unpresenter:   name,
			UnpresenterID: id,
		}
		//考勤结果存入mysql
		err2 := model.DB.Table("attend_table").Create(&at).Error
		if err2 != nil {
			fmt.Println("err:", err2)
			req.Ack(false)
			continue
		}
		result, _ := json.Marshal(stus)
		// fmt.Println("res:", string(result))
		//将考勤结果存入redis
		err = model.RDB6.HSet(model.RDB6.Context(), strconv.Itoa(int(cr.Id)), strconv.Itoa(int(week)), string(result)).Err()
		if err != nil {
			fmt.Println("err343:", err)
			req.Ack(false)
			continue
		}
		fmt.Println("测试6")
		//将本周考勤状态存入redis
		err3 := model.RDB7.HSet(model.RDB7.Context(), strconv.Itoa(int(cr.Id)), strconv.Itoa(int(week)), "2").Err() //点名结束 本周课程点名状态设置为1
		if err3 != nil {
			fmt.Println("err3", err3)
			req.Ack(false)
			continue
		}
		fmt.Println("RDB3:", model.RDB3)
		_ = model.RDB3.ZAdd(context.Background(), mp.University, &redis.Z{
			Score:  2,
			Member: strconv.Itoa(int(cr.Id)),
		}).Err()
		err5 := model.RDB1.HSet(context.Background(), mp.University, strconv.Itoa(int(cr.Id)), strconv.Itoa(int(week))).Err()
		if err5 != nil {
			fmt.Println("err5", err5)
			req.Ack(false)
			continue
		}
		// s, _ := model.RDB1.HGet(context.Background(), mp.University, strconv.Itoa(int(cr.Id))).Result()
		// fmt.Println("s", s)
		fmt.Println("测试7")
		req.Ack(false)
	}
}

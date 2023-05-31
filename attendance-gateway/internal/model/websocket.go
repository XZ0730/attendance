package model

import (
	"attendance-gateway/internal/logic"
	"attendance-gateway/internal/svc"
	"attendance-gateway/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID        string
	Socket    *websocket.Conn
	HeartBeat int64
	Send      chan []byte

	Done chan struct{}
}
type ClientManager struct {
	Clients sync.Map
	// BroadCast  chan *BroadCast
	ReplyMsg   chan *Client
	Register   chan *Client
	Unregister chan *Client
}

// type BroadCast struct {
// 	Client  *Client
// 	Message []byte
// 	Type    string
// }

var Manager = ClientManager{ //链接管理员
	// BroadCast:  make(chan *BroadCast),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

type Message struct {
	University string `json:"university"`
	CourseID   string `json:"course_id"`
	StudentID  string `json:"student_id"`
	CourseMain int64  `json:"course_main"`
	Type       uint   `json:"type"`
}

func (manager *ClientManager) Start() {
	for {
		select {
		case client := <-Manager.Register: // 建立连接
			Manager.Clients.Store(client.ID, client)
			msg := []byte("链接成功")
			_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
		case client := <-Manager.Unregister:
			msg := []byte("链接断开")
			fmt.Println("链接断开")
			_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
			_ = client.Socket.Close()
			close(client.Send)
			Manager.Clients.Delete(client.ID)
			// case broadcast := <-Manager.BroadCast:
			// 	//
			// 	broadcast.Client.Send <- broadcast.Message
		}
	}
}

// 读取点名信息，然后调用微服务
func (c *Client) Readmsg(req *types.NormalAttReq, svcCtx *svc.ServiceContext) {
	defer func() {
		fmt.Println("俺结束哩")
	}()
	for {
		select {
		case <-c.Done:
			fmt.Println("read出来了")
			return
		default:
			vo := &Message{}
			fmt.Println("-----------测试1")
			err := c.Socket.ReadJSON(&vo)
			if err != nil {
				fmt.Println("err:", err)
				c.Done <- struct{}{}
				return
			}
			fmt.Println("-----------测试2")
			if vo.Type == 3 {
				c.UpdateHeartBeat()
				continue
			}
			fmt.Println("vo:", vo)
			req.Type = 1
			req.CourseID = vo.CourseID
			req.CourseMain = vo.CourseMain
			req.StudentID = vo.StudentID
			fmt.Println("-----------测试3")
			nal := logic.NewNormalAttendLogic(context.Background(), svcCtx)
			ar, _ := nal.NormalAttend(req)
			ar1, _ := json.Marshal(ar)
			c.Socket.WriteMessage(websocket.TextMessage, ar1)
		}

	}

}

// 每隔三秒回写实时考勤信息
func (c *Client) Writemsg(req *types.NormalAttReq, svcCtx *svc.ServiceContext) {
	defer func() {
		fmt.Println("俺结束哩")
		Manager.Unregister <- c
	}()
	//发送初始名单
	req.Type = 0
	nal := logic.NewNormalAttendLogic(context.Background(), svcCtx)
	ar, err := nal.NormalAttend(req)
	if err != nil {
		fmt.Println("err:", err)
	}
	//从数据库中获取选课名单 生成map
	//获取redis中缺勤的同学匹配生成实时考勤名单
	//marshal传入
	b, err2 := json.Marshal(ar)
	if err2 != nil {
		fmt.Println("err2:", err2)
	}
	fmt.Println("ar:", ar)
	c.Socket.WriteMessage(websocket.TextMessage, b)

	flag := false
	// msg := []byte("")
	for {
		select {
		// case msg <- c.Send:
		// 	c.Socket.WriteMessage(websocket.TextMessage, msg)
		// case <-time.After(3 * time.Second):
		// 	//调用接口

		// 	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "消息")
		case <-c.Done:
			fmt.Println("俺结束哩")
			flag = true
			break
		}
		if flag {
			break
		}
	}
}
func (c *Client) LinkWith() {
	defer func() {
		fmt.Println("俺也结束哩")
	}()
	flag1 := false
	for {
		select {
		case <-time.After(5 * time.Second):
			currentTime := time.Now().Unix()
			if c.IsHeartBeatTimeout(currentTime) {
				for i := 0; i < 2; i++ {
					fmt.Println("range2:", i)
					c.Done <- struct{}{}
				}
				fmt.Println("range2")
				flag1 = true
			}
			if flag1 {
				break
			}
		}
		if flag1 {
			break
		}
	}
}
func (c *Client) IsHeartBeatTimeout(cur int64) bool {
	// add := time.Duration.Seconds(5)
	// more := int64(add)
	if c.HeartBeat+25 < cur {
		fmt.Println(c.HeartBeat+10, ":", cur)
		return true
	}
	return false
}
func (c *Client) UpdateHeartBeat() {
	c.HeartBeat = time.Now().Unix()
	return
}

// func ProtectRoutine() {

//		for {
//			fmt.Println("hello")
//			RangeManager()
//		}
//	}
// func RangeManager() {

//		Manager.Clients.Range(func(key, value any) bool {
//			currentTime := time.Now().Unix()
//			if value.(*Client).IsHeartBeatTimeout(currentTime) {
//				for i := 0; i < 2; i++ {
//					fmt.Println("range2:", i)
//					value.(*Client).Done <- struct{}{}
//				}
//				fmt.Println("range2")
//			}
//			<-time.After(10 * time.Second)
//			fmt.Println("i am here")
//			return true
//		})
//	}

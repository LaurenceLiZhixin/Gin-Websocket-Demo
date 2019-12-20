package notification

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

//Hub 中央通知机构
/*
负责时刻监听后端的调用，后端将信息从合适身份channel输入，由此
数据结构转发给所有对应身份已连接的用户。
未连接的用户需要将消息放入缓存。这只是一个实时的demo
*/
type Hub struct {
	TeacherNoteChannel chan []byte
	StudentNoteChannel chan []byte
	SigninClient       chan *Client
	SignoutClient      chan *Client
	IsOnline           map[*Client]bool
}

//NewInstance 获取单例对象
func NewInstance() *Hub {
	return &Hub{
		TeacherNoteChannel: make(chan []byte),
		StudentNoteChannel: make(chan []byte),
		SigninClient:       make(chan *Client),
		SignoutClient:      make(chan *Client),
		IsOnline:           make(map[*Client]bool),
	}
}

//SendTeacherMsg 老师通知发送接口
func (h *Hub) SendTeacherMsg(msg string) {
	h.TeacherNoteChannel <- []byte(msg)
}

//SendStudentMsg 学生通知发送接口
func (h *Hub) SendStudentMsg(msg string) {
	h.StudentNoteChannel <- []byte(msg)
}

//Run 开启中央模块服务
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.SigninClient:
			log.Println("用户登入:" + client.Identity)
			h.IsOnline[client] = true
		case client := <-h.SignoutClient:
			log.Println("用户登出:" + client.Identity)
			delete(h.IsOnline, client)
			client.Conn.Close()
		case teacherMsg := <-h.TeacherNoteChannel:
			for teacher := range h.IsOnline {
				if teacher.Identity == "teacher" {
					teacher.Conn.WriteMessage(websocket.TextMessage, teacherMsg)
				}
			}
		case studentMsg := <-h.StudentNoteChannel:
			for student := range h.IsOnline {
				if student.Identity == "student" {
					log.Println("将要给学生发信息")
					if err := student.Conn.WriteMessage(websocket.TextMessage, studentMsg); err != nil {
						h.SignoutClient <- student
						break
					}
				}
			}
		}
	}
}

//Client 用户连接服务
type Client struct {
	Identity string
	Conn     *websocket.Conn
	HubPtr   *Hub
}

//Start 确保用户离开后关闭连接
func (c *Client) Start() {
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			err := c.Conn.WriteMessage(websocket.BinaryMessage, nil)
			if err != nil {
				log.Println("发现用户离开")
				c.HubPtr.SignoutClient <- c
				return
			}
		}
	}
}

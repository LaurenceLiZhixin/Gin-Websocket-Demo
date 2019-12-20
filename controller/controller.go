package controller

import (
	"gin-websocket-demo/notification"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//MsgSender 不停为老师和学生发送消息
func MsgSender() {
	studentTimeTick := time.NewTicker(1 * time.Second)
	teacherTimeTick := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-studentTimeTick.C:
			hub.SendStudentMsg("给学生发的消息：1s一次")
		case <-teacherTimeTick.C:
			hub.SendTeacherMsg("给老师发的消息：5s一次")
		}
	}
}

//LoginHandler 处理用户的登录请求
func LoginHandler(c *gin.Context) {
	useremail := c.Param("email") //URL内部参数使用
	log.Println(useremail + "登录")
	c.HTML(http.StatusOK, "home.html", gin.H{})
}

//NotificationHandler 构造通知对象Client并在中央注册
func NotificationHandler(c *gin.Context) {
	identity := c.DefaultQuery("identity", "student") // URL尾部参数
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("websocket 创建错误")
		return
	}
	defer ws.Close()
	//构造通知对象
	clientPtr := &notification.Client{
		Conn:     ws,
		Identity: identity,
		HubPtr:   hub,
	}
	//中央注册
	hub.SigninClient <- clientPtr
	//用户连接维护
	clientPtr.Start()
}

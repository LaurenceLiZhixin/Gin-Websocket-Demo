package controller

import (
	"gin-websocket-demo/notification"
)

var hub *notification.Hub

func init() {
	hub = notification.NewInstance()
	//开启中央通知服务
	go hub.Run()
}

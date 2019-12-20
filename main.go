package main

import (
	"gin-websocket-demo/controller"
	"gin-websocket-demo/router"
)

//URL 监听端口
const URL = "localhost:8082"

func main() {
	//用于demo测试
	go controller.MsgSender()
	router.Serve(URL)
}

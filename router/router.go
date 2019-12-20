package router

import (
	"gin-websocket-demo/controller"

	"github.com/gin-gonic/gin"
)

//Serve 提供路由
func Serve(URL string) {
	router := gin.Default()
	router.LoadHTMLGlob("view/*")
	router.GET("/user/:email/login", controller.LoginHandler)
	router.GET("/ws", controller.NotificationHandler)
	router.Run(URL)
}

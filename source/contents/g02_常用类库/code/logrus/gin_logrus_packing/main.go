package main

import (
	"fmt"
	"gin_logrus_packing/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func helloWorld(c *gin.Context) {
	// 测试写入日志
	middleware.Logger.WithFields(logrus.Fields{
		"data": "访问/hello",
	}).Info("测试写入info")

	// c.JSON：返回JSON格式的数据
	c.JSON(200, gin.H{
		"message": "Hello world!",
	})
}

func main() {
	r := gin.Default()
	r.Use(middleware.LoggerMiddleware())
	r.GET("/hello", helloWorld)
	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	fmt.Println(`http://127.0.0.1:8080/hello`)
	r.Run(":8080")
}

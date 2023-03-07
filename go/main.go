package main

import (
	"errors"
	"io"
	"os"
	"tiktok/go/config"
	"tiktok/go/util"

	"github.com/gin-gonic/gin"

	// "tiktok/go/config"
	"tiktok/go/route"
)

/*
启动类
*/
func main() {
	//初始化项目
	initProject()
	// 1.创建路由
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	route.LoadRouter(r)
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	err := r.Run(":8000")
	if err != nil {
		panic(errors.New("项目启动失败"))
	}
}

func initProject() {
	// mysql 初始化
	config.InitDataSource()
	// redis 初始化
	config.InitRedisClient()
	// 过滤器
	util.InitSensitiveFilter()

	// 设置日志 --取消注释即可创建日志文件
	f, _ := os.Create("resources/gin.log") // // 如果文件已存在，会将文件清空。
	gin.DefaultWriter = io.MultiWriter(f)
	//gin.DebugPrintRouteFunc()

	// 如果需要同时将日志写入文件和控制台，请使用以下代码。
	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	util.Log("服务器开启成功!")
}

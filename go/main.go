package main

import (
	"fmt"
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
		fmt.Print("项目启动失败")
	}
}

func initProject() {
	// config.Init()
	config.Init()
	// 过滤器
	util.InitSensitiveFilter()
	// redis

}

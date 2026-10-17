package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"tiktok/go/config"
	"tiktok/go/model"
	"tiktok/go/route"
	"tiktok/go/util"
)

// main 启动类：初始化 DB/Redis、自动建表、注册路由、优雅关闭。
func main() {
	// MySQL 初始化（失败 fail-fast）
	db := config.InitDataSource()
	// 注入到 model 包，供各数据访问函数使用
	model.InitDB(db)
	// 自动建表（替代手跑 initial.sql）
	if err := model.AutoMigrateAll(db); err != nil {
		util.LogFatal("auto migrate failed: " + err.Error())
		panic(err)
	}
	// Redis 探活（按需，失败仅告警，不阻断启动）
	if err := config.PingAll(context.Background()); err != nil {
		util.Log("redis 探活失败（部分功能不可用）: " + err.Error())
	}
	// 过滤器
	util.InitSensitiveFilter()

	// 本地存储目录预创建（/data/{videos,covers,tmp}），确保上传/抽帧可读写
	if err := config.Storage.MkdirAll(); err != nil {
		util.Log("创建本地存储目录失败: " + err.Error())
	}

	// 配置自检：JWT 密钥若仍是代码内弱默认值，提示生产环境用环境变量覆盖
	if config.AppConfig.JWT.Secret == "this_is_a_secret_key" {
		util.Log("警告: JWT_SECRET 使用代码内弱默认值，生产环境请通过环境变量 JWT_SECRET 覆盖")
	}

	// 日志：写入 resources/gin.log（如需同时输出控制台，改为 MultiWriter(f, os.Stdout)）
	if err := os.MkdirAll("resources", 0o755); err != nil {
		util.Log("创建 resources 目录失败: " + err.Error())
	} else if f, err := os.Create("resources/gin.log"); err != nil {
		util.Log("无法打开 gin 日志文件 resources/gin.log: " + err.Error())
	} else {
		gin.DefaultWriter = io.MultiWriter(f)
	}

	r := gin.Default()
	route.LoadRouter(r)
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	// 静态资源：视频/封面经 /static/videos/*、/static/covers/* 对外提供（本地存储卷挂载到 /data）
	r.Static("/static", "/data")

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(config.AppConfig.Server.Port),
		Handler: r,
	}

	go func() {
		util.Log("服务器开启成功!")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			util.LogFatal("server error: " + err.Error())
			panic(err)
		}
	}()

	// 优雅关闭：捕获 SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	util.Log("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		util.LogFatal("forced shutdown: " + err.Error())
	}
	util.Log("server exited")
}

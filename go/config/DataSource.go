package config

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tiktok/go/util"
)

// InitDataSource 使用外部化配置初始化并返回 *gorm.DB。
// 连接具备重试能力，避免 MySQL 尚未就绪时直接 panic 导致容器崩溃。
func InitDataSource() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s",
		AppConfig.MySQL.User,
		AppConfig.MySQL.Password,
		AppConfig.MySQL.Host,
		AppConfig.MySQL.Port,
		AppConfig.MySQL.DBName,
		AppConfig.MySQL.Charset,
		AppConfig.MySQL.Loc,
	)
	const maxRetry = 10
	var db *gorm.DB
	var err error
	for i := 1; i <= maxRetry; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			PrepareStmt: true, // 缓存预编译语句
		})
		if err == nil {
			util.Log("连接成功")
			return db
		}
		util.LogError(fmt.Sprintf("MySQL 连接失败(第 %d/%d 次): %s", i, maxRetry, err.Error()))
		if i < maxRetry {
			time.Sleep(2 * time.Second)
		}
	}
	util.LogFatal("MySQL 连接重试耗尽: " + err.Error())
	panic(err)
}

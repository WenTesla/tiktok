package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tiktok/go/util"
)

//var Db *gorm.DB

// 初始化并返回数据链接
func InitDataSource() *gorm.DB {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	// 本地数据库-推荐-自己建表 -自己替换
	dsn := "root:zhang134679@tcp(127.0.0.1:3306)/tiktok?charset=utf8mb4&parseTime=True&loc=Local"
	// 阿里云数据库 非常慢
	//dsn := "tiktok:tiktok@tcp(rm-2ze62585lf96k7285mo.mysql.rds.aliyuncs.com)/tiktok?charset=utf8mb4&parseTime=True&loc=Local"
	Db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true, //缓存预编译语句
	})
	if err != nil {
		util.LogFatal(err.Error())
		panic(err)
		return nil
	}
	println("连接成功")
	return Db
}

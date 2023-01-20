package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:zhang134679@tcp(127.0.0.1:3306)/tiktok?charset=utf8mb4&parseTime=True&loc=Local"
	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true, //缓存预编译语句
	})
	if err != nil {
		//fmt.Println("连接失败")
		println("连接失败")
		return
	}
	println("连接成功")
}

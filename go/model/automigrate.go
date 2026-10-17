package model

import "gorm.io/gorm"

// AutoMigrateAll 在开发/容器启动时自动创建并增量更新所有表结构，
// 替代手跑 resources/initial.sql（SQL 文件保留作参考与种子数据）。
// 表名/列名通过各结构体的 gorm tag 与 GORM 命名约定对齐 initial.sql。
func AutoMigrateAll(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&TableVideo{},
		&Like{},
		&Comment{},
		&Follow{},
		&Message{},
	)
}

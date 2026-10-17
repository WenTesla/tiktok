package service

import (
	"errors"
	"tiktok/go/config"
)

// dataSourceErr 自定义数据库错误，用于向上层返回。
var dataSourceErr = errors.New(config.DatabaseError)

// 原有的 7 个包级 redis client 变量（var x, _ = config.RedisClient(i)）已移除：
// 它们在导入期即 Ping/panic 且吞掉错误（F4/F24）。现统一按需调用 config.RedisClient(db) 单例：
//   - 点赞计数/是否点赞：config.RedisClient(4)
//   - 消息记录：        config.RedisClient(6)
//   - 限速器：          config.RedisClient(10)/(11)（见 middle/jwt/Verify.go）

package service

import (
	"errors"
	"tiktok/go/config"
)

// 自定义数据库错误 用于返回上层
var dataSourceErr = errors.New(config.DatabaseError)

// 提供基本服务

var redisDb, _ = config.InitRedisClient()

var userRedisDb, _ = config.RedisClient(1)

var videoRedisDb, _ = config.RedisClient(2)

var followRedisDb, _ = config.RedisClient(3)

var likeRedisDb, _ = config.RedisClient(4)

var commentRedisDb, _ = config.RedisClient(5)

var messageRedisDb, _ = config.RedisClient(6)

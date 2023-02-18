package service

import "tiktok/go/config"

// 提供基本服务

var redisDb, _ = config.InitRedisClient()

var likeRedisDb, _ = config.RedisClient(3)

var commentRedisDb, _ = config.RedisClient(4)

var messageRedisDb, _ = config.RedisClient(5)

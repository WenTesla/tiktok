package middle

import "tiktok/go/config"

// 用于限速和防止大量请求
var FlowLimitRedisDbByIp, _ = config.RedisClient(10)

var FlowLimitRedisDbByUserId, _ = config.RedisClient(11)

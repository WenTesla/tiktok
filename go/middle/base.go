package middle

// 限速器 Redis 客户端改为按需 config.RedisClient(10)/(11) 获取，见 middle/jwt/Verify.go。
// 原先在此处用 var x, _ = config.RedisClient(i) 在导入期 Ping/panic（F4/F24），已移除。

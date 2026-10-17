package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"tiktok/go/util"
)

var (
	redisMu      sync.Mutex
	redisClients = map[int]*redis.Client{}
)

// RedisClient 返回绑定到逻辑库 db 的 *redis.Client 单例（懒加载并缓存，不在创建时 Ping）。
// 取代原先每个包各 new 一个 client 且导入期 Ping/panic 的写法（F4/F24）。
func RedisClient(db int) *redis.Client {
	redisMu.Lock()
	defer redisMu.Unlock()
	if c, ok := redisClients[db]; ok {
		return c
	}
	c := redis.NewClient(&redis.Options{
		Addr:     AppConfig.Redis.Addr,
		Password: AppConfig.Redis.Password,
		DB:       db,
	})
	redisClients[db] = c
	return c
}

// InitRedisClient 初始化默认库(db 0)并探活，保留以兼容 main.go。
func InitRedisClient() (*redis.Client, error) {
	c := RedisClient(0)
	if err := c.Ping(context.Background()).Err(); err != nil {
		util.LogFatal("redis init failed: " + err.Error())
		return nil, err
	}
	util.Log("redis 初始化成功")
	return c, nil
}

// PingAll 校验所有已创建的 client 连通性（供 main 启动时 fail-fast 使用）。
func PingAll(ctx context.Context) error {
	redisMu.Lock()
	defer redisMu.Unlock()
	for db, c := range redisClients {
		if err := c.Ping(ctx).Err(); err != nil {
			return fmt.Errorf("redis db %d ping failed: %w", db, err)
		}
	}
	return nil
}

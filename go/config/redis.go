package config

import (
	"github.com/go-redis/redis"
	"tiktok/go/util"
)

// 声明一个全局的redisDb变量
//var redisDb *redis.Client

// 根据redis配置初始化一个客户端

func InitRedisClient() (redisDb *redis.Client, err error) {
	redisDb = redis.NewClient(&redis.Options{
		Addr:     "43.138.126.75:6388", // redis地址
		Password: "redis",              // redis密码，没有则留空
		DB:       0,                    // 默认数据库，默认是0
	})

	//通过 *redis.Client.Ping() 来检查是否成功连接到了redis服务器
	_, err = redisDb.Ping().Result()
	if err != nil {
		util.LogFatal(err.Error())
		panic(err)
		return nil, err
	}
	util.Log("redis 初始化成功")
	return redisDb, nil
}

//  选择redis的数据库 -> i

func RedisClient(i int) (*redis.Client, error) {
	redisDb := redis.NewClient(&redis.Options{
		Addr:     "43.138.126.75:6388", // redis地址
		Password: "redis",              // redis密码，没有则留空
		DB:       i,                    // 默认数据库，默认是0
	})
	//通过 *redis.Client.Ping() 来检查是否成功连接到了redis服务器
	_, err := redisDb.Ping().Result()
	if err != nil {
		util.LogFatal(err.Error())
		panic(err)
		return nil, err
	}
	return redisDb, nil
}

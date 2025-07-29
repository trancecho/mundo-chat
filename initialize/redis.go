package initialize

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"log"
)

var (
	Client *redis.Client
)

func InitRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConns"),
	})
	if pong, err := Client.Ping().Result(); err != nil {
		log.Println("初始化redis失败:", err, "响应：", pong)
	}
}

func GetRedisClient() *redis.Client {
	return Client
}

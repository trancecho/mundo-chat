package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/trancecho/mundo-chat/initialize"
	"github.com/trancecho/mundo-chat/server"
	"log"
)

const (
	userOnlinePrefix    = "acc:user:online:" // 用户在线状态
	userOnlineCacheTime = 24 * 60 * 60
)

func getUserOnlineKey(userKey string) (key string) {
	key = fmt.Sprintf("%s%s", userOnlinePrefix, userKey)
	return
}

// GetUserOnlineInfo 获取用户在线信息
func GetUserOnlineInfo(userKey string) (userOnline *server.Client) {
	redisClient := initialize.GetRedisClient()
	key := getUserOnlineKey(userKey)
	data, err := redisClient.Get(key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Println("GetUserOnlineInfo:", userKey, err)
			return nil
		}
		log.Println("GetUserOnlineInfo:", userKey, err)
		return nil
	}
	userOnline = &server.Client{}
	err = json.Unmarshal(data, userOnline)
	if err != nil {
		log.Println("获取用户在线数据 json Unmarshal 错误:", userKey, err)
		return nil
	}
	log.Println("获取用户在线数据:", userKey, "time", userOnline.LoginTime, userOnline.HeartBeatTime)
	return userOnline
}

// SetUserOnlineInfo 设置用户在线数据
func SetUserOnlineInfo(userKey string, userOnline *server.Client) (err error) {
	redisClient := initialize.GetRedisClient()
	key := getUserOnlineKey(userKey)
	dataByte, err := json.Marshal(userOnline)
	if err != nil {
		log.Println("SetUserOnlineInfo:", key, "json Marshal 错误:", err)
		return err
	}
	_, err = redisClient.Do(context.Background(), "setEx", userOnlineCacheTime, string(dataByte)).Result()
	if err != nil {
		log.Println("SetUserOnlineInfo:", key, "设置在线数据错误:", err)
		return
	}
	return
}

package cache

import (
	"encoding/json"
	"github.com/trancecho/mundo-chat/initialize"
	"github.com/trancecho/mundo-chat/models"
	"log"
	"time"
)

const (
	//submitAgainPrefix = "acc:submit:again:" // 数据不重复提交
	chatHistoryExPire = 24 * time.Hour
	chatHistoryLimit  = 100             // 聊天记录最大保存100条
	chatHistoryKey    = "chat:history:" // 聊天记录key前缀
)

// getSubmitKey 获取数据提交去除key
//func getSubmitKey(from string, value string) (key string) {
//	key = fmt.Sprintf("%s%s:%s", submitAgainPrefix, from, value)
//	return
//}

func getChatHistoryKey(roomID string) (key string) {
	return chatHistoryKey + roomID
}

func SaveChatMessage(roomID string, message *models.Message) error {
	redisClient := initialize.GetRedisClient()
	key := getChatHistoryKey(roomID)

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Println("json化聊天内容失败:", err)
		return err
	}

	pipe := redisClient.Pipeline()
	pipe.RPush(key, string(jsonData))
	pipe.LTrim(key, -chatHistoryLimit, -1)
	pipe.Expire(key, chatHistoryExPire)

	_, err = pipe.Exec()
	if err != nil {
		log.Println("保存聊天记录到Redis失败:", err)
		return err
	}
	return nil
}

func GetChatHistory(roomID string, limit int) ([]models.Message, error) {
	if limit <= 0 {
		limit = chatHistoryLimit
	}

	redisClient := initialize.GetRedisClient()
	key := getChatHistoryKey(roomID)

	results, err := redisClient.LRange(key, 0, int64(limit-1)).Result()
	if err != nil {
		log.Println("获取聊天记录失败:", err)
		return nil, err
	}

	messages := make([]models.Message, 0)
	for _, result := range results {
		var message models.Message
		err = json.Unmarshal([]byte(result), &message)
		if err != nil {
			log.Println("解析聊天记录失败:", err)
			continue
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// submitAgain 重复提交
// return true:重复提交 false:第一次提交
//func submitAgain(from string, second int, value string) (isSubmitAgain bool) {
//	//默认重复提交
//	isSubmitAgain = true
//	key := getSubmitAgainKey(from, value)
//	redisClient := initialize.GetRedisClient()
//	number, err := redisClient.Do(context.Background(), "setNx", key, "1").Int()
//	if err != nil {
//		log.Println("submit again:", key, number, err)
//		return
//	}
//	if number == 1 {
//		return
//	}
//
//	isSubmitAgain = false
//	redisClient.Do(context.Background(), "Expire", key, second)
//	return
//}

// SeqDuplicates Seq 重复提交
//func SeqDuplicates(seq string) (result bool) {
//	result = submitAgain("seq", 12*60*60, seq)
//	return
//}

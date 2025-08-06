package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/trancecho/mundo-chat/models"
	"github.com/trancecho/mundo-chat/server/cache"
	"log"
	"net/http"
	"time"
)

var (
	Managers = NewRoomManager()
)

const (
	heartbeatExpirationTime = 6 * 60
	MaxPeople               = 10 // 房间最大人数
	CharPrefix              = "mundo-chat-"
)

func GenerateRoomKey(RoomName string) string {
	return CharPrefix + RoomName
}

func WsPage(w http.ResponseWriter, req *http.Request, RoomID string, UserID string, Username string) {
	if !Managers.RoomExists(RoomID) {
		log.Println("房间不存在:", RoomID)
		http.Error(w, "房间不存在", http.StatusNotFound)
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			log.Println("升级协议", "ua:", r.Header["User-Agent"], "referer:", r.Header["Referer"])
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	// 升级协议
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("升级协议失败:", err)
		http.Error(w, "无法建立WebSocket连接", http.StatusInternalServerError)
		return
	}
	log.Println("webSocket 建立连接:", conn.RemoteAddr().String())
	currentTime := uint64(time.Now().Unix())
	client := NewClient(conn.RemoteAddr().String(), conn, currentTime, UserID, Username)

	manager := Managers.GetRoom(RoomID)
	if manager == nil {
		log.Println("获取房间管理器失败:", RoomID)
		http.Error(w, "无法获取房间信息", http.StatusInternalServerError)
		_ = conn.Close()
		return
	}
	// 用户连接事件
	manager.Register <- client

	//获取历史记录
	history, err := cache.GetChatHistory(RoomID, 50)
	if err == nil && len(history) > 0 {
		for _, msg := range history {
			msgData, e := json.Marshal(msg)
			if e == nil {
				client.SendMsg(msgData)
			}
		}
	}

	//开启协程
	go client.Read(RoomID)
	go client.Write(RoomID)

	//发送加入通知
	joinMsg := models.NewMsg(Username, Username+"进入房间")
	joinData, _ := json.Marshal(joinMsg)
	manager.Broadcast <- joinData
}

package server

import (
	"github.com/gorilla/websocket"
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

func WsPage(w http.ResponseWriter, req *http.Request, RoomID string, UserID string) {

	// 升级协议
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		log.Println("升级协议", "ua:", r.Header["User-Agent"], "referer:", r.Header["Referer"])
		return true
	}}).Upgrade(w, req, nil)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	log.Println("webSocket 建立连接:", conn.RemoteAddr().String())
	currentTime := uint64(time.Now().Unix())
	client := NewClient(conn.RemoteAddr().String(), conn, currentTime, UserID)
	go client.read()
	go client.write(RoomID)

	// 用户连接事件
	Managers.Rooms[RoomID].Register <- client
}

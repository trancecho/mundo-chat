package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/trancecho/mundo-chat/models"
	"github.com/trancecho/mundo-chat/server/cache"
	"log"
	"runtime/debug"
)

type Client struct {
	Addr          string          // 客户端地址
	Socket        *websocket.Conn // WebSocket连接
	UserID        string          // 用户ID
	Username      string          // 用户名
	SendChan      chan []byte     // 发送消息的通道
	FirstTime     uint64          // 第一次连接时间
	HeartBeatTime uint64          // 用户上次心跳时间
	LoginTime     uint64          // 用户登录时间
}

func NewClient(addr string, socket *websocket.Conn, firstTime uint64, UserID, Username string) *Client {
	return &Client{
		UserID:        UserID,
		Username:      Username,
		Addr:          addr,
		Socket:        socket,
		SendChan:      make(chan []byte, 100), // 缓冲区大小为100
		FirstTime:     firstTime,
		HeartBeatTime: firstTime,
	}
}

func (c *Client) SendMsg(message []byte) {
	if c.Socket == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("SendMsg stop", string(debug.Stack()), r)
		}
	}()
	c.SendChan <- message
}

func (c *Client) Read(roomID string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("white stop", string(debug.Stack()), r)
		}
	}()
	defer func() {
		log.Println("用户读取通道关闭", c)
		manager := Managers.GetRoom(roomID)
		if manager != nil {
			manager.Unregister <- c
		}
		close(c.SendChan)
	}()
	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			log.Printf("用户读取消息失败: %s, 错误: %v", c.Addr, err)
			return
		}
		log.Printf("用户读取消息: %s, 内容: %s", c.Addr, message)

		manager := Managers.GetRoom(roomID)
		if manager == nil {
			log.Println("获取房间管理器失败:", roomID)
			continue
		}

		m := models.NewMsg(c.Username, string(message))
		jsonMessage, err := json.Marshal(m)
		if err != nil {
			log.Println("消息序列化失败:", err)
			continue
		}
		//保存消息到redis缓存
		if err = cache.SaveChatMessage(roomID, m); err != nil {
			log.Println("保存聊天记录到Redis失败:", err)
		}

		manager.Broadcast <- jsonMessage
	}
}

func (c *Client) Login(uid string, loginTime uint64) {
	c.LoginTime = loginTime
	c.UserID = uid
	c.HeartBeatTime = loginTime
}

func (c *Client) IsLogin() bool {
	return c.UserID != ""
}

func (c *Client) Write(roomID string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("write stop", string(debug.Stack()), r)
		}
	}()
	defer func() {
		Managers.Rooms[roomID].Unregister <- c
		_ = c.Socket.Close()
		log.Println("用户写入通道关闭defer", c)
	}()
	for {
		select {
		case message, ok := <-c.SendChan:
			if !ok {
				log.Println("Client发送数据关闭连接:", c.Addr, "ok", ok)
				return
			}
			err := c.Socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("用户发送消息失败: %s, 错误: %v", c.Addr, err)
				return
			}
		}
	}
}

func (c *Client) IsHeartBeatTimeout(currentTime uint64) bool {
	return c.HeartBeatTime+uint64(heartbeatExpirationTime) <= currentTime
}

func (c *Client) HeartBeat(currentTime uint64) {
	c.HeartBeatTime = currentTime
	return
}

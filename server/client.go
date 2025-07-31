package server

import (
	"github.com/gorilla/websocket"
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

func NewClient(addr string, socket *websocket.Conn, firstTime uint64, UserID string) *Client {
	return &Client{
		UserID:        UserID,
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

func (c *Client) read() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("white stop", string(debug.Stack()), r)
		}
	}()
	defer func() {
		log.Println("用户读取通道关闭", c)
		close(c.SendChan)
	}()
	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			log.Printf("用户读取消息失败: %s, 错误: %v", c.Addr, err)
			return
		}
		log.Printf("用户读取消息: %s, 内容: %s", c.Addr, message)
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

func (c *Client) write(roomID string) {
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
			_ = c.Socket.WriteMessage(websocket.TextMessage, message)
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

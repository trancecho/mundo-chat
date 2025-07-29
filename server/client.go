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
	SendChan      chan []byte     // 发送消息的通道
	FirstTime     uint64          // 第一次连接时间
	HeartBeatTime uint64          // 用户上次心跳时间
	LoginTime     uint64          // 用户登录时间
}

func NewClient(socket *websocket.Conn, firstTime uint64) *Client {
	return &Client{
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

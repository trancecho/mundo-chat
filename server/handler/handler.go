package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/trancecho/mundo-chat/server"
	"github.com/trancecho/mundo-chat/server/common"
	"github.com/trancecho/mundo-chat/server/libx"
	"time"
)

func Chat(c *gin.Context) {
	uid, code := libx.Uid(c), common.OK
	currentTime := uint64(time.Now().Unix())
	var conn *websocket.Conn
	Client := server.NewClient(conn, currentTime)
}

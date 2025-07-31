package router

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-chat/server/handler"
	"github.com/trancecho/mundo-chat/server/middleware"
)

func GenerateRouter(r *gin.Engine) {

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	v1 := r.Group("/api", middleware.JWTAuthMiddleware())
	{
		v1.POST("/add_room", handler.AddRoom)         // 添加房间
		v1.GET("/get_all_rooms", handler.GetAllRooms) // 获取所有房间
		v1.GET("/room_exists", handler.RoomExists)    // 检查房间是否存在
		v1.GET("/chat", handler.Chat)                 // 进入聊天室
	}
}

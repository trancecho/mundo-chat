package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-chat/models"
	"github.com/trancecho/mundo-chat/server"
	"github.com/trancecho/mundo-chat/server/common"
	"github.com/trancecho/mundo-chat/server/libx"
	"strconv"
)

func AddRoom(c *gin.Context) {
	var request models.AddRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		libx.Err(c, common.ParameterIllegal, "Invalid request parameters", err)
		return
	}
	if request.RoomName == "" {
		libx.Err(c, common.ParameterIllegal, "Invalid request parameters", nil)
		return
	}
	if request.MaxPeople == 0 {
		request.MaxPeople = server.MaxPeople
	}
	roomID := server.GenerateRoomKey(request.RoomName)
	manager := server.NewManagerWithInfo(roomID, request.RoomName, request.Description, request.MaxPeople)
	server.Managers.Register <- manager
	libx.Ok(c, "Room created successfully", roomID)
}

func GetAllRooms(c *gin.Context) {
	rooms := server.Managers.GetAllRooms()
	if len(rooms) == 0 {
		libx.Err(c, common.NotData, "No rooms available", nil)
		return
	}
	libx.Ok(c, rooms)
}

func RoomExists(c *gin.Context) {
	roomName := c.Query("room_name")
	if roomName == "" {
		libx.Err(c, common.ParameterIllegal, "Room name is required", nil)
		return
	}
	roomID := server.GenerateRoomKey(roomName)
	exists := server.Managers.RoomExists(roomID)
	if exists {
		libx.Ok(c, "Room exists:"+roomID)
	} else {
		libx.Err(c, common.NotData, "Room does not exist", nil)
	}
}

func Chat(c *gin.Context) {
	roomName := c.Query("room_name")
	uid := libx.Uid(c)
	server.WsPage(c.Writer, c.Request, server.GenerateRoomKey(roomName), strconv.Itoa(int(uid)))
}

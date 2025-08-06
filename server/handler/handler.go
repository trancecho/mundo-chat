package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-chat/models"
	"github.com/trancecho/mundo-chat/server"
	"github.com/trancecho/mundo-chat/server/libx"
	"net/http"
	"strconv"
)

func AddRoom(c *gin.Context) {
	var request models.AddRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		libx.Err(c, http.StatusBadRequest, "Invalid request parameters", err)
		return
	}
	if request.RoomName == "" {
		libx.Err(c, http.StatusBadRequest, "Invalid request parameters", nil)
		return
	}
	if request.MaxPeople == 0 {
		request.MaxPeople = server.MaxPeople
	}
	roomID := server.GenerateRoomKey(request.RoomName)
	manager := server.NewManagerWithInfo(roomID, request.RoomName, request.Description, request.MaxPeople)
	server.Managers.Register <- manager
	libx.Ok(c, "Room added successfully", roomID)
}

func GetAllRooms(c *gin.Context) {
	rooms := server.Managers.GetAllRooms()
	if len(rooms) == 0 {
		libx.Err(c, http.StatusInternalServerError, "No rooms available", nil)
		return
	}
	libx.Ok(c, "Rooms retrieved successfully", rooms)
}

func RoomExists(c *gin.Context) {
	roomName := c.Query("room_name")
	if roomName == "" {
		libx.Err(c, http.StatusBadRequest, "Room name is required", nil)
		return
	}
	roomID := server.GenerateRoomKey(roomName)
	exists := server.Managers.RoomExists(roomID)
	if exists {
		libx.Ok(c, "Room exists:"+roomID)
	} else {
		libx.Err(c, http.StatusInternalServerError, "Room does not exist", nil)
	}
}

func GetRoomInfo(c *gin.Context) {
	roomName := c.Query("room_name")
	if roomName == "" {
		libx.Err(c, http.StatusBadRequest, "Room name is required", nil)
		return
	}
	roomID := server.GenerateRoomKey(roomName)
	roomInfo := server.Managers.GetRoomInfo(roomID)
	if roomInfo == nil {
		libx.Err(c, http.StatusInternalServerError, "Room does not exist", nil)
		return
	}
	libx.Ok(c, "Room info retrieved successfully", roomInfo)
}

func Chat(c *gin.Context) {
	roomName := c.Query("room_name")
	if roomName == "" {
		libx.Err(c, http.StatusBadRequest, "Room name is required", nil)
		return
	}
	uid := libx.Uid(c)
	if uid <= 0 {
		libx.Err(c, http.StatusUnauthorized, "Invalid user ID", nil)
		return
	}
	roomID := server.GenerateRoomKey(roomName)
	if !server.Managers.RoomExists(roomID) {
		libx.Err(c, http.StatusInternalServerError, "Room does not exist", nil)
		return
	}
	name := libx.GetUsername(c)
	server.WsPage(c.Writer, c.Request, roomID, strconv.Itoa(int(uid)), name)
}

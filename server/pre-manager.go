package server

import (
	"context"
	"sync"
)

// RoomManager 聊天室房间管理器
type RoomManager struct {
	Rooms      map[string]*Manager // 房间映射，key为房间ID，value为Manager实例
	RoomsLock  sync.RWMutex        // 房间锁
	Register   chan *Manager       // 注册新房间
	Unregister chan string         // 注销房间
	Broadcast  chan []byte         // 广播消息
}

// NewRoomManager 创建新的房间管理器
func NewRoomManager() *RoomManager {
	return &RoomManager{
		Rooms:      make(map[string]*Manager),
		RoomsLock:  sync.RWMutex{},
		Register:   make(chan *Manager, 100),
		Unregister: make(chan string, 100),
		Broadcast:  make(chan []byte, 1000),
	}
}

// RemoveRoom 删除房间
func (rm *RoomManager) RemoveRoom(roomID string) bool {
	rm.RoomsLock.Lock()
	defer rm.RoomsLock.Unlock()

	manager, exists := rm.Rooms[roomID]
	if !exists {
		return false
	}

	// 断开所有用户连接
	clients := manager.AllClients()
	for client := range clients {
		client.Socket.Close()
	}

	// 从映射中删除
	delete(rm.Rooms, roomID)

	return true
}

// GetRoom 获取房间Manager
func (rm *RoomManager) GetRoom(roomID string) *Manager {
	rm.RoomsLock.RLock()
	defer rm.RoomsLock.RUnlock()

	manager, exists := rm.Rooms[roomID]
	if !exists {
		return nil
	}
	return manager
}

// GetAllRooms 获取所有房间
func (rm *RoomManager) GetAllRooms() map[string]*Manager {
	rm.RoomsLock.RLock()
	defer rm.RoomsLock.RUnlock()

	roomsCopy := make(map[string]*Manager, len(rm.Rooms))
	for key, value := range rm.Rooms {
		roomsCopy[key] = value
	}
	return roomsCopy
}

// GetRoomCount 获取房间数量
func (rm *RoomManager) GetRoomCount() int {
	rm.RoomsLock.RLock()
	defer rm.RoomsLock.RUnlock()
	return len(rm.Rooms)
}

// RoomExists 检查房间是否存在
func (rm *RoomManager) RoomExists(roomID string) bool {
	rm.RoomsLock.RLock()
	defer rm.RoomsLock.RUnlock()
	_, exists := rm.Rooms[roomID]
	return exists
}

// AddUserToRoom 添加用户到房间
func (rm *RoomManager) AddUserToRoom(roomID string, client *Client) bool {
	manager := rm.GetRoom(roomID)
	if manager == nil {
		return false
	}

	// 检查房间是否已满
	if manager.IsRoomFull() {
		return false
	}

	// 检查用户是否已在房间中
	if manager.GetUserByID(client.UserID) != nil {
		return false
	}

	// 添加用户到Manager
	manager.AddClient(client)

	return true
}

// RemoveUserFromRoom 从房间移除用户
func (rm *RoomManager) RemoveUserFromRoom(roomID string, userID string) bool {
	manager := rm.GetRoom(roomID)
	if manager == nil {
		return false
	}

	client := manager.GetUserByID(userID)
	if client == nil {
		return false
	}

	// 从Manager移除用户
	manager.RemoveClient(client)

	return true
}

// GetRoomUsers 获取房间内的用户
func (rm *RoomManager) GetRoomUsers(roomID string) map[*Client]bool {
	manager := rm.GetRoom(roomID)
	if manager == nil {
		return nil
	}
	return manager.AllClients()
}

// BroadcastToRoom 向房间广播消息
func (rm *RoomManager) BroadcastToRoom(roomID string, message []byte) {
	manager := rm.GetRoom(roomID)
	if manager == nil {
		return
	}

	// 使用Manager的广播功能
	manager.Broadcast <- message
}

// BroadcastToAllRooms 向所有房间广播消息
func (rm *RoomManager) BroadcastToAllRooms(message []byte) {
	rm.RoomsLock.RLock()
	defer rm.RoomsLock.RUnlock()

	for _, manager := range rm.Rooms {
		select {
		case manager.Broadcast <- message:
		default:
			// 如果广播通道满了，跳过
		}
	}
}

// GetRoomInfo 获取房间信息
func (rm *RoomManager) GetRoomInfo(roomID string) map[string]interface{} {
	manager := rm.GetRoom(roomID)
	if manager == nil {
		return nil
	}
	return manager.GetRoomInfo()
}

// GetAllRoomsInfo 获取所有房间信息
func (rm *RoomManager) GetAllRoomsInfo() map[string]map[string]interface{} {
	rm.RoomsLock.RLock()
	defer rm.RoomsLock.RUnlock()

	roomsInfo := make(map[string]map[string]interface{})
	for roomID, manager := range rm.Rooms {
		roomsInfo[roomID] = manager.GetRoomInfo()
	}
	return roomsInfo
}

func (rm *RoomManager) Start(ctx context.Context) {
	// 使用context来控制管理器的生命周期
	for {
		select {
		case manager := <-rm.Register:
			// 使用接收到的manager进行房间注册
			rm.RoomsLock.Lock()
			if _, exists := rm.Rooms[manager.ID]; exists {
				// 房间ID已存在，可以记录日志或通知
				rm.RoomsLock.Unlock()
				continue
			}
			rm.Rooms[manager.ID] = manager
			rm.RoomsLock.Unlock()

			// 启动房间管理
			go manager.Start()

		case roomID := <-rm.Unregister:
			rm.RemoveRoom(roomID)

		case message := <-rm.Broadcast:
			// 全局广播到所有房间
			rm.BroadcastToAllRooms(message)

		case <-ctx.Done():
			// 优雅退出
			rm.shutdown()
			return
		}
	}
}

// shutdown 关闭所有房间连接
func (rm *RoomManager) shutdown() {
	rm.RoomsLock.Lock()
	defer rm.RoomsLock.Unlock()

	// 关闭所有房间
	for roomID, manager := range rm.Rooms {
		// 断开所有用户连接
		clients := manager.AllClients()
		for client := range clients {
			client.Socket.Close()
		}

		// 从映射中删除
		delete(rm.Rooms, roomID)
	}
}

package server

import (
	"sync"
	"time"
)

type Manager struct {
	// 房间的信息
	ID          string    `json:"id"`          // 房间ID
	Name        string    `json:"name"`        // 房间的名称
	Description string    `json:"description"` // 房间描述
	MaxUsers    int       `json:"maxUsers"`    // 最大用户数
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`   // 更新时间
	IsActive    bool      `json:"isActive"`    // 是否活跃

	// 客户端管理
	Clients     map[*Client]bool   //全部连接
	ClientsLock sync.RWMutex       //读写锁
	UserMap     map[string]*Client //用户ID到客户端的映射
	UserMapLock sync.RWMutex       //读写锁
	Register    chan *Client       //注册新连接
	Unregister  chan *Client       //注销连接
	Broadcast   chan []byte        //广播消息
}

func NewManager() *Manager {
	return &Manager{
		Clients:     make(map[*Client]bool),
		ClientsLock: sync.RWMutex{},
		UserMap:     make(map[string]*Client),
		UserMapLock: sync.RWMutex{},
		Register:    make(chan *Client, 1000),
		Unregister:  make(chan *Client, 1000),
		Broadcast:   make(chan []byte, 1000), // 缓冲区大小为100
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		MaxUsers:    MaxPeople, // 默认最大用户数为10
	}
}

// NewManagerWithInfo 创建带房间信息的Manager
func NewManagerWithInfo(id, name, description string, maxUsers int) *Manager {
	manager := NewManager()
	manager.ID = id
	manager.Name = name
	manager.Description = description
	manager.MaxUsers = maxUsers
	return manager
}

// GetUserCount 获取当前用户数
func (m *Manager) GetUserCount() int {
	return m.GetClientLen()
}

// IsRoomFull 检查房间是否已满
func (m *Manager) IsRoomFull() bool {
	return m.GetClientLen() >= m.MaxUsers
}

// GetRoomInfo 获取房间信息
func (m *Manager) GetRoomInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":          m.ID,
		"name":        m.Name,
		"description": m.Description,
		"maxUsers":    m.MaxUsers,
		"userCount":   m.GetClientLen(),
		"createdAt":   m.CreatedAt,
		"updatedAt":   m.UpdatedAt,
		"isActive":    m.IsActive,
	}
}

func (m *Manager) ClientExist(client *Client) bool {
	m.ClientsLock.RLock()
	defer m.ClientsLock.RUnlock()
	_, exists := m.Clients[client]
	return exists
}

func (m *Manager) AllClients() map[*Client]bool {
	m.ClientsLock.RLock()
	defer m.ClientsLock.RUnlock()
	clientsCopy := make(map[*Client]bool, len(m.Clients))
	for key, value := range m.Clients {
		clientsCopy[key] = value
	}
	return clientsCopy
}

func (m *Manager) AddClient(client *Client) {
	m.ClientsLock.Lock()
	defer m.ClientsLock.Unlock()
	m.Clients[client] = true
	m.UserMapLock.Lock()
	defer m.UserMapLock.Unlock()
	m.UserMap[client.UserID] = client
	m.UpdatedAt = time.Now()
}

func (m *Manager) RemoveClient(client *Client) {
	m.ClientsLock.Lock()
	defer m.ClientsLock.Unlock()
	if _, exists := m.Clients[client]; exists {
		delete(m.Clients, client)
	}
	m.UserMapLock.Lock()
	defer m.UserMapLock.Unlock()
	if _, exists := m.UserMap[client.UserID]; exists {
		delete(m.UserMap, client.UserID)
	}
	m.UpdatedAt = time.Now()
}

func (m *Manager) GetUserByID(userID string) *Client {
	m.UserMapLock.RLock()
	defer m.UserMapLock.RUnlock()
	client, exists := m.UserMap[userID]
	if !exists {
		return nil
	}
	return client
}

func (m *Manager) GetClientLen() int {
	m.ClientsLock.RLock()
	defer m.ClientsLock.RUnlock()
	return len(m.Clients)
}

func (m *Manager) Start() {
	for {
		select {
		case client := <-m.Register:
			if !m.ClientExist(client) {
				m.AddClient(client)
			}
		case client := <-m.Unregister:
			if m.ClientExist(client) {
				m.RemoveClient(client)
			}
		case message := <-m.Broadcast:
			clients := m.AllClients()
			for conn := range clients {
				select {
				case conn.SendChan <- message:
				default:
					close(conn.SendChan)
				}
			}
		}
	}
}

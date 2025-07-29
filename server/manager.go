package server

import "sync"

type Manager struct {
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

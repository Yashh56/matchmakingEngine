package clients

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Manager struct {
	mu      sync.RWMutex
	clients map[string]*websocket.Conn
}

func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]*websocket.Conn),
	}
}

func (m *Manager) Add(id string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[id] = conn
}

func (m *Manager) Get(id string) *websocket.Conn {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.clients[id]
}

func (m *Manager) Remove(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clients, id)
}

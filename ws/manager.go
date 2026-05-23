package ws

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/Sna-ken/videoweb/biz/service/chat"
	"github.com/Sna-ken/videoweb/internal/model"
)

type Manager struct {
	Clients      map[string]*Client
	Lock         sync.RWMutex
	OnNewMessage func(msg *model.Message)
}

var ManagerInstance = &Manager{
	Clients: make(map[string]*Client),
}

func (m *Manager) AddClient(client *Client) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.Clients[client.UserID] = client
}

func (m *Manager) RemoveClient(userID string) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	delete(m.Clients, userID)
}

func (m *Manager) BroadcastMessage(message []byte, sender *Client) {
	ctx := context.Background()

	msg, err := chat.SaveMessage(ctx, sender.UserID, sender.Username, string(message))
	if err != nil {
		log.Printf("save massage error:%v", err)
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("marshal error:%v", err)
		return
	}

	m.Lock.RLock()
	targets := make([]*Client, 0, len(message))
	for _, client := range m.Clients {
		targets = append(targets, client)
	}
	m.Lock.RUnlock()

	for _, client := range m.Clients {
		client.Send <- data
	}

	if m.OnNewMessage != nil {
		go m.OnNewMessage(msg)
	} else {
		log.Printf("OnNewMessage is nil")
	}
}

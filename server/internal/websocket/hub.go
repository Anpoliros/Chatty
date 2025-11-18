package websocket

import (
	"sync"

	"github.com/Anpoliros/Chatty/server/pkg/logger"
)

// Hub 维护所有活跃的客户端连接和广播消息
type Hub struct {
	// 已注册的客户端
	clients map[uint]*Client // key是用户ID

	// 从客户端接收的消息
	broadcast chan *Message

	// 注册请求
	register chan *Client

	// 注销请求
	unregister chan *Client

	// 用户ID到客户端的映射锁
	mu sync.RWMutex

	log *logger.Logger
}

// NewHub 创建新的Hub实例
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]*Client),
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		log:        logger.New(),
	}
}

// Run 运行Hub,处理注册、注销和广播
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			h.log.Info("Client registered:", client.UserID)

			// 发送在线状态更新
			// TODO: 广播用户上线通知

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.send)
				h.log.Info("Client unregistered:", client.UserID)
			}
			h.mu.Unlock()

			// 发送离线状态更新
			// TODO: 广播用户下线通知

		case message := <-h.broadcast:
			h.handleBroadcast(message)
		}
	}
}

// handleBroadcast 处理消息广播
func (h *Hub) handleBroadcast(message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	switch message.Type {
	case MessageTypePrivate:
		// 私聊消息:只发送给接收者
		if client, ok := h.clients[message.ReceiverID]; ok {
			select {
			case client.send <- message:
			default:
				// 发送失败,关闭连接
				close(client.send)
				delete(h.clients, client.UserID)
			}
		}

	case MessageTypeGroup:
		// 群聊消息:发送给所有群成员
		// TODO: 从数据库获取群成员列表
		// 暂时简单广播给所有在线用户
		for userID, client := range h.clients {
			if userID != message.SenderID {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, userID)
				}
			}
		}

	case MessageTypeSystem:
		// 系统消息:广播给所有在线用户
		for userID, client := range h.clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, userID)
			}
		}
	}
}

// SendToUser 发送消息给指定用户
func (h *Hub) SendToUser(userID uint, message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.clients[userID]; ok {
		select {
		case client.send <- message:
		default:
			h.log.Warn("Failed to send message to user:", userID)
		}
	}
}

// GetOnlineUsers 获取在线用户列表
func (h *Hub) GetOnlineUsers() []uint {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]uint, 0, len(h.clients))
	for userID := range h.clients {
		users = append(users, userID)
	}
	return users
}

// IsUserOnline 检查用户是否在线
func (h *Hub) IsUserOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, ok := h.clients[userID]
	return ok
}

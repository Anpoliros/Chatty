package websocket

import (
	"encoding/json"
	"time"

	"github.com/Anpoliros/Chatty/server/pkg/logger"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 8192
)

// Client 表示一个WebSocket客户端连接
type Client struct {
	// WebSocket连接
	conn *websocket.Conn

	// Hub引用
	hub *Hub

	// 用户ID
	UserID uint

	// 发送消息的通道
	send chan *Message

	log *logger.Logger
}

// NewClient 创建新的客户端
func NewClient(conn *websocket.Conn, hub *Hub, userID uint) *Client {
	return &Client{
		conn:   conn,
		hub:    hub,
		UserID: userID,
		send:   make(chan *Message, 256),
		log:    logger.New(),
	}
}

// readPump 从WebSocket连接读取消息
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	c.conn.SetReadLimit(maxMessageSize)

	for {
		_, messageData, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.log.Error("WebSocket error:", err)
			}
			break
		}

		// 解析消息
		var msg Message
		if err := json.Unmarshal(messageData, &msg); err != nil {
			c.log.Error("Failed to parse message:", err)
			continue
		}

		// 设置发送者ID
		msg.SenderID = c.UserID
		msg.Timestamp = time.Now()

		// TODO: 这里应该先保存消息到数据库,然后再广播
		// 暂时直接广播
		c.hub.broadcast <- &msg
	}
}

// writePump 向WebSocket连接写入消息
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub关闭了通道
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 发送消息
			data, err := json.Marshal(message)
			if err != nil {
				c.log.Error("Failed to marshal message:", err)
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				c.log.Error("Failed to write message:", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Start 启动客户端读写循环
func (c *Client) Start() {
	go c.writePump()
	go c.readPump()
}

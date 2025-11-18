package websocket

import "time"

// MessageType WebSocket消息类型
type MessageType string

const (
	MessageTypePrivate MessageType = "private" // 私聊消息
	MessageTypeGroup   MessageType = "group"   // 群聊消息
	MessageTypeSystem  MessageType = "system"  // 系统消息
	MessageTypeTyping  MessageType = "typing"  // 正在输入
	MessageTypeRead    MessageType = "read"    // 已读回执
	MessageTypePing    MessageType = "ping"    // 心跳
	MessageTypePong    MessageType = "pong"    // 心跳响应
)

// Message WebSocket消息结构
type Message struct {
	Type           MessageType `json:"type"`
	SenderID       uint        `json:"sender_id"`
	ReceiverID     uint        `json:"receiver_id,omitempty"`      // 私聊接收者ID
	ConversationID uint        `json:"conversation_id,omitempty"`  // 会话ID
	Content        string      `json:"content,omitempty"`          // 消息内容
	MediaURL       string      `json:"media_url,omitempty"`        // 媒体URL
	MessageID      uint        `json:"message_id,omitempty"`       // 消息ID(数据库ID)
	Timestamp      time.Time   `json:"timestamp"`                  // 时间戳
	Extra          interface{} `json:"extra,omitempty"`            // 额外数据
}

// TypingMessage 正在输入消息
type TypingMessage struct {
	ConversationID uint `json:"conversation_id"`
	IsTyping       bool `json:"is_typing"`
}

// ReadMessage 已读回执消息
type ReadMessage struct {
	ConversationID uint `json:"conversation_id"`
	MessageID      uint `json:"message_id"`
}

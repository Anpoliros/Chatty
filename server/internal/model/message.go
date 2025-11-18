package model

import (
	"time"
)

// MessageType 消息类型
type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeFile  MessageType = "file"
	MessageTypeAudio MessageType = "audio"
	MessageTypeVideo MessageType = "video"
)

// Message 消息模型
type Message struct {
	ID             uint        `json:"id" gorm:"primaryKey"`
	ConversationID uint        `json:"conversation_id" gorm:"index;not null"`
	SenderID       uint        `json:"sender_id" gorm:"index;not null"`
	Type           MessageType `json:"type" gorm:"default:'text'"`
	Content        string      `json:"content" gorm:"type:text"`
	MediaURL       string      `json:"media_url,omitempty"` // 图片/文件/音视频URL
	IsRead         bool        `json:"is_read" gorm:"default:false"`
	IsDeleted      bool        `json:"is_deleted" gorm:"default:false"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`

	// 关联
	Sender *User `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}

// Conversation 会话模型
type Conversation struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Type          string    `json:"type" gorm:"default:'private'"` // private, group
	Name          string    `json:"name"`                          // 群聊名称
	Avatar        string    `json:"avatar"`                        // 群聊头像
	LastMessageID uint      `json:"last_message_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// 关联
	LastMessage  *Message               `json:"last_message,omitempty" gorm:"foreignKey:LastMessageID"`
	Participants []*ConversationMember  `json:"participants,omitempty" gorm:"foreignKey:ConversationID"`
}

// TableName 指定表名
func (Conversation) TableName() string {
	return "conversations"
}

// ConversationMember 会话成员
type ConversationMember struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	ConversationID uint      `json:"conversation_id" gorm:"index;not null"`
	UserID         uint      `json:"user_id" gorm:"index;not null"`
	UnreadCount    int       `json:"unread_count" gorm:"default:0"`
	JoinedAt       time.Time `json:"joined_at"`
	LastReadAt     time.Time `json:"last_read_at"`

	// 关联
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (ConversationMember) TableName() string {
	return "conversation_members"
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	ConversationID uint        `json:"conversation_id" binding:"required"`
	Type           MessageType `json:"type" binding:"required"`
	Content        string      `json:"content"`
	MediaURL       string      `json:"media_url,omitempty"`
}

// MessageResponse 消息响应
type MessageResponse struct {
	ID             uint         `json:"id"`
	ConversationID uint         `json:"conversation_id"`
	SenderID       uint         `json:"sender_id"`
	Type           MessageType  `json:"type"`
	Content        string       `json:"content"`
	MediaURL       string       `json:"media_url,omitempty"`
	IsRead         bool         `json:"is_read"`
	CreatedAt      time.Time    `json:"created_at"`
	Sender         *UserResponse `json:"sender,omitempty"`
}

// ToResponse 转换为响应格式
func (m *Message) ToResponse() *MessageResponse {
	resp := &MessageResponse{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		SenderID:       m.SenderID,
		Type:           m.Type,
		Content:        m.Content,
		MediaURL:       m.MediaURL,
		IsRead:         m.IsRead,
		CreatedAt:      m.CreatedAt,
	}

	if m.Sender != nil {
		resp.Sender = m.Sender.ToResponse()
	}

	return resp
}

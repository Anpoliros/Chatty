package repository

import (
	"errors"

	"github.com/Anpoliros/Chatty/server/internal/model"
	"gorm.io/gorm"
)

// MessageRepository 消息数据访问层
type MessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository 创建消息仓库
func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// CreateMessage 创建消息
func (r *MessageRepository) CreateMessage(message *model.Message) error {
	return r.db.Create(message).Error
}

// GetMessageByID 根据ID获取消息
func (r *MessageRepository) GetMessageByID(id uint) (*model.Message, error) {
	var message model.Message
	err := r.db.Preload("Sender").First(&message, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("message not found")
		}
		return nil, err
	}
	return &message, nil
}

// GetMessages 获取会话的消息列表
func (r *MessageRepository) GetMessages(conversationID uint, limit, offset int) ([]*model.Message, error) {
	var messages []*model.Message
	err := r.db.Where("conversation_id = ? AND is_deleted = ?", conversationID, false).
		Preload("Sender").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

// UpdateMessage 更新消息
func (r *MessageRepository) UpdateMessage(message *model.Message) error {
	return r.db.Save(message).Error
}

// DeleteMessage 删除消息(软删除)
func (r *MessageRepository) DeleteMessage(id uint) error {
	return r.db.Model(&model.Message{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// MarkAsRead 标记消息为已读
func (r *MessageRepository) MarkAsRead(messageID uint) error {
	return r.db.Model(&model.Message{}).Where("id = ?", messageID).Update("is_read", true).Error
}

// MarkConversationAsRead 标记会话中所有消息为已读
func (r *MessageRepository) MarkConversationAsRead(conversationID, userID uint) error {
	return r.db.Model(&model.Message{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = ?", conversationID, userID, false).
		Update("is_read", true).Error
}

// GetUnreadCount 获取未读消息数
func (r *MessageRepository) GetUnreadCount(conversationID, userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Message{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = ?", conversationID, userID, false).
		Count(&count).Error
	return count, err
}

// CreateConversation 创建会话
func (r *MessageRepository) CreateConversation(conversation *model.Conversation) error {
	return r.db.Create(conversation).Error
}

// GetConversationByID 根据ID获取会话
func (r *MessageRepository) GetConversationByID(id uint) (*model.Conversation, error) {
	var conversation model.Conversation
	err := r.db.Preload("LastMessage").
		Preload("Participants").
		Preload("Participants.User").
		First(&conversation, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("conversation not found")
		}
		return nil, err
	}
	return &conversation, nil
}

// GetUserConversations 获取用户的会话列表
func (r *MessageRepository) GetUserConversations(userID uint) ([]*model.Conversation, error) {
	var conversations []*model.Conversation
	err := r.db.
		Joins("JOIN conversation_members ON conversation_members.conversation_id = conversations.id").
		Where("conversation_members.user_id = ?", userID).
		Preload("LastMessage").
		Preload("LastMessage.Sender").
		Preload("Participants").
		Preload("Participants.User").
		Order("conversations.updated_at DESC").
		Find(&conversations).Error
	return conversations, err
}

// FindPrivateConversation 查找两个用户之间的私聊会话
func (r *MessageRepository) FindPrivateConversation(userID1, userID2 uint) (*model.Conversation, error) {
	var conversation model.Conversation
	err := r.db.
		Joins("JOIN conversation_members cm1 ON cm1.conversation_id = conversations.id AND cm1.user_id = ?", userID1).
		Joins("JOIN conversation_members cm2 ON cm2.conversation_id = conversations.id AND cm2.user_id = ?", userID2).
		Where("conversations.type = ?", "private").
		First(&conversation).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 未找到会话返回nil而不是错误
		}
		return nil, err
	}
	return &conversation, nil
}

// AddConversationMember 添加会话成员
func (r *MessageRepository) AddConversationMember(member *model.ConversationMember) error {
	return r.db.Create(member).Error
}

// UpdateConversation 更新会话
func (r *MessageRepository) UpdateConversation(conversation *model.Conversation) error {
	return r.db.Save(conversation).Error
}

// UpdateMemberUnreadCount 更新成员未读数
func (r *MessageRepository) UpdateMemberUnreadCount(conversationID, userID uint, increment bool) error {
	if increment {
		return r.db.Model(&model.ConversationMember{}).
			Where("conversation_id = ? AND user_id = ?", conversationID, userID).
			UpdateColumn("unread_count", gorm.Expr("unread_count + ?", 1)).Error
	}
	return r.db.Model(&model.ConversationMember{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Update("unread_count", 0).Error
}

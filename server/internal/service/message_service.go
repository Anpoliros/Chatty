package service

import (
	"errors"
	"time"

	"github.com/Anpoliros/Chatty/server/internal/model"
	"github.com/Anpoliros/Chatty/server/internal/repository"
)

// MessageService 消息服务
type MessageService struct {
	messageRepo *repository.MessageRepository
	userRepo    *repository.UserRepository
}

// NewMessageService 创建消息服务
func NewMessageService(messageRepo *repository.MessageRepository, userRepo *repository.UserRepository) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
	}
}

// SendMessage 发送消息
func (s *MessageService) SendMessage(senderID, conversationID uint, msgType model.MessageType, content, mediaURL string) (*model.Message, error) {
	// 验证会话是否存在
	conversation, err := s.messageRepo.GetConversationByID(conversationID)
	if err != nil {
		return nil, err
	}

	// TODO: 验证用户是否是会话成员

	// 创建消息
	message := &model.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Type:           msgType,
		Content:        content,
		MediaURL:       mediaURL,
		IsRead:         false,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
	}

	if err := s.messageRepo.CreateMessage(message); err != nil {
		return nil, err
	}

	// 更新会话的最后一条消息
	conversation.LastMessageID = message.ID
	conversation.UpdatedAt = time.Now()
	if err := s.messageRepo.UpdateConversation(conversation); err != nil {
		return nil, err
	}

	// 加载发送者信息
	sender, _ := s.userRepo.GetByID(senderID)
	message.Sender = sender

	// TODO: 增加其他成员的未读数

	return message, nil
}

// GetMessages 获取消息列表
func (s *MessageService) GetMessages(conversationID uint, page, pageSize int) ([]*model.Message, error) {
	offset := (page - 1) * pageSize
	return s.messageRepo.GetMessages(conversationID, pageSize, offset)
}

// GetConversations 获取用户的会话列表
func (s *MessageService) GetConversations(userID uint) ([]*model.Conversation, error) {
	return s.messageRepo.GetUserConversations(userID)
}

// CreatePrivateConversation 创建私聊会话
func (s *MessageService) CreatePrivateConversation(userID1, userID2 uint) (*model.Conversation, error) {
	// 检查是否已存在私聊会话
	existingConv, err := s.messageRepo.FindPrivateConversation(userID1, userID2)
	if err != nil {
		return nil, err
	}
	if existingConv != nil {
		return existingConv, nil
	}

	// 验证用户是否存在
	user1, err := s.userRepo.GetByID(userID1)
	if err != nil {
		return nil, errors.New("user1 not found")
	}
	user2, err := s.userRepo.GetByID(userID2)
	if err != nil {
		return nil, errors.New("user2 not found")
	}

	// 创建会话
	conversation := &model.Conversation{
		Type:      "private",
		Name:      user2.Nickname, // 对于私聊,名称可以是对方的昵称
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.messageRepo.CreateConversation(conversation); err != nil {
		return nil, err
	}

	// 添加两个成员
	member1 := &model.ConversationMember{
		ConversationID: conversation.ID,
		UserID:         userID1,
		UnreadCount:    0,
		JoinedAt:       time.Now(),
	}
	member2 := &model.ConversationMember{
		ConversationID: conversation.ID,
		UserID:         userID2,
		UnreadCount:    0,
		JoinedAt:       time.Now(),
	}

	if err := s.messageRepo.AddConversationMember(member1); err != nil {
		return nil, err
	}
	if err := s.messageRepo.AddConversationMember(member2); err != nil {
		return nil, err
	}

	// 加载会话信息
	conversation.Participants = []*model.ConversationMember{member1, member2}
	member1.User = user1
	member2.User = user2

	return conversation, nil
}

// CreateGroupConversation 创建群聊会话
func (s *MessageService) CreateGroupConversation(creatorID uint, name string, memberIDs []uint) (*model.Conversation, error) {
	// 创建群聊
	conversation := &model.Conversation{
		Type:      "group",
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.messageRepo.CreateConversation(conversation); err != nil {
		return nil, err
	}

	// 添加创建者
	creatorMember := &model.ConversationMember{
		ConversationID: conversation.ID,
		UserID:         creatorID,
		UnreadCount:    0,
		JoinedAt:       time.Now(),
	}
	if err := s.messageRepo.AddConversationMember(creatorMember); err != nil {
		return nil, err
	}

	// 添加其他成员
	for _, memberID := range memberIDs {
		if memberID == creatorID {
			continue // 跳过创建者
		}
		member := &model.ConversationMember{
			ConversationID: conversation.ID,
			UserID:         memberID,
			UnreadCount:    0,
			JoinedAt:       time.Now(),
		}
		if err := s.messageRepo.AddConversationMember(member); err != nil {
			return nil, err
		}
	}

	return conversation, nil
}

// MarkAsRead 标记消息为已读
func (s *MessageService) MarkAsRead(conversationID, userID uint) error {
	// 标记会话中所有消息为已读
	if err := s.messageRepo.MarkConversationAsRead(conversationID, userID); err != nil {
		return err
	}

	// 清空未读数
	return s.messageRepo.UpdateMemberUnreadCount(conversationID, userID, false)
}

// DeleteMessage 删除消息
func (s *MessageService) DeleteMessage(messageID, userID uint) error {
	// 获取消息
	message, err := s.messageRepo.GetMessageByID(messageID)
	if err != nil {
		return err
	}

	// 验证是否是发送者
	if message.SenderID != userID {
		return errors.New("permission denied")
	}

	// 软删除
	return s.messageRepo.DeleteMessage(messageID)
}

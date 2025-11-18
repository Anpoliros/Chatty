package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetConversations 获取会话列表
func GetConversations(c *gin.Context) {
	// TODO: 从上下文中获取当前用户ID
	// userID := c.GetUint("user_id")

	// TODO: 获取用户的会话列表
	// conversations, err := messageService.GetConversations(userID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Get conversations",
		// "conversations": conversations,
	})
}

// GetMessages 获取消息列表
func GetMessages(c *gin.Context) {
	conversationID := c.Param("conversationId")

	// 分页参数
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "50")

	// TODO: 获取消息列表
	// messages, err := messageService.GetMessages(conversationID, page, pageSize)

	c.JSON(http.StatusOK, gin.H{
		"message":         "Get messages",
		"conversation_id": conversationID,
		"page":            page,
		"page_size":       pageSize,
		// "messages": messages,
	})
}

// SendMessage 发送消息(HTTP方式,也可以通过WebSocket)
func SendMessage(c *gin.Context) {
	// TODO: 从上下文中获取当前用户ID
	// userID := c.GetUint("user_id")

	var req struct {
		ConversationID uint   `json:"conversation_id" binding:"required"`
		Type           string `json:"type" binding:"required"`
		Content        string `json:"content"`
		MediaURL       string `json:"media_url,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// TODO: 保存消息到数据库
	// message, err := messageService.SendMessage(userID, req)

	// TODO: 通过WebSocket推送消息

	c.JSON(http.StatusOK, gin.H{
		"message": "Message sent",
		// "data": message.ToResponse(),
	})
}

// MarkAsRead 标记消息为已读
func MarkAsRead(c *gin.Context) {
	// TODO: 从上下文中获取当前用户ID
	// userID := c.GetUint("user_id")

	var req struct {
		ConversationID uint `json:"conversation_id" binding:"required"`
		MessageID      uint `json:"message_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// TODO: 标记消息为已读
	// err := messageService.MarkAsRead(userID, req.ConversationID, req.MessageID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Marked as read",
	})
}

// DeleteMessage 删除消息
func DeleteMessage(c *gin.Context) {
	messageID := c.Param("id")

	// TODO: 从上下文中获取当前用户ID
	// userID := c.GetUint("user_id")

	// TODO: 删除消息(软删除)
	// err := messageService.DeleteMessage(userID, messageID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Message deleted",
		"id":      messageID,
	})
}

// CreateConversation 创建会话(单聊或群聊)
func CreateConversation(c *gin.Context) {
	// TODO: 从上下文中获取当前用户ID
	// userID := c.GetUint("user_id")

	var req struct {
		Type          string `json:"type" binding:"required"` // private or group
		ParticipantID uint   `json:"participant_id,omitempty"` // 私聊对方ID
		Name          string `json:"name,omitempty"`          // 群聊名称
		MemberIDs     []uint `json:"member_ids,omitempty"`    // 群聊成员ID列表
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// TODO: 创建会话
	// conversation, err := messageService.CreateConversation(userID, req)

	c.JSON(http.StatusOK, gin.H{
		"message": "Conversation created",
		// "conversation": conversation,
	})
}

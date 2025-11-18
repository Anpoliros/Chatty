package api

import (
	"net/http"
	"strconv"

	"github.com/Anpoliros/Chatty/server/internal/model"
	"github.com/Anpoliros/Chatty/server/internal/service"
	"github.com/Anpoliros/Chatty/server/internal/websocket"
	"github.com/gin-gonic/gin"
)

// MessageHandler 消息处理器
type MessageHandler struct {
	messageService *service.MessageService
	hub            *websocket.Hub
}

// NewMessageHandler 创建消息处理器
func NewMessageHandler(messageService *service.MessageService, hub *websocket.Hub) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		hub:            hub,
	}
}

// GetConversations 获取会话列表
// @Summary 获取会话列表
// @Description 获取当前用户的所有会话
// @Tags 会话
// @Security Bearer
// @Produce json
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/v1/conversations [get]
func (h *MessageHandler) GetConversations(c *gin.Context) {
	// 从上下文中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// 获取用户的会话列表
	conversations, err := h.messageService.GetConversations(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get conversations: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conversations": conversations,
	})
}

// GetMessages 获取消息列表
// @Summary 获取消息列表
// @Description 获取指定会话的消息历史
// @Tags 消息
// @Security Bearer
// @Produce json
// @Param conversationId path int true "会话ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(50)
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/v1/messages/{conversationId} [get]
func (h *MessageHandler) GetMessages(c *gin.Context) {
	conversationIDStr := c.Param("conversationId")
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid conversation ID",
		})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	// 获取消息列表
	messages, err := h.messageService.GetMessages(uint(conversationID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get messages: " + err.Error(),
		})
		return
	}

	// 转换为响应格式
	messageResponses := make([]interface{}, len(messages))
	for i, msg := range messages {
		messageResponses[i] = msg.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"messages":  messageResponses,
		"page":      page,
		"page_size": pageSize,
	})
}

// SendMessage 发送消息(HTTP方式,也可以通过WebSocket)
// @Summary 发送消息
// @Description 发送文本或媒体消息
// @Tags 消息
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body object true "消息内容"
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/v1/messages [post]
func (h *MessageHandler) SendMessage(c *gin.Context) {
	// 从上下文中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var req struct {
		ConversationID uint              `json:"conversation_id" binding:"required"`
		Type           model.MessageType `json:"type" binding:"required"`
		Content        string            `json:"content"`
		MediaURL       string            `json:"media_url,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// 保存消息到数据库
	message, err := h.messageService.SendMessage(
		userID.(uint),
		req.ConversationID,
		req.Type,
		req.Content,
		req.MediaURL,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send message: " + err.Error(),
		})
		return
	}

	// TODO: 通过WebSocket推送消息给其他用户
	// 这里需要找到会话的其他成员并推送消息

	c.JSON(http.StatusOK, gin.H{
		"message": "Message sent",
		"data":    message.ToResponse(),
	})
}

// MarkAsRead 标记消息为已读
// @Summary 标记消息为已读
// @Description 标记会话中的消息为已读
// @Tags 消息
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body object true "标记信息"
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/v1/messages/read [post]
func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	// 从上下文中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

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

	// 标记消息为已读
	if err := h.messageService.MarkAsRead(req.ConversationID, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to mark as read: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Marked as read",
	})
}

// DeleteMessage 删除消息
// @Summary 删除消息
// @Description 删除指定消息(软删除)
// @Tags 消息
// @Security Bearer
// @Produce json
// @Param id path int true "消息ID"
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 403 {object} map[string]interface{} "权限不足"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/v1/messages/{id} [delete]
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	messageIDStr := c.Param("id")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid message ID",
		})
		return
	}

	// 从上下文中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// 删除消息(软删除)
	if err := h.messageService.DeleteMessage(uint(messageID), userID.(uint)); err != nil {
		if err.Error() == "permission denied" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete message: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Message deleted",
		"id":      messageID,
	})
}

// CreateConversation 创建会话(单聊或群聊)
// @Summary 创建会话
// @Description 创建私聊或群聊会话
// @Tags 会话
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body object true "会话信息"
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/v1/conversations [post]
func (h *MessageHandler) CreateConversation(c *gin.Context) {
	// 从上下文中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

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

	var conversation *model.Conversation
	var err error

	// 根据类型创建不同的会话
	if req.Type == "private" {
		if req.ParticipantID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "participant_id is required for private conversation",
			})
			return
		}
		conversation, err = h.messageService.CreatePrivateConversation(userID.(uint), req.ParticipantID)
	} else if req.Type == "group" {
		if req.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "name is required for group conversation",
			})
			return
		}
		conversation, err = h.messageService.CreateGroupConversation(userID.(uint), req.Name, req.MemberIDs)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid conversation type",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create conversation: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Conversation created",
		"conversation": conversation,
	})
}

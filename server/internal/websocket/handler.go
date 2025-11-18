package websocket

import (
	"fmt"
	"net/http"

	"github.com/Anpoliros/Chatty/server/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	// 允许所有来源(生产环境应该限制)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var log = logger.New()

// HandleWebSocket 处理WebSocket连接
func HandleWebSocket(hub *Hub, c *gin.Context) {
	// TODO: 从token或query参数中获取用户ID
	// 这里暂时从query参数获取
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	// 转换用户ID(简化处理,生产环境应该验证token)
	var userID uint
	if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error("Failed to upgrade connection:", err)
		return
	}

	// 创建客户端
	client := NewClient(conn, hub, userID)

	// 注册到Hub
	hub.register <- client

	// 启动客户端读写循环
	client.Start()
}

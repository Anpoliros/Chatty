package api

import (
	"net/http"

	"github.com/Anpoliros/Chatty/server/internal/model"
	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req model.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// TODO: 实现注册逻辑
	// 1. 检查用户名/邮箱是否已存在
	// 2. 加密密码
	// 3. 创建用户
	// 4. 生成JWT token

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
		"user": gin.H{
			"username": req.Username,
			"email":    req.Email,
		},
		// "token": token,
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req model.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// TODO: 实现登录逻辑
	// 1. 查找用户
	// 2. 验证密码
	// 3. 生成JWT token
	// 4. 更新用户在线状态

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		// "token": token,
		// "user": userResponse,
	})
}

// Logout 用户登出
func Logout(c *gin.Context) {
	// TODO: 实现登出逻辑
	// 1. 获取用户ID
	// 2. 更新用户离线状态
	// 3. 清除token(如果使用Redis存储)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCurrentUser 获取当前用户信息
func GetCurrentUser(c *gin.Context) {
	// TODO: 从上下文中获取当前用户ID(通过JWT中间件设置)
	// userID := c.GetUint("user_id")

	// TODO: 查询用户信息
	// user, err := userService.GetUserByID(userID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Get current user",
		// "user": user.ToResponse(),
	})
}

// GetUserByID 根据ID获取用户信息
func GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	// TODO: 查询用户信息
	// user, err := userService.GetUserByID(userID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Get user by ID: " + userID,
		// "user": user.ToResponse(),
	})
}

// SearchUsers 搜索用户
func SearchUsers(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "keyword is required",
		})
		return
	}

	// TODO: 搜索用户
	// users, err := userService.SearchUsers(keyword)

	c.JSON(http.StatusOK, gin.H{
		"message": "Search users",
		"keyword": keyword,
		// "users": users,
	})
}

// UpdateProfile 更新用户资料
func UpdateProfile(c *gin.Context) {
	// TODO: 获取当前用户ID
	// userID := c.GetUint("user_id")

	var req struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// TODO: 更新用户信息
	// user, err := userService.UpdateProfile(userID, req)

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated",
		// "user": user.ToResponse(),
	})
}

package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"uniqueIndex;not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	Password     string    `json:"-" gorm:"not null"` // 密码不返回给前端
	Nickname     string    `json:"nickname"`
	Avatar       string    `json:"avatar"`
	Status       string    `json:"status" gorm:"default:'offline'"` // online, offline, away
	LastSeenAt   time.Time `json:"last_seen_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserRegisterRequest 注册请求
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

// UserLoginRequest 登录请求
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse 用户响应(不包含敏感信息)
type UserResponse struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Nickname   string    `json:"nickname"`
	Avatar     string    `json:"avatar"`
	Status     string    `json:"status"`
	LastSeenAt time.Time `json:"last_seen_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		Nickname:   u.Nickname,
		Avatar:     u.Avatar,
		Status:     u.Status,
		LastSeenAt: u.LastSeenAt,
		CreatedAt:  u.CreatedAt,
	}
}

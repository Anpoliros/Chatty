package repository

import (
	"context"
	"fmt"

	"github.com/Anpoliros/Chatty/server/pkg/utils"
	"github.com/redis/go-redis/v9"
)

// RedisClient Redis客户端管理器
type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

// NewRedis 创建新的Redis连接
func NewRedis(config utils.RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	ctx := context.Background()

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	return &RedisClient{
		Client: client,
		Ctx:    ctx,
	}, nil
}

// Close 关闭Redis连接
func (r *RedisClient) Close() error {
	return r.Client.Close()
}

// SetUserOnline 设置用户在线状态
func (r *RedisClient) SetUserOnline(userID uint) error {
	key := fmt.Sprintf("user:online:%d", userID)
	return r.Client.Set(r.Ctx, key, "1", 0).Err()
}

// SetUserOffline 设置用户离线状态
func (r *RedisClient) SetUserOffline(userID uint) error {
	key := fmt.Sprintf("user:online:%d", userID)
	return r.Client.Del(r.Ctx, key).Err()
}

// IsUserOnline 检查用户是否在线
func (r *RedisClient) IsUserOnline(userID uint) (bool, error) {
	key := fmt.Sprintf("user:online:%d", userID)
	result, err := r.Client.Exists(r.Ctx, key).Result()
	return result > 0, err
}

// GetOnlineUsers 获取所有在线用户
func (r *RedisClient) GetOnlineUsers() ([]string, error) {
	keys, err := r.Client.Keys(r.Ctx, "user:online:*").Result()
	return keys, err
}

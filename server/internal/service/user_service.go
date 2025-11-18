package service

import (
	"errors"
	"time"

	"github.com/Anpoliros/Chatty/server/internal/model"
	"github.com/Anpoliros/Chatty/server/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务
type UserService struct {
	userRepo    *repository.UserRepository
	redisClient *repository.RedisClient
}

// NewUserService 创建用户服务
func NewUserService(userRepo *repository.UserRepository, redisClient *repository.RedisClient) *UserService {
	return &UserService{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}

// Register 用户注册
func (s *UserService) Register(req *model.UserRegisterRequest) (*model.User, error) {
	// 检查用户是否已存在
	exists, err := s.userRepo.Exists(req.Username, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username or email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 创建用户
	user := &model.User{
		Username:   req.Username,
		Email:      req.Email,
		Password:   string(hashedPassword),
		Nickname:   req.Nickname,
		Status:     "offline",
		LastSeenAt: time.Now(),
	}

	if user.Nickname == "" {
		user.Nickname = req.Username
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(username, password string) (*model.User, error) {
	// 查找用户
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// 更新用户状态为在线
	user.Status = "online"
	user.LastSeenAt = time.Now()
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// 设置Redis在线状态
	if s.redisClient != nil {
		s.redisClient.SetUserOnline(user.ID)
	}

	return user, nil
}

// Logout 用户登出
func (s *UserService) Logout(userID uint) error {
	// 更新用户状态为离线
	if err := s.userRepo.UpdateStatus(userID, "offline"); err != nil {
		return err
	}

	// 设置Redis离线状态
	if s.redisClient != nil {
		return s.redisClient.SetUserOffline(userID)
	}

	return nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(userID uint, nickname, avatar string) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if nickname != "" {
		user.Nickname = nickname
	}
	if avatar != "" {
		user.Avatar = avatar
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// SearchUsers 搜索用户
func (s *UserService) SearchUsers(keyword string) ([]*model.User, error) {
	return s.userRepo.Search(keyword, 20)
}

// IsUserOnline 检查用户是否在线
func (s *UserService) IsUserOnline(userID uint) (bool, error) {
	if s.redisClient != nil {
		return s.redisClient.IsUserOnline(userID)
	}
	// fallback到数据库查询
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}
	return user.Status == "online", nil
}

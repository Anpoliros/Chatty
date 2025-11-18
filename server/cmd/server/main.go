package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Anpoliros/Chatty/server/internal/api"
	"github.com/Anpoliros/Chatty/server/internal/middleware"
	"github.com/Anpoliros/Chatty/server/internal/repository"
	"github.com/Anpoliros/Chatty/server/internal/service"
	"github.com/Anpoliros/Chatty/server/internal/websocket"
	"github.com/Anpoliros/Chatty/server/pkg/logger"
	"github.com/Anpoliros/Chatty/server/pkg/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志
	log := logger.New()
	log.Info("Starting Chatty Server...")

	// 加载配置
	cfg, err := utils.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化数据库连接
	db, err := repository.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	defer db.Close()

	// 自动迁移数据库
	if err := db.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Info("Database migrated successfully")

	// 初始化Redis连接
	redisClient, err := repository.NewRedis(cfg.Redis)
	if err != nil {
		log.Warn("Failed to connect redis (will continue without cache):", err)
		redisClient = nil // 可以在没有Redis的情况下运行
	} else {
		defer redisClient.Close()
		log.Info("Redis connected successfully")
	}

	// 创建Repository层
	userRepo := repository.NewUserRepository(db.DB)
	messageRepo := repository.NewMessageRepository(db.DB)

	// 创建Service层
	userService := service.NewUserService(userRepo, redisClient)
	messageService := service.NewMessageService(messageRepo, userRepo)

	// 创建JWT工具
	jwtUtil := utils.NewJWTUtil(cfg.JWT.Secret, cfg.JWT.ExpireTime)

	// 创建WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	// 设置Gin模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 设置Gin路由
	router := gin.Default()

	// 添加CORS中间件
	router.Use(middleware.CORSMiddleware())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"time":    time.Now().Unix(),
			"version": "v0.1-alpha",
		})
	})

	// 创建API处理器集合
	authHandler := api.NewAuthHandler(userService, jwtUtil)
	userHandler := api.NewUserHandler(userService)
	messageHandler := api.NewMessageHandler(messageService, hub)

	// API路由组
	v1 := router.Group("/api/v1")
	{
		// 认证路由(无需token)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// WebSocket连接(暂时不需要认证,用于测试)
		v1.GET("/ws", func(c *gin.Context) {
			websocket.HandleWebSocket(hub, c)
		})

		// 需要认证的路由
		authorized := v1.Group("/")
		authorized.Use(middleware.AuthMiddleware(jwtUtil))
		{
			// 用户相关
			authorized.POST("/auth/logout", authHandler.Logout)
			authorized.GET("/users/me", userHandler.GetCurrentUser)
			authorized.GET("/users/search", userHandler.SearchUsers)
			authorized.PUT("/users/profile", userHandler.UpdateProfile)

			// 会话相关
			authorized.GET("/conversations", messageHandler.GetConversations)
			authorized.POST("/conversations", messageHandler.CreateConversation)

			// 消息相关
			authorized.GET("/messages/:conversationId", messageHandler.GetMessages)
			authorized.POST("/messages", messageHandler.SendMessage)
			authorized.POST("/messages/read", messageHandler.MarkAsRead)
			authorized.DELETE("/messages/:id", messageHandler.DeleteMessage)
		}
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// 启动服务器
	go func() {
		log.Info("Server is running on :" + cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Info("Server exited")
}

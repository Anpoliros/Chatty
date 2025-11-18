package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Anpoliros/Chatty/server/internal/api"
	"github.com/Anpoliros/Chatty/server/internal/websocket"
	"github.com/Anpoliros/Chatty/server/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志
	log := logger.New()
	log.Info("Starting Chatty Server...")

	// TODO: 加载配置
	// cfg := config.Load()

	// TODO: 初始化数据库连接
	// db := repository.NewDatabase(cfg.Database)
	// defer db.Close()

	// TODO: 初始化Redis连接
	// redisClient := repository.NewRedis(cfg.Redis)
	// defer redisClient.Close()

	// 创建服务层
	// userService := service.NewUserService(db, redisClient)
	// messageService := service.NewMessageService(db, redisClient)

	// 创建WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	// 设置Gin路由
	router := gin.Default()

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	// API路由组
	v1 := router.Group("/api/v1")
	{
		// 认证路由
		auth := v1.Group("/auth")
		{
			auth.POST("/register", api.Register)
			auth.POST("/login", api.Login)
			auth.POST("/logout", api.Logout)
		}

		// WebSocket连接
		v1.GET("/ws", func(c *gin.Context) {
			websocket.HandleWebSocket(hub, c)
		})

		// 需要认证的路由
		authorized := v1.Group("/")
		// authorized.Use(middleware.AuthMiddleware())
		{
			authorized.GET("/users/me", api.GetCurrentUser)
			authorized.GET("/conversations", api.GetConversations)
			authorized.GET("/messages/:conversationId", api.GetMessages)
		}
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// 启动服务器
	go func() {
		log.Info("Server is running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server: ", err)
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
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Info("Server exited")
}

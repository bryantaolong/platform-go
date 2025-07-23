package main

import (
	"log"
	"net/http"
	"time"

	"github.com/bryantaolong/platform/internal/config"
	"github.com/bryantaolong/platform/internal/handler"
	"github.com/bryantaolong/platform/internal/middleware"
	"github.com/bryantaolong/platform/internal/repository"
	"github.com/bryantaolong/platform/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db := config.InitDB(cfg)

	// 初始化依赖
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg)
	authHandler := handler.NewAuthHandler(authService)

	// 初始化Gin
	r := gin.Default()

	// 添加CORS中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 设置路由
	setupRoutes(r, authHandler)

	// 启动服务器
	log.Printf("服务器启动在端口 %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}

func setupRoutes(r *gin.Engine, authHandler *handler.AuthHandler) {
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetCurrentUser)
		}
	}
}

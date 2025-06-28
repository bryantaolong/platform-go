package main

import (
	"github.com/bryantaolong/platform/internal/config"
	"github.com/bryantaolong/platform/internal/handler"
	"github.com/bryantaolong/platform/internal/middleware"
	"github.com/bryantaolong/platform/internal/repository"
	"github.com/bryantaolong/platform/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

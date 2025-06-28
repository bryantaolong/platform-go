package middleware

import (
	"github.com/bryantaolong/platform/internal/config"
	"github.com/bryantaolong/platform/internal/util"
	"github.com/bryantaolong/platform/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	cfg := config.Load()
	jwtUtil := util.NewJWTUtil(cfg.JWTSecret)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.Unauthorized("请求头中缺少 Authorization Token"))
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, response.Unauthorized("Token 格式不正确"))
			c.Abort()
			return
		}

		token := authHeader[7:] // 移除 "Bearer " 前缀

		claims, err := jwtUtil.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Unauthorized("Token 解析失败或无效"))
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("roles", claims.Roles)

		c.Next()
	}
}

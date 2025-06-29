package handler

import (
	"github.com/bryantaolong/platform/internal/model"
	"github.com/bryantaolong/platform/internal/service"
	"github.com/bryantaolong/platform/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("请求参数验证失败: "+err.Error()))
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success(user))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("请求参数验证失败: "+err.Error()))
		return
	}

	token, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Unauthorized(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(token))
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("未找到用户信息"))
		return
	}

	user, err := h.authService.GetUserByID(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "用户不存在"))
		return
	}

	c.JSON(http.StatusOK, response.Success(user))
}

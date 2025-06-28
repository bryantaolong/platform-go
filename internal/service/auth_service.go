package service

import (
	"errors"
	"github.com/bryantaolong/platform/internal/config"
	"github.com/bryantaolong/platform/internal/model"
	"github.com/bryantaolong/platform/internal/repository"
	"github.com/bryantaolong/platform/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
	jwtUtil  *util.JWTUtil
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtUtil:  util.NewJWTUtil(cfg.JWTSecret),
	}
}

func (s *AuthService) Register(req *model.RegisterRequest) (*model.User, error) {
	// 检查用户名是否已存在
	if s.userRepo.ExistsByUsername(req.Username) {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if req.Email != "" && s.userRepo.ExistsByEmail(req.Email) {
		return nil, errors.New("邮箱已被使用")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Status:   0,
		Roles:    "ROLE_USER",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("用户创建失败")
	}

	return user, nil
}

func (s *AuthService) Login(req *model.LoginRequest) (string, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return "", errors.New("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 0 {
		return "", errors.New("用户已被封禁")
	}

	// 生成JWT token
	token, err := s.jwtUtil.GenerateToken(user.ID, user.Roles)
	if err != nil {
		return "", errors.New("Token生成失败")
	}

	return token, nil
}

func (s *AuthService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

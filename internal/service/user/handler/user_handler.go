package handler

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"adult-short-videos/internal/pkg/cache"
	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/pkg/response"
	"adult-short-videos/internal/pkg/utils"
	"adult-short-videos/internal/service/user/logic"
	"adult-short-videos/internal/service/user/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserHandler 用户处理器
// 包含所有依赖项
type UserHandler struct {
	db        *gorm.DB
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExpire int64
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(db *gorm.DB, jwtSecret string, jwtExpire int64) *UserHandler {
	return &UserHandler{
		db:        db,
		userRepo:  repository.NewUserRepository(db),
		jwtSecret: jwtSecret,
		jwtExpire: jwtExpire,
	}
}

// Register 处理用户注册请求
// 路由: POST /api/user/register
func (h *UserHandler) Register(c *gin.Context) {
	// 1: 解析请求参数
	var req logic.RegisterReq

	// ShouldBindJSON 会自动解析 JSON 请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数解析失败
		response.HandleError(c, err)
		return
	}

	// 2: 参数验证
	// 验证用户名长度
	if len(req.Username) < 3 || len(req.Username) > 20 {
		response.Error(c, errors.CodeInvalidParam, "用户名长度必须在3-20个字符之间")
		return
	}

	// 验证密码长度
	if len(req.Password) < 8 || len(req.Password) > 50 {
		response.Error(c, errors.CodeInvalidParam, "密码长度必须在8-50个字符之间")
		return
	}

	// 验证邮箱格式（简单验证）
	if req.Email == "" {
		response.Error(c, errors.CodeInvalidParam, "邮箱不能为空")
		return
	}

	// 3: 调用业务逻辑
	// 创建业务逻辑实例
	l := logic.NewRegisterLogic(
		c.Request.Context(),
		h.userRepo,
		h.jwtSecret,
		h.jwtExpire,
	)

	// 执行注册
	resp, err := l.Register(&req)
	if err != nil {
		// 业务逻辑返回错误
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMsg(c, "注册成功", resp)
}

// Login 处理用户登录请求
// 路由: POST /api/user/login
func (h *UserHandler) Login(c *gin.Context) {
	// 1: 解析请求参数
	var req logic.LoginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeInvalidParam, "参数格式错误")
		return
	}

	// 2: 参数验证
	if req.Username == "" {
		response.Error(c, errors.CodeInvalidParam, "用户名不能为空")
		return
	}

	if req.Password == "" {
		response.Error(c, errors.CodeInvalidParam, "密码不能为空")
		return
	}

	// 3: 调用业务逻辑
	l := logic.NewLoginLogic(
		c.Request.Context(),
		h.userRepo,
		h.jwtSecret,
		h.jwtExpire,
		h.db,
	)

	resp, err := l.Login(&req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMsg(c, "登录成功", resp)
}

// GetUserInfo 获取用户信息
// 路由: GET /api/user/info
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	// 1: 从上下文获取用户 ID
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	// 2: 查询用户信息
	user, err := h.userRepo.FindByID(c.Request.Context(), userId.(int64))
	if err != nil {
		response.HandleError(c, err)
		return
	}

	// 3: 返回用户信息
	// 构建响应数据（不包含敏感信息）
	userInfo := map[string]interface{}{
		"user_id":       user.UserId,
		"username":      user.Username,
		"email":         user.Email,
		"avatar":        user.Avatar,
		"vip_level":     user.VipLevel,
		"vip_expire_at": user.VipExpireAt.Unix(),
		"created_at":    user.CreatedAt.Unix(),
	}

	response.Success(c, userInfo)
}

// Logout 处理退出登录
// 路由: POST /api/user/logout
// 将当前 Token 加入 Redis 黑名单，TTL = Token 剩余有效期
func (h *UserHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	var tokenString string
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	} else {
		tokenString = strings.TrimSpace(authHeader)
	}

	if tokenString == "" {
		response.Error(c, errors.CodeInvalidParam, "未提供令牌")
		return
	}

	claims, err := utils.ParseToken(tokenString, h.jwtSecret)
	if err != nil {
		response.Error(c, errors.CodeInvalidParam, "令牌无效")
		return
	}

	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining > 0 && cache.Client != nil {
		hash := sha256.Sum256([]byte(tokenString))
		key := fmt.Sprintf("blacklist:%x", hash)
		cache.Client.Set(c.Request.Context(), key, "1", remaining)
	}

	response.SuccessWithMsg(c, "已退出登录", nil)
}

package logic

import (
	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/pkg/logger"
	"adult-short-videos/internal/pkg/metrics"
	"adult-short-videos/internal/pkg/utils"
	"adult-short-videos/internal/pkg/validator"
	"adult-short-videos/internal/service/user/dto"
	"adult-short-videos/internal/service/user/model"
	"adult-short-videos/internal/service/user/repository"
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"adult-short-videos/internal/pkg/cache" // 导入 cache 包

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UserService 用户业务服务
type UserService struct {
	userRepo  repository.UserRepository
	db        *gorm.DB
	jwtSecret string
	jwtExpire int64
}

// NewUserService 创建 UserService 实例
func NewUserService(userRepo repository.UserRepository, db *gorm.DB, jwtSecret string, jwtExpire int64) *UserService {
	return &UserService{
		userRepo:  userRepo,
		db:        db,
		jwtSecret: jwtSecret,
		jwtExpire: jwtExpire,
	}
}

// Login 执行登录逻辑
func (s *UserService) Login(ctx context.Context, req *dto.LoginReq) (*dto.LoginResp, error) {
	// 1: 查找用户
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		metrics.UserLoginsTotal.WithLabelValues("false").Inc()
		return nil, errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	if user == nil {
		return nil, errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	// 2: 验证密码
	if !utils.PasswordVerify(user.Password, req.Password) {
		s.recordLoginLog(user.UserId, "password", false, "密码错误", req.ClientIP, req.UserAgent)
		metrics.UserLoginsTotal.WithLabelValues("false").Inc()
		return nil, errors.New(errors.CodePasswordError, "密码错误")
	}

	// 3: 检查账号状态
	if user.Status != 1 {
		metrics.UserLoginsTotal.WithLabelValues("false").Inc()
		return nil, errors.New(errors.CodeUserDisabled, "账号已被禁用")
	}

	// 4: 更新最后登录信息
	if err := s.userRepo.UpdateLastLogin(ctx, user.UserId, req.ClientIP); err != nil {
		// 非致命错误，忽略
	}

	// 5: 记录登录成功日志
	s.recordLoginLog(user.UserId, "password", true, "", req.ClientIP, req.UserAgent)
	metrics.UserLoginsTotal.WithLabelValues("true").Inc()

	// 6: 生成 JWT Token
	accessToken, err := utils.GenerateToken(user.UserId, user.Username, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, errors.New(errors.CodeTokenInvalid, "Token生成失败")
	}

	refreshToken, err := utils.GenerateToken(user.UserId, user.Username, s.jwtSecret, s.jwtExpire*7)
	if err != nil {
		return nil, errors.New(errors.CodeTokenInvalid, "Token生成失败")
	}

	return &dto.LoginResp{
		UserId:       user.UserId,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.jwtExpire,
	}, nil
}

// Register 执行注册逻辑
func (s *UserService) Register(ctx context.Context, req *dto.RegisterReq) (*dto.RegisterResp, error) {
	// 验证请求参数
	if err := validator.Validate(req); err != nil {
		return nil, errors.ErrInvalidParam
	}

	// 检查用户名是否已存在
	existingUser, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New(errors.CodeDatabaseError, "数据库查询失败")
	}
	if existingUser != nil {
		return nil, errors.New(errors.CodeUserExists, "用户已存在")
	}

	// 检查邮箱是否已被使用
	existingEmail, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New(errors.CodeDatabaseError, "数据库查询失败")
	}
	if existingEmail != nil {
		return nil, errors.New(errors.CodeEmailExists, "邮箱已被注册")
	}

	// 密码加密
	hashedPassword, err := utils.PasswordHash(req.Password)
	if err != nil {
		return nil, errors.New(errors.CodeEncryptError, "密码加密失败")
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Status:   1, // 默认正常状态
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.New(errors.CodeDatabaseError, "创建用户失败")
	}

	// 生成 JWT Token
	accessToken, err := utils.GenerateToken(user.UserId, user.Username, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, errors.New(errors.CodeServerError, "生成令牌失败")
	}
	refreshToken, err := utils.GenerateToken(user.UserId, user.Username, s.jwtSecret, s.jwtExpire*7)
	if err != nil {
		return nil, errors.New(errors.CodeServerError, "生成令牌失败")
	}

	logger.Info("用户注册成功", zap.Int64("user_id", user.UserId), zap.String("username", user.Username))
	metrics.UserRegistrationsTotal.Inc()

	return &dto.RegisterResp{
		UserId:       user.UserId,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.jwtExpire,
	}, nil
}

// GetInfo 获取用户信息
func (s *UserService) GetInfo(ctx context.Context, userId int64) (*dto.UserInfoResp, error) {
	user, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeUserNotFound, "用户不存在")
		}
		return nil, errors.New(errors.CodeDatabaseError, "数据库查询失败")
	}

	return &dto.UserInfoResp{
		UserId:   user.UserId,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Status:   user.Status,
	}, nil
}

// Logout 处理退出登录
func (s *UserService) Logout(ctx context.Context, tokenString string) error {
	claims, err := utils.ParseToken(tokenString, s.jwtSecret)
	if err != nil {
		return errors.New(errors.CodeInvalidParam, "令牌无效")
	}

	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining > 0 && cache.Client != nil {
		hash := sha256.Sum256([]byte(tokenString))
		key := fmt.Sprintf("blacklist:%x", hash)
		cache.Client.Set(ctx, key, "1", remaining)
	}
	return nil
}

// GetTopCreators 获取创作者榜单
func (s *UserService) GetTopCreators(ctx context.Context, limit int) ([]dto.CreatorItem, error) {
	if limit < 1 || limit > 50 {
		limit = 20
	}

	type row struct {
		UserId     int64
		Username   string
		Avatar     string
		VideoCount int64
		TotalPlays int64
	}

	var rows []row
	err := s.db.WithContext(ctx).
		Table("users u").
		Select("u.user_id, u.username, u.avatar, COUNT(v.video_id) AS video_count, COALESCE(SUM(v.play_count),0) AS total_plays").
		Joins("LEFT JOIN videos v ON v.user_id = u.user_id AND v.status = 1").
		Group("u.user_id, u.username, u.avatar").
		Having("COUNT(v.video_id) > 0").
		Order("(COUNT(v.video_id) * 10 + COALESCE(SUM(v.play_count),0) * 0.001) DESC").
		Limit(limit).
		Scan(&rows).Error

	if err != nil {
		return nil, errors.New(errors.CodeDatabaseError, "查询创作者榜单失败")
	}

	items := make([]dto.CreatorItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, dto.CreatorItem{
			UserId:     r.UserId,
			Username:   r.Username,
			Avatar:     r.Avatar,
			VideoCount: r.VideoCount,
			TotalPlays: r.TotalPlays,
			Score:      float64(r.VideoCount)*10 + float64(r.TotalPlays)*0.001,
		})
	}
	return items, nil
}

// UpdateProfile 更新用户昵称和个人简介，昵称 30 天内只能修改一次
func (s *UserService) UpdateProfile(ctx context.Context, userId int64, req *dto.UpdateProfileReq) (*dto.UpdateProfileResp, error) {
	user, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		return nil, errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	updates := map[string]interface{}{}

	if req.Username != "" && req.Username != user.Username {
		if user.UsernameChangedAt != nil && time.Since(*user.UsernameChangedAt) < 30*24*time.Hour {
			return nil, errors.New(errors.CodeInvalidParam, "昵称30天内只能修改一次")
		}
		existing, _ := s.userRepo.FindByUsername(ctx, req.Username)
		if existing != nil && existing.UserId != userId {
			return nil, errors.New(errors.CodeUserExists, "用户名已被使用")
		}
		updates["username"] = req.Username
		now := time.Now()
		updates["username_changed_at"] = now
		user.Username = req.Username
	}

	// bio 始终更新，允许置空
	updates["bio"] = req.Bio
	user.Bio = req.Bio

	if err := s.userRepo.UpdateByMap(ctx, userId, updates); err != nil {
		return nil, errors.New(errors.CodeDatabaseError, "更新失败")
	}

	return &dto.UpdateProfileResp{Username: user.Username, Bio: user.Bio}, nil
}

// UpdateAvatar 更新用户头像，30 天内只能修改一次
func (s *UserService) UpdateAvatar(ctx context.Context, userId int64, relativePath string) error {
	user, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		return errors.New(errors.CodeUserNotFound, "用户不存在")
	}
	if user.AvatarChangedAt != nil && time.Since(*user.AvatarChangedAt) < 30*24*time.Hour {
		return errors.New(errors.CodeInvalidParam, "头像30天内只能修改一次")
	}
	now := time.Now()
	if err := s.userRepo.UpdateByMap(ctx, userId, map[string]interface{}{
		"avatar":            relativePath,
		"avatar_changed_at": now,
	}); err != nil {
		return errors.New(errors.CodeDatabaseError, "更新头像失败")
	}
	return nil
}

func (s *UserService) recordLoginLog(userId int64, loginType string, success bool, failReason, clientIP, userAgent string) {
	status := int32(1)
	if !success {
		status = 0
	}
	loginLog := &model.LoginLog{
		UserId:     userId,
		LoginType:  loginType,
		LoginIp:    clientIP,
		UserAgent:  userAgent,
		Status:     status,
		FailReason: failReason,
	}
	go func() {
		s.db.Create(loginLog)
	}()
}

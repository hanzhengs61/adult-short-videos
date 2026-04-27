package logic

import (
	"adult-short-videos/internal/pkg/metrics"
	"adult-short-videos/internal/pkg/utils"
	"adult-short-videos/internal/service/user/model"
	"adult-short-videos/internal/service/user/repository"
	"context"
	"errors"

	"gorm.io/gorm"
)

// LoginLogic 登录业务逻辑
type LoginLogic struct {
	ctx       context.Context
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExpire int64
	db        *gorm.DB // 用于记录登录日志
}

// NewLoginLogic 创建登录逻辑实例
func NewLoginLogic(ctx context.Context, userRepo repository.UserRepository, jwtSecret string, jwtExpire int64, db *gorm.DB) *LoginLogic {
	return &LoginLogic{
		ctx:       ctx,
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpire: jwtExpire,
		db:        db,
	}
}

// LoginReq 登录请求参数
type LoginReq struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	ClientIP  string `json:"-"`
	UserAgent string `json:"-"`
}

// LoginResp 登录响应
type LoginResp struct {
	UserId       int64  `json:"user_id"`
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Login 执行登录逻辑
// 步骤：
// 1. 查找用户
// 2. 验证密码
// 3. 检查账号状态
// 4. 更新最后登录信息
// 5. 记录登录日志
// 6. 生成 JWT Token
func (l *LoginLogic) Login(req *LoginReq) (*LoginResp, error) {
	// 1: 查找用户
	user, err := l.userRepo.FindByUsername(l.ctx, req.Username)
	if err != nil {
		// 如果是"记录不存在"错误，说明用户名不存在
		if errors.Is(err, gorm.ErrRecordNotFound) {
			metrics.UserLoginsTotal.WithLabelValues("false").Inc()
			return nil, errors.New("用户名或密码错误")
		}
		return nil, errors.New("数据库查询失败")
	}

	// 2: 验证密码
	// 使用 bcrypt 比对密码
	if !utils.PasswordVerify(user.Password, req.Password) {
		// 记录登录失败日志
		l.recordLoginLog(user.UserId, "password", false, "密码错误", req.ClientIP, req.UserAgent)
		metrics.UserLoginsTotal.WithLabelValues("false").Inc()
		return nil, errors.New("用户名或密码错误")
	}

	// 3: 检查账号状态
	if user.Status != 1 {
		metrics.UserLoginsTotal.WithLabelValues("false").Inc()
		return nil, errors.New("账号已被禁用")
	}

	// 4: 更新最后登录信息
	if err := l.userRepo.UpdateLastLogin(l.ctx, user.UserId, req.ClientIP); err != nil {
		// 非致命错误，忽略
	}

	// 5: 记录登录成功日志
	l.recordLoginLog(user.UserId, "password", true, "", req.ClientIP, req.UserAgent)
	metrics.UserLoginsTotal.WithLabelValues("true").Inc()

	// 6: 生成 JWT Token
	accessToken, err := utils.GenerateToken(
		user.UserId,
		user.Username,
		l.jwtSecret,
		l.jwtExpire,
	)
	if err != nil {
		return nil, errors.New("生成访问令牌失败")
	}

	refreshToken, err := utils.GenerateToken(
		user.UserId,
		user.Username,
		l.jwtSecret,
		l.jwtExpire*7,
	)
	if err != nil {
		return nil, errors.New("生成刷新令牌失败")
	}

	return &LoginResp{
		UserId:       user.UserId,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    l.jwtExpire,
	}, nil
}

func (l *LoginLogic) recordLoginLog(userId int64, loginType string, success bool, failReason, clientIP, userAgent string) {
	status := int32(1)
	if !success {
		status = 0
	}

	// 创建登录日志记录
	loginLog := &model.LoginLog{
		UserId:     userId,
		LoginType:  loginType,
		LoginIp:    clientIP,
		UserAgent:  userAgent,
		Status:     status,
		FailReason: failReason,
	}

	// 异步记录日志，不影响登录主流程
	go func() {
		l.db.Create(loginLog)
	}()
}

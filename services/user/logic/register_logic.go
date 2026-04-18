package logic

import (
	"adult-short-videos/common/utils"
	"adult-short-videos/services/user/model"
	"adult-short-videos/services/user/repository"
	"context"
	"errors"

	"gorm.io/gorm"
)

// RegisterLogic 注册业务逻辑
type RegisterLogic struct {
	ctx       context.Context           // 请求上下文
	userRepo  repository.UserRepository // 用户仓储
	jwtSecret string                    // JWT 密钥
	jwtExpire int64                     // JWT 过期时间（秒）
}

// NewRegisterLogic 创建注册逻辑实例
func NewRegisterLogic(ctx context.Context, userRepo repository.UserRepository, jwtSecret string, jwtExpire int64) *RegisterLogic {
	return &RegisterLogic{
		ctx:       ctx,
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpire: jwtExpire,
	}
}

// RegisterReq 注册请求参数
type RegisterReq struct {
	Username   string `json:"username"`    // 用户名
	Password   string `json:"password"`    // 密码
	Email      string `json:"email"`       // 邮箱
	InviteCode string `json:"invite_code"` // 邀请码（可选）
}

// RegisterResp 注册响应
type RegisterResp struct {
	UserId       int64  `json:"user_id"`
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Register 执行注册逻辑
// 步骤：
// 1. 检查用户名是否已存在
// 2. 检查邮箱是否已被使用
// 3. 密码加密
// 4. 创建用户
// 5. 生成邀请码
// 6. 生成 JWT Token
func (l *RegisterLogic) Register(req *RegisterReq) (*RegisterResp, error) {
	// 1: 检查用户名是否已存在
	existingUser, err := l.userRepo.FindByUsername(l.ctx, req.Username)

	// 如果查询出错，但不是"记录不存在"错误，说明数据库有问题
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("数据库查询失败")
	}

	// 如果找到了用户，说明用户名已存在
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 2: 检查邮箱是否已被使用
	existingEmail, err := l.userRepo.FindByEmail(l.ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("数据库查询失败")
	}
	if existingEmail != nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 3: 密码加密
	hashedPassword, err := utils.PasswordHash(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 4: 创建用户
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		VipLevel: 0, // 默认普通用户
		Status:   1, // 默认正常状态
	}

	// 保存到数据库
	if err := l.userRepo.Create(l.ctx, user); err != nil {
		return nil, errors.New("创建用户失败")
	}

	// 5: 生成邀请码
	// 用户创建成功后，数据库会自动分配 UserId
	user.InviteCode = utils.GenerateInviterCode(user.UserId)

	// 更新邀请码到数据库
	if err := l.userRepo.Update(l.ctx, user); err != nil {
		// 这不是致命错误，记录日志即可
	}

	// 6: 生成 JWT Token
	// 生成 Access Token（用于日常请求）
	accessToken, err := utils.GenerateToken(
		user.UserId,
		user.Username,
		l.jwtSecret,
		l.jwtExpire,
	)
	if err != nil {
		return nil, errors.New("生成访问令牌失败")
	}

	// 生成 Refresh Token（用于刷新 Access Token）
	// Refresh Token 的有效期通常更长（这里是 7 天）
	refreshToken, err := utils.GenerateToken(
		user.UserId,
		user.Username,
		l.jwtSecret,
		l.jwtExpire*7, // 7 倍的 Access Token 时间
	)
	if err != nil {
		return nil, errors.New("生成刷新令牌失败")
	}

	return &RegisterResp{
		UserId:       user.UserId,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    l.jwtExpire,
	}, nil
}

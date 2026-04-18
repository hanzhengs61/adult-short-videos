package repository

import (
	"adult-short-videos/services/user/model"
	"context"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *model.User) error

	// FindByID 根据 ID 查找用户
	FindByID(ctx context.Context, userId int64) (*model.User, error)

	// FindByUsername 根据用户名查找用户
	FindByUsername(ctx context.Context, username string) (*model.User, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (*model.User, error)

	// Update 更新用户信息
	Update(ctx context.Context, user *model.User) error

	// UpdateLastLogin 更新最后登录信息
	UpdateLastLogin(ctx context.Context, userId int64, ip string) error
}

// userRepository 用户仓储实现
// 实现了 UserRepository 接口
type userRepository struct {
	db *gorm.DB // GORM 数据库连接
}

// NewUserRepository 创建用户仓储实例
// 这是一个工厂函数，返回接口类型
func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	// WithContext 支持上下文取消
	// Create 会自动设置 CreatedAt 和 UpdatedAt
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID 根据 ID 查找用户
func (r *userRepository) FindByID(ctx context.Context, userId int64) (*model.User, error) {
	var user model.User

	// 查询条件：user_id = ? AND status = 1
	// First 查询第一条记录，如果不存在返回 gorm.ErrRecordNotFound
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = 1", userId).
		First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("username = ? AND status = 1", username).
		First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("email = ? AND status = 1", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	// Updates 只更新非零值字段
	// UpdatedAt 会自动更新
	return r.db.WithContext(ctx).
		Model(user).
		Updates(user).Error
}

// UpdateLastLogin 更新最后登录信息
func (r *userRepository) UpdateLastLogin(ctx context.Context, userId int64, ip string) error {
	// 使用 map 可以更新零值
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("user_id = ?", userId).
		Updates(map[string]interface{}{
			"last_login_at": time.Now(),
			"last_login_ip": ip,
		}).Error
}

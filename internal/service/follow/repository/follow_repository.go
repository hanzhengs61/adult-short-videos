package repository

import (
	"adult-short-videos/internal/service/follow/model"
	"context"

	"gorm.io/gorm"
)

// AuthorProfile 关注列表中的作者信息，从 users 表实时 JOIN 获取
type AuthorProfile struct {
	Author string // 关注时记录的作者名
	Avatar string // users.avatar 相对路径，未注册用户为空
}

// FollowRepository 关注仓储接口
type FollowRepository interface {
	Create(ctx context.Context, f *model.AuthorFollow) error
	Delete(ctx context.Context, userID int64, author string) error
	Exists(ctx context.Context, userID int64, author string) (bool, error)
	ListAuthors(ctx context.Context, userID int64) ([]string, error)
	ListAuthorProfiles(ctx context.Context, userID int64) ([]AuthorProfile, error)
}

type followRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &followRepository{db: db}
}

func (r *followRepository) Create(ctx context.Context, f *model.AuthorFollow) error {
	return r.db.WithContext(ctx).Create(f).Error
}

func (r *followRepository) Delete(ctx context.Context, userID int64, author string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND author = ?", userID, author).
		Delete(&model.AuthorFollow{}).Error
}

func (r *followRepository) Exists(ctx context.Context, userID int64, author string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.AuthorFollow{}).
		Where("user_id = ? AND author = ?", userID, author).
		Count(&count).Error
	return count > 0, err
}

func (r *followRepository) ListAuthors(ctx context.Context, userID int64) ([]string, error) {
	var authors []string
	err := r.db.WithContext(ctx).
		Model(&model.AuthorFollow{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Pluck("author", &authors).Error
	return authors, err
}

// ListAuthorProfiles LEFT JOIN users 表，带回实时头像；未注册作者 Avatar 为空
func (r *followRepository) ListAuthorProfiles(ctx context.Context, userID int64) ([]AuthorProfile, error) {
	type row struct {
		Author string
		Avatar string
	}
	var rows []row
	err := r.db.WithContext(ctx).
		Table("author_follows af").
		Select("af.author, COALESCE(u.avatar, '') AS avatar").
		Joins("LEFT JOIN users u ON u.username = af.author AND u.status = 1").
		Where("af.user_id = ?", userID).
		Order("af.created_at DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	profiles := make([]AuthorProfile, 0, len(rows))
	for _, r := range rows {
		profiles = append(profiles, AuthorProfile{Author: r.Author, Avatar: r.Avatar})
	}
	return profiles, nil
}

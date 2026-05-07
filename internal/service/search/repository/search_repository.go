package repository

import (
	"context"
	"strings"

	videoModel "adult-short-videos/internal/service/video/model"

	"gorm.io/gorm"
)

// SearchRepository 搜索仓储接口
type SearchRepository interface {
	SearchVideos(ctx context.Context, keyword string, offset, limit int, filters map[string]interface{}) ([]*videoModel.VideoWithAuthor, int64, error)
}

type searchRepository struct {
	db *gorm.DB
}

// NewSearchRepository 创建搜索仓储实例
func NewSearchRepository(db *gorm.DB) SearchRepository {
	return &searchRepository{db: db}
}

// SearchVideos 按标题搜索视频，LEFT JOIN users 获取实时作者信息
func (r *searchRepository) SearchVideos(ctx context.Context, keyword string, offset, limit int, filters map[string]interface{}) ([]*videoModel.VideoWithAuthor, int64, error) {
	var total int64
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, 0, nil
	}
	like := "%" + keyword + "%"

	if err := r.db.WithContext(ctx).Model(&videoModel.Video{}).
		Where("status = 1 AND title ILIKE ?", like).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderCol := "v.created_at"
	if col, ok := filters["order_by"].(string); ok {
		switch col {
		case "play_count", "favorite_count":
			orderCol = "v." + col
		case "created_at":
			orderCol = "v.created_at"
		}
	}

	var videos []*videoModel.VideoWithAuthor
	err := r.db.WithContext(ctx).
		Table("videos v").
		Select(`v.*,
			COALESCE(u.username, '') AS author_name,
			COALESCE(u.avatar, '')   AS author_avatar`).
		Joins(`LEFT JOIN users u ON v.user_id > 0 AND u.user_id = v.user_id AND u.status = 1`).
		Where("v.status = 1 AND v.title ILIKE ?", like).
		Order(orderCol + " DESC").
		Offset(offset).
		Limit(limit).
		Scan(&videos).Error

	return videos, total, err
}

package repository

import (
	"context"
	"strings"

	videoModel "adult-short-videos/internal/service/video/model"

	"gorm.io/gorm"
)

// SearchRepository 搜索仓储接口
type SearchRepository interface {
	// SearchVideos 搜索视频
	SearchVideos(ctx context.Context, keyword string, offset, limit int) ([]*videoModel.Video, int64, error)
}

type searchRepository struct {
	db *gorm.DB
}

// NewSearchRepository 创建搜索仓储实例
func NewSearchRepository(db *gorm.DB) SearchRepository {
	return &searchRepository{db: db}
}

// SearchVideos 按标题或标签搜索视频
func (r *searchRepository) SearchVideos(ctx context.Context, keyword string, offset, limit int) ([]*videoModel.Video, int64, error) {
	var videos []*videoModel.Video
	var total int64

	// 关键词预处理
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return videos, 0, nil
	}

	// 模糊匹配关键词
	like := "%" + keyword + "%"

	query := r.db.WithContext(ctx).
		Model(&videoModel.Video{}).
		Where("status = 1").
		Where("title ILIKE ?", like)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("play_count DESC").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error

	return videos, total, err
}

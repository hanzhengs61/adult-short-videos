package repository

import (
	"adult-short-videos/internal/service/video/model"
	"context"

	"gorm.io/gorm"
)

// VideoRepository 视频仓储接口
type VideoRepository interface {
	// FindByID 根据 ID 查找视频
	FindByID(ctx context.Context, videoId int64) (*model.Video, error)
	// List 分页查询视频列表
	List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.Video, int64, error)
	// GetHotVideos 获取热门视频
	GetHotVideos(ctx context.Context, limit int) ([]*model.Video, error)
	// IncrementPlayCount 增加播放次数
	IncrementPlayCount(ctx context.Context, videoId int64) error
}

// VideoRepository 视频仓储实现
type videoRepository struct {
	db *gorm.DB
}

// NewVideoRepository 创建视频仓储实例
func NewVideoRepository(db *gorm.DB) *videoRepository {
	return &videoRepository{db: db}
}

// FindByID 根据 ID 查找视频
func (r *videoRepository) FindByID(ctx context.Context, videoId int64) (*model.Video, error) {
	var video model.Video
	// 查询条件：video_id = ? AND status = 1（只查正常状态的视频）
	err := r.db.WithContext(ctx).
		Select(`
			videos.*,
			(
				SELECT COUNT(*)
				FROM comments c
				WHERE c.video_id = videos.video_id
					AND c.parent_id = 0
					AND c.status = 1
			) AS comment_count_alias
		`). // 给子查询的 comment_count 添加别名
		Where("video_id = ? AND status = 1", videoId).
		First(&video).Error

	if err != nil {
		return nil, err
	}
	return &video, nil
}

// List 分页查询视频列表
func (r *videoRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.Video, int64, error) {
	var videos []*model.Video
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Video{}).Where("status = 1")

	orderCol := "created_at"
	if col, ok := filters["order_by"].(string); ok {
		switch col {
		case "play_count", "favorite_count", "created_at":
			orderCol = col
		case "comment_count":
			orderCol = "comment_count_alias"
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Select(`
			videos.*,
			(
				SELECT COUNT(*)
				FROM comments c
				WHERE c.video_id = videos.video_id
					AND c.parent_id = 0
					AND c.status = 1
			) AS comment_count_alias
		`).
		Order(orderCol + " DESC").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error

	return videos, total, err
}

func (r *videoRepository) GetHotVideos(ctx context.Context, limit int) ([]*model.Video, error) {
	var videos []*model.Video

	err := r.db.WithContext(ctx).
		Table("videos v").
		Select(`
			v.*,
			(
				SELECT COUNT(*)
				FROM comments c
				WHERE c.video_id = v.video_id
					AND c.parent_id = 0
					AND c.status = 1
			) AS comment_count_alias
		`).
		Joins("INNER JOIN video_heat_stats h ON v.video_id = h.video_id").
		Where("h.is_hot = ? AND v.status = 1", true).
		Order("h.heat_score DESC").
		Limit(limit).
		Find(&videos).Error

	return videos, err
}

func (r *videoRepository) IncrementPlayCount(ctx context.Context, videoId int64) error {
	return r.db.WithContext(ctx).
		Model(&model.Video{}).
		Where("video_id = ?", videoId).
		UpdateColumn("play_count", gorm.Expr("play_count + ?", 1)).Error
}

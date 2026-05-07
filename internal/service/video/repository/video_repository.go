package repository

import (
	"adult-short-videos/internal/service/video/model"
	"context"

	"gorm.io/gorm"
)

// authorSelect 所有视频查询统一的 SELECT 片段：
// user_id > 0 时实时从 users 表取昵称和头像
const authorSelect = `
	v.*,
	COALESCE(u.username, '') AS author_name,
	COALESCE(u.avatar, '')   AS author_avatar,
	(
		SELECT COUNT(*)
		FROM comments c
		WHERE c.video_id = v.video_id
			AND c.parent_id = 0
			AND c.status = 1
	) AS comment_count`

const authorJoin = `LEFT JOIN users u ON v.user_id > 0 AND u.user_id = v.user_id AND u.status = 1`

// VideoRepository 视频仓储接口
type VideoRepository interface {
	FindByID(ctx context.Context, videoId int64) (*model.VideoWithAuthor, error)
	List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.VideoWithAuthor, int64, error)
	GetHotVideos(ctx context.Context, limit int) ([]*model.VideoWithAuthor, error)
	IncrementPlayCount(ctx context.Context, videoId int64) error
}

type videoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *videoRepository {
	return &videoRepository{db: db}
}

func (r *videoRepository) FindByID(ctx context.Context, videoId int64) (*model.VideoWithAuthor, error) {
	var v model.VideoWithAuthor
	err := r.db.WithContext(ctx).
		Table("videos v").
		Select(authorSelect).
		Joins(authorJoin).
		Where("v.video_id = ? AND v.status = 1", videoId).
		Scan(&v).Error
	if err != nil {
		return nil, err
	}
	if v.VideoId == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &v, nil
}

func (r *videoRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.VideoWithAuthor, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Video{}).Where("status = 1").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderCol := "v.created_at"
	if col, ok := filters["order_by"].(string); ok {
		switch col {
		case "play_count", "favorite_count":
			orderCol = "v." + col
		case "comment_count":
			orderCol = "comment_count"
		case "created_at":
			orderCol = "v.created_at"
		}
	}

	var videos []*model.VideoWithAuthor
	err := r.db.WithContext(ctx).
		Table("videos v").
		Select(authorSelect).
		Joins(authorJoin).
		Where("v.status = 1").
		Order(orderCol + " DESC").
		Offset(offset).
		Limit(limit).
		Scan(&videos).Error

	return videos, total, err
}

func (r *videoRepository) GetHotVideos(ctx context.Context, limit int) ([]*model.VideoWithAuthor, error) {
	var videos []*model.VideoWithAuthor
	err := r.db.WithContext(ctx).
		Table("videos v").
		Select(authorSelect).
		Joins(authorJoin).
		Joins("INNER JOIN video_heat_stats h ON v.video_id = h.video_id").
		Where("h.is_hot = ? AND v.status = 1", true).
		Order("h.heat_score DESC").
		Limit(limit).
		Scan(&videos).Error
	return videos, err
}

func (r *videoRepository) IncrementPlayCount(ctx context.Context, videoId int64) error {
	return r.db.WithContext(ctx).
		Model(&model.Video{}).
		Where("video_id = ?", videoId).
		UpdateColumn("play_count", gorm.Expr("play_count + ?", 1)).Error
}

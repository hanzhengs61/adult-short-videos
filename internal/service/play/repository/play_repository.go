package repository

import (
	"adult-short-videos/internal/service/play/model"
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PlayHistoryRepository 播放历史仓储接口
type PlayHistoryRepository interface {
	// CreateOrUpdate 创建或更新播放历史（同一用户同一视频只保留一条记录）
	CreateOrUpdate(ctx context.Context, history *model.PlayHistory) error

	// ListWithVideo 获取用户播放历史列表
	ListWithVideo(ctx context.Context, userId int64, offset, limit int) ([]*model.PlayHistoryWithVideo, int64, error)

	// Delete 删除播放历史
	Delete(ctx context.Context, userId, videoId int64) error

	// DeleteAll 清空用户所有播放历史
	DeleteAll(ctx context.Context, userId int64) error
}

type playHistoryRepository struct {
	db *gorm.DB
}

func NewPlayHistoryRepository(db *gorm.DB) *playHistoryRepository {
	return &playHistoryRepository{db: db}
}

// CreateOrUpdate 创建或更新播放历史
// 使用 ON CONFLICT 实现 Upsert（如果存在则更新）
func (r *playHistoryRepository) CreateOrUpdate(ctx context.Context, history *model.PlayHistory) error {
	// 设置最后播放时间
	history.LastPlayAt = time.Now()

	// 使用 Clauses 实现 Upsert
	// 如果 user_id + video_id 组合已存在，则更新
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}, {Name: "video_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"play_duration",
				"play_progress",
				"last_play_at",
				"updated_at",
			}),
		}).
		Create(history).Error
}

// ListWithVideo 获取用户播放历史列表
func (r *playHistoryRepository) ListWithVideo(ctx context.Context, userId int64, offset, limit int) ([]*model.PlayHistoryWithVideo, int64, error) {
	var results []*model.PlayHistoryWithVideo
	var total int64

	// 获取总数
	err := r.db.WithContext(ctx).
		Model(&model.PlayHistory{}).
		Where("user_id = ?", userId).
		Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	// 联表查询
	err = r.db.WithContext(ctx).
		Table("play_history h").
		Select(`
			h.id as history_id,
			v.video_id,
			v.title,
			v.cover_url,
			v.preview_url,
			v.duration,
			h.play_duration,
			h.play_progress,
			v.fanhao,
			v.region,
			v.category,
			EXTRACT(EPOCH FROM h.last_play_at)::bigint as last_play_at
		`).
		Joins("INNER JOIN videos v ON h.video_id = v.video_id").
		Where("h.user_id = ? AND v.status = 1", userId).
		Order("h.last_play_at DESC").
		Offset(offset).
		Limit(limit).
		Scan(&results).Error

	return results, total, err
}

// Delete 删除播放历史
func (r *playHistoryRepository) Delete(ctx context.Context, userId, videoId int64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND video_id = ?", userId, videoId).
		Delete(&model.PlayHistory{}).Error
}

// DeleteAll 清空用户所有播放历史
func (r *playHistoryRepository) DeleteAll(ctx context.Context, userId int64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Delete(&model.PlayHistory{}).Error
}

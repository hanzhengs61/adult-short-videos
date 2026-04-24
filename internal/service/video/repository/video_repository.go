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
	// offset: 偏移量, limit: 每页数量, filters: 筛选条件
	List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.Video, int64, error)

	// GetHotVideos 获取热门视频
	GetHotVideos(ctx context.Context, limit int) ([]*model.Video, error)

	// IncrementPlayCount 增加播放次数
	IncrementPlayCount(ctx context.Context, videoId int64) error

	// GetVideoActors 查询视频的演员列表
	GetVideoActors(ctx context.Context, videoId int64) ([]*model.Actor, error)
}

// VideoRepository 视频仓储实现
type videoRepository struct {
	db *gorm.DB
}

// NewVideoRepository 创建视频仓储实例
func NewVideoRepository(db *gorm.DB) *videoRepository {
	return &videoRepository{
		db: db,
	}
}

// FindByID 根据 ID 查找视频
func (r *videoRepository) FindByID(ctx context.Context, videoId int64) (*model.Video, error) {
	var video model.Video

	// 查询条件：video_id = ? AND status = 1（只查正常状态的视频）
	err := r.db.WithContext(ctx).
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

	// 构建查询
	query := r.db.WithContext(ctx).Model(&model.Video{}).Where("status = 1")

	// filters 是一个 map，可以包含：region, category, is_vip_only 等
	for key, value := range filters {
		switch key {
		case "region":
			// 按地区筛选：日本、欧美、国产等
			query = query.Where("region = ?", value)
		case "category":
			// 按分类筛选：职场、制服、学生等
			query = query.Where("category = ?", value)
		case "is_vip_only":
			// VIP 专属筛选
			query = query.Where("is_vip_only = ?", value)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按创建时间倒序排列，最新的在前面
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error

	return videos, total, err
}

// GetHotVideos 获取热门视频
func (r *videoRepository) GetHotVideos(ctx context.Context, limit int) ([]*model.Video, error) {
	var videos []*model.Video

	// 联表查询：video 视频表 + video_heat_stats 视频热度统计表
	// 只查询 is_hot = true 的视频，按热度分数排序
	err := r.db.WithContext(ctx).
		Table("videos v").
		Select("v.*").
		Joins("INNER JOIN video_heat_stats h ON v.video_id = h.video_id").
		Where("h.is_hot = ? AND v.status = 1", true).
		Order("h.heat_score DESC").
		Limit(limit).
		Find(&videos).Error

	return videos, err
}

// IncrementPlayCount 增加播放次数
func (r *videoRepository) IncrementPlayCount(ctx context.Context, videoId int64) error {
	// 使用 SQL 表达式：play_count = play_count + 1
	return r.db.WithContext(ctx).
		Model(&model.Video{}).
		Where("video_id = ?", videoId).
		UpdateColumn("play_count", gorm.Expr("play_count + ?", 1)).Error
}

// GetVideoActors 查询视频的演员列表
func (r *videoRepository) GetVideoActors(ctx context.Context, videoId int64) ([]*model.Actor, error) {
	var actors []*model.Actor

	// 联表查询：actors 演员表 + video_actors 视频-演员关联表
	err := r.db.WithContext(ctx).
		Table("actors a").
		Select("a.*").
		Joins("INNER JOIN video_actors va ON a.actor_id = va.actor_id").
		Where("va.video_id = ?", videoId).
		Find(&actors).Error

	return actors, err
}

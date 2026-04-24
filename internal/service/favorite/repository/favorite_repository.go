package repository

import (
	"adult-short-videos/internal/service/favorite/model"
	"context"

	"gorm.io/gorm"
)

// FavoriteRepository 收藏仓储接口
type FavoriteRepository interface {
	// Create 添加收藏
	// 返回：收藏记录ID, 错误
	Create(ctx context.Context, favorite *model.Favorite) error

	// Delete 删除收藏
	// 参数：用户ID, 视频ID
	Delete(ctx context.Context, userId, videoId int64) error

	// Exists 检查是否已收藏
	// 返回：true=已收藏, false=未收藏
	Exists(ctx context.Context, userId, videoId int64) (bool, error)

	// ListWithVideo 获取用户的收藏列表（包含视频信息）
	// 参数：用户ID, offset（偏移量）, limit（数量）
	// 返回：收藏列表, 总数, 错误
	ListWithVideo(ctx context.Context, userId int64, offset, limit int) ([]*model.FavoriteWithVideo, int64, error)

	// GetUserFavoriteVideoIds 获取用户收藏的视频ID列表
	// 用于批量查询时判断哪些视频已被收藏
	GetUserFavoriteVideoIds(ctx context.Context, userId int64) ([]int64, error)
}

// favoriteRepository 收藏仓储实现
type favoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository 创建收藏仓储实例
func NewFavoriteRepository(db *gorm.DB) *favoriteRepository {
	return &favoriteRepository{
		db: db,
	}
}

// Create 添加收藏
func (r *favoriteRepository) Create(ctx context.Context, favorite *model.Favorite) error {
	// 使用 Create 会自动设置 created_at
	return r.db.WithContext(ctx).Create(favorite).Error
}

// Delete 删除收藏
func (r *favoriteRepository) Delete(ctx context.Context, userId, videoId int64) error {
	// 删除条件：user_id = ? AND video_id = ?
	return r.db.WithContext(ctx).
		Where("user_id = ? AND video_id = ?", userId, videoId).
		Delete(&model.Favorite{}).Error
}

// Exists 检查是否已收藏
func (r *favoriteRepository) Exists(ctx context.Context, userId, videoId int64) (bool, error) {
	var count int64

	// 查询符合条件的记录数
	err := r.db.WithContext(ctx).
		Model(&model.Favorite{}).
		Where("user_id = ? AND video_id = ?", userId, videoId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	// count > 0 说明已收藏
	return count > 0, nil
}

// ListWithVideo 获取用户的收藏列表（包含视频信息）
func (r *favoriteRepository) ListWithVideo(ctx context.Context, userId int64, offset, limit int) ([]*model.FavoriteWithVideo, int64, error) {
	var results []*model.FavoriteWithVideo
	var total int64

	// 获取收藏总数
	err := r.db.WithContext(ctx).
		Model(&model.Favorite{}).
		Where("user_id = ?", userId).
		Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	// 联表：favorites 收藏表 + videos 视频表
	// 查询收藏记录，并获取对应的视频信息
	err = r.db.WithContext(ctx).
		Table("favorites f").
		Select(`
			f.id as favorite_id,
			v.video_id,
			v.title,
			v.cover_url,
			v.preview_url,
			v.duration,
			v.play_count,
			v.fanhao,
			v.region,
			v.category,
			v.is_vip_only,
			EXTRACT(EPOCH FROM f.created_at)::bigint as favorited_at
		`).
		Joins("INNER JOIN videos v ON f.video_id = v.video_id").
		Where("f.user_id = ? AND v.status = 1", userId). // 只查询正常状态的视频
		Order("f.created_at DESC").                      // 按收藏时间倒序
		Offset(offset).
		Limit(limit).
		Scan(&results).Error

	return results, total, err
}

// GetUserFavoriteVideoIds 获取用户收藏的视频ID列表
func (r *favoriteRepository) GetUserFavoriteVideoIds(ctx context.Context, userId int64) ([]int64, error) {
	var videoIds []int64

	// 只查询 video_id 字段
	err := r.db.WithContext(ctx).
		Model(&model.Favorite{}).
		Where("user_id = ?", userId).
		Pluck("video_id", &videoIds).Error

	return videoIds, err
}

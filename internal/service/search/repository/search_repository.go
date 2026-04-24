package repository

import (
	"context"
	"strings"

	videoModel "adult-short-videos/internal/service/video/model"

	"gorm.io/gorm"
)

// SearchRepository 搜索仓储接口
type SearchRepository interface {
	// SearchVideos 搜索视频（支持标题、番号、演员名）
	SearchVideos(ctx context.Context, keyword string, offset, limit int) ([]*videoModel.Video, int64, error)

	// SearchActors 搜索演员
	SearchActors(ctx context.Context, keyword string, offset, limit int) ([]*videoModel.Actor, int64, error)

	// SearchByFanhao 按番号精确查找
	SearchByFanhao(ctx context.Context, fanhao string) (*videoModel.Video, error)

	// AdvancedSearch 高级搜索（支持多条件组合）
	AdvancedSearch(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*videoModel.Video, int64, error)
}

type searchRepository struct {
	db *gorm.DB
}

// NewSearchRepository 创建搜索仓储实例
func NewSearchRepository(db *gorm.DB) SearchRepository {
	return &searchRepository{db: db}
}

// SearchVideos 搜索视频
// 支持：标题、番号、演员名
func (r *searchRepository) SearchVideos(ctx context.Context, keyword string, offset, limit int) ([]*videoModel.Video, int64, error) {
	var videos []*videoModel.Video
	var total int64

	// 关键词预处理
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return videos, 0, nil
	}

	// 模糊匹配关键词
	likeKeyword := "%" + keyword + "%"

	// 搜索条件：标题包含关键词 OR 番号包含关键词 OR 演员名包含关键词
	query := r.db.WithContext(ctx).
		Model(&videoModel.Video{}).
		Where("status = 1").
		Where(r.db.Where("title ILIKE ?", likeKeyword).
			Or("fanhao ILIKE ?", likeKeyword).
			Or(`video_id IN (
				SELECT va.video_id 
				FROM video_actors va 
				INNER JOIN actors a ON va.actor_id = a.actor_id 
				WHERE a.actor_name ILIKE ?
			)`, likeKeyword))

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 按相关度排序：番号完全匹配 > 标题匹配 > 播放次数
	err := query.
		Order(gorm.Expr(
			"CASE WHEN fanhao ILIKE ? THEN 1 WHEN title ILIKE ? THEN 2 ELSE 3 END, play_count DESC",
			likeKeyword, likeKeyword,
		)).
		Offset(offset).
		Limit(limit).
		Find(&videos).Error

	return videos, total, err
}

// SearchActors 搜索演员
func (r *searchRepository) SearchActors(ctx context.Context, keyword string, offset, limit int) ([]*videoModel.Actor, int64, error) {
	var actors []*videoModel.Actor
	var total int64

	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return actors, 0, nil
	}

	likeKeyword := "%" + keyword + "%"

	// 按演员名模糊搜索
	query := r.db.WithContext(ctx).
		Model(&videoModel.Actor{}).
		Where("actor_name ILIKE ?", likeKeyword)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按粉丝数排序
	err := query.
		Order("follower_count DESC").
		Offset(offset).
		Limit(limit).
		Find(&actors).Error

	return actors, total, err
}

// SearchByFanhao 按番号精确查找
func (r *searchRepository) SearchByFanhao(ctx context.Context, fanhao string) (*videoModel.Video, error) {
	var video videoModel.Video

	// 番号转大写，精确匹配
	fanhao = strings.ToUpper(strings.TrimSpace(fanhao))

	err := r.db.WithContext(ctx).
		Where("UPPER(fanhao) = ? AND status = 1", fanhao).
		First(&video).Error

	if err != nil {
		return nil, err
	}

	return &video, nil
}

// AdvancedSearch 高级搜索
// 支持组合条件：关键词 + 地区 + 分类 + VIP筛选 + 时长筛选
func (r *searchRepository) AdvancedSearch(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*videoModel.Video, int64, error) {
	var videos []*videoModel.Video
	var total int64

	// 构建基础查询
	query := r.db.WithContext(ctx).Model(&videoModel.Video{}).Where("status = 1")

	// 1. 关键词搜索
	if keyword, ok := filters["keyword"].(string); ok && keyword != "" {
		likeKeyword := "%" + keyword + "%"
		query = query.Where(r.db.Where("title ILIKE ?", likeKeyword).
			Or("fanhao ILIKE ?", likeKeyword))
	}

	// 2. 地区筛选
	if region, ok := filters["region"].(string); ok && region != "" {
		query = query.Where("region = ?", region)
	}

	// 3. 分类筛选
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}

	// 4. VIP筛选
	if isVipOnly, ok := filters["is_vip_only"].(bool); ok {
		query = query.Where("is_vip_only = ?", isVipOnly)
	}

	// 5. 分辨率筛选
	if resolution, ok := filters["resolution"].(string); ok && resolution != "" {
		query = query.Where("resolution = ?", resolution)
	}

	// 6. 时长筛选（分钟）
	if minDuration, ok := filters["min_duration"].(int); ok && minDuration > 0 {
		query = query.Where("duration >= ?", minDuration*60)
	}
	if maxDuration, ok := filters["max_duration"].(int); ok && maxDuration > 0 {
		query = query.Where("duration <= ?", maxDuration*60)
	}

	// 7. 演员筛选
	if actorId, ok := filters["actor_id"].(int64); ok && actorId > 0 {
		query = query.Where(`video_id IN (
			SELECT video_id FROM video_actors WHERE actor_id = ?
		)`, actorId)
	}

	// ========== 排序 ==========
	sortBy := "created_at DESC" // 默认按创建时间倒序
	if sort, ok := filters["sort"].(string); ok {
		switch sort {
		case "hot":
			sortBy = "play_count DESC"
		case "favorite":
			sortBy = "favorite_count DESC"
		case "newest":
			sortBy = "created_at DESC"
		case "oldest":
			sortBy = "created_at ASC"
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order(sortBy).
		Offset(offset).
		Limit(limit).
		Find(&videos).Error

	return videos, total, err
}

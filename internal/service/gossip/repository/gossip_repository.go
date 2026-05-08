package repository

import (
	"adult-short-videos/internal/service/gossip/model"
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// GossipRepository 吃瓜文章仓储接口
type GossipRepository interface {
	List(ctx context.Context, offset, limit int, tag string) ([]*model.GossipPost, int64, error)
	FindByID(ctx context.Context, id int64) (*model.GossipPost, error)
	IncrementViewCount(ctx context.Context, id int64) error
	// FindByTags 用于详情页关联视频查询时取同类标签的最新文章（目前未使用，预留）
	FindByTags(ctx context.Context, tags []string, excludeID int64, limit int) ([]*model.GossipPost, error)
}

type gossipRepository struct {
	db *gorm.DB
}

func NewGossipRepository(db *gorm.DB) GossipRepository {
	return &gossipRepository{db: db}
}

func (r *gossipRepository) List(ctx context.Context, offset, limit int, tag string) ([]*model.GossipPost, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.GossipPost{}).Where("status = 1")
	if tag != "" {
		// tags 字段为逗号分隔字符串，用 LIKE 过滤
		q = q.Where("tags LIKE ?", "%"+tag+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var posts []*model.GossipPost
	err := q.Order("published_at DESC").Offset(offset).Limit(limit).Find(&posts).Error
	return posts, total, err
}

func (r *gossipRepository) FindByID(ctx context.Context, id int64) (*model.GossipPost, error) {
	var post model.GossipPost
	err := r.db.WithContext(ctx).Where("id = ? AND status = 1", id).First(&post).Error
	return &post, err
}

func (r *gossipRepository) IncrementViewCount(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&model.GossipPost{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (r *gossipRepository) FindByTags(ctx context.Context, tags []string, excludeID int64, limit int) ([]*model.GossipPost, error) {
	if len(tags) == 0 {
		return nil, nil
	}
	// 任意一个标签匹配即可
	conditions := make([]string, len(tags))
	args := make([]interface{}, len(tags))
	for i, t := range tags {
		conditions[i] = "tags LIKE ?"
		args[i] = fmt.Sprintf("%%%s%%", t)
	}
	args = append(args, excludeID)

	var posts []*model.GossipPost
	err := r.db.WithContext(ctx).
		Where("status = 1 AND ("+strings.Join(conditions, " OR ")+") AND id != ?", args...).
		Order("published_at DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

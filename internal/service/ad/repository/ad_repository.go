package repository

import (
	"context"

	"adult-short-videos/internal/service/ad/model"

	"gorm.io/gorm"
)

type AdRepository interface {
	ListBanners(ctx context.Context) ([]*model.AdBanner, error)
	ListApps(ctx context.Context, category string) ([]*model.AdApp, error)
}

type adRepository struct {
	db *gorm.DB
}

func NewAdRepository(db *gorm.DB) AdRepository {
	return &adRepository{db: db}
}

func (r *adRepository) ListBanners(ctx context.Context) ([]*model.AdBanner, error) {
	var list []*model.AdBanner
	err := r.db.WithContext(ctx).
		Where("status = 1").
		Order("sort DESC, id ASC").
		Find(&list).Error
	return list, err
}

func (r *adRepository) ListApps(ctx context.Context, category string) ([]*model.AdApp, error) {
	q := r.db.WithContext(ctx).Where("status = 1")
	if category != "" && category != "all" {
		// categories 字段为逗号分隔，用 LIKE 匹配分类标识
		q = q.Where("categories LIKE ?", "%"+category+"%")
	}
	var list []*model.AdApp
	err := q.Order("sort DESC, id ASC").Find(&list).Error
	return list, err
}

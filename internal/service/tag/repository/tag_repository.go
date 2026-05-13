package repository

import (
	"adult-short-videos/internal/service/tag/model"
	"strings"

	"gorm.io/gorm"
)

// TagRepository handles database operations for Tag model.
type TagRepository struct {
	db *gorm.DB
}

// NewTagRepository creates a new TagRepository.
func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

// Create saves a new tag to the database.
func (r *TagRepository) Create(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

// FindByName finds a tag by its name.
func (r *TagRepository) FindByName(name string) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	return &tag, err
}

// GetAllActive 获取所有启用状态的标签
func (r *TagRepository) GetAllActive() ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Where("status = 1").Order("tag_id").Find(&tags).Error
	return tags, err
}

// GetHot 按点击数降序返回前 limit 个启用标签
func (r *TagRepository) GetHot(limit int) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Where("status = 1").Order("click_count DESC, tag_id ASC").Limit(limit).Find(&tags).Error
	return tags, err
}

// IncrClickCount 原子自增指定标签的点击数（标签不存在时静默忽略）
func (r *TagRepository) IncrClickCount(name string) error {
	return r.db.Model(&model.Tag{}).
		Where("name = ? AND status = 1", name).
		UpdateColumn("click_count", gorm.Expr("click_count + 1")).Error
}

// MatchTags 用标签库对标题做子串匹配，返回命中的标签名列表
func MatchTags(tags []model.Tag, title string) []string {
	result := make([]string, 0)
	for _, t := range tags {
		if strings.Contains(title, t.Name) {
			result = append(result, t.Name)
		}
	}
	return result
}

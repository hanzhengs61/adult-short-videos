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

package model

import "time"

// GossipPost 吃瓜文章表
type GossipPost struct {
	ID      int64  `gorm:"primaryKey;autoIncrement;comment:文章 ID" json:"id"`
	Title   string `gorm:"type:varchar(500);not null;comment:标题" json:"title"`
	Cover   string `gorm:"type:varchar(500);comment:封面图 URL" json:"cover"`
	Summary string `gorm:"type:varchar(1000);comment:摘要（列表卡片展示）" json:"summary"`
	Content string `gorm:"type:text;comment:正文 HTML" json:"content"`
	Source  string `gorm:"type:varchar(200);default:'';comment:来源标识（如外站域名）" json:"source"`
	// 逗号分隔的标签，如 "爆料,女优,新人"；用于关联推荐视频
	Tags      string `gorm:"type:varchar(500);default:'';index;comment:标签，逗号分隔" json:"tags"`
	ViewCount int64  `gorm:"default:0;comment:浏览次数" json:"view_count"`
	Status    int32  `gorm:"default:1;index;comment:0:下架 1:正常" json:"status"`

	PublishedAt time.Time `gorm:"type:timestamp;comment:发布时间" json:"published_at"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (GossipPost) TableName() string {
	return "gossip_posts"
}

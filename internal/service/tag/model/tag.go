package model

import "time"

// Tag 标签库：预置关键词，视频入库时对标题做子串匹配
type Tag struct {
	TagId     int64     `gorm:"primaryKey;autoIncrement" json:"tag_id"`
	Name      string    `gorm:"type:varchar(100);uniqueIndex;not null;comment:标签名，如「深喉」「国产」" json:"name"`
	Status    int32     `gorm:"default:1;comment:0:禁用 1:启用" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Tag) TableName() string {
	return "tags"
}

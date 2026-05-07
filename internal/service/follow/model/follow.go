package model

import "time"

// AuthorFollow 用户关注作者记录
type AuthorFollow struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId    int64     `json:"user_id" gorm:"not null;index:idx_user_author,unique;comment:用户ID"`
	Author    string    `json:"author" gorm:"type:varchar(100);not null;index:idx_user_author,unique;comment:作者名"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;comment:关注时间"`
}

func (AuthorFollow) TableName() string {
	return "author_follows"
}

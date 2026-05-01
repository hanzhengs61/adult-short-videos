package model

import (
	"time"
)

// Comment 评论表
type Comment struct {
	CommentId int64     `gorm:"primaryKey;autoIncrement;comment:评论 ID" json:"comment_id"`
	VideoId   int64     `gorm:"index;not null;comment:视频 ID" json:"video_id"`
	UserId    int64     `gorm:"index;not null;comment:用户 ID" json:"user_id"`
	Content   string    `gorm:"type:text;not null;comment:评论内容" json:"content"`
	ParentId  int64     `gorm:"index;default:0;comment:父评论ID（0表示顶级评论）" json:"parent_id"`
	LikeCount int64     `gorm:"default:0;comment:点赞数" json:"like_count"`
	Status    int32     `gorm:"default:1;index;comment:0:已删除 1:正常 2:审核中" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime;index;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;comment:修改时间" json:"updated_at"`
}

func (Comment) TableName() string {
	return "comments"
}

// CommentWithUser 评论（包含用户信息）
type CommentWithUser struct {
	CommentId int64  `json:"comment_id"` // 评论 ID
	VideoId   int64  `json:"video_id"`   // 视频 ID
	UserId    int64  `json:"user_id"`    // 用户 ID
	Username  string `json:"username"`   // 用户名
	Avatar    string `json:"avatar"`     // 头像
	Content   string `json:"content"`    // 评论内容
	ParentId  int64  `json:"parent_id"`  // 父评论 ID
	LikeCount int64  `json:"like_count"` // 点赞数
	Status    int32  `json:"status"`     // 状态
	CreatedAt int64  `json:"created_at"` // 创建时间
	IsLiked   bool   `json:"is_liked"`   // 当前用户是否点赞
}

// CommentLike 评论点赞表
type CommentLike struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;comment:主键 ID" json:"id"`
	CommentId int64     `gorm:"index;not null;comment:评论 ID" json:"comment_id"`
	UserId    int64     `gorm:"index;not null;comment:用户 ID" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
}

func (CommentLike) TableName() string {
	return "comment_likes"
}

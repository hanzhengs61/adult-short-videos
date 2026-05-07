package model

import "time"

// Video 视频表
type Video struct {
	VideoId  int64  `gorm:"primaryKey;autoIncrement;comment:视频 ID" json:"video_id"`
	Title    string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	CoverURL string `gorm:"type:varchar(500);comment:封面图 URL" json:"cover_url"`
	Duration int32  `gorm:"default:0;comment:时长（秒）" json:"duration"`
	RemoteId string `gorm:"type:varchar(100);uniqueIndex;default:'';comment:外站唯一 ID，用于去重" json:"remote_id"`

	// ========== 冷热分离 ==========
	StorageType string `gorm:"type:varchar(10);default:'cold';comment:cold/hot" json:"storage_type"`
	SourceURL   string `gorm:"type:varchar(500);comment:源站 URL（冷数据）" json:"source_url"`
	LocalURL    string `gorm:"type:varchar(500);comment:本地 URL（热数据）" json:"local_url"`

	// ========== 统计 ==========
	PlayCount     int64 `gorm:"default:0;index:idx_play_count;comment:播放次数" json:"play_count"`
	FavoriteCount int64 `gorm:"default:0;comment:收藏数" json:"favorite_count"`
	CommentCount  int64 `gorm:"default:0;comment:评论数" json:"comment_count"`
	ShareCount    int64 `gorm:"default:0;comment:分享数" json:"share_count"`

	// ========== 媒体属性 ==========
	IsPortrait bool  `gorm:"default:false;comment:是否竖屏视频" json:"is_portrait"`
	UserId     int64 `gorm:"default:0;index;comment:关联用户ID，0=未关联" json:"user_id"`

	// ========== 状态 ==========
	Status int32 `gorm:"default:1;index;comment:0:下架 1:正常" json:"status"`

	// ========== 时间 ==========
	PublishedAt time.Time `gorm:"type:timestamp;comment:发布时间" json:"published_at"`
	CreatedAt   time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;comment:修改时间" json:"updated_at"`
}

func (Video) TableName() string {
	return "videos"
}

// VideoWithAuthor repository 层 JOIN 查询的扫描目标
// AuthorName/AuthorAvatar 由 COALESCE(u.username/avatar, ”) 填充
type VideoWithAuthor struct {
	Video
	AuthorName   string
	AuthorAvatar string
}

// VideoHeatStats 视频热度统计表
type VideoHeatStats struct {
	VideoId          int64     `gorm:"primaryKey;comment:视频 ID" json:"video_id"`
	PlayCount24h     int64     `gorm:"default:0;comment:24小时播放次数" json:"play_count_24h"`
	PlayCount7d      int64     `gorm:"default:0;comment:7天播放次数" json:"play_count_7d"`
	HeatScore        float64   `gorm:"default:0;index;comment:热度分数" json:"heat_score"`
	IsHot            bool      `gorm:"default:false;index;comment:是否为热门视频" json:"is_hot"`
	LastCalculatedAt time.Time `gorm:"type:timestamp;comment:最后计算时间" json:"last_calculated_at"`
}

func (VideoHeatStats) TableName() string {
	return "video_heat_stats"
}

package model

import "time"

// Video 视频表
type Video struct {
	VideoId     int64  `gorm:"primaryKey;autoIncrement;comment:视频 ID" json:"video_id"`
	Title       string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Description string `gorm:"type:text;comment:描述" json:"description"`
	CoverURL    string `gorm:"type:varchar(500);comment:封面图 URL" json:"cover_url"`
	PreviewURL  string `gorm:"type:varchar(500);comment:WebP 预览图" json:"preview_url"`
	Duration    int32  `gorm:"default:0;comment:时长（秒）" json:"duration"`

	// ========== 冷热分离核心字段 ==========
	StorageType string `gorm:"type:varchar(10);default:'cold';comment:cold/hot" json:"storage_type"`
	SourceURL   string `gorm:"type:varchar(500);comment:源站 URL（冷数据）" json:"source_url"`
	LocalURL    string `gorm:"type:varchar(500);comment:本地 URL（热数据）" json:"local_url"`

	// ========== 统计数据 ==========
	PlayCount     int64 `gorm:"default:0;index:idx_play_count;comment:播放次数" json:"play_count"`
	FavoriteCount int64 `gorm:"default:0" json:"favorite_count;comment:收藏数"`
	CommentCount  int64 `gorm:"default:0" json:"comment_count;comment:评论数"`
	ShareCount    int64 `gorm:"default:0" json:"share_count;comment:分享数"`

	// ========== 分类信息 ==========
	Fanhao   string `gorm:"type:varchar(50);uniqueIndex;comment:番号" json:"fanhao"`
	Region   string `gorm:"type:varchar(50);index;comment:地区：日本、欧美、国产" json:"region"`
	Category string `gorm:"type:varchar(50);index;comment:分类：职场、制服、学生等" json:"category"`

	// ========== 质量信息 ==========
	Resolution string `gorm:"type:varchar(20);comment:分辨率：480p, 720p, 1080p, 4k" json:"resolution"`
	FileSize   int64  `gorm:"default:0;comment:文件大小（字节）" json:"file_size"`

	// ========== 状态信息 ==========
	Status    int32 `gorm:"default:1;index;comment:0:下架 1:正常 2:审核中" json:"status"`
	IsVipOnly bool  `gorm:"default:false;comment:是否 VIP 专属" json:"is_vip_only"`

	// ========== 时间字段 ==========
	PublishedAt time.Time `gorm:"type:timestamp;comment:发布时间" json:"published_at"`
	CreatedAt   time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;comment:修改时间" json:"updated_at"`
}

func (Video) TableName() string {
	return "videos"
}

// VideoHeatStats 视频热度统计表
// 用于判断视频是否需要从冷数据转为热数据
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

// Actor 演员表
type Actor struct {
	ActorId       int64     `gorm:"primaryKey;autoIncrement;comment:演员 ID" json:"actor_id"`
	ActorName     string    `gorm:"type:varchar(100);uniqueIndex;not null;comment:演员名字" json:"actor_name"`
	AvatarURL     string    `gorm:"type:varchar(500);comment:演员头像 URL" json:"avatar_url"`
	Description   string    `gorm:"type:text;comment:描述" json:"description"`
	Region        string    `gorm:"type:varchar(50);comment:地区" json:"region"`
	VideoCount    int64     `gorm:"default:0;index;comment:作品数量" json:"video_count"`
	FollowerCount int64     `gorm:"default:0;index;comment:粉丝数" json:"follower_count"`
	CreatedAt     time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;comment:修改时间" json:"updated_at"`
}

func (Actor) TableName() string {
	return "actors"
}

// VideoActor 视频-演员关联表
type VideoActor struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;comment:主键 ID" json:"id"`
	VideoId   int64     `gorm:"index;not null;comment:视频 ID" json:"video_id"`
	ActorId   int64     `gorm:"index;not null;comment:演员 ID" json:"actor_id"`
	RoleType  string    `gorm:"type:varchar(20);comment:main=主演，supporting=配角" json:"role_type"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
}

func (VideoActor) TableName() string {
	return "video_actors"
}

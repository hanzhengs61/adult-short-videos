package model

import "time"

// Favorite 收藏表
// 记录用户收藏了哪些视频
type Favorite struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;comment:主键 ID" json:"id"`
	UserId    int64     `gorm:"index;not null;comment:用户 ID" json:"user_id"`
	VideoId   int64     `gorm:"index;not null;comment:视频 ID" json:"video_id"`
	CreatedAt time.Time `gorm:"autoCreateTime;index;comment:收藏时间" json:"created_at"`
}

func (Favorite) TableName() string {
	return "favorites"
}

// FavoriteWithVideo 收藏记录（包含视频信息）
// 用于返回给前端，让用户看到收藏的视频详情
type FavoriteWithVideo struct {
	FavoriteId  int64  `json:"favorite_id"`  // 收藏记录 ID
	VideoId     int64  `json:"video_id"`     // 视频 ID
	Title       string `json:"title"`        // 视频标题
	CoverURL    string `json:"cover_url"`    // 封面
	Duration    int32  `json:"duration"`     // 时长
	IsPortrait  bool   `json:"is_portrait"`  // 是否竖屏
	PlayCount   int64  `json:"play_count"`   // 播放次数
	FavoritedAt int64  `json:"favorited_at"` // 收藏时间
}

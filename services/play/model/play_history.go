package model

import "time"

// PlayHistory 播放历史表
type PlayHistory struct {
	ID           int64     `gorm:"primaryKey;autoIncrement;comment:主键 ID" json:"id"`
	UserId       int64     `gorm:"index;not null;comment:用户 ID" json:"user_id"`
	VideoId      int64     `gorm:"index;not null;comment:视频 ID" json:"video_id"`
	PlayDuration int32     `gorm:"default:0;comment:观看时长（秒）" json:"play_duration"`
	PlayProgress float32   `gorm:"default:0;comment:播放进度（0-100）" json:"play_progress"`
	LastPlayAt   time.Time `gorm:"type:timestamp;index;comment:最后播放时间" json:"last_play_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;comment:修改时间" json:"updated_at"`
}

func (PlayHistory) TableName() string {
	return "play_history"
}

// PlayHistoryWithVideo 播放历史（包含视频信息）
type PlayHistoryWithVideo struct {
	HistoryId    int64   `json:"history_id"`    // 播放历史 ID
	VideoId      int64   `json:"video_id"`      // 视频 ID
	Title        string  `json:"title"`         // 视频标题
	CoverURL     string  `json:"cover_url"`     // 封面
	PreviewURL   string  `json:"preview_url"`   // 预览图
	Duration     int32   `json:"duration"`      // 时长
	PlayDuration int32   `json:"play_duration"` // 观看时长
	PlayProgress float32 `json:"play_progress"` // 播放进度
	Fanhao       string  `json:"fanhao"`        // 番号
	Region       string  `json:"region"`        // 地区
	Category     string  `json:"category"`      // 分类
	LastPlayAt   int64   `json:"last_play_at"`  // 最后播放时间
}

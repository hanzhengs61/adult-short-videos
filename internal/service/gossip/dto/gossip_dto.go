package dto

import "adult-short-videos/internal/service/video/dto"

// GossipListReq 吃瓜列表请求
type GossipListReq struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Tag      string `form:"tag"` // 按标签过滤，空表示全部
}

// GossipItem 列表卡片数据
type GossipItem struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Cover       string   `json:"cover"`
	Summary     string   `json:"summary"`
	Tags        []string `json:"tags"`
	ViewCount   int64    `json:"view_count"`
	PublishedAt int64    `json:"published_at"` // Unix 时间戳
}

// GossipListResp 吃瓜列表响应
type GossipListResp struct {
	Total int64        `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
	Posts []GossipItem `json:"posts"`
}

// GossipDetail 文章详情（含关联视频）
type GossipDetail struct {
	ID            int64           `json:"id"`
	Title         string          `json:"title"`
	Cover         string          `json:"cover"`
	Content       string          `json:"content"`
	Tags          []string        `json:"tags"`
	ViewCount     int64           `json:"view_count"`
	PublishedAt   int64           `json:"published_at"`
	RelatedVideos []dto.VideoItem `json:"related_videos"` // 底部推荐视频
}

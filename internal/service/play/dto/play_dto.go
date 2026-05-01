package dto

import "adult-short-videos/internal/service/play/model"

// RecordPlayReq 记录播放请求
type RecordPlayReq struct {
	VideoId      int64   `json:"video_id"`
	PlayDuration int32   `json:"play_duration"` // 本次播放时长，秒
	PlayProgress float32 `json:"play_progress"` // 播放进度百分比 (0-100)
}

// PlayHistoryListReq 播放历史列表请求
type PlayHistoryListReq struct {
	Page int `form:"page"`
	Size int `form:"page_size"`
}

// PlayHistoryListResp 播放历史列表响应
type PlayHistoryListResp struct {
	Total   int64                         `json:"total"`
	Page    int                           `json:"page"`
	Size    int                           `json:"size"`
	History []*model.PlayHistoryWithVideo `json:"history"` // 包含视频信息的播放历史
}

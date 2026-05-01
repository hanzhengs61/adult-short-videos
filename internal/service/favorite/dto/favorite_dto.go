package dto

import "adult-short-videos/internal/service/favorite/model"

// AddFavoriteReq 添加收藏请求
type AddFavoriteReq struct {
	VideoId int64 `json:"video_id"`
}

// FavoriteListReq 收藏列表请求
type FavoriteListReq struct {
	Page int `form:"page"`
	Size int `form:"page_size"`
}

// FavoriteListResp 收藏列表响应
type FavoriteListResp struct {
	Total     int64                      `json:"total"`
	Page      int                        `json:"page"`
	Size      int                        `json:"size"`
	Favorites []*model.FavoriteWithVideo `json:"favorites"` // 包含视频详情
}

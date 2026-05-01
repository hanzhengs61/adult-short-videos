package dto

// SearchReq 搜索请求
type SearchReq struct {
	Keyword string `form:"keyword"`
	Page    int    `form:"page"`
	Size    int    `form:"page_size"`
	OrderBy string `form:"order_by"` // created_at, play_count
}

// VideoSearchItem 搜索结果项
type VideoSearchItem struct {
	VideoId       int64  `json:"video_id"`
	Title         string `json:"title"`
	CoverURL      string `json:"cover_url"`
	Duration      int32  `json:"duration"`
	IsPortrait    bool   `json:"is_portrait"`
	Author        string `json:"author"`
	PlayURL       string `json:"play_url"`
	PlayCount     int64  `json:"play_count"`
	FavoriteCount int64  `json:"favorite_count"`
	PublishedAt   int64  `json:"published_at"`
}

// SearchResp 搜索响应
type SearchResp struct {
	Total  int64             `json:"total"`
	Page   int               `json:"page"`
	Size   int               `json:"size"`
	Videos []VideoSearchItem `json:"videos"`
}

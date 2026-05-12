package dto

// VideoItem 视频列表项
type VideoItem struct {
	VideoId       int64    `json:"video_id"`
	Title         string   `json:"title"`
	CoverURL      string   `json:"cover_url"`
	Duration      int32    `json:"duration"` // 秒
	IsPortrait    bool     `json:"is_portrait"`
	AuthorId      int64    `json:"author_id"`
	AuthorName    string   `json:"author_name"`
	AuthorAvatar  string   `json:"author_avatar"`
	PlayURL       string   `json:"play_url"`
	PlayCount     int64    `json:"play_count"`
	FavoriteCount int64    `json:"favorite_count"`
	CommentCount  int64    `json:"comment_count"`
	PublishedAt   int64    `json:"published_at"` // Unix 时间戳
	Tags          []string `json:"tags"`
}

// VideoDetail 视频详情
type VideoDetail struct {
	VideoId       int64    `json:"video_id"`
	Title         string   `json:"title"`
	CoverURL      string   `json:"cover_url"`
	Duration      int32    `json:"duration"`
	IsPortrait    bool     `json:"is_portrait"`
	AuthorId      int64    `json:"author_id"`
	AuthorName    string   `json:"author_name"`
	AuthorAvatar  string   `json:"author_avatar"`
	StorageType   string   `json:"storage_type"`
	PlayURL       string   `json:"play_url"`
	PlayCount     int64    `json:"play_count"`
	FavoriteCount int64    `json:"favorite_count"`
	CommentCount  int64    `json:"comment_count"`
	PublishedAt   int64    `json:"published_at"`
	Tags          []string `json:"tags"`
}

// VideoListReq 视频列表请求
type VideoListReq struct {
	Page       int    `form:"page"`
	Size       int    `form:"page_size"`
	OrderBy    string `form:"order_by"`    // created_at, play_count, comment_count, favorite_count
	AuthorName string `form:"author_name"` // 按作者用户名过滤
}

// VideoListResp 视频列表响应
type VideoListResp struct {
	Total  int64       `json:"total"`
	Page   int         `json:"page"`
	Size   int         `json:"size"`
	Videos []VideoItem `json:"videos"`
}

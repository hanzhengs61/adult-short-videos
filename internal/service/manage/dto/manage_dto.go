package dto

// PageReq 通用分页请求
type PageReq struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

func (p *PageReq) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 || p.PageSize > 100 {
		p.PageSize = 20
	}
	return (p.Page - 1) * p.PageSize
}

// ── Dashboard ────────────────────────────────────────

type DashboardResp struct {
	TotalVideos   int64       `json:"total_videos"`
	TotalUsers    int64       `json:"total_users"`
	TotalComments int64       `json:"total_comments"`
	TodayPlays    int64       `json:"today_plays"`
	Trend         []TrendItem `json:"trend"` // 近 7 天播放趋势
}

type TrendItem struct {
	Date  string `json:"date"`
	Plays int64  `json:"plays"`
}

// ── 视频管理 ─────────────────────────────────────────

type VideoListReq struct {
	PageReq
	Status  *int32 `form:"status"`  // 不传则全部
	Keyword string `form:"keyword"` // 标题模糊搜索
}

type VideoItem struct {
	VideoId   int64  `json:"video_id"`
	Title     string `json:"title"`
	CoverURL  string `json:"cover_url"`
	Duration  int32  `json:"duration"`
	PlayCount int64  `json:"play_count"`
	Status    int32  `json:"status"`
	CreatedAt string `json:"created_at"`
}

type VideoListResp struct {
	Total  int64       `json:"total"`
	Videos []VideoItem `json:"videos"`
}

type VideoStatusReq struct {
	// 0=下架 1=正常
	Status int32 `json:"status" binding:"oneof=0 1"`
}

// ── 评论管理 ─────────────────────────────────────────

type CommentListReq struct {
	PageReq
	VideoId int64  `form:"video_id"` // 0 则不过滤
	Keyword string `form:"keyword"`
}

type CommentItem struct {
	CommentId  int64  `json:"comment_id"`
	VideoId    int64  `json:"video_id"`
	VideoTitle string `json:"video_title"`
	UserId     int64  `json:"user_id"`
	Username   string `json:"username"`
	Content    string `json:"content"`
	LikeCount  int64  `json:"like_count"`
	CreatedAt  string `json:"created_at"`
}

type CommentListResp struct {
	Total    int64         `json:"total"`
	Comments []CommentItem `json:"comments"`
}

type BatchDeleteReq struct {
	IDs []int64 `json:"ids" binding:"required,min=1"`
}

// ── 用户管理 ─────────────────────────────────────────

type UserListReq struct {
	PageReq
	Keyword string `form:"keyword"` // 用户名/邮箱
	Status  *int32 `form:"status"`
}

type UserItem struct {
	UserId    int64  `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    int32  `json:"status"`
	CreatedAt string `json:"created_at"`
}

type UserListResp struct {
	Total int64      `json:"total"`
	Users []UserItem `json:"users"`
}

type UserUpdateReq struct {
	// 不传的字段不更新，用指针区分"未传"和"传了零值"
	Status *int32  `json:"status"` // 0=禁用 1=正常
	Role   *string `json:"role"`   // user / creator / admin
}

// ── 广告 Banner ──────────────────────────────────────

type BannerSaveReq struct {
	Title    string `json:"title"    binding:"required"`
	SubTitle string `json:"sub_title"`
	BgStyle  string `json:"bg_style"  binding:"required"`
	Tag      string `json:"tag"`
	TagClass string `json:"tag_class"`
	URL      string `json:"url"`
	Sort     int32  `json:"sort"`
	Status   int32  `json:"status"`
}

// ── 广告 App ─────────────────────────────────────────

type AppSaveReq struct {
	Name       string `json:"name"        binding:"required"`
	Icon       string `json:"icon"`
	IconText   string `json:"icon_text"`
	BgStyle    string `json:"bg_style"    binding:"required"`
	Official   bool   `json:"official"`
	Hot        bool   `json:"hot"`
	Categories string `json:"categories"  binding:"required"`
	URL        string `json:"url"`
	Sort       int32  `json:"sort"`
	Status     int32  `json:"status"`
}

// ── 吃瓜文章 ─────────────────────────────────────────

type GossipListReq struct {
	PageReq
	Keyword string `form:"keyword"`
	Status  *int32 `form:"status"`
}

type GossipItem struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	ViewCount   int64  `json:"view_count"`
	Status      int32  `json:"status"`
	PublishedAt string `json:"published_at"`
}

type GossipListResp struct {
	Total int64        `json:"total"`
	Posts []GossipItem `json:"posts"`
}

type GossipSaveReq struct {
	Title       string `json:"title"       binding:"required"`
	Cover       string `json:"cover"`
	Summary     string `json:"summary"`
	Content     string `json:"content"`
	Tags        string `json:"tags"`
	Status      int32  `json:"status"`
	PublishedAt string `json:"published_at"` // RFC3339
}

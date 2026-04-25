package logic

import (
	"context"
	"errors"

	"adult-short-videos/internal/pkg/metrics"
	"adult-short-videos/internal/service/search/repository"
	videoModel "adult-short-videos/internal/service/video/model"

	"gorm.io/gorm"
)

// SearchLogic 搜索业务逻辑
type SearchLogic struct {
	ctx        context.Context
	searchRepo repository.SearchRepository
}

// NewSearchLogic 创建搜索逻辑实例
func NewSearchLogic(ctx context.Context, searchRepo repository.SearchRepository) *SearchLogic {
	return &SearchLogic{
		ctx:        ctx,
		searchRepo: searchRepo,
	}
}

// SearchReq 搜索请求
type SearchReq struct {
	Keyword string `form:"q"`    // 搜索关键词
	Page    int    `form:"page"` // 页码
	Size    int    `form:"size"` // 每页数量
}

// VideoSearchItem 视频搜索结果项
type VideoSearchItem struct {
	VideoId       int64  `json:"video_id"`
	Title         string `json:"title"`
	CoverURL      string `json:"cover_url"`
	PreviewURL    string `json:"preview_url"`
	Duration      int32  `json:"duration"`
	PlayCount     int64  `json:"play_count"`
	FavoriteCount int64  `json:"favorite_count"`
	Fanhao        string `json:"fanhao"`
	Region        string `json:"region"`
	Category      string `json:"category"`
	IsVipOnly     bool   `json:"is_vip_only"`
	PublishedAt   int64  `json:"published_at"`
}

// SearchResp 搜索响应
type SearchResp struct {
	Total  int64             `json:"total"`
	Page   int               `json:"page"`
	Size   int               `json:"size"`
	Videos []VideoSearchItem `json:"videos"`
}

// Search 搜索视频
func (l *SearchLogic) Search(req *SearchReq) (*SearchResp, error) {

	if req.Keyword == "" {
		return nil, errors.New("搜索关键词不能为空")
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Size < 1 {
		req.Size = 20
	}
	if req.Size > 100 {
		req.Size = 100
	}

	metrics.SearchRequestsTotal.WithLabelValues("video").Inc()
	offset := (req.Page - 1) * req.Size

	videos, total, err := l.searchRepo.SearchVideos(l.ctx, req.Keyword, offset, req.Size)
	if err != nil {
		return nil, errors.New("搜索失败")
	}

	items := make([]VideoSearchItem, 0, len(videos))
	for _, video := range videos {
		items = append(items, VideoSearchItem{
			VideoId:       video.VideoId,
			Title:         video.Title,
			CoverURL:      video.CoverURL,
			PreviewURL:    video.PreviewURL,
			Duration:      video.Duration,
			PlayCount:     video.PlayCount,
			FavoriteCount: video.FavoriteCount,
			Fanhao:        video.Fanhao,
			Region:        video.Region,
			Category:      video.Category,
			IsVipOnly:     video.IsVipOnly,
			PublishedAt:   video.PublishedAt.Unix(),
		})
	}

	return &SearchResp{
		Total:  total,
		Page:   req.Page,
		Size:   req.Size,
		Videos: items,
	}, nil
}

// ActorSearchLogic 演员搜索逻辑
type ActorSearchLogic struct {
	ctx        context.Context
	searchRepo repository.SearchRepository
}

func NewActorSearchLogic(ctx context.Context, searchRepo repository.SearchRepository) *ActorSearchLogic {
	return &ActorSearchLogic{
		ctx:        ctx,
		searchRepo: searchRepo,
	}
}

// ActorSearchItem 演员搜索结果项
type ActorSearchItem struct {
	ActorId       int64  `json:"actor_id"`
	ActorName     string `json:"actor_name"`
	AvatarURL     string `json:"avatar_url"`
	Region        string `json:"region"`
	VideoCount    int64  `json:"video_count"`
	FollowerCount int64  `json:"follower_count"`
}

// ActorSearchResp 演员搜索响应
type ActorSearchResp struct {
	Total  int64             `json:"total"`
	Page   int               `json:"page"`
	Size   int               `json:"size"`
	Actors []ActorSearchItem `json:"actors"`
}

// SearchActors 搜索演员
func (l *ActorSearchLogic) SearchActors(req *SearchReq) (*ActorSearchResp, error) {
	if req.Keyword == "" {
		return nil, errors.New("搜索关键词不能为空")
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Size < 1 {
		req.Size = 20
	}
	if req.Size > 100 {
		req.Size = 100
	}

	metrics.SearchRequestsTotal.WithLabelValues("actor").Inc()
	offset := (req.Page - 1) * req.Size

	actors, total, err := l.searchRepo.SearchActors(l.ctx, req.Keyword, offset, req.Size)
	if err != nil {
		return nil, errors.New("搜索失败")
	}

	items := make([]ActorSearchItem, 0, len(actors))
	for _, actor := range actors {
		items = append(items, ActorSearchItem{
			ActorId:       actor.ActorId,
			ActorName:     actor.ActorName,
			AvatarURL:     actor.AvatarURL,
			Region:        actor.Region,
			VideoCount:    actor.VideoCount,
			FollowerCount: actor.FollowerCount,
		})
	}

	return &ActorSearchResp{
		Total:  total,
		Page:   req.Page,
		Size:   req.Size,
		Actors: items,
	}, nil
}

// FanhaoSearchLogic 番号搜索逻辑
type FanhaoSearchLogic struct {
	ctx        context.Context
	searchRepo repository.SearchRepository
}

func NewFanhaoSearchLogic(ctx context.Context, searchRepo repository.SearchRepository) *FanhaoSearchLogic {
	return &FanhaoSearchLogic{
		ctx:        ctx,
		searchRepo: searchRepo,
	}
}

type SearchFanhaoReq struct {
	Fanhao string `json:"fanhao"` // 番号
}

// SearchByFanhao 按番号搜索
func (l *FanhaoSearchLogic) SearchByFanhao(req SearchFanhaoReq) (*videoModel.Video, error) {
	if req.Fanhao == "" {
		return nil, errors.New("番号不能为空")
	}

	metrics.SearchRequestsTotal.WithLabelValues("fanhao").Inc()
	video, err := l.searchRepo.SearchByFanhao(l.ctx, req.Fanhao)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("未找到该番号的视频")
		}
		return nil, errors.New("搜索失败")
	}

	return video, nil
}

// AdvancedSearchLogic 高级搜索逻辑
type AdvancedSearchLogic struct {
	ctx        context.Context
	searchRepo repository.SearchRepository
}

func NewAdvancedSearchLogic(ctx context.Context, searchRepo repository.SearchRepository) *AdvancedSearchLogic {
	return &AdvancedSearchLogic{
		ctx:        ctx,
		searchRepo: searchRepo,
	}
}

// AdvancedSearchReq 高级搜索请求
type AdvancedSearchReq struct {
	Keyword     string `form:"q"`            // 关键词
	Region      string `form:"region"`       // 地区
	Category    string `form:"category"`     // 分类
	Resolution  string `form:"resolution"`   // 分辨率
	MinDuration int    `form:"min_duration"` // 最小时长（分钟）
	MaxDuration int    `form:"max_duration"` // 最大时长（分钟）
	ActorId     int64  `form:"actor_id"`     // 演员 ID
	Sort        string `form:"sort"`         // 排序：hot, favorite, newest, oldest
	Page        int    `form:"page"`
	Size        int    `form:"size"`
}

// AdvancedSearch 高级搜索
func (l *AdvancedSearchLogic) AdvancedSearch(req *AdvancedSearchReq) (*SearchResp, error) {
	// 参数校验
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Size < 1 {
		req.Size = 20
	}
	if req.Size > 100 {
		req.Size = 100
	}

	metrics.SearchRequestsTotal.WithLabelValues("advanced").Inc()
	offset := (req.Page - 1) * req.Size

	// 构建过滤条件
	filters := make(map[string]interface{})
	if req.Keyword != "" {
		filters["keyword"] = req.Keyword
	}
	if req.Region != "" {
		filters["region"] = req.Region
	}
	if req.Category != "" {
		filters["category"] = req.Category
	}
	if req.Resolution != "" {
		filters["resolution"] = req.Resolution
	}
	if req.MinDuration > 0 {
		filters["min_duration"] = req.MinDuration
	}
	if req.MaxDuration > 0 {
		filters["max_duration"] = req.MaxDuration
	}
	if req.ActorId > 0 {
		filters["actor_id"] = req.ActorId
	}
	if req.Sort != "" {
		filters["sort"] = req.Sort
	}

	// 执行搜索
	videos, total, err := l.searchRepo.AdvancedSearch(l.ctx, filters, offset, req.Size)
	if err != nil {
		return nil, errors.New("搜索失败")
	}

	// 转换结果
	items := make([]VideoSearchItem, 0, len(videos))
	for _, video := range videos {
		items = append(items, VideoSearchItem{
			VideoId:       video.VideoId,
			Title:         video.Title,
			CoverURL:      video.CoverURL,
			PreviewURL:    video.PreviewURL,
			Duration:      video.Duration,
			PlayCount:     video.PlayCount,
			FavoriteCount: video.FavoriteCount,
			Fanhao:        video.Fanhao,
			Region:        video.Region,
			Category:      video.Category,
			IsVipOnly:     video.IsVipOnly,
			PublishedAt:   video.PublishedAt.Unix(),
		})
	}

	return &SearchResp{
		Total:  total,
		Page:   req.Page,
		Size:   req.Size,
		Videos: items,
	}, nil
}

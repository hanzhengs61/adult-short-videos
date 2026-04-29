package logic

import (
	"context"
	"errors"

	"adult-short-videos/internal/pkg/metrics"
	"adult-short-videos/internal/service/search/repository"
)

// SearchLogic 搜索业务逻辑
type SearchLogic struct {
	ctx        context.Context
	searchRepo repository.SearchRepository
}

// NewSearchLogic 创建搜索逻辑实例
func NewSearchLogic(ctx context.Context, searchRepo repository.SearchRepository) *SearchLogic {
	return &SearchLogic{ctx: ctx, searchRepo: searchRepo}
}

// SearchReq 搜索请求
type SearchReq struct {
	Keyword string `form:"keyword"`   // 搜索关键词
	Page    int    `form:"page"`      // 页码
	Size    int    `form:"page_size"` // 每页数量
}

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

type SearchResp struct {
	Total  int64             `json:"total"`
	Page   int               `json:"page"`
	Size   int               `json:"size"`
	Videos []VideoSearchItem `json:"videos"`
}

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
	for _, v := range videos {
		playURL := v.SourceURL
		if v.StorageType != "hot" || v.LocalURL == "" {
			playURL = "/api/storage/proxy?url=" + v.SourceURL
		} else {
			playURL = v.LocalURL
		}
		items = append(items, VideoSearchItem{
			VideoId:       v.VideoId,
			Title:         v.Title,
			CoverURL:      v.CoverURL,
			Duration:      v.Duration,
			IsPortrait:    v.IsPortrait,
			Author:        v.Author,
			PlayURL:       playURL,
			PlayCount:     v.PlayCount,
			FavoriteCount: v.FavoriteCount,
			PublishedAt:   v.PublishedAt.Unix(),
		})
	}

	return &SearchResp{
		Total:  total,
		Page:   req.Page,
		Size:   req.Size,
		Videos: items,
	}, nil
}

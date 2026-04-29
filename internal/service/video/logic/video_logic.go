package logic

import (
	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/service/video/repository"
	"context"

	"gorm.io/gorm"
)

// VideoItem 视频列表项
type VideoItem struct {
	VideoId    int64  `json:"video_id"`
	Title      string `json:"title"`
	CoverURL   string `json:"cover_url"`
	Duration   int32  `json:"duration"`
	IsPortrait bool   `json:"is_portrait"`
	Author     string `json:"author"`

	PlayURL       string `json:"play_url"`
	PlayCount     int64  `json:"play_count"`
	FavoriteCount int64  `json:"favorite_count"`
	PublishedAt   int64  `json:"published_at"`
}

func BuildPlayURL(storageType, sourceURL, localURL string) string {
	if storageType == "hot" && localURL != "" {
		return localURL
	}
	return "/api/storage/proxy?url=" + sourceURL
}

// VideoDetail 视频详情
type VideoDetail struct {
	VideoId    int64  `json:"video_id"`
	Title      string `json:"title"`
	CoverURL   string `json:"cover_url"`
	Duration   int32  `json:"duration"`
	IsPortrait bool   `json:"is_portrait"`
	Author     string `json:"author"`

	StorageType   string `json:"storage_type"`
	PlayURL       string `json:"play_url"`
	PlayCount     int64  `json:"play_count"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	PublishedAt   int64  `json:"published_at"`
}

// VideoListReq 视频列表请求
type VideoListReq struct {
	Page    int
	Size    int
	OrderBy string
}

// VideoListResp 视频列表响应
type VideoListResp struct {
	Total  int64       `json:"total"`
	Page   int         `json:"page"`
	Size   int         `json:"size"`
	Videos []VideoItem `json:"videos"`
}

// VideoListLogic 视频列表逻辑
type VideoListLogic struct {
	ctx       context.Context
	videoRepo repository.VideoRepository
}

func NewVideoListLogic(ctx context.Context, videoRepo repository.VideoRepository) *VideoListLogic {
	return &VideoListLogic{ctx: ctx, videoRepo: videoRepo}
}

func (l *VideoListLogic) GetVideoList(req *VideoListReq) (*VideoListResp, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Size < 1 {
		req.Size = 20
	}
	if req.Size > 100 {
		req.Size = 100
	}

	filters := map[string]interface{}{}
	if req.OrderBy != "" {
		filters["order_by"] = req.OrderBy
	}

	offset := (req.Page - 1) * req.Size
	videos, total, err := l.videoRepo.List(l.ctx, offset, req.Size, filters)
	if err != nil {
		return nil, errors.New(errors.CodeVideoNotFound, "查询视频列表失败")
	}

	items := make([]VideoItem, 0, len(videos))
	for _, v := range videos {
		items = append(items, VideoItem{
			VideoId:    v.VideoId,
			Title:      v.Title,
			CoverURL:   v.CoverURL,
			Duration:   v.Duration,
			IsPortrait: v.IsPortrait,
			Author:     v.Author,

			PlayURL:       BuildPlayURL(v.StorageType, v.SourceURL, v.LocalURL),
			PlayCount:     v.PlayCount,
			FavoriteCount: v.FavoriteCount,
			PublishedAt:   v.PublishedAt.Unix(),
		})
	}

	return &VideoListResp{
		Total:  total,
		Page:   req.Page,
		Size:   req.Size,
		Videos: items,
	}, nil
}

// VideoDetailLogic 视频详情逻辑
type VideoDetailLogic struct {
	ctx       context.Context
	videoRepo repository.VideoRepository
}

func NewVideoDetailLogic(ctx context.Context, videoRepo repository.VideoRepository) *VideoDetailLogic {
	return &VideoDetailLogic{ctx: ctx, videoRepo: videoRepo}
}

func (l *VideoDetailLogic) GetVideoDetail(videoId int64) (*VideoDetail, error) {
	video, err := l.videoRepo.FindByID(l.ctx, videoId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New(errors.CodeVideoNotFound, "查询视频失败")
	}

	playURL := BuildPlayURL(video.StorageType, video.SourceURL, video.LocalURL)

	go func() {
		l.videoRepo.IncrementPlayCount(context.Background(), videoId)
	}()

	return &VideoDetail{
		VideoId:       video.VideoId,
		Title:         video.Title,
		CoverURL:      video.CoverURL,
		Duration:      video.Duration,
		IsPortrait:    video.IsPortrait,
		Author:        video.Author,
		StorageType:   video.StorageType,
		PlayURL:       playURL,
		PlayCount:     video.PlayCount,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		PublishedAt:   video.PublishedAt.Unix(),
	}, nil
}

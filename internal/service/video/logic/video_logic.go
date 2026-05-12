package logic

import (
	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/pkg/logger"
	"adult-short-videos/internal/service/video/dto"
	"adult-short-videos/internal/service/video/repository"
	"context"
	errs "errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// VideoService 视频服务
type VideoService struct {
	videoRepo repository.VideoRepository
}

// NewVideoService 创建 VideoService 实例
func NewVideoService(videoRepo repository.VideoRepository) *VideoService {
	return &VideoService{videoRepo: videoRepo}
}

// GetVideoList 获取视频列表
func (s *VideoService) GetVideoList(ctx context.Context, req *dto.VideoListReq) (*dto.VideoListResp, error) {
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
	if req.AuthorName != "" {
		filters["author_name"] = req.AuthorName
	}

	offset := (req.Page - 1) * req.Size
	videos, total, err := s.videoRepo.List(ctx, offset, req.Size, filters)
	if err != nil {
		if errs.Is(err, context.Canceled) {
			return nil, errors.New(errors.CodeVideoNotFound, "请求已取消")
		}
		logger.Error("查询视频列表失败", zap.Error(err), zap.Any("filters", filters))
		return nil, errors.New(errors.CodeVideoNotFound, "查询视频列表失败")
	}

	items := make([]dto.VideoItem, 0, len(videos))
	for _, v := range videos {
		tags := v.Tags
		if tags == nil {
			tags = []string{}
		}
		items = append(items, dto.VideoItem{
			VideoId:       v.VideoId,
			Title:         v.Title,
			CoverURL:      v.CoverURL,
			Duration:      v.Duration,
			IsPortrait:    v.IsPortrait,
			AuthorId:      v.UserId,
			AuthorName:    v.AuthorName,
			AuthorAvatar:  v.AuthorAvatar,
			PlayURL:       BuildPlayURL(v.StorageType, v.SourceURL, v.LocalURL),
			PlayCount:     v.PlayCount,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			PublishedAt:   v.PublishedAt.Unix(),
			Tags:          tags,
		})
	}

	return &dto.VideoListResp{
		Total:  total,
		Page:   req.Page,
		Size:   req.Size,
		Videos: items,
	}, nil
}

// GetVideoDetail 获取视频详情
func (s *VideoService) GetVideoDetail(ctx context.Context, videoId int64) (*dto.VideoDetail, error) {
	video, err := s.videoRepo.FindByID(ctx, videoId)
	if err != nil {
		if errs.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(errors.CodeVideoNotFound, "视频不存在")
		}
		logger.Error("查询视频详情失败", zap.Error(err), zap.Int64("videoId", videoId))
		return nil, errors.New(errors.CodeDatabaseError, "查询视频失败")
	}

	playURL := BuildPlayURL(video.StorageType, video.SourceURL, video.LocalURL)

	// 异步增加播放数
	go func() {
		_ = s.videoRepo.IncrementPlayCount(context.Background(), videoId)
	}()

	tags := video.Tags
	if tags == nil {
		tags = []string{}
	}

	return &dto.VideoDetail{
		VideoId:       video.VideoId,
		Title:         video.Title,
		CoverURL:      video.CoverURL,
		Duration:      video.Duration,
		IsPortrait:    video.IsPortrait,
		AuthorId:      video.UserId,
		AuthorName:    video.AuthorName,
		AuthorAvatar:  video.AuthorAvatar,
		StorageType:   video.StorageType,
		PlayURL:       playURL,
		PlayCount:     video.PlayCount,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		PublishedAt:   video.PublishedAt.Unix(),
		Tags:          tags,
	}, nil
}

// GetHotVideos 获取热门视频列表
func (s *VideoService) GetHotVideos(ctx context.Context, limit int) ([]dto.VideoItem, error) {
	videos, err := s.videoRepo.GetHotVideos(ctx, limit)
	if err != nil {
		logger.Error("查询热门视频失败", zap.Error(err), zap.Int("limit", limit))
		return nil, errors.New(errors.CodeDatabaseError, "查询热门视频失败")
	}

	items := make([]dto.VideoItem, 0, len(videos))
	for _, v := range videos {
		tags := v.Tags
		if tags == nil {
			tags = []string{}
		}
		items = append(items, dto.VideoItem{
			VideoId:       v.VideoId,
			Title:         v.Title,
			CoverURL:      v.CoverURL,
			Duration:      v.Duration,
			IsPortrait:    v.IsPortrait,
			AuthorId:      v.UserId,
			AuthorName:    v.AuthorName,
			AuthorAvatar:  v.AuthorAvatar,
			PlayURL:       BuildPlayURL(v.StorageType, v.SourceURL, v.LocalURL),
			PlayCount:     v.PlayCount,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			PublishedAt:   v.PublishedAt.Unix(),
			Tags:          tags,
		})
	}
	return items, nil
}

// BuildPlayURL 构建播放链接
func BuildPlayURL(storageType, sourceURL, localURL string) string {
	if storageType == "hot" && localURL != "" {
		return localURL
	}
	return "/api/storage/proxy?url=" + sourceURL
}

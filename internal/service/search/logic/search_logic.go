package logic

import (
	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/pkg/logger"
	"context"

	"adult-short-videos/internal/pkg/metrics"
	"adult-short-videos/internal/service/search/dto"
	"adult-short-videos/internal/service/search/repository"
	videoLogic "adult-short-videos/internal/service/video/logic" // 导入 video logic 包

	"go.uber.org/zap"
)

// SearchService 搜索业务服务
type SearchService struct {
	searchRepo repository.SearchRepository
}

// NewSearchService 创建搜索服务实例
func NewSearchService(searchRepo repository.SearchRepository) *SearchService {
	return &SearchService{searchRepo: searchRepo}
}

// Search 执行搜索
func (s *SearchService) Search(ctx context.Context, req *dto.SearchReq) (*dto.SearchResp, error) {
	if req.Keyword == "" {
		logger.Error("搜索关键词为空")
		return nil, errors.New(errors.CodeInvalidParam, "搜索关键词为空")
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

	filters := map[string]interface{}{}
	if req.OrderBy != "" {
		filters["order_by"] = req.OrderBy
	}

	metrics.SearchRequestsTotal.WithLabelValues("video").Inc()
	offset := (req.Page - 1) * req.Size

	videos, total, err := s.searchRepo.SearchVideos(ctx, req.Keyword, offset, req.Size, filters)
	if err != nil {
		logger.Error("搜索失败", zap.Error(err))
		return nil, errors.New(errors.CodeDatabaseError, "搜索失败")
	}

	items := make([]dto.VideoSearchItem, 0, len(videos))
	for _, v := range videos {
		items = append(items, dto.VideoSearchItem{
			VideoId:       v.VideoId,
			Title:         v.Title,
			CoverURL:      v.CoverURL,
			Duration:      v.Duration,
			IsPortrait:    v.IsPortrait,
			AuthorId:      v.UserId,
			AuthorName:    v.AuthorName,
			AuthorAvatar:  v.AuthorAvatar,
			PlayURL:       videoLogic.BuildPlayURL(v.StorageType, v.SourceURL, v.LocalURL),
			PlayCount:     v.PlayCount,
			FavoriteCount: v.FavoriteCount,
			PublishedAt:   v.PublishedAt.Unix(),
		})
	}

	return &dto.SearchResp{
		Total:  total,
		Page:   req.Page,
		Size:   req.Size,
		Videos: items,
	}, nil
}

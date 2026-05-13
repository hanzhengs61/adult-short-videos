package logic

import (
	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/pkg/logger"
	gossipDTO "adult-short-videos/internal/service/gossip/dto"
	"adult-short-videos/internal/service/gossip/repository"
	videoDTO "adult-short-videos/internal/service/video/dto"
	videoLogic "adult-short-videos/internal/service/video/logic"
	videoRepo "adult-short-videos/internal/service/video/repository"
	"context"
	"strings"

	errs "errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GossipService 吃瓜文章服务
type GossipService struct {
	gossipRepo repository.GossipRepository
	videoRepo  videoRepo.VideoRepository
}

func NewGossipService(gr repository.GossipRepository, vr videoRepo.VideoRepository) *GossipService {
	return &GossipService{gossipRepo: gr, videoRepo: vr}
}

func (s *GossipService) GetList(ctx context.Context, req *gossipDTO.GossipListReq) (*gossipDTO.GossipListResp, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	offset := (req.Page - 1) * req.PageSize
	posts, total, err := s.gossipRepo.List(ctx, offset, req.PageSize, req.Tag, req.Keyword)
	if err != nil {
		logger.Error("查询吃瓜列表失败", zap.Error(err))
		return nil, errors.New(errors.CodeDatabaseError, "查询失败")
	}

	items := make([]gossipDTO.GossipItem, 0, len(posts))
	for _, p := range posts {
		items = append(items, gossipDTO.GossipItem{
			ID:          p.ID,
			Title:       p.Title,
			Cover:       p.Cover,
			Summary:     p.Summary,
			Tags:        splitTags(p.Tags),
			ViewCount:   p.ViewCount,
			PublishedAt: p.PublishedAt.Unix(),
		})
	}

	return &gossipDTO.GossipListResp{
		Total: total,
		Page:  req.Page,
		Size:  req.PageSize,
		Posts: items,
	}, nil
}

func (s *GossipService) GetDetail(ctx context.Context, id int64) (*gossipDTO.GossipDetail, error) {
	post, err := s.gossipRepo.FindByID(ctx, id)
	if err != nil {
		if errs.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(errors.CodeNotFound, "文章不存在")
		}
		logger.Error("查询吃瓜详情失败", zap.Error(err), zap.Int64("id", id))
		return nil, errors.New(errors.CodeDatabaseError, "查询失败")
	}

	// 异步增加浏览数
	go func() {
		_ = s.gossipRepo.IncrementViewCount(context.Background(), id)
	}()

	// 取最新 6 条视频作为关联推荐（后续可按标签精准匹配）
	relatedVideos, _, err := s.videoRepo.List(ctx, 0, 6, map[string]interface{}{
		"order_by": "created_at",
	})
	if err != nil {
		logger.Warn("查询关联视频失败", zap.Error(err))
		relatedVideos = nil
	}

	related := make([]videoDTO.VideoItem, 0, len(relatedVideos))
	for _, v := range relatedVideos {
		related = append(related, videoDTO.VideoItem{
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
			CommentCount:  v.CommentCount,
			PublishedAt:   v.PublishedAt.Unix(),
		})
	}

	return &gossipDTO.GossipDetail{
		ID:            post.ID,
		Title:         post.Title,
		Cover:         post.Cover,
		Content:       post.Content,
		Tags:          splitTags(post.Tags),
		ViewCount:     post.ViewCount,
		PublishedAt:   post.PublishedAt.Unix(),
		RelatedVideos: related,
	}, nil
}

// splitTags 将逗号分隔的标签字符串拆成切片，过滤空串
func splitTags(raw string) []string {
	if raw == "" {
		return []string{}
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

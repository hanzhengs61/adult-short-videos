package logic

import (
	"context"

	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/pkg/logger"
	"adult-short-videos/internal/service/follow/model"
	"adult-short-videos/internal/service/follow/repository"
	videoDto "adult-short-videos/internal/service/video/dto"
	videoLogic "adult-short-videos/internal/service/video/logic"
	videoModel "adult-short-videos/internal/service/video/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// FollowService 关注服务
type FollowService struct {
	repo repository.FollowRepository
	db   *gorm.DB
}

func NewFollowService(repo repository.FollowRepository, db *gorm.DB) *FollowService {
	return &FollowService{repo: repo, db: db}
}

// Follow 关注作者
func (s *FollowService) Follow(ctx context.Context, userID int64, author string) error {
	if author == "" {
		return errors.New(errors.CodeInvalidParam, "作者名不能为空")
	}
	exists, err := s.repo.Exists(ctx, userID, author)
	if err != nil {
		logger.Error("检查关注状态失败", zap.Error(err))
		return errors.New(errors.CodeDatabaseError, "操作失败")
	}
	if exists {
		return errors.New(errors.CodeAlreadyFollowed, "已关注该作者")
	}
	if err := s.repo.Create(ctx, &model.AuthorFollow{UserId: userID, Author: author}); err != nil {
		logger.Error("关注失败", zap.Error(err))
		return errors.New(errors.CodeDatabaseError, "关注失败")
	}
	return nil
}

// Unfollow 取消关注
func (s *FollowService) Unfollow(ctx context.Context, userID int64, author string) error {
	if author == "" {
		return errors.New(errors.CodeInvalidParam, "作者名不能为空")
	}
	if err := s.repo.Delete(ctx, userID, author); err != nil {
		logger.Error("取消关注失败", zap.Error(err))
		return errors.New(errors.CodeDatabaseError, "取消关注失败")
	}
	return nil
}

// CheckFollow 检查是否已关注
func (s *FollowService) CheckFollow(ctx context.Context, userID int64, author string) (bool, error) {
	return s.repo.Exists(ctx, userID, author)
}

// ListFollowing 获取已关注作者列表
func (s *FollowService) ListFollowing(ctx context.Context, userID int64) ([]string, error) {
	return s.repo.ListAuthors(ctx, userID)
}

// AuthorProfileItem 返回给前端的关注作者信息
type AuthorProfileItem struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"` // 完整 URL，未注册作者为空字符串
}

// ListFollowingProfiles 返回带实时头像的关注列表
func (s *FollowService) ListFollowingProfiles(ctx context.Context, userID int64, baseURL string) ([]AuthorProfileItem, error) {
	profiles, err := s.repo.ListAuthorProfiles(ctx, userID)
	if err != nil {
		return nil, errors.New(errors.CodeDatabaseError, "查询失败")
	}
	items := make([]AuthorProfileItem, 0, len(profiles))
	for _, p := range profiles {
		avatar := ""
		if p.Avatar != "" {
			avatar = baseURL + p.Avatar
		}
		items = append(items, AuthorProfileItem{Name: p.Author, Avatar: avatar})
	}
	return items, nil
}

// GetFeed 获取已关注作者的视频列表
func (s *FollowService) GetFeed(ctx context.Context, userID int64, page, pageSize int) ([]videoDto.VideoItem, bool, error) {
	authors, err := s.repo.ListAuthors(ctx, userID)
	if err != nil {
		logger.Error("获取关注列表失败", zap.Error(err))
		return nil, false, errors.New(errors.CodeDatabaseError, "查询失败")
	}
	if len(authors) == 0 {
		return []videoDto.VideoItem{}, false, nil
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var videos []*videoModel.VideoWithAuthor
	err = s.db.WithContext(ctx).
		Table("videos v").
		Select(`v.*,
			COALESCE(u.username, v.author, '') AS author_name,
			COALESCE(u.avatar, '')             AS author_avatar`).
		Joins(`LEFT JOIN users u ON v.user_id > 0 AND u.user_id = v.user_id AND u.status = 1`).
		Where("v.author IN ? AND v.status = 1", authors).
		Order("v.published_at DESC").
		Offset(offset).Limit(pageSize).
		Scan(&videos).Error
	if err != nil {
		logger.Error("查询订阅视频失败", zap.Error(err))
		return nil, false, errors.New(errors.CodeDatabaseError, "查询失败")
	}

	items := make([]videoDto.VideoItem, 0, len(videos))
	for _, v := range videos {
		items = append(items, videoDto.VideoItem{
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
	return items, len(videos) >= pageSize, nil
}

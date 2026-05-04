package logic

import (
	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/pkg/logger"
	"adult-short-videos/internal/pkg/metrics"
	"adult-short-videos/internal/service/favorite/dto"
	"adult-short-videos/internal/service/favorite/model"
	"adult-short-videos/internal/service/favorite/repository"
	videoRepo "adult-short-videos/internal/service/video/repository"
	"context"
	errs "errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// FavoriteService 收藏服务
type FavoriteService struct {
	favoriteRepo repository.FavoriteRepository
	videoRepo    videoRepo.VideoRepository
	db           *gorm.DB
}

// NewFavoriteService 创建 FavoriteService 实例
func NewFavoriteService(favoriteRepo repository.FavoriteRepository, videoRepo videoRepo.VideoRepository, db *gorm.DB) *FavoriteService {
	return &FavoriteService{
		favoriteRepo: favoriteRepo,
		videoRepo:    videoRepo,
		db:           db,
	}
}

// AddFavorite 添加收藏
func (s *FavoriteService) AddFavorite(ctx context.Context, userId int64, req *dto.AddFavoriteReq) error {
	// 检查视频是否存在
	video, err := s.videoRepo.FindByID(ctx, req.VideoId)
	if err != nil {
		if errs.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(errors.CodeVideoNotFound, "视频不存在")
		}
		logger.Error("查询视频失败", zap.Error(err))
		return errors.New(errors.CodeDatabaseError, "数据库查询失败")
	}

	if video.Status != 1 {
		return errors.New(errors.CodeVideoOffline, "视频已下架")
	}

	// 检查是否已收藏
	exists, err := s.favoriteRepo.Exists(ctx, userId, req.VideoId)
	if err != nil {
		logger.Error("检查收藏状态失败", zap.Error(err))
		return errors.New(errors.CodeDatabaseError, "检查收藏状态失败")
	}
	if exists {
		logger.Warn("用户已经收藏过此视频")
		return errors.New(errors.CodeAlreadyFavorite, "已经收藏过了")
	}

	// 创建收藏记录
	favorite := &model.Favorite{
		UserId:  userId,
		VideoId: req.VideoId,
	}
	if err := s.favoriteRepo.Create(ctx, favorite); err != nil {
		logger.Error("创建收藏记录失败", zap.Error(err))
		return errors.New(errors.CodeDatabaseError, "收藏失败")
	}

	// 异步更新视频的收藏数
	go func() {
		s.db.Exec(`UPDATE videos SET favorite_count = favorite_count + 1 WHERE video_id = ?`, req.VideoId)
	}()

	metrics.FavoritesTotal.WithLabelValues("add").Inc()
	return nil
}

// RemoveFavorite 取消收藏
func (s *FavoriteService) RemoveFavorite(ctx context.Context, userId, videoId int64) error {
	exists, err := s.favoriteRepo.Exists(ctx, userId, videoId)
	if err != nil {
		logger.Error("检查收藏状态失败", zap.Error(err))
		return errors.New(errors.CodeDatabaseError, "检查收藏状态失败")
	}
	if !exists {
		logger.Warn("用户未收藏过此视频")
		return errors.New(errors.CodeNotFavorite, "未收藏过此视频")
	}

	if err := s.favoriteRepo.Delete(ctx, userId, videoId); err != nil {
		logger.Error("取消收藏失败", zap.Error(err))
		return errors.New(errors.CodeDatabaseError, "取消收藏失败")
	}

	go func() {
		s.db.Exec(`UPDATE videos SET favorite_count = GREATEST(favorite_count - 1, 0) WHERE video_id = ?`, videoId)
	}()

	metrics.FavoritesTotal.WithLabelValues("remove").Inc()
	return nil
}

// GetFavoriteList 获取收藏列表
func (s *FavoriteService) GetFavoriteList(ctx context.Context, userId int64, req *dto.FavoriteListReq) (*dto.FavoriteListResp, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Size < 1 {
		req.Size = 20
	}
	if req.Size > 100 {
		req.Size = 100
	}

	offset := (req.Page - 1) * req.Size
	favorites, total, err := s.favoriteRepo.ListWithVideo(ctx, userId, offset, req.Size)
	if err != nil {
		logger.Error("查询收藏列表失败", zap.Error(err))
		return nil, errors.New(errors.CodeDatabaseError, "查询收藏列表失败")
	}

	return &dto.FavoriteListResp{
		Total:     total,
		Page:      req.Page,
		Size:      req.Size,
		Favorites: favorites,
	}, nil
}

// CheckFavorite 检查是否收藏
func (s *FavoriteService) CheckFavorite(ctx context.Context, userId, videoId int64) (bool, error) {
	return s.favoriteRepo.Exists(ctx, userId, videoId)
}

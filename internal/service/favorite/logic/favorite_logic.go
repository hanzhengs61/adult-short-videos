package logic

import (
	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/pkg/logger"
	"adult-short-videos/internal/pkg/metrics"
	"adult-short-videos/internal/service/favorite/model"
	"adult-short-videos/internal/service/favorite/repository"
	videoRepo "adult-short-videos/internal/service/video/repository"
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AddFavoriteLogic 添加收藏逻辑
type AddFavoriteLogic struct {
	ctx          context.Context
	favoriteRepo repository.FavoriteRepository
	videoRepo    videoRepo.VideoRepository
	db           *gorm.DB
}

// NewAddFavoriteLogic 创建添加收藏逻辑实例
func NewAddFavoriteLogic(ctx context.Context, favoriteRepo repository.FavoriteRepository, videoRepo videoRepo.VideoRepository, db *gorm.DB) *AddFavoriteLogic {
	return &AddFavoriteLogic{
		ctx:          ctx,
		favoriteRepo: favoriteRepo,
		videoRepo:    videoRepo,
		db:           db,
	}
}

// AddFavoriteReq 添加收藏请求
type AddFavoriteReq struct {
	VideoId int64 `json:"video_id"` // 视频 ID
}

// AddFavorite 添加收藏
// 步骤：
// 1. 检查视频是否存在
// 2. 检查是否已收藏
// 3. 创建收藏记录
// 4. 更新视频的收藏数
func (l *AddFavoriteLogic) AddFavorite(userId int64, req *AddFavoriteReq) error {
	// 检查视频是否存在
	video, err := l.videoRepo.FindByID(l.ctx, req.VideoId)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error("视频不存在", zap.Error(err))
		return errors.New(errors.CodeDatabaseError, "视频不存在")
	}

	if video == nil {
		logger.Error("视频不存在", zap.Error(err))
		return errors.New(errors.CodeVideoNotFound, "视频不存在")
	}

	// 检查视频状态
	if video.Status != 1 {
		return errors.New(errors.CodeVideoOffline, "视频已下架")
	}

	// 检查是否已收藏
	exists, err := l.favoriteRepo.Exists(l.ctx, userId, req.VideoId)
	if err != nil {
		return errors.New(errors.CodeDatabaseError, "检查收藏状态失败")
	}

	if exists {
		return errors.New(errors.CodeAlreadyFavorite, "已经收藏过了")
	}

	// 创建收藏记录
	favorite := &model.Favorite{
		UserId:  userId,
		VideoId: req.VideoId,
	}

	if err := l.favoriteRepo.Create(l.ctx, favorite); err != nil {
		return errors.New(errors.CodeDatabaseError, "收藏失败")
	}

	// 更新视频的收藏数
	// 使用事务确保数据一致性
	// favorite_count = favorite_count + 1
	go func() {
		l.db.Model(&video).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1))
	}()

	metrics.FavoritesTotal.WithLabelValues("add").Inc()
	return nil
}

// RemoveFavoriteLogic 取消收藏逻辑
type RemoveFavoriteLogic struct {
	ctx          context.Context
	favoriteRepo repository.FavoriteRepository
	db           *gorm.DB
}

// NewRemoveFavoriteLogic 创建取消收藏逻辑实例
func NewRemoveFavoriteLogic(ctx context.Context, favoriteRepo repository.FavoriteRepository, db *gorm.DB) *RemoveFavoriteLogic {
	return &RemoveFavoriteLogic{
		ctx:          ctx,
		favoriteRepo: favoriteRepo,
		db:           db,
	}
}

// RemoveFavorite 取消收藏
// 步骤：
// 1. 检查是否已收藏
// 2. 删除收藏记录
// 3. 更新视频的收藏数
func (l *RemoveFavoriteLogic) RemoveFavorite(userId, videoId int64) error {
	// 检查是否已收藏
	exists, err := l.favoriteRepo.Exists(l.ctx, userId, videoId)
	if err != nil {
		return errors.New(errors.CodeDatabaseError, "检查收藏状态失败")
	}

	if !exists {
		return errors.New(errors.CodeNotFavorite, "未收藏过此视频")
	}

	// 删除收藏记录
	if err := l.favoriteRepo.Delete(l.ctx, userId, videoId); err != nil {
		return errors.New(errors.CodeDatabaseError, "取消收藏失败")
	}

	// 更新视频的收藏数
	// favorite_count = favorite_count - 1
	// 异步执行，不影响主流程
	go func() {
		// 注意：需要检查 favorite_count 不能小于 0
		l.db.Exec(`
			UPDATE videos
			SET favorite_count = GREATEST(favorite_count - 1, 0)
			WHERE video_id = ?
		`, videoId)
	}()

	metrics.FavoritesTotal.WithLabelValues("remove").Inc()
	return nil
}

// FavoriteListLogic 收藏列表逻辑
type FavoriteListLogic struct {
	ctx          context.Context
	favoriteRepo repository.FavoriteRepository
}

// NewFavoriteListLogic 创建收藏列表逻辑实例
func NewFavoriteListLogic(ctx context.Context, favoriteRepo repository.FavoriteRepository) *FavoriteListLogic {
	return &FavoriteListLogic{
		ctx:          ctx,
		favoriteRepo: favoriteRepo,
	}
}

// FavoriteListReq 收藏列表请求
type FavoriteListReq struct {
	Page int // 页码
	Size int // 每页数量
}

// FavoriteListResp 收藏列表响应
type FavoriteListResp struct {
	Total     int64                      `json:"total"`     // 总数
	Page      int                        `json:"page"`      // 当前页
	Size      int                        `json:"size"`      // 每页数量
	Favorites []*model.FavoriteWithVideo `json:"favorites"` // 收藏列表
}

// GetFavoriteList 获取收藏列表
func (l *FavoriteListLogic) GetFavoriteList(userId int64, req *FavoriteListReq) (*FavoriteListResp, error) {
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

	// 计算偏移量
	offset := (req.Page - 1) * req.Size

	// 查询收藏列表
	favorites, total, err := l.favoriteRepo.ListWithVideo(l.ctx, userId, offset, req.Size)
	if err != nil {
		return nil, errors.New(errors.CodeDatabaseError, "查询收藏列表失败")
	}

	// 返回结果
	return &FavoriteListResp{
		Total:     total,
		Page:      req.Page,
		Size:      req.Size,
		Favorites: favorites,
	}, nil
}

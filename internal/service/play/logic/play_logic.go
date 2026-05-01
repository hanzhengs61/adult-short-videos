package logic

import (
	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/pkg/metrics"
	"adult-short-videos/internal/service/play/dto"
	"adult-short-videos/internal/service/play/model"
	"adult-short-videos/internal/service/play/repository"
	videoRepo "adult-short-videos/internal/service/video/repository"
	"context"
	"time"

	"gorm.io/gorm"
)

// PlayService 播放历史服务
type PlayService struct {
	playRepo  repository.PlayHistoryRepository
	videoRepo videoRepo.VideoRepository
	db        *gorm.DB
}

// NewPlayService 创建 PlayService 实例
func NewPlayService(playRepo repository.PlayHistoryRepository, videoRepo videoRepo.VideoRepository, db *gorm.DB) *PlayService {
	return &PlayService{
		playRepo:  playRepo,
		videoRepo: videoRepo,
		db:        db,
	}
}

// RecordPlay 记录播放历史
func (s *PlayService) RecordPlay(ctx context.Context, userId int64, req *dto.RecordPlayReq) error {
	// 检查视频是否存在
	_, err := s.videoRepo.FindByID(ctx, req.VideoId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeVideoNotFound, "视频不存在")
		}
		return errors.New(errors.CodeDatabaseError, "数据库查询失败")
	}

	// 参数校验
	if req.PlayProgress < 0 || req.PlayProgress > 100 {
		req.PlayProgress = 0
	}

	// 创建或更新播放历史
	history := &model.PlayHistory{
		UserId:       userId,
		VideoId:      req.VideoId,
		PlayDuration: req.PlayDuration,
		PlayProgress: req.PlayProgress,
	}

	if err := s.playRepo.CreateOrUpdate(ctx, history); err != nil {
		return errors.New(errors.CodeDatabaseError, "记录播放历史失败")
	}

	// 异步更新视频播放次数
	go func() {
		_ = s.videoRepo.IncrementPlayCount(context.Background(), req.VideoId)
	}()

	// 异步更新用户统计
	go func() {
		s.db.Exec(`
			UPDATE user_statistics
			SET play_count = play_count + 1,
				total_watch_time = total_watch_time + ?,
				last_active_at = ?
			WHERE user_id = ?
		`, req.PlayDuration, time.Now(), userId)
	}()

	metrics.VideoPlaysTotal.Inc()
	return nil
}

// GetPlayHistory 获取播放历史列表
func (s *PlayService) GetPlayHistory(ctx context.Context, userId int64, req *dto.PlayHistoryListReq) (*dto.PlayHistoryListResp, error) {
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
	history, total, err := s.playRepo.ListWithVideo(ctx, userId, offset, req.Size)
	if err != nil {
		return nil, errors.New(errors.CodeDatabaseError, "查询播放历史失败")
	}

	return &dto.PlayHistoryListResp{
		Total:   total,
		Page:    req.Page,
		Size:    req.Size,
		History: history,
	}, nil
}

// ClearPlayHistory 清空播放历史
func (s *PlayService) ClearPlayHistory(ctx context.Context, i int64) error {
	err := s.playRepo.DeleteAll(ctx, i)
	if err != nil {
		return errors.New(errors.CodeDatabaseError, "清空播放历史失败")
	}
	return nil
}

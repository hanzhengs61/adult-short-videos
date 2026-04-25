package logic

import (
	"adult-short-videos/internal/pkg/metrics"
	"adult-short-videos/internal/service/play/model"
	"adult-short-videos/internal/service/play/repository"
	videoRepo "adult-short-videos/internal/service/video/repository"
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

// RecordPlayLogic 记录播放逻辑
type RecordPlayLogic struct {
	ctx       context.Context
	playRepo  repository.PlayHistoryRepository
	videoRepo videoRepo.VideoRepository
	db        *gorm.DB
}

func NewRecordPlayLogic(ctx context.Context, playRepo repository.PlayHistoryRepository, videoRepo videoRepo.VideoRepository, db *gorm.DB) *RecordPlayLogic {
	return &RecordPlayLogic{
		ctx:       ctx,
		playRepo:  playRepo,
		videoRepo: videoRepo,
		db:        db,
	}
}

// RecordPlayReq 记录播放请求
type RecordPlayReq struct {
	VideoId      int64   `json:"video_id"`      // 视频 ID
	PlayDuration int32   `json:"play_duration"` // 观看时长（秒）
	PlayProgress float32 `json:"play_progress"` // 播放进度（0-100）
}

// RecordPlay 记录播放历史
func (l *RecordPlayLogic) RecordPlay(userId int64, req *RecordPlayReq) error {
	// 检查视频是否存在
	_, err := l.videoRepo.FindByID(l.ctx, req.VideoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("视频不存在")
		}
		return errors.New("查询视频失败")
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

	if err := l.playRepo.CreateOrUpdate(l.ctx, history); err != nil {
		return errors.New("记录播放历史失败")
	}

	// 异步更新视频播放次数
	go func() {
		err := l.videoRepo.IncrementPlayCount(context.Background(), req.VideoId)
		if err != nil {
			log.Println("更新视频播放次数失败:", err)
		}
	}()

	// 异步更新用户统计
	go func() {
		l.db.Exec(`
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

// PlayHistoryListLogic 播放历史列表逻辑
type PlayHistoryListLogic struct {
	ctx      context.Context
	playRepo repository.PlayHistoryRepository
}

func NewPlayHistoryListLogic(ctx context.Context, playRepo repository.PlayHistoryRepository) *PlayHistoryListLogic {
	return &PlayHistoryListLogic{
		ctx:      ctx,
		playRepo: playRepo,
	}
}

// PlayHistoryListReq 播放历史列表请求
type PlayHistoryListReq struct {
	Page int
	Size int
}

// PlayHistoryListResp 播放历史列表响应
type PlayHistoryListResp struct {
	Total   int64                         `json:"total"`
	Page    int                           `json:"page"`
	Size    int                           `json:"size"`
	History []*model.PlayHistoryWithVideo `json:"history"`
}

// GetPlayHistoryList 获取播放历史列表
func (l *PlayHistoryListLogic) GetPlayHistoryList(userId int64, req *PlayHistoryListReq) (*PlayHistoryListResp, error) {
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

	offset := (req.Page - 1) * req.Size

	// 查询播放历史
	history, total, err := l.playRepo.ListWithVideo(l.ctx, userId, offset, req.Size)
	if err != nil {
		return nil, errors.New("查询播放历史失败")
	}

	return &PlayHistoryListResp{
		Total:   total,
		Page:    req.Page,
		Size:    req.Size,
		History: history,
	}, nil
}

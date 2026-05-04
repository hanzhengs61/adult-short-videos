package logic

import (
	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/pkg/logger"
	"adult-short-videos/internal/pkg/metrics"
	"context"
	errs "errors"
	"strings"

	"adult-short-videos/internal/service/comment/dto"
	"adult-short-videos/internal/service/comment/model"
	"adult-short-videos/internal/service/comment/repository"
	videoModel "adult-short-videos/internal/service/video/model"
	videoRepo "adult-short-videos/internal/service/video/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CommentService 评论服务
type CommentService struct {
	commentRepo repository.CommentRepository
	videoRepo   videoRepo.VideoRepository
	db          *gorm.DB
}

// NewCommentService 创建 CommentService 实例
func NewCommentService(commentRepo repository.CommentRepository, videoRepo videoRepo.VideoRepository, db *gorm.DB) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		videoRepo:   videoRepo,
		db:          db,
	}
}

// AddComment 发表评论
func (s *CommentService) AddComment(ctx context.Context, userId int64, req *dto.AddCommentReq) (*model.Comment, error) {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		logger.Error("评论内容不能为空")
		return nil, errors.New(errors.CodeCommentEmpty, "评论内容不能为空")
	}

	if len(content) > 500 {
		logger.Error("评论内容不能超过500字")
		return nil, errors.New(errors.CodeCommentTooLong, "评论内容不能超过500字")
	}

	// 检查视频是否存在
	_, err := s.videoRepo.FindByID(ctx, req.VideoId)
	if err != nil && !errs.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("查询视频失败", zap.Error(err))
		return nil, errors.New(errors.CodeDatabaseError, "数据库查询失败")
	}

	// 创建评论
	comment := &model.Comment{
		VideoId:  req.VideoId,
		UserId:   userId,
		Content:  content,
		ParentId: req.ParentId,
		Status:   1, // 正常状态
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		logger.Error("发表评论失败", zap.Error(err))
		return nil, errors.New(errors.CodeCommentFailed, "发表评论失败")
	}

	// 异步更新视频评论数
	go func() {
		s.db.Model(&videoModel.Video{}).
			Where("video_id = ?", req.VideoId).
			UpdateColumn("comment_count", gorm.Expr("comment_count + 1"))
	}()

	metrics.CommentsTotal.WithLabelValues("add").Inc()
	return comment, nil
}

// GetCommentList 获取评论列表
// currentUserId: 当前登录用户ID，如果未登录则为0
func (s *CommentService) GetCommentList(ctx context.Context, req *dto.CommentListReq, currentUserId int64) (*dto.CommentListResp, error) {
	// 参数验证
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Size < 1 {
		req.Size = 10
	}
	if req.Size > 100 {
		req.Size = 100
	}

	offset := (req.Page - 1) * req.Size

	// 查询评论列表
	comments, total, err := s.commentRepo.ListWithUser(ctx, req.VideoId, currentUserId, offset, req.Size)
	if err != nil {
		return nil, errors.New(errors.CodeDatabaseError, "查询评论列表失败")
	}

	return &dto.CommentListResp{
		Total:    total,
		Page:     req.Page,
		Size:     req.Size,
		Comments: comments,
	}, nil
}

// LikeComment 点赞评论
func (s *CommentService) LikeComment(ctx context.Context, userId, commentId int64) error {
	// 检查是否已点赞
	isLiked, err := s.commentRepo.IsLiked(ctx, commentId, userId)
	if err != nil {
		return errors.New(errors.CodeDatabaseError, "检查点赞状态失败")
	}

	if isLiked {
		return errors.New(errors.CodeAlreadyLiked, "已经点赞过了")
	}

	// 点赞
	if err := s.commentRepo.LikeComment(ctx, commentId, userId); err != nil {
		return errors.New(errors.CodeDatabaseError, "点赞失败")
	}

	// 增加点赞数
	go func() {
		_ = s.commentRepo.IncrementLikeCount(context.Background(), commentId)
	}()

	metrics.CommentsTotal.WithLabelValues("like").Inc()
	return nil
}

// UnlikeComment 取消点赞
func (s *CommentService) UnlikeComment(ctx context.Context, userId, commentId int64) error {
	// 检查是否已点赞
	isLiked, err := s.commentRepo.IsLiked(ctx, commentId, userId)
	if err != nil {
		return errors.New(errors.CodeDatabaseError, "检查点赞状态失败")
	}

	if !isLiked {
		return errors.New(errors.CodeNotLiked, "未点赞过该评论")
	}

	// 取消点赞
	if err := s.commentRepo.UnlikeComment(ctx, commentId, userId); err != nil {
		return errors.New(errors.CodeDatabaseError, "取消点赞失败")
	}

	// 减少点赞数
	go func() {
		_ = s.commentRepo.DecrementLikeCount(context.Background(), commentId)
	}()

	metrics.CommentsTotal.WithLabelValues("unlike").Inc()
	return nil
}

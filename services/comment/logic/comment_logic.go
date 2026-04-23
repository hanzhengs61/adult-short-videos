package logic

import (
	"context"
	"errors"
	"strings"

	"adult-short-videos/services/comment/model"
	"adult-short-videos/services/comment/repository"
	videoRepo "adult-short-videos/services/video/repository"

	"gorm.io/gorm"
)

// AddCommentLogic 发表评论逻辑
type AddCommentLogic struct {
	ctx         context.Context
	commentRepo repository.CommentRepository
	videoRepo   videoRepo.VideoRepository
	db          *gorm.DB
}

func NewAddCommentLogic(ctx context.Context, commentRepo repository.CommentRepository, videoRepo videoRepo.VideoRepository, db *gorm.DB) *AddCommentLogic {
	return &AddCommentLogic{
		ctx:         ctx,
		commentRepo: commentRepo,
		videoRepo:   videoRepo,
		db:          db,
	}
}

// AddCommentReq 发表评论请求
type AddCommentReq struct {
	VideoId  int64  `json:"video_id"`  // 视频 ID
	Content  string `json:"content"`   // 评论内容
	ParentId int64  `json:"parent_id"` // 父评论ID（0表示顶级评论）
}

// AddComment 发表评论
func (l *AddCommentLogic) AddComment(userId int64, req *AddCommentReq) (*model.Comment, error) {

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errors.New("评论内容不能为空")
	}

	if len(content) > 500 {
		return nil, errors.New("评论内容不能超过500字")
	}

	// 检查视频是否存在
	_, err := l.videoRepo.FindByID(l.ctx, req.VideoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("视频不存在")
		}
		return nil, errors.New("查询视频失败")
	}

	// 创建评论
	comment := &model.Comment{
		VideoId:  req.VideoId,
		UserId:   userId,
		Content:  content,
		ParentId: req.ParentId,
		Status:   1, // 正常状态
	}

	if err := l.commentRepo.Create(l.ctx, comment); err != nil {
		return nil, errors.New("发表评论失败")
	}

	// 异步更新视频评论数
	go func() {
		l.db.Model(&model.Comment{}).
			Where("video_id = ?", req.VideoId).
			UpdateColumn("comment_count", gorm.Expr("comment_count + 1"))
	}()

	return comment, nil
}

// CommentListLogic 评论列表逻辑
type CommentListLogic struct {
	ctx         context.Context
	commentRepo repository.CommentRepository
}

func NewCommentListLogic(ctx context.Context, commentRepo repository.CommentRepository) *CommentListLogic {
	return &CommentListLogic{
		ctx:         ctx,
		commentRepo: commentRepo,
	}
}

// CommentListReq 评论列表请求
type CommentListReq struct {
	VideoId int64 // 视频 ID
	Page    int
	Size    int
}

// CommentListResp 评论列表响应
type CommentListResp struct {
	Total    int64                    `json:"total"`
	Page     int                      `json:"page"`
	Size     int                      `json:"size"`
	Comments []*model.CommentWithUser `json:"comments"`
}

// GetCommentList 获取评论列表
func (l *CommentListLogic) GetCommentList(req *CommentListReq) (*CommentListResp, error) {
	// 参数验证
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

	// 查询评论列表
	comments, total, err := l.commentRepo.ListWithUser(l.ctx, req.VideoId, offset, req.Size)
	if err != nil {
		return nil, errors.New("查询评论列表失败")
	}

	return &CommentListResp{
		Total:    total,
		Page:     req.Page,
		Size:     req.Size,
		Comments: comments,
	}, nil
}

// LikeCommentLogic 点赞评论逻辑
type LikeCommentLogic struct {
	ctx         context.Context
	commentRepo repository.CommentRepository
}

func NewLikeCommentLogic(ctx context.Context, commentRepo repository.CommentRepository) *LikeCommentLogic {
	return &LikeCommentLogic{
		ctx:         ctx,
		commentRepo: commentRepo,
	}
}

// LikeComment 点赞评论
func (l *LikeCommentLogic) LikeComment(userId, commentId int64) error {
	// 检查是否已点赞
	isLiked, err := l.commentRepo.IsLiked(l.ctx, commentId, userId)
	if err != nil {
		return errors.New("检查点赞状态失败")
	}

	if isLiked {
		return errors.New("已经点赞过了")
	}

	// 点赞
	if err := l.commentRepo.LikeComment(l.ctx, commentId, userId); err != nil {
		return errors.New("点赞失败")
	}

	// 增加点赞数
	go func() {
		l.commentRepo.IncrementLikeCount(context.Background(), commentId)
	}()

	return nil
}

// UnlikeComment 取消点赞
func (l *LikeCommentLogic) UnlikeComment(userId, commentId int64) error {
	// 检查是否已点赞
	isLiked, err := l.commentRepo.IsLiked(l.ctx, commentId, userId)
	if err != nil {
		return errors.New("检查点赞状态失败")
	}

	if !isLiked {
		return errors.New("未点赞过该评论")
	}

	// 取消点赞
	if err := l.commentRepo.UnlikeComment(l.ctx, commentId, userId); err != nil {
		return errors.New("取消点赞失败")
	}

	// 减少点赞数
	go func() {
		l.commentRepo.DecrementLikeCount(context.Background(), commentId)
	}()

	return nil
}

package repository

import (
	"adult-short-videos/internal/service/comment/model"
	"context"

	"gorm.io/gorm"
)

// CommentRepository 评论仓储接口
type CommentRepository interface {
	// Create 创建评论
	Create(ctx context.Context, comment *model.Comment) error

	// ListWithUser 获取视频的评论列表（包含用户信息）
	ListWithUser(ctx context.Context, videoId int64, offset, limit int) ([]*model.CommentWithUser, int64, error)

	// GetReplies 获取评论的回复列表
	GetReplies(ctx context.Context, parentId int64, offset, limit int) ([]*model.CommentWithUser, error)

	// Delete 删除评论
	Delete(ctx context.Context, commentId, userId int64) error

	// LikeComment 点赞评论
	LikeComment(ctx context.Context, commentId, userId int64) error

	// UnlikeComment 取消点赞
	UnlikeComment(ctx context.Context, commentId, userId int64) error

	// IsLiked 检查是否已点赞
	IsLiked(ctx context.Context, commentId, userId int64) (bool, error)

	// IncrementLikeCount 增加点赞数
	IncrementLikeCount(ctx context.Context, commentId int64) error

	// DecrementLikeCount 减少点赞数
	DecrementLikeCount(ctx context.Context, commentId int64) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

// Create 创建评论
func (r *commentRepository) Create(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

// ListWithUser 获取视频的评论列表
func (r *commentRepository) ListWithUser(ctx context.Context, videoId int64, offset, limit int) ([]*model.CommentWithUser, int64, error) {
	var results []*model.CommentWithUser
	var total int64

	// 获取总数（只统计顶级评论）
	err := r.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("video_id = ? AND parent_id = 0 AND status = 1", videoId).
		Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	// 联表查询评论和用户信息
	err = r.db.WithContext(ctx).
		Table("comments c").
		Select(`
			c.comment_id,
			c.video_id,
			c.user_id,
			u.username,
			u.avatar,
			c.content,
			c.parent_id,
			c.like_count,
			c.status,
			EXTRACT(EPOCH FROM c.created_at)::bigint as created_at
		`).
		Joins("INNER JOIN users u ON c.user_id = u.user_id").
		Where("c.video_id = ? AND c.parent_id = 0 AND c.status = 1", videoId).
		Order("c.created_at DESC").
		Offset(offset).
		Limit(limit).
		Scan(&results).Error

	return results, total, err
}

// GetReplies 获取评论的回复列表
func (r *commentRepository) GetReplies(ctx context.Context, parentId int64, offset, limit int) ([]*model.CommentWithUser, error) {
	var results []*model.CommentWithUser

	err := r.db.WithContext(ctx).
		Table("comments c").
		Select(`
			c.comment_id,
			c.video_id,
			c.user_id,
			u.username,
			u.avatar,
			c.content,
			c.parent_id,
			c.like_count,
			c.status,
			EXTRACT(EPOCH FROM c.created_at)::bigint as created_at
		`).
		Joins("INNER JOIN users u ON c.user_id = u.user_id").
		Where("c.parent_id = ? AND c.status = 1", parentId).
		Order("c.created_at ASC"). // 回复按时间正序
		Offset(offset).
		Limit(limit).
		Scan(&results).Error

	return results, err
}

// Delete 删除评论（软删除）
func (r *commentRepository) Delete(ctx context.Context, commentId, userId int64) error {
	// 只能删除自己的评论
	return r.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("comment_id = ? AND user_id = ?", commentId, userId).
		Update("status", 0).Error
}

// LikeComment 点赞评论
func (r *commentRepository) LikeComment(ctx context.Context, commentId, userId int64) error {
	like := &model.CommentLike{
		CommentId: commentId,
		UserId:    userId,
	}
	return r.db.WithContext(ctx).Create(like).Error
}

// UnlikeComment 取消点赞
func (r *commentRepository) UnlikeComment(ctx context.Context, commentId, userId int64) error {
	return r.db.WithContext(ctx).
		Where("comment_id = ? AND user_id = ?", commentId, userId).
		Delete(&model.CommentLike{}).Error
}

// IsLiked 检查是否已点赞
func (r *commentRepository) IsLiked(ctx context.Context, commentId, userId int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.CommentLike{}).
		Where("comment_id = ? AND user_id = ?", commentId, userId).
		Count(&count).Error

	return count > 0, err
}

// IncrementLikeCount 增加点赞数
func (r *commentRepository) IncrementLikeCount(ctx context.Context, commentId int64) error {
	return r.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("comment_id = ?", commentId).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

// DecrementLikeCount 减少点赞数
func (r *commentRepository) DecrementLikeCount(ctx context.Context, commentId int64) error {
	return r.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("comment_id = ?", commentId).
		UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error
}

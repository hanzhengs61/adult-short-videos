package handler

import (
	"strconv"

	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/pkg/response"
	"adult-short-videos/internal/service/comment/logic"
	"adult-short-videos/internal/service/comment/repository"
	videoRepo "adult-short-videos/internal/service/video/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CommentHandler 评论处理器
type CommentHandler struct {
	db          *gorm.DB
	commentRepo repository.CommentRepository
	videoRepo   videoRepo.VideoRepository
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{
		db:          db,
		commentRepo: repository.NewCommentRepository(db),
		videoRepo:   videoRepo.NewVideoRepository(db),
	}
}

// AddComment 发表评论
// 路由: POST /api/comment/add
func (h *CommentHandler) AddComment(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req logic.AddCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	l := logic.NewAddCommentLogic(c.Request.Context(), h.commentRepo, h.videoRepo, h.db)
	comment, err := l.AddComment(userId.(int64), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMsg(c, "评论成功", map[string]interface{}{
		"comment_id": comment.CommentId,
	})
}

// GetCommentList 获取评论列表
// 路由: GET /api/comment/list
func (h *CommentHandler) GetCommentList(c *gin.Context) {
	videoIdStr := c.Query("video_id")
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "20")

	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	req := &logic.CommentListReq{
		VideoId: videoId,
		Page:    page,
		Size:    size,
	}

	l := logic.NewCommentListLogic(c.Request.Context(), h.commentRepo)
	resp, err := l.GetCommentList(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

// LikeComment 点赞评论
// 路由: POST /api/comment/like
func (h *CommentHandler) LikeComment(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var body struct {
		CommentId int64 `json:"comment_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.CommentId <= 0 {
		response.Error(c, errors.CodeInvalidParam, "评论 ID 无效")
		return
	}

	l := logic.NewLikeCommentLogic(c.Request.Context(), h.commentRepo)
	if err := l.LikeComment(userId.(int64), body.CommentId); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMsg(c, "点赞成功", nil)
}

// DeleteComment 删除评论
// 路由: DELETE /api/comment/delete/:id
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	commentId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || commentId <= 0 {
		response.Error(c, errors.CodeInvalidParam, "评论 ID 无效")
		return
	}

	if err := h.commentRepo.Delete(c.Request.Context(), commentId, userId.(int64)); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}

// UnlikeComment 取消点赞
// 路由: DELETE /api/comment/like/:id
func (h *CommentHandler) UnlikeComment(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	commentIdStr := c.Param("id")
	commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	l := logic.NewLikeCommentLogic(c.Request.Context(), h.commentRepo)
	if err := l.UnlikeComment(userId.(int64), commentId); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMsg(c, "已取消点赞", nil)
}

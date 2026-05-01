package handler

import (
	"strconv"

	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/pkg/response"
	"adult-short-videos/internal/service/comment/dto"
	"adult-short-videos/internal/service/comment/logic"
	"adult-short-videos/internal/service/comment/repository"
	videoRepo "adult-short-videos/internal/service/video/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CommentHandler 评论处理器
type CommentHandler struct {
	commentService *logic.CommentService
	commentRepo    repository.CommentRepository
}

// NewCommentHandler 创建 CommentHandler 评论处理器
func NewCommentHandler(db *gorm.DB) *CommentHandler {
	commentRepo := repository.NewCommentRepository(db)
	vRepo := videoRepo.NewVideoRepository(db)
	return &CommentHandler{
		commentRepo:    commentRepo,
		commentService: logic.NewCommentService(commentRepo, vRepo, db),
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

	var req dto.AddCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	comment, err := h.commentService.AddComment(c.Request.Context(), userId.(int64), &req)
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

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	req := &dto.CommentListReq{
		VideoId: videoId,
		Page:    page,
		Size:    size,
	}

	// 获取当前用户ID，如果未登录则为0
	var currentUserId int64
	if userIdVal, exists := c.Get("user_id"); exists {
		currentUserId = userIdVal.(int64)
	}

	resp, err := h.commentService.GetCommentList(c.Request.Context(), req, currentUserId)
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

	if err := h.commentService.LikeComment(c.Request.Context(), userId.(int64), body.CommentId); err != nil {
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

	// 注意：这里仍然直接使用了 commentRepo，如果 Delete 逻辑也需要迁移到 service 层，则需要进一步重构
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

	if err := h.commentService.UnlikeComment(c.Request.Context(), userId.(int64), commentId); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMsg(c, "已取消点赞", nil)
}

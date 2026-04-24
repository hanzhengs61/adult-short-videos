package handler

import (
	"strconv"

	"adult-short-videos/common/response"
	"adult-short-videos/services/play/logic"
	"adult-short-videos/services/play/repository"
	videoRepo "adult-short-videos/services/video/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PlayHandler 播放处理器
type PlayHandler struct {
	db        *gorm.DB
	playRepo  repository.PlayHistoryRepository
	videoRepo videoRepo.VideoRepository
}

func NewPlayHandler(db *gorm.DB) *PlayHandler {
	return &PlayHandler{
		db:        db,
		playRepo:  repository.NewPlayHistoryRepository(db),
		videoRepo: videoRepo.NewVideoRepository(db),
	}
}

// RecordPlay 记录播放历史
// 路由: POST /api/play/record
func (h *PlayHandler) RecordPlay(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req logic.RecordPlayReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	l := logic.NewRecordPlayLogic(c.Request.Context(), h.playRepo, h.videoRepo, h.db)
	if err := l.RecordPlay(userId.(int64), &req); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMsg(c, "记录成功", nil)
}

// GetPlayHistory 获取播放历史
// 路由: GET /api/play/history
func (h *PlayHandler) GetPlayHistory(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "20")

	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	req := &logic.PlayHistoryListReq{
		Page: page,
		Size: size,
	}

	l := logic.NewPlayHistoryListLogic(c.Request.Context(), h.playRepo)
	resp, err := l.GetPlayHistoryList(userId.(int64), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

// DeletePlayHistory 删除播放历史
// 路由: DELETE /api/play/history/:videoId
func (h *PlayHandler) DeletePlayHistory(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	videoIdStr := c.Param("videoId")
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	if err := h.playRepo.Delete(c.Request.Context(), userId.(int64), videoId); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMsg(c, "已删除", nil)
}

// ClearPlayHistory 清空播放历史
// 路由: DELETE /api/play/history
func (h *PlayHandler) ClearPlayHistory(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	if err := h.playRepo.DeleteAll(c.Request.Context(), userId.(int64)); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMsg(c, "已清空播放历史", nil)
}

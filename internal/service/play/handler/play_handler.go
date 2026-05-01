package handler

import (
	"adult-short-videos/internal/pkg/response"
	"adult-short-videos/internal/service/play/dto"
	"adult-short-videos/internal/service/play/logic"
	"adult-short-videos/internal/service/play/repository"
	videoRepo "adult-short-videos/internal/service/video/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PlayHandler 播放历史处理器
type PlayHandler struct {
	playService *logic.PlayService
}

func NewPlayHandler(db *gorm.DB) *PlayHandler {
	pRepo := repository.NewPlayHistoryRepository(db)
	vRepo := videoRepo.NewVideoRepository(db)
	return &PlayHandler{
		playService: logic.NewPlayService(pRepo, vRepo, db),
	}
}

// RecordPlay 记录播放
// 路由: POST /api/play/record
func (h *PlayHandler) RecordPlay(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.RecordPlayReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	if err := h.playService.RecordPlay(c.Request.Context(), userId.(int64), &req); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMsg(c, "播放历史已记录", nil)
}

// GetPlayHistory 获取播放历史
// 路由: GET /api/play/history
func (h *PlayHandler) GetPlayHistory(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.PlayHistoryListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	resp, err := h.playService.GetPlayHistory(c.Request.Context(), userId.(int64), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

// ClearPlayHistory 清空播放历史
// 路由: DELETE /api/play/history
func (h *PlayHandler) ClearPlayHistory(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	err := h.playService.ClearPlayHistory(c.Request.Context(), userId.(int64))

	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMsg(c, "已清空播放历史", nil)
}

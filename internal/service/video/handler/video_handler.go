package handler

import (
	"strconv"

	"adult-short-videos/internal/pkg/response"
	"adult-short-videos/internal/service/video/dto"
	"adult-short-videos/internal/service/video/logic"
	"adult-short-videos/internal/service/video/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// VideoHandler 视频处理器
type VideoHandler struct {
	videoService *logic.VideoService
}

func NewVideoHandler(db *gorm.DB) *VideoHandler {
	vRepo := repository.NewVideoRepository(db)
	return &VideoHandler{
		videoService: logic.NewVideoService(vRepo),
	}
}

// GetVideoList 获取视频列表
// 路由: GET /api/video/list
func (h *VideoHandler) GetVideoList(c *gin.Context) {
	var req dto.VideoListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	resp, err := h.videoService.GetVideoList(c.Request.Context(), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

// GetVideoDetail 获取视频详情
// 路由: GET /api/video/detail/:id
func (h *VideoHandler) GetVideoDetail(c *gin.Context) {
	videoIdStr := c.Param("id")
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	resp, err := h.videoService.GetVideoDetail(c.Request.Context(), videoId)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

// GetHotVideos 获取热门视频
// 路由: GET /api/video/popular
func (h *VideoHandler) GetHotVideos(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	resp, err := h.videoService.GetHotVideos(c.Request.Context(), limit)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, map[string]interface{}{
		"videos": resp,
	})
}

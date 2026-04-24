package handler

import (
	"adult-short-videos/common/response"
	"adult-short-videos/services/video/logic"
	"adult-short-videos/services/video/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// VideoHandler 视频处理器
type VideoHandler struct {
	db        *gorm.DB
	videoRepo repository.VideoRepository
}

// NewVideoHandler 创建视频处理器实例
func NewVideoHandler(db *gorm.DB) *VideoHandler {
	return &VideoHandler{
		db:        db,
		videoRepo: repository.NewVideoRepository(db),
	}
}

// GetVideoList 获取视频列表
// 路由: GET /api/video/list
// 参数: page, size, region, category
func (h *VideoHandler) GetVideoList(c *gin.Context) {
	// c.DefaultQuery 可以设置默认值
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "20")
	region := c.Query("region")
	category := c.Query("category")

	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	req := &logic.VideoListReq{
		Page:     page,
		Size:     size,
		Region:   region,
		Category: category,
	}

	l := logic.NewVideoListLogic(c.Request.Context(), h.videoRepo)
	// 获取视频列表
	resp, err := l.GetVideoList(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

// GetVideoDetail 获取视频详情
// 路由: GET /api/video/detail/:id
// 参数: id（视频ID）
func (h *VideoHandler) GetVideoDetail(c *gin.Context) {

	// :id 是路由中定义的参数
	idStr := c.Param("id")

	videoId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	l := logic.NewVideoDetailLogic(c.Request.Context(), h.videoRepo)
	// 获取视频详情
	detail, err := l.GetVideoDetail(videoId)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, detail)
}

// GetHotVideos 获取热门视频
// 路由: GET /api/video/hot
// 参数: limit（可选，默认20）
func (h *VideoHandler) GetHotVideos(c *gin.Context) {

	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)

	// 限制范围
	if limit < 1 || limit > 50 {
		limit = 20
	}

	// 查询热门视频
	videos, err := h.videoRepo.GetHotVideos(c.Request.Context(), limit)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	items := make([]logic.VideoItem, 0, len(videos))
	for _, video := range videos {
		items = append(items, logic.VideoItem{
			VideoId:       video.VideoId,
			Title:         video.Title,
			CoverURL:      video.CoverURL,
			PreviewURL:    video.PreviewURL,
			Duration:      video.Duration,
			PlayCount:     video.PlayCount,
			FavoriteCount: video.FavoriteCount,
			Fanhao:        video.Fanhao,
			Region:        video.Region,
			Category:      video.Category,
			IsVipOnly:     video.IsVipOnly,
			PublishedAt:   video.PublishedAt.Unix(),
		})
	}

	response.Success(c, map[string]interface{}{
		"videos": items,
		"total":  len(items),
	})
}

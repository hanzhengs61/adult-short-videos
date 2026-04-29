package handler

import (
	"adult-short-videos/internal/pkg/response"
	"adult-short-videos/internal/service/video/logic"
	"adult-short-videos/internal/service/video/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VideoHandler struct {
	videoRepo repository.VideoRepository
}

func NewVideoHandler(db *gorm.DB) *VideoHandler {
	return &VideoHandler{videoRepo: repository.NewVideoRepository(db)}
}

// GetVideoList GET /api/video/list
func (h *VideoHandler) GetVideoList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	req := &logic.VideoListReq{
		Page:    page,
		Size:    size,
		OrderBy: c.DefaultQuery("order_by", "created_at"),
	}

	l := logic.NewVideoListLogic(c.Request.Context(), h.videoRepo)
	resp, err := l.GetVideoList(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, resp)
}

// GetVideoDetail GET /api/video/detail/:id
func (h *VideoHandler) GetVideoDetail(c *gin.Context) {
	videoId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	l := logic.NewVideoDetailLogic(c.Request.Context(), h.videoRepo)
	detail, err := l.GetVideoDetail(videoId)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, detail)
}

// GetHotVideos GET /api/video/hot
func (h *VideoHandler) GetHotVideos(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 50 {
		limit = 20
	}

	videos, err := h.videoRepo.GetHotVideos(c.Request.Context(), limit)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	items := make([]logic.VideoItem, 0, len(videos))
	for _, v := range videos {
		items = append(items, logic.VideoItem{
			VideoId:       v.VideoId,
			Title:         v.Title,
			CoverURL:      v.CoverURL,
			Duration:      v.Duration,
			IsPortrait:    v.IsPortrait,
			PlayURL:       logic.BuildPlayURL(v.StorageType, v.SourceURL, v.LocalURL),
			PlayCount:     v.PlayCount,
			FavoriteCount: v.FavoriteCount,
			PublishedAt:   v.PublishedAt.Unix(),
		})
	}

	response.Success(c, map[string]interface{}{
		"videos": items,
		"total":  len(items),
	})
}

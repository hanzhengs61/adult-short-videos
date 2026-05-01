package handler

import (
	"adult-short-videos/internal/pkg/response"
	"adult-short-videos/internal/service/search/dto"
	"adult-short-videos/internal/service/search/logic"
	"adult-short-videos/internal/service/search/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SearchHandler 搜索处理器
type SearchHandler struct {
	searchService *logic.SearchService
}

// NewSearchHandler 创建搜索处理器实例
func NewSearchHandler(db *gorm.DB) *SearchHandler {
	searchRepo := repository.NewSearchRepository(db)
	return &SearchHandler{
		searchService: logic.NewSearchService(searchRepo),
	}
}

// SearchVideos 搜索视频
// 路由: GET /api/search/videos
func (h *SearchHandler) SearchVideos(c *gin.Context) {
	var req dto.SearchReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	resp, err := h.searchService.Search(c.Request.Context(), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

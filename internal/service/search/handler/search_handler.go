package handler

import (
	"adult-short-videos/internal/pkg/response"
	"adult-short-videos/internal/service/search/logic"
	"adult-short-videos/internal/service/search/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SearchHandler 搜索处理器
type SearchHandler struct {
	searchRepo repository.SearchRepository
}

func NewSearchHandler(db *gorm.DB) *SearchHandler {
	return &SearchHandler{searchRepo: repository.NewSearchRepository(db)}
}

// SearchVideos GET /api/search/videos?keyword=xxx&page=1&page_size=20
func (h *SearchHandler) SearchVideos(c *gin.Context) {
	var req logic.SearchReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	l := logic.NewSearchLogic(c.Request.Context(), h.searchRepo)
	resp, err := l.Search(&req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, resp)
}

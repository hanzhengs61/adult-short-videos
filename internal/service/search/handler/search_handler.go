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
	db         *gorm.DB
	searchRepo repository.SearchRepository
}

func NewSearchHandler(db *gorm.DB) *SearchHandler {
	return &SearchHandler{
		db:         db,
		searchRepo: repository.NewSearchRepository(db),
	}
}

// SearchVideos 搜索视频
// 路由: GET /api/search/videos
// 参数: q (关键词), page, size
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

// SearchActors 搜索演员
// 路由: GET /api/search/actors
// 参数: q (关键词), page, size
func (h *SearchHandler) SearchActors(c *gin.Context) {
	var req logic.SearchReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	l := logic.NewActorSearchLogic(c.Request.Context(), h.searchRepo)
	resp, err := l.SearchActors(&req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

// SearchByFanhao 按番号搜索
// 路由: GET /api/search/fanhao/:fanhao
func (h *SearchHandler) SearchByFanhao(c *gin.Context) {
	var req logic.SearchFanhaoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	l := logic.NewFanhaoSearchLogic(c.Request.Context(), h.searchRepo)
	video, err := l.SearchByFanhao(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	// 返回视频详情
	response.Success(c, map[string]interface{}{
		"video_id":       video.VideoId,
		"title":          video.Title,
		"description":    video.Description,
		"cover_url":      video.CoverURL,
		"preview_url":    video.PreviewURL,
		"duration":       video.Duration,
		"fanhao":         video.Fanhao,
		"region":         video.Region,
		"category":       video.Category,
		"play_count":     video.PlayCount,
		"favorite_count": video.FavoriteCount,
		"is_vip_only":    video.IsVipOnly,
		"published_at":   video.PublishedAt.Unix(),
	})
}

// AdvancedSearch 高级搜索
// 路由: GET /api/search/advanced
// 参数: q, region, category, resolution, min_duration, max_duration, actor_id, sort, page, size
func (h *SearchHandler) AdvancedSearch(c *gin.Context) {
	var req logic.AdvancedSearchReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	l := logic.NewAdvancedSearchLogic(c.Request.Context(), h.searchRepo)
	resp, err := l.AdvancedSearch(&req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

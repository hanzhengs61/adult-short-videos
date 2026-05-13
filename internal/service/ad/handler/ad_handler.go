package handler

import (
	"adult-short-videos/internal/pkg/response"
	"adult-short-videos/internal/service/ad/dto"
	"adult-short-videos/internal/service/ad/logic"
	"adult-short-videos/internal/service/ad/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdHandler struct {
	svc *logic.AdService
}

func NewAdHandler(db *gorm.DB) *AdHandler {
	repo := repository.NewAdRepository(db)
	return &AdHandler{svc: logic.NewAdService(repo)}
}

// GetBanners 获取广告轮播图列表
// 路由: GET /api/ad/banners
func (h *AdHandler) GetBanners(c *gin.Context) {
	list, err := h.svc.GetBanners(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{"banners": list})
}

// GetApps 获取广告 App 列表，支持按分类筛选
// 路由: GET /api/ad/apps?category=video
func (h *AdHandler) GetApps(c *gin.Context) {
	var req dto.AppListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	list, err := h.svc.GetApps(c.Request.Context(), req.Category)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{"apps": list})
}

func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB) {
	h := NewAdHandler(db)
	g := api.Group("/ad")
	{
		g.GET("/banners", h.GetBanners)
		g.GET("/apps", h.GetApps)
	}
}

package handler

import (
	"net/http"
	"strconv"

	"adult-short-videos/internal/pkg/response"
	"adult-short-videos/internal/service/tag/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes 注册标签相关路由
func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB) {
	repo := repository.NewTagRepository(db)
	h := &tagHandler{repo: repo}

	g := api.Group("/tags")
	g.GET("/hot", h.GetHot)
	g.POST("/:name/click", h.Click)
}

type tagHandler struct {
	repo *repository.TagRepository
}

// GetHot GET /api/tags/hot?limit=10
// 返回点击数最高的标签列表
func (h *tagHandler) GetHot(c *gin.Context) {
	limit := 10
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 50 {
			limit = n
		}
	}
	tags, err := h.repo.GetHot(limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}
	items := make([]gin.H, 0, len(tags))
	for _, t := range tags {
		items = append(items, gin.H{
			"name":        t.Name,
			"click_count": t.ClickCount,
		})
	}
	response.Success(c, gin.H{"tags": items})
}

// Click POST /api/tags/:name/click
// 记录用户点击热搜标签（异步自增，无需登录）
func (h *tagHandler) Click(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.Error(c, http.StatusBadRequest, "标签名不能为空")
		return
	}
	// 异步写库，不阻塞响应
	go h.repo.IncrClickCount(name)
	response.Success(c, nil)
}

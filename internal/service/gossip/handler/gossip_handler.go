package handler

import (
	"strconv"

	"adult-short-videos/internal/pkg/response"
	"adult-short-videos/internal/service/gossip/dto"
	"adult-short-videos/internal/service/gossip/logic"
	"adult-short-videos/internal/service/gossip/repository"
	videoRepo "adult-short-videos/internal/service/video/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GossipHandler 吃瓜文章处理器
type GossipHandler struct {
	svc *logic.GossipService
}

func NewGossipHandler(db *gorm.DB) *GossipHandler {
	gr := repository.NewGossipRepository(db)
	vr := videoRepo.NewVideoRepository(db)
	return &GossipHandler{svc: logic.NewGossipService(gr, vr)}
}

// GetList 获取吃瓜列表
// 路由: GET /api/gossip/list
func (h *GossipHandler) GetList(c *gin.Context) {
	var req dto.GossipListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	resp, err := h.svc.GetList(c.Request.Context(), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

// GetDetail 获取吃瓜文章详情
// 路由: GET /api/gossip/detail/:id
func (h *GossipHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	resp, err := h.svc.GetDetail(c.Request.Context(), id)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

package handler

import (
	"strconv"

	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/pkg/response"
	"adult-short-videos/internal/service/favorite/dto"
	"adult-short-videos/internal/service/favorite/logic"
	"adult-short-videos/internal/service/favorite/repository"
	videoRepo "adult-short-videos/internal/service/video/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FavoriteHandler 收藏处理器
type FavoriteHandler struct {
	favoriteService *logic.FavoriteService
}

// NewFavoriteHandler 创建收藏处理器实例
func NewFavoriteHandler(db *gorm.DB) *FavoriteHandler {
	fRepo := repository.NewFavoriteRepository(db)
	vRepo := videoRepo.NewVideoRepository(db)
	return &FavoriteHandler{
		favoriteService: logic.NewFavoriteService(fRepo, vRepo, db),
	}
}

// AddFavorite 添加收藏
// 路由: POST /api/favorite/add
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.AddFavoriteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	if err := h.favoriteService.AddFavorite(c.Request.Context(), userId.(int64), &req); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMsg(c, "收藏成功", nil)
}

// RemoveFavorite 取消收藏
// 路由: DELETE /api/favorite/remove
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req struct {
		VideoId int64 `json:"video_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.VideoId <= 0 {
		response.Error(c, errors.CodeInvalidParam, "视频 ID 无效")
		return
	}

	if err := h.favoriteService.RemoveFavorite(c.Request.Context(), userId.(int64), req.VideoId); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMsg(c, "取消收藏成功", nil)
}

// GetFavoriteList 获取收藏列表
// 路由: GET /api/favorite/list
func (h *FavoriteHandler) GetFavoriteList(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.FavoriteListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	resp, err := h.favoriteService.GetFavoriteList(c.Request.Context(), userId.(int64), &req)

	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

// CheckFavorite 检查是否收藏
// 路由: GET /api/favorite/check
func (h *FavoriteHandler) CheckFavorite(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	videoIdStr := c.Query("video_id")
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil || videoId <= 0 {
		response.Error(c, errors.CodeInvalidParam, "视频 ID 无效")
		return
	}

	isFavorite, err := h.favoriteService.CheckFavorite(c.Request.Context(), userId.(int64), videoId)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, map[string]bool{"is_favorite": isFavorite})
}

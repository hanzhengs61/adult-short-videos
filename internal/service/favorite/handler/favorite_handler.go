package handler

import (
	"adult-short-videos/internal/pkg/errors"
	"strconv"

	"adult-short-videos/internal/pkg/response"
	"adult-short-videos/internal/service/favorite/logic"
	"adult-short-videos/internal/service/favorite/repository"
	videoRepo "adult-short-videos/internal/service/video/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FavoriteHandler 收藏处理器
type FavoriteHandler struct {
	db           *gorm.DB
	favoriteRepo repository.FavoriteRepository
	videoRepo    videoRepo.VideoRepository
}

// NewFavoriteHandler 创建收藏处理器实例
func NewFavoriteHandler(db *gorm.DB) *FavoriteHandler {
	return &FavoriteHandler{
		db:           db,
		favoriteRepo: repository.NewFavoriteRepository(db),
		videoRepo:    videoRepo.NewVideoRepository(db),
	}
}

// AddFavorite 添加收藏
// 路由: POST /api/favorite/add
// 需要认证
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	// 获取用户 ID
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req logic.AddFavoriteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
	}

	if req.VideoId <= 0 {
		response.Error(c, errors.CodeInvalidParam, "视频 ID无效")
		return
	}

	l := logic.NewAddFavoriteLogic(
		c.Request.Context(),
		h.favoriteRepo,
		h.videoRepo,
		h.db,
	)

	if err := l.AddFavorite(userId.(int64), &req); err != nil {
		response.HandleError(c, err)
		return
	}

	// 返回成功响应
	response.SuccessWithMsg(c, "收藏成功", nil)
}

// RemoveFavorite 取消收藏
// 路由: DELETE /api/favorite/remove
// 需要认证
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	// 获取用户 ID
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var body struct {
		VideoId int64 `json:"video_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.VideoId <= 0 {
		response.Error(c, errors.CodeInvalidParam, "视频 ID 无效")
		return
	}
	videoId := body.VideoId

	l := logic.NewRemoveFavoriteLogic(
		c.Request.Context(),
		h.favoriteRepo,
		h.db,
	)

	// 删除收藏
	if err := l.RemoveFavorite(userId.(int64), videoId); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMsg(c, "已取消收藏", nil)
}

// GetFavoriteList 获取收藏列表
// 路由: GET /api/favorite/list
// 需要认证
func (h *FavoriteHandler) GetFavoriteList(c *gin.Context) {
	// 获取用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "20")

	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	req := &logic.FavoriteListReq{
		Page: page,
		Size: size,
	}

	l := logic.NewFavoriteListLogic(
		c.Request.Context(),
		h.favoriteRepo,
	)

	// 获取收藏列表
	resp, err := l.GetFavoriteList(userId.(int64), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

// CheckFavorite 检查是否已收藏
// 路由: GET /api/favorite/check?video_id=
// 需要认证
func (h *FavoriteHandler) CheckFavorite(c *gin.Context) {
	// 获取用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil || videoId <= 0 {
		response.Error(c, errors.CodeInvalidParam, "视频 ID 无效")
		return
	}

	// 检查是否已收藏
	isFavorited, err := h.favoriteRepo.Exists(c.Request.Context(), userId.(int64), videoId)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, map[string]interface{}{
		"is_favorited": isFavorited,
	})
}

package handler

import (
	"strconv"

	"adult-short-videos/common/response"
	"adult-short-videos/services/favorite/logic"
	"adult-short-videos/services/favorite/repository"
	videoRepo "adult-short-videos/services/video/repository"

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
		response.Error(c, response.INVALID_PARAM, "参数格式错误")
		return
	}

	if req.VideoId <= 0 {
		response.Error(c, response.INVALID_PARAM, "视频 ID无效")
		return
	}

	l := logic.NewAddFavoriteLogic(
		c.Request.Context(),
		h.favoriteRepo,
		h.videoRepo,
		h.db,
	)

	if err := l.AddFavorite(userId.(int64), &req); err != nil {
		response.Error(c, response.ERROR, err.Error())
		return
	}

	// 返回成功响应
	response.SuccessWithMsg(c, "收藏成功", nil)
}

// RemoveFavorite 取消收藏
// 路由: DELETE /api/favorite/remove/:videoId
// 需要认证
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	// 获取用户 ID
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	// 获取视频 ID
	videoIdStr := c.Param("videoId")
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		response.Error(c, response.INVALID_PARAM, "视频ID格式错误")
		return
	}

	l := logic.NewRemoveFavoriteLogic(
		c.Request.Context(),
		h.favoriteRepo,
		h.db,
	)

	// 删除收藏
	if err := l.RemoveFavorite(userId.(int64), videoId); err != nil {
		response.Error(c, response.ERROR, err.Error())
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
		response.Error(c, response.ERROR, err.Error())
		return
	}

	response.Success(c, resp)
}

// CheckFavorite 检查是否已收藏
// 路由: GET /api/favorite/check/:videoId
// 需要认证
func (h *FavoriteHandler) CheckFavorite(c *gin.Context) {
	// 获取用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	// 获取视频ID
	videoIdStr := c.Param("videoId")
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		response.Error(c, response.INVALID_PARAM, "视频 ID格式错误")
		return
	}

	// 检查是否已收藏
	isFavorited, err := h.favoriteRepo.Exists(c.Request.Context(), userId.(int64), videoId)
	if err != nil {
		response.Error(c, response.ERROR, "检查收藏状态失败")
		return
	}

	response.Success(c, map[string]interface{}{
		"is_favorited": isFavorited,
	})
}

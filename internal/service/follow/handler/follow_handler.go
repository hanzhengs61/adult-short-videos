package handler

import (
	"strconv"

	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/pkg/response"
	"adult-short-videos/internal/service/follow/logic"
	"adult-short-videos/internal/service/follow/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FollowHandler struct {
	svc                 *logic.FollowService
	staticAssetsBaseURL string
}

func NewFollowHandler(db *gorm.DB, staticAssetsBaseURL string) *FollowHandler {
	repo := repository.NewFollowRepository(db)
	return &FollowHandler{
		svc:                 logic.NewFollowService(repo, db),
		staticAssetsBaseURL: staticAssetsBaseURL,
	}
}

// Follow POST /api/follow
func (h *FollowHandler) Follow(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		Author string `json:"author"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Author == "" {
		response.Error(c, errors.CodeInvalidParam, "作者名无效")
		return
	}
	if err := h.svc.Follow(c.Request.Context(), userID.(int64), req.Author); err != nil {
		response.HandleError(c, err)
		return
	}
	response.SuccessWithMsg(c, "关注成功", nil)
}

// Unfollow DELETE /api/follow
func (h *FollowHandler) Unfollow(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		Author string `json:"author"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Author == "" {
		response.Error(c, errors.CodeInvalidParam, "作者名无效")
		return
	}
	if err := h.svc.Unfollow(c.Request.Context(), userID.(int64), req.Author); err != nil {
		response.HandleError(c, err)
		return
	}
	response.SuccessWithMsg(c, "取消关注成功", nil)
}

// CheckFollow GET /api/follow/check?author=xxx
func (h *FollowHandler) CheckFollow(c *gin.Context) {
	userID, _ := c.Get("user_id")
	author := c.Query("author")
	if author == "" {
		response.Error(c, errors.CodeInvalidParam, "作者名无效")
		return
	}
	isFollowed, err := h.svc.CheckFollow(c.Request.Context(), userID.(int64), author)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, map[string]bool{"is_followed": isFollowed})
}

// ListFollowing GET /api/follow/list
func (h *FollowHandler) ListFollowing(c *gin.Context) {
	userID, _ := c.Get("user_id")
	profiles, err := h.svc.ListFollowingProfiles(c.Request.Context(), userID.(int64), h.staticAssetsBaseURL)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	if profiles == nil {
		profiles = []logic.AuthorProfileItem{}
	}
	response.Success(c, map[string]interface{}{"authors": profiles})
}

// GetFeed GET /api/follow/feed?page=1&page_size=20
func (h *FollowHandler) GetFeed(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	videos, hasMore, err := h.svc.GetFeed(c.Request.Context(), userID.(int64), page, pageSize)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, map[string]interface{}{
		"videos":   videos,
		"has_more": hasMore,
	})
}

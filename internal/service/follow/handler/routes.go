package handler

import (
	"adult-short-videos/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, jwtSecret string, staticAssetsBaseURL string) {
	h := NewFollowHandler(db, staticAssetsBaseURL)
	g := rg.Group("/follow")
	g.Use(middleware.AuthMiddleware(jwtSecret))
	g.POST("", h.Follow)
	g.DELETE("", h.Unfollow)
	g.GET("/check", h.CheckFollow)
	g.GET("/list", h.ListFollowing)
	g.GET("/feed", h.GetFeed)
}

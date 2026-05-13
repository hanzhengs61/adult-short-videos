package handler

import (
	"adult-short-videos/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB, jwtSecret string) {
	h := NewManageHandler(db)

	g := api.Group("/manage")
	g.Use(middleware.AuthMiddleware(jwtSecret), middleware.AdminMiddleware())
	{
		g.GET("/dashboard", h.Dashboard)

		// 视频
		g.GET("/videos", h.VideoList)
		g.PUT("/videos/:id/status", h.VideoUpdateStatus)
		g.DELETE("/videos/:id", h.VideoDelete)

		// 评论
		g.GET("/comments", h.CommentList)
		g.DELETE("/comments/:id", h.CommentDelete)
		g.DELETE("/comments", h.CommentBatchDelete)

		// 用户
		g.GET("/users", h.UserList)
		g.PUT("/users/:id", h.UserUpdate)

		// 广告 Banner
		g.GET("/banners", h.BannerList)
		g.POST("/banners", h.BannerCreate)
		g.PUT("/banners/:id", h.BannerUpdate)
		g.DELETE("/banners/:id", h.BannerDelete)

		// 广告 App
		g.GET("/apps", h.AppList)
		g.POST("/apps", h.AppCreate)
		g.PUT("/apps/:id", h.AppUpdate)
		g.DELETE("/apps/:id", h.AppDelete)

		// 吃瓜文章
		g.GET("/gossip", h.GossipList)
		g.POST("/gossip", h.GossipCreate)
		g.PUT("/gossip/:id", h.GossipUpdate)
		g.DELETE("/gossip/:id", h.GossipDelete)
	}
}

package handler

import (
	"adult-short-videos/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, jwtSecret string) {
	h := NewCommentHandler(db)
	g := rg.Group("/comment")
	g.GET("/list", middleware.OptionalAuthMiddleware(jwtSecret), h.GetCommentList)

	auth := g.Group("")
	auth.Use(middleware.AuthMiddleware(jwtSecret))
	auth.POST("/add", h.AddComment)
	auth.POST("/like", h.LikeComment)
	auth.DELETE("/like/:id", h.UnlikeComment)
	auth.DELETE("/delete/:id", h.DeleteComment)
}

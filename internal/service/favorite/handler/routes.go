package handler

import (
	"adult-short-videos/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, jwtSecret string) {
	h := NewFavoriteHandler(db)
	g := rg.Group("/favorite")
	g.Use(middleware.AuthMiddleware(jwtSecret))
	g.POST("/add", h.AddFavorite)
	g.DELETE("/remove", h.RemoveFavorite)
	g.GET("/list", h.GetFavoriteList)
	g.GET("/check", h.CheckFavorite)
}

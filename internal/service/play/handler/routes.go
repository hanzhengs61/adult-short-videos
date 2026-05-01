package handler

import (
	"adult-short-videos/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, jwtSecret string) {
	h := NewPlayHandler(db)
	g := rg.Group("/play")
	g.Use(middleware.AuthMiddleware(jwtSecret))
	g.POST("/record", h.RecordPlay)
	g.GET("/history", h.GetPlayHistory)
	g.DELETE("/history", h.ClearPlayHistory)
}

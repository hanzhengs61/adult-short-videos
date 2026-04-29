package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	h := NewSearchHandler(db)
	g := rg.Group("/search")
	g.GET("/videos", h.SearchVideos)
}

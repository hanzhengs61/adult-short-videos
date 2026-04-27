package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	h := NewVideoHandler(db)
	g := rg.Group("/video")
	g.GET("/list", h.GetVideoList)
	g.GET("/detail/:id", h.GetVideoDetail)
	g.GET("/hot", h.GetHotVideos)
	g.GET("/popular", h.GetHotVideos)
}

package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	h := NewGossipHandler(db)
	g := rg.Group("/gossip")
	g.GET("/list", h.GetList)
	g.GET("/detail/:id", h.GetDetail)
}

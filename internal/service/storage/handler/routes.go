package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	h := NewStorageHandler()
	g := rg.Group("/storage")
	g.GET("/proxy", h.ProxyPlay)
}

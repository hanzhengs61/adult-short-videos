package handler

import (
	"adult-short-videos/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, jwtSecret string, jwtExpire int64) {
	h := NewUserHandler(db, jwtSecret, jwtExpire)
	g := rg.Group("/user")
	g.POST("/register", h.Register)
	g.POST("/login", h.Login)
	g.GET("/creators", h.GetTopCreators)

	auth := g.Group("")
	auth.Use(middleware.AuthMiddleware(jwtSecret))
	auth.GET("/info", h.GetUserInfo)
	auth.POST("/logout", h.Logout)
}

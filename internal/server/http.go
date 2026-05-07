package server

import (
	"net/http"
	"strings"
	"time"

	"adult-short-videos/internal/config"
	"adult-short-videos/internal/pkg/cache"
	"adult-short-videos/internal/pkg/logger"
	"adult-short-videos/internal/pkg/middleware"
	commentHandler "adult-short-videos/internal/service/comment/handler"
	favoriteHandler "adult-short-videos/internal/service/favorite/handler"
	followHandler "adult-short-videos/internal/service/follow/handler"
	playHandler "adult-short-videos/internal/service/play/handler"
	searchHandler "adult-short-videos/internal/service/search/handler"
	storageHandler "adult-short-videos/internal/service/storage/handler"
	userHandler "adult-short-videos/internal/service/user/handler"
	videoHandler "adult-short-videos/internal/service/video/handler"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// New 创建并返回配置好的 HTTP 服务器。
// 使用标准库 *http.Server 而非 r.Run()，以支持 Shutdown() 优雅退出。
func New(cfg *config.Config, db *gorm.DB) *http.Server {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := buildRouter(cfg, db)

	return &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}
}

func buildRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(requestLogger())
	r.Use(gin.Recovery())
	r.Use(middleware.Metrics("/health", cfg.Metrics.Path))

	var corsOrigins []string
	if cfg.Server.CORSOrigins != "" {
		for _, o := range strings.Split(cfg.Server.CORSOrigins, ",") {
			if o = strings.TrimSpace(o); o != "" {
				corsOrigins = append(corsOrigins, o)
			}
		}
	}
	r.Use(middleware.CORS(corsOrigins))
	//r.Use(middleware.RateLimit(middleware.NewIPRateLimiter(rate.Limit(10), 20)))

	r.GET("/health", healthCheck(db))
	if cfg.Metrics.Enabled {
		r.GET(cfg.Metrics.Path, gin.WrapH(promhttp.Handler()))
	}

	// 用户上传的静态文件（头像等）；生产环境可换 Nginx/CDN 直接服务，此行可注释掉
	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	userHandler.RegisterRoutes(api, db, cfg.JWT.Secret, cfg.JWT.AccessExpire, cfg.Server.StaticAssetsBaseURL)
	videoHandler.RegisterRoutes(api, db)
	favoriteHandler.RegisterRoutes(api, db, cfg.JWT.Secret)
	playHandler.RegisterRoutes(api, db, cfg.JWT.Secret)
	storageHandler.RegisterRoutes(api)
	searchHandler.RegisterRoutes(api, db)
	commentHandler.RegisterRoutes(api, db, cfg.JWT.Secret)
	followHandler.RegisterRoutes(api, db, cfg.JWT.Secret, cfg.Server.StaticAssetsBaseURL)

	return r
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		logger.Info("API 请求",
			zap.String("request_id", c.GetString(middleware.RequestIDKey)),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("cost", time.Since(start)),
		)
	}
}

func healthCheck(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sqlDB, _ := db.DB()
		dbErr := sqlDB.Ping()

		var redisErr error
		if cache.Client != nil {
			redisErr = cache.Client.Ping(c.Request.Context()).Err()
		}

		status := "healthy"
		if dbErr != nil || redisErr != nil {
			status = "unhealthy"
		}

		c.JSON(http.StatusOK, gin.H{
			"status": status,
			"checks": gin.H{
				"database": dbErr == nil,
				"redis":    redisErr == nil,
			},
			"timestamp": time.Now().Unix(),
		})
	}
}

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"adult-short-videos/internal/config"
	"adult-short-videos/internal/pkg/cache"
	"adult-short-videos/internal/pkg/logger"
	"adult-short-videos/internal/pkg/middleware"
	commentHandler "adult-short-videos/internal/service/comment/handler"
	commentModel "adult-short-videos/internal/service/comment/model"
	favoriteHandler "adult-short-videos/internal/service/favorite/handler"
	favoriteModel "adult-short-videos/internal/service/favorite/model"
	playHandler "adult-short-videos/internal/service/play/handler"
	playModel "adult-short-videos/internal/service/play/model"
	searchHandler "adult-short-videos/internal/service/search/handler"
	storageHandler "adult-short-videos/internal/service/storage/handler"
	"adult-short-videos/internal/service/storage/scheduler"
	"adult-short-videos/internal/service/user/handler"
	userModel "adult-short-videos/internal/service/user/model"
	videoHandler "adult-short-videos/internal/service/video/handler"
	videoModel "adult-short-videos/internal/service/video/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	// ========== 第1步: 加载配置 ==========
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		fmt.Printf("❌ 加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// ========== 第2步: 初始化日志系统 ==========
	if err := logger.Init(
		cfg.Log.Level,
		cfg.Log.FilePath,
		cfg.Log.MaxSize,
		cfg.Log.MaxBackups,
		cfg.Log.MaxAge,
	); err != nil {
		fmt.Printf("❌ 初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("🚀 应用启动", zap.String("mode", cfg.Server.Mode))

	// ========== 第3步: 连接数据库 ==========
	db, err := initDatabase(cfg)
	if err != nil {
		logger.Fatal("数据库初始化失败", zap.Error(err))
	}
	logger.Info("✅ 数据库连接成功")

	// ========== 第4步: 初始化 Redis ==========
	if err := cache.Init(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB); err != nil {
		logger.Warn("Redis连接失败，缓存功能将不可用", zap.Error(err))
	} else {
		logger.Info("✅ Redis连接成功")
	}

	// ========== 第5步: 自动迁移数据库表 ==========
	if err := migrateDatabase(db); err != nil {
		logger.Fatal("数据库迁移失败", zap.Error(err))
	}
	logger.Info("✅ 数据库迁移完成")

	// ========== 第6步: 创建索引 ==========
	createIndexes(db)

	// ========== 第7步: 设置 Gin 模式 ==========
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// ========== 第8步: 创建 Gin 引擎 ==========
	r := setupRouter(cfg, db)

	// ========== 第9步: 启动热度计算器 ==========
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	heatCalc := scheduler.NewHeatCalculator(db, 5*time.Minute)
	go heatCalc.Start(ctx)
	logger.Info("🔥 热度计算器已启动")

	// ========== 第10步: 启动服务器 ==========
	printStartupInfo(cfg)

	// 优雅关闭
	srv := startServer(r, cfg.Server.Port)
	gracefulShutdown(srv, cancel)
}

// initDatabase 初始化数据库
func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	// GORM 日志配置
	logLevel := gormLogger.Info
	if cfg.Server.Mode == "release" {
		logLevel = gormLogger.Warn
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{
		Logger: gormLogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	// 最大空闲连接数
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	// 设置连接最大空闲时间
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	logger.Info("数据库连接池配置",
		zap.Int("max_idle", cfg.Database.MaxIdleConns),
		zap.Int("max_open", cfg.Database.MaxOpenConns),
		zap.Int("max_lifetime", cfg.Database.ConnMaxLifetime),
	)

	return db, nil
}

// migrateDatabase 数据库迁移
func migrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		// 用户相关
		&userModel.User{},
		&userModel.UserStatistics{},
		&userModel.LoginLog{},
		// 视频相关
		&videoModel.Video{},
		&videoModel.VideoHeatStats{},
		&videoModel.Actor{},
		&videoModel.VideoActor{},
		// 收藏
		&favoriteModel.Favorite{},
		// 播放历史
		&playModel.PlayHistory{},
		// 评论
		&commentModel.Comment{},
		&commentModel.CommentLike{},
	)
}

// createIndexes 创建数据库索引
func createIndexes(db *gorm.DB) {
	logger.Info("📊 开始创建数据库索引...")

	indexes := []string{
		// 视频表索引
		"CREATE INDEX IF NOT EXISTS idx_videos_status_region ON videos(status, region)",
		"CREATE INDEX IF NOT EXISTS idx_videos_status_category ON videos(status, category)",
		"CREATE INDEX IF NOT EXISTS idx_videos_play_count ON videos(play_count DESC)",
		"CREATE INDEX IF NOT EXISTS idx_videos_created_at ON videos(created_at DESC)",

		// 收藏表索引
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_favorites_user_video ON favorites(user_id, video_id)",

		// 播放历史表索引
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_play_history_user_video ON play_history(user_id, video_id)",
		"CREATE INDEX IF NOT EXISTS idx_play_history_last_play ON play_history(last_play_at DESC)",

		// 评论表索引
		"CREATE INDEX IF NOT EXISTS idx_comments_video_id ON comments(video_id, status, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id)",
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_comment_likes_user_comment ON comment_likes(user_id, comment_id)",
	}

	for _, sql := range indexes {
		if err := db.Exec(sql).Error; err != nil {
			logger.Warn("创建索引失败", zap.Error(err), zap.String("sql", sql))
		}
	}

	logger.Info("✅ 数据库索引创建完成")
}

// setupRouter 设置路由
func setupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	r := gin.New()

	// ========== 全局中间件 ==========
	// 日志中间件
	r.Use(ginLogger())
	// 恢复中间件
	r.Use(gin.Recovery())
	// CORS 中间件
	var corsOrigins []string
	if cfg.Server.CORSOrigins != "" {
		for _, o := range strings.Split(cfg.Server.CORSOrigins, ",") {
			if o = strings.TrimSpace(o); o != "" {
				corsOrigins = append(corsOrigins, o)
			}
		}
	}
	r.Use(middleware.CORS(corsOrigins))
	// 限流中间件（每秒10个请求，突发20个）
	limiter := middleware.NewIPRateLimiter(rate.Limit(10), 20)
	r.Use(middleware.RateLimit(limiter))

	// ========== 健康检查 ==========
	r.GET("/health", healthCheck(db))

	// ========== 创建处理器 ==========
	userHandler := handler.NewUserHandler(db, cfg.JWT.Secret, cfg.JWT.AccessExpire)
	videoHandlerInstance := videoHandler.NewVideoHandler(db)
	favoriteHandlerInstance := favoriteHandler.NewFavoriteHandler(db)
	playHandlerInstance := playHandler.NewPlayHandler(db)
	storageHandlerInstance := storageHandler.NewStorageHandler()
	searchHandlerInstance := searchHandler.NewSearchHandler(db)
	commentHandlerInstance := commentHandler.NewCommentHandler(db)

	// ========== 注册路由 ==========
	api := r.Group("/api")
	{
		// 用户路由
		user := api.Group("/user")
		{
			user.POST("/register", userHandler.Register)
			user.POST("/login", userHandler.Login)

			auth := user.Group("")
			auth.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
			{
				auth.GET("/info", userHandler.GetUserInfo)
				auth.POST("/logout", userHandler.Logout)
			}
		}

		// 视频路由
		video := api.Group("/video")
		{
			video.GET("/list", videoHandlerInstance.GetVideoList)
			video.GET("/detail/:id", videoHandlerInstance.GetVideoDetail)
			video.GET("/hot", videoHandlerInstance.GetHotVideos)
		}

		// 收藏路由
		favorite := api.Group("/favorite")
		favorite.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			favorite.POST("/add", favoriteHandlerInstance.AddFavorite)
			favorite.DELETE("/remove/:videoId", favoriteHandlerInstance.RemoveFavorite)
			favorite.GET("/list", favoriteHandlerInstance.GetFavoriteList)
			favorite.GET("/check/:videoId", favoriteHandlerInstance.CheckFavorite)
		}

		// 播放历史路由
		play := api.Group("/play")
		play.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			play.POST("/record", playHandlerInstance.RecordPlay)
			play.GET("/history", playHandlerInstance.GetPlayHistory)
			play.DELETE("/history/:videoId", playHandlerInstance.DeletePlayHistory)
			play.DELETE("/history", playHandlerInstance.ClearPlayHistory)
		}

		// 存储服务路由
		storage := api.Group("/storage")
		{
			storage.GET("/proxy", storageHandlerInstance.ProxyPlay)
		}

		// 搜索路由
		search := api.Group("/search")
		{
			search.GET("/videos", searchHandlerInstance.SearchVideos)
			search.GET("/actors", searchHandlerInstance.SearchActors)
			search.GET("/fanhao/:fanhao", searchHandlerInstance.SearchByFanhao)
			search.GET("/advanced", searchHandlerInstance.AdvancedSearch)
		}

		// 评论路由
		comment := api.Group("/comment")
		{
			comment.GET("/list", commentHandlerInstance.GetCommentList)

			commentAuth := comment.Group("")
			commentAuth.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
			{
				commentAuth.POST("/add", commentHandlerInstance.AddComment)
				commentAuth.POST("/like/:id", commentHandlerInstance.LikeComment)
				commentAuth.DELETE("/like/:id", commentHandlerInstance.UnlikeComment)
			}
		}
	}

	return r
}

// ginLogger Gin 日志中间件
func ginLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		cost := time.Since(start)

		logger.Info("API请求",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("cost", cost),
		)
	}
}

// healthCheck 健康检查
func healthCheck(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查数据库
		sqlDB, _ := db.DB()
		dbErr := sqlDB.Ping()

		// 检查 Redis
		var redisErr error
		if cache.Client != nil {
			redisErr = cache.Client.Ping(c.Request.Context()).Err()
		}

		status := "healthy"
		if dbErr != nil || redisErr != nil {
			status = "unhealthy"
		}

		c.JSON(200, gin.H{
			"status": status,
			"checks": gin.H{
				"database": dbErr == nil,
				"redis":    redisErr == nil,
			},
			"timestamp": time.Now().Unix(),
		})
	}
}

// startServer 启动服务器
func startServer(r *gin.Engine, port string) *gin.Engine {
	go func() {
		if err := r.Run(port); err != nil {
			logger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()
	return r
}

// gracefulShutdown 优雅关闭
func gracefulShutdown(r *gin.Engine, cancel context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 正在关闭服务器...")

	// 取消所有 context
	cancel()

	// 关闭 Redis
	if cache.Client != nil {
		cache.Client.Close()
	}

	logger.Info("✅ 服务器已关闭")
}

// printStartupInfo 打印启动信息
func printStartupInfo(cfg *config.Config) {
	banner := `
╔════════════════════════════════════════════════════════════════╗
║                                                                ║
║    🎬 短视频平台 - Adult Short Videos Platform           		 ║
║                                                                ║
╚════════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)

	fmt.Printf("🚀 服务器启动在: http://localhost%s\n", cfg.Server.Port)
	fmt.Printf("📝 运行模式: %s\n", cfg.Server.Mode)
	fmt.Printf("📊 日志级别: %s\n", cfg.Log.Level)
	fmt.Println()

	fmt.Println(strings.Repeat("=", 64))
	fmt.Println("📝 API 接口列表")
	fmt.Println(strings.Repeat("=", 64))

	endpoints := []struct {
		group string
		apis  []string
	}{
		{
			group: "用户服务",
			apis: []string{
				"POST /api/user/register      - 用户注册",
				"POST /api/user/login         - 用户登录",
				"GET  /api/user/info          - 获取用户信息（需认证）",
				"POST /api/user/logout        - 退出登录（需认证）",
			},
		},
		{
			group: "视频服务",
			apis: []string{
				"GET  /api/video/list         - 视频列表",
				"GET  /api/video/detail/:id   - 视频详情",
				"GET  /api/video/hot          - 热门视频",
			},
		},
		{
			group: "收藏服务",
			apis: []string{
				"POST   /api/favorite/add         - 添加收藏（需认证）",
				"DELETE /api/favorite/remove/:id  - 取消收藏（需认证）",
				"GET    /api/favorite/list        - 收藏列表（需认证）",
				"GET    /api/favorite/check/:id   - 检查收藏状态（需认证）",
			},
		},
		{
			group: "播放历史",
			apis: []string{
				"POST   /api/play/record          - 记录播放（需认证）",
				"GET    /api/play/history         - 播放历史（需认证）",
				"DELETE /api/play/history/:id     - 删除单条（需认证）",
				"DELETE /api/play/history         - 清空历史（需认证）",
			},
		},
		{
			group: "搜索服务",
			apis: []string{
				"GET /api/search/videos      - 搜索视频",
				"GET /api/search/actors      - 搜索演员",
				"GET /api/search/fanhao/:id  - 按番号搜索",
				"GET /api/search/advanced    - 高级搜索",
			},
		},
		{
			group: "评论服务",
			apis: []string{
				"GET    /api/comment/list         - 评论列表",
				"POST   /api/comment/add          - 发表评论（需认证）",
				"POST   /api/comment/like/:id     - 点赞评论（需认证）",
				"DELETE /api/comment/like/:id     - 取消点赞（需认证）",
			},
		},
		{
			group: "系统服务",
			apis: []string{
				"GET /health                  - 健康检查",
				"GET /api/storage/proxy      - 代理播放",
			},
		},
	}

	for _, ep := range endpoints {
		fmt.Printf("\n【%s】\n", ep.group)
		for _, api := range ep.apis {
			fmt.Printf("  %s\n", api)
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 64))
	fmt.Println()
}

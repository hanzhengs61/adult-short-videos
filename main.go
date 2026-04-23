package main

import (
	"adult-short-videos/common/middleware"
	commentHandler "adult-short-videos/services/comment/handler"
	commentModel "adult-short-videos/services/comment/model"
	favoriteHandler "adult-short-videos/services/favorite/handler"
	favoriteModel "adult-short-videos/services/favorite/model"
	playHandler "adult-short-videos/services/play/handler"
	playModel "adult-short-videos/services/play/model"
	searchHandler "adult-short-videos/services/search/handler"
	storageHandler "adult-short-videos/services/storage/handler"
	"adult-short-videos/services/storage/scheduler"
	userHandler "adult-short-videos/services/user/handler"
	userModel "adult-short-videos/services/user/model"
	videoHandler "adult-short-videos/services/video/handler"
	videoModel "adult-short-videos/services/video/model"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// ========== 第1步: 连接数据库 ==========
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=adult_videos sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ 数据库连接失败:", err)
	}

	fmt.Println("✅ 数据库连接成功")

	// ========== 第2步: 自动迁移数据库表 ==========
	err = db.AutoMigrate(
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
	if err != nil {
		log.Fatal("❌ 数据库迁移失败:", err)
	}

	fmt.Println("✅ 数据库表创建成功")

	// ========== 第3步: 创建数据库索引 ==========
	createIndexes(db)

	// ========== 第4步: 配置参数 ==========
	jwtSecret := "your-secret-key-change-in-production"
	jwtExpire := int64(86400) // 24小时

	// ========== 第5步: 创建 Gin 引擎 ==========
	r := gin.Default()
	r.Use(middleware.CORS())

	// ========== 第6步: 创建处理器 ==========
	userHandlerInstance := userHandler.NewUserHandler(db, jwtSecret, jwtExpire)
	videoHandlerInstance := videoHandler.NewVideoHandler(db)
	favoriteHandlerInstance := favoriteHandler.NewFavoriteHandler(db)
	playHandlerInstance := playHandler.NewPlayHandler(db)
	storageHandlerInstance := storageHandler.NewStorageHandler()
	searchHandlerInstance := searchHandler.NewSearchHandler(db)
	commentHandlerInstance := commentHandler.NewCommentHandler(db)

	// ========== 第7步: 注册路由 ==========
	api := r.Group("/api")
	{
		// 用户路由
		user := api.Group("/user")
		{
			user.POST("/register", userHandlerInstance.Register)
			user.POST("/login", userHandlerInstance.Login)

			auth := user.Group("")
			auth.Use(middleware.AuthMiddleware(jwtSecret))
			{
				auth.GET("/info", userHandlerInstance.GetUserInfo)
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
		favorite.Use(middleware.AuthMiddleware(jwtSecret))
		{
			favorite.POST("/add", favoriteHandlerInstance.AddFavorite)
			favorite.DELETE("/remove/:videoId", favoriteHandlerInstance.RemoveFavorite)
			favorite.GET("/list", favoriteHandlerInstance.GetFavoriteList)
			favorite.GET("/check/:videoId", favoriteHandlerInstance.CheckFavorite)
		}

		// 播放历史路由
		play := api.Group("/play")
		play.Use(middleware.AuthMiddleware(jwtSecret))
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
			// 公开路由
			comment.GET("/list", commentHandlerInstance.GetCommentList)

			// 需要认证的路由
			commentAuth := comment.Group("")
			commentAuth.Use(middleware.AuthMiddleware(jwtSecret))
			{
				commentAuth.POST("/add", commentHandlerInstance.AddComment)
				commentAuth.POST("/like/:id", commentHandlerInstance.LikeComment)
				commentAuth.DELETE("/like/:id", commentHandlerInstance.UnlikeComment)
			}
		}
	}

	// ========== 第8步: 启动热度计算器 ==========
	heatCalc := scheduler.NewHeatCalculator(db, 5*time.Minute)
	go heatCalc.Start(context.Background())

	// ========== 第9步: 启动服务器 ==========
	printAPIInfo()

	port := ":8080"
	fmt.Printf("\n🚀 服务器启动在 http://localhost%s\n\n", port)

	if err := r.Run(port); err != nil {
		log.Fatal("❌ 服务器启动失败:", err)
	}
}

// createIndexes 创建数据库索引
func createIndexes(db *gorm.DB) {
	fmt.Println("📊 开始创建数据库索引...")

	// 视频表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_videos_status_region ON videos(status, region)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_videos_status_category ON videos(status, category)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_videos_play_count ON videos(play_count DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_videos_created_at ON videos(created_at DESC)")

	// 收藏表索引
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_favorites_user_video ON favorites(user_id, video_id)")

	// 播放历史表索引
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_play_history_user_video ON play_history(user_id, video_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_play_history_last_play ON play_history(last_play_at DESC)")

	// 评论表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_video_id ON comments(video_id, status, created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id)")
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_comment_likes_user_comment ON comment_likes(user_id, comment_id)")

	fmt.Println("✅ 数据库索引创建完成")
}

// printAPIInfo 打印 API 信息
func printAPIInfo() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📝 API 接口列表")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n【用户服务】")
	fmt.Println("  POST http://localhost:8080/api/user/register      - 用户注册")
	fmt.Println("  POST http://localhost:8080/api/user/login         - 用户登录")
	fmt.Println("  GET  http://localhost:8080/api/user/info          - 获取用户信息（需认证）")

	fmt.Println("\n【视频服务】")
	fmt.Println("  GET  http://localhost:8080/api/video/list         - 视频列表")
	fmt.Println("       参数: page=1&size=20&region=日本&category=职场")
	fmt.Println("  GET  http://localhost:8080/api/video/detail/:id   - 视频详情")
	fmt.Println("  GET  http://localhost:8080/api/video/hot          - 热门视频")

	fmt.Println("\n【收藏服务】（需认证）")
	fmt.Println("  POST   http://localhost:8080/api/favorite/add         - 添加收藏")
	fmt.Println("  DELETE http://localhost:8080/api/favorite/remove/:id  - 取消收藏")
	fmt.Println("  GET    http://localhost:8080/api/favorite/list        - 收藏列表")
	fmt.Println("  GET    http://localhost:8080/api/favorite/check/:id   - 检查收藏状态")

	fmt.Println("\n【播放历史】（需认证）")
	fmt.Println("  POST   http://localhost:8080/api/play/record          - 记录播放")
	fmt.Println("  GET    http://localhost:8080/api/play/history         - 播放历史")
	fmt.Println("  DELETE http://localhost:8080/api/play/history/:id     - 删除单条")
	fmt.Println("  DELETE http://localhost:8080/api/play/history         - 清空历史")

	fmt.Println("\n【存储服务】")
	fmt.Println("  GET  http://localhost:8080/api/storage/proxy      - 代理播放")
	fmt.Println("       参数: url=源站URL")

	fmt.Println("\n【搜索服务】")
	fmt.Println("  GET  http://localhost:8080/api/search/videos      - 搜索视频")
	fmt.Println("       参数: q=关键词&page=1&size=20")
	fmt.Println("  GET  http://localhost:8080/api/search/actors      - 搜索演员")
	fmt.Println("  GET  http://localhost:8080/api/search/fanhao/:id  - 按番号搜索")
	fmt.Println("  GET  http://localhost:8080/api/search/advanced    - 高级搜索")
	fmt.Println("       参数: q=关键词&region=日本&category=职场&sort=hot")

	fmt.Println("\n【评论服务】")
	fmt.Println("  GET    http://localhost:8080/api/comment/list         - 评论列表")
	fmt.Println("         参数: video_id=1&page=1&size=20")
	fmt.Println("  POST   http://localhost:8080/api/comment/add          - 发表评论（需认证）")
	fmt.Println("  POST   http://localhost:8080/api/comment/like/:id     - 点赞评论（需认证）")
	fmt.Println("  DELETE http://localhost:8080/api/comment/like/:id     - 取消点赞（需认证）")

	fmt.Println("\n" + strings.Repeat("=", 60))
}

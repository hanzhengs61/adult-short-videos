package main

import (
	"adult-short-videos/common/middleware"
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
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// ========== 第1步: 连接数据库 ==========
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=adult_short_videos sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

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
	)
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	fmt.Println("✅ 数据库表创建成功")

	// ========== 第3步: 配置参数 ==========
	jwtSecret := "your-secret-key-change-in-production" // JWT 密钥
	jwtExpire := int64(86400)                           // 24小时

	// ========== 第4步: 创建 Gin 引擎 ==========
	r := gin.Default()

	// ========== 第5步: 注册全局中间件 ==========
	r.Use(middleware.CORS()) // 跨域中间件

	// ========== 第6步: 创建处理器 ==========
	userHandlerInstance := userHandler.NewUserHandler(db, jwtSecret, jwtExpire)
	videoHandlerInstance := videoHandler.NewVideoHandler(db)
	favoriteHandlerInstance := favoriteHandler.NewFavoriteHandler(db)
	playHandlerInstance := playHandler.NewPlayHandler(db)
	storageHandlerInstance := storageHandler.NewStorageHandler()
	searchHandlerInstance := searchHandler.NewSearchHandler(db)

	// ========== 第7步: 注册路由 ==========
	api := r.Group("/api")
	{
		// ========== 用户路由 ==========
		user := api.Group("/user")
		{
			// 公开路由（不需要认证）
			user.POST("/register", userHandlerInstance.Register) // 注册
			user.POST("/login", userHandlerInstance.Login)       // 登录

			// 需要认证的路由
			auth := user.Group("")
			auth.Use(middleware.AuthMiddleware(jwtSecret)) // 使用认证中间件
			{
				auth.GET("/info", userHandlerInstance.GetUserInfo) // 获取用户信息
			}
		}
		// ========== 视频路由 ==========
		video := api.Group("/video")
		{
			// 公开路由（不需要认证就能看视频列表）
			video.GET("/list", videoHandlerInstance.GetVideoList)         // 视频列表
			video.GET("/detail/:id", videoHandlerInstance.GetVideoDetail) // 视频详情
			video.GET("/hot", videoHandlerInstance.GetHotVideos)          // 热门视频
		}

		// ========== 收藏路由 ==========
		favorite := api.Group("/favorite")
		favorite.Use(middleware.AuthMiddleware(jwtSecret)) // 所有收藏接口都需要认证
		{
			favorite.POST("/add", favoriteHandlerInstance.AddFavorite)               // 添加收藏
			favorite.GET("/remove/:videoId", favoriteHandlerInstance.RemoveFavorite) // 取消收藏
			favorite.GET("/list", favoriteHandlerInstance.GetFavoriteList)           // 收藏列表
			favorite.GET("/check/:videoId", favoriteHandlerInstance.CheckFavorite)   // 检查收藏状态
		}

		// ========== 播放历史路由 ==========
		play := api.Group("/play")
		play.Use(middleware.AuthMiddleware(jwtSecret))
		{
			play.POST("/record", playHandlerInstance.RecordPlay)                 // 记录播放
			play.GET("/history", playHandlerInstance.GetPlayHistory)             // 播放历史列表
			play.GET("/history/:videoId", playHandlerInstance.DeletePlayHistory) // 删除单条
			play.GET("/clerHistory", playHandlerInstance.ClearPlayHistory)       // 清空历史
		}

		// ========== 存储服务路由 ==========
		storage := api.Group("/storage")
		{
			storage.GET("/proxy", storageHandlerInstance.ProxyPlay) // 代理播放
		}

		// ========== 搜索路由 ==========
		search := api.Group("/search")
		{
			search.GET("/videos", searchHandlerInstance.SearchVideos)     // 搜索视频
			search.GET("/actors", searchHandlerInstance.SearchActors)     // 搜索演员
			search.POST("/fanhao/", searchHandlerInstance.SearchByFanhao) // 按番号搜索
			search.GET("/advanced", searchHandlerInstance.AdvancedSearch) // 高级搜索
		}
	}

	// ========== 启动热度计算器 ==========
	heatCalc := scheduler.NewHeatCalculator(db, 5*time.Minute) // 每5分钟计算一次
	go heatCalc.Start(context.Background())

	// ========== 第8步: 启动服务器 ==========
	port := ":8080"
	fmt.Printf("🚀 服务器启动在 http://localhost%s\n", port)
	fmt.Println("\n📝 可用接口:")
	fmt.Println("\n【用户服务】")
	fmt.Println("  POST http://localhost:8080/api/user/register      - 用户注册")
	fmt.Println("  POST http://localhost:8080/api/user/login         - 用户登录")
	fmt.Println("  GET  http://localhost:8080/api/user/info          - 获取用户信息（需认证）")
	fmt.Println("\n【视频服务】")
	fmt.Println("  GET  http://localhost:8080/api/video/list         - 视频列表")
	fmt.Println("       参数: page=1&size=20&region=日本&category=职场")
	fmt.Println("  GET  http://localhost:8080/api/video/detail/:id   - 视频详情")
	fmt.Println("  GET  http://localhost:8080/api/video/hot          - 热门视频")
	fmt.Println("\n【收藏服务】（需要认证）")
	fmt.Println("  POST   http://localhost:8080/api/favorite/add           - 添加收藏")
	fmt.Println("  GET    http://localhost:8080/api/favorite/remove/:id    - 取消收藏")
	fmt.Println("  GET    http://localhost:8080/api/favorite/list          - 收藏列表")
	fmt.Println("  GET    http://localhost:8080/api/favorite/check/:id     - 检查收藏状态")
	fmt.Println("\n【播放历史】")
	fmt.Println("  POST   http://localhost:8080/api/play/record				- 记录播放")
	fmt.Println("  GET    http://localhost:8080/api/play/history			- 播放历史列表")
	fmt.Println("  GET    http://localhost:8080/api/play/history/:id		- 删除单条")
	fmt.Println("  GET    http://localhost:8080/api/play/history			- 清空历史")
	fmt.Println("\n【存储服务】")
	fmt.Println("  GET    http://localhost:8080/api/storage/proxy      - 代理播放")
	fmt.Println("\n【搜索服务】")
	fmt.Println("  GET  http://localhost:8080/api/search/videos       - 搜索视频")
	fmt.Println("       参数: q=关键词&page=1&size=20")
	fmt.Println("  GET  http://localhost:8080/api/search/actors       - 搜索演员")
	fmt.Println("  GET  http://localhost:8080/api/search/fanhao/:id   - 按番号搜索")
	fmt.Println("  GET  http://localhost:8080/api/search/advanced     - 高级搜索")
	fmt.Println("       参数: q=关键词&region=日本&category=职场&sort=hot")
	fmt.Println()

	if err := r.Run(port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

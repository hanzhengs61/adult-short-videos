package main

import (
	"adult-short-videos/common/middleware"
	favoriteHandler "adult-short-videos/services/favorite/handler"
	favoriteModel "adult-short-videos/services/favorite/model"
	userHandler "adult-short-videos/services/user/handler"
	userModel "adult-short-videos/services/user/model"
	videoHandler "adult-short-videos/services/video/handler"
	"adult-short-videos/services/video/model"
	videoModel "adult-short-videos/services/video/model"
	"fmt"
	"log"

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

	// 用户相关表
	err = db.AutoMigrate(
		&userModel.User{},
		&userModel.UserStatistics{},
		&userModel.LoginLog{},
	)
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 视频相关表
	err = db.AutoMigrate(
		&model.Video{},
		&videoModel.VideoHeatStats{},
		&videoModel.Actor{},
		&videoModel.VideoActor{},
	)
	if err != nil {
		log.Fatal("视频表迁移失败:", err)
	}

	fmt.Println("✅ 数据库表创建成功")

	// 收藏表
	err = db.AutoMigrate(
		&favoriteModel.Favorite{},
	)
	if err != nil {
		log.Fatal("收藏表迁移失败:", err)
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
	userHandler := userHandler.NewUserHandler(db, jwtSecret, jwtExpire)
	videoHandler := videoHandler.NewVideoHandler(db)
	favoriteHandler := favoriteHandler.NewFavoriteHandler(db)

	// ========== 第7步: 注册路由 ==========
	api := r.Group("/api")
	{
		// ========== 用户路由 ==========
		user := api.Group("/user")
		{
			// 公开路由（不需要认证）
			user.POST("/register", userHandler.Register) // 注册
			user.POST("/login", userHandler.Login)       // 登录

			// 需要认证的路由
			auth := user.Group("")
			auth.Use(middleware.AuthMiddleware(jwtSecret)) // 使用认证中间件
			{
				auth.GET("/info", userHandler.GetUserInfo) // 获取用户信息
			}
		}
		// ========== 视频路由 ==========
		video := api.Group("/video")
		{
			// 公开路由（不需要认证就能看视频列表）
			video.GET("/list", videoHandler.GetVideoList)         // 视频列表
			video.GET("/detail/:id", videoHandler.GetVideoDetail) // 视频详情
			video.GET("/hot", videoHandler.GetHotVideos)          // 热门视频
		}

		// ========== 收藏路由 ==========
		favorite := api.Group("/favorite")
		favorite.Use(middleware.AuthMiddleware(jwtSecret)) // 所有收藏接口都需要认证
		{
			favorite.POST("/add", favoriteHandler.AddFavorite)               // 添加收藏
			favorite.GET("/remove/:videoId", favoriteHandler.RemoveFavorite) // 取消收藏
			favorite.GET("/list", favoriteHandler.GetFavoriteList)           // 收藏列表
			favorite.GET("/check/:videoId", favoriteHandler.CheckFavorite)   // 检查收藏状态
		}
	}

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
	fmt.Println()

	if err := r.Run(port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

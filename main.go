package main

import (
	"adult-short-videos/common/middleware"
	"adult-short-videos/services/user/handler"
	"adult-short-videos/services/user/model"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 配置数据库连接字符串
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=adult_short_videos sslmode=disable TimeZone=Asia/Shanghai"
	// 连接 PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移数据库表
	err = db.AutoMigrate(
		&model.User{},
		&model.UserStatistics{},
		&model.LoginLog{},
	)
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	fmt.Println("✅ 数据库表创建成功")

	// 配置参数
	jwtSecret := "your-secret-key-change-in-production" // JWT 密钥
	jwtExpire := int64(86400)                           // 24小时

	r := gin.Default()
	r.Use(middleware.CORS()) // 跨域中间件
	userHandler := handler.NewUserHandler(db, jwtSecret, jwtExpire)

	// 注册路由
	api := r.Group("/api")
	{
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
	}

	port := ":8080"
	fmt.Printf("🚀 服务器启动在 http://localhost%s\n", port)
	fmt.Println("📝 可用接口:")
	fmt.Println("  POST http://localhost:8080/api/user/register - 用户注册")
	fmt.Println("  POST http://localhost:8080/api/user/login    - 用户登录")
	fmt.Println("  GET  http://localhost:8080/api/user/info     - 获取用户信息（需认证）")

	if err := r.Run(port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

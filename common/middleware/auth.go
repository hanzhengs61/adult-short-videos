package middleware

import (
	"adult-short-videos/common/response"
	"adult-short-videos/common/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 认证中间件
// 参数: jwtSecret - JWT 密钥
// 返回: Gin 中间件函数
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1: 从 Header 获取 Token
		// 格式: Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		authHeader = strings.TrimSpace(authHeader)

		if authHeader == "" {
			response.Unauthorized(c, "未提供认证令牌")
			return
		}

		// 2: 检查 Token 格式
		var tokenString string
		if strings.HasPrefix(authHeader, "Bearer ") {
			// 提取 Bearer 后的 Token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Unauthorized(c, "认证令牌格式错误")
				return
			}
			tokenString = strings.TrimSpace(parts[1])
		} else {
			// 直接使用整个字符串作为 Token
			tokenString = authHeader
		}

		// 3: 解析和验证 Token
		claims, err := utils.ParseToken(tokenString, jwtSecret)
		if err != nil {
			response.Unauthorized(c, "认证令牌无效或已过期")
			return
		}

		// 4: 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// CORS 跨域中间件
// 允许前端跨域访问 API
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有来源访问（生产环境应该限制具体域名）
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

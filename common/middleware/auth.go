package middleware

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"adult-short-videos/common/cache"
	"adult-short-videos/common/response"
	"adult-short-videos/common/utils"

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

		// 4: 检查 Token 是否在黑名单（已退出登录）
		if cache.Client != nil {
			hash := sha256.Sum256([]byte(tokenString))
			blacklistKey := fmt.Sprintf("blacklist:%x", hash)
			if exists, _ := cache.Exists(c.Request.Context(), blacklistKey); exists {
				response.Unauthorized(c, "Token 已失效，请重新登录")
				return
			}
		}

		// 5: 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// CORS 跨域中间件，只允许白名单中的 Origin 跨域访问
func CORS(allowedOrigins []string) gin.HandlerFunc {
	originSet := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		if o = strings.TrimSpace(o); o != "" {
			originSet[o] = true
		}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := origin != "" && originSet[origin]

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}

		if c.Request.Method == "OPTIONS" {
			if allowed {
				c.AbortWithStatus(204)
			} else {
				c.AbortWithStatus(403)
			}
			return
		}

		c.Next()
	}
}

package middleware

import (
	"sync"

	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/pkg/logger"
	"adult-short-videos/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// IPRateLimiter IP限流器
type IPRateLimiter struct {
	ips map[string]*rate.Limiter // IP 地址对应的限流器
	mu  *sync.RWMutex            // 读写锁
	r   rate.Limit               // 每秒允许的请求数
	b   int                      // 突发请求数
}

// NewIPRateLimiter 创建限流器
// r: 每秒允许的请求数
// b: 突发请求数
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// GetLimiter 获取IP对应的限流器
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

// RateLimit 限流中间件
func RateLimit(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// 获取该IP的限流器
		l := limiter.GetLimiter(ip)

		// 检查是否允许通过
		if !l.Allow() {
			logger.Warn("请求被限流",
				zap.String("ip", ip),
				zap.String("path", c.Request.URL.Path),
			)

			response.Error(c, errors.CodeTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

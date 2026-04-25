package middleware

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	RequestIDKey    = "request_id"
	maxRequestIDLen = 128
)

// RequestID 为每个请求生成唯一 ID，写入 Context、响应头和日志字段。
// 若客户端已在 X-Request-ID 头中传入合法 ID 则复用，方便端到端追踪。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if !isValidRequestID(id) {
			id = generateRequestID()
		}
		c.Set(RequestIDKey, id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

func isValidRequestID(id string) bool {
	if len(id) == 0 || len(id) > maxRequestIDLen {
		return false
	}
	for _, ch := range id {
		if ch < 0x20 || ch > 0x7e {
			return false
		}
	}
	return true
}

func generateRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant RFC 4122
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

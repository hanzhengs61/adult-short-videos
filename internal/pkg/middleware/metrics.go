package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "HTTP 请求总数",
	}, []string{"method", "path", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP 请求耗时（秒）",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})
)

// Metrics 记录每个请求的耗时和状态码。
// path 使用 c.FullPath() 返回路由模板（如 /api/video/detail/:id），
// 避免高基数 cardinality 问题。
// excludedPaths 为需要跳过统计的路径（如健康检查、Prometheus 指标端点），
// 调用方可通过 cfg.Metrics.Path 动态传入，避免硬编码。
func Metrics(excludedPaths ...string) gin.HandlerFunc {
	excluded := make(map[string]struct{}, len(excludedPaths))
	for _, p := range excludedPaths {
		excluded[p] = struct{}{}
	}

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if _, skip := excluded[path]; skip {
			return
		}
		if path == "" {
			path = "unknown"
		}
		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
	}
}

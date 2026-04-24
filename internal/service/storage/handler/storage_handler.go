package handler

import (
	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/pkg/response"
	"adult-short-videos/internal/service/storage/proxy"

	"github.com/gin-gonic/gin"
)

// StorageHandler 存储处理器
type StorageHandler struct {
	proxy *proxy.RefererProxy
}

// NewStorageHandler 创建存储处理器
func NewStorageHandler() *StorageHandler {
	return &StorageHandler{
		proxy: proxy.NewRefererProxy(),
	}
}

// ProxyPlay 代理播放
// 路由: GET /api/storage/proxy
// 参数: url（源站视频URL）
func (h *StorageHandler) ProxyPlay(c *gin.Context) {
	// 获取目标 URL
	targetURL := c.Query("url")
	if targetURL == "" {
		response.Error(c, errors.CodeInvalidParam, "缺少 url 参数")
		return
	}

	// 代理请求
	err := h.proxy.ProxyRequest(targetURL, c.Writer, c.Request)
	if err != nil {
		// 代理失败
		response.HandleError(c, err)
		return
	}

	// 注意：这里不需要调用 response.Success
	// 因为 ProxyRequest 已经直接写入了响应
}

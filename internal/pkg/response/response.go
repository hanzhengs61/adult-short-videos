package response

import (
	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Response 响应结构体
type Response struct {
	Code int         `json:"code"`           // 状态码
	Msg  string      `json:"msg"`            // 提示信息
	Data interface{} `json:"data,omitempty"` // 数据（可选）
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: errors.CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

// SuccessWithMsg 成功响应（自定义消息）
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: errors.CodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// HandleError 处理错误（自动识别业务错误）
func HandleError(c *gin.Context, err error) {
	// 业务错误
	if bizErr := errors.GetBizError(err); bizErr != nil {
		Error(c, bizErr.Code, bizErr.Msg)
		return
	}

	// 系统错误
	logger.Error("系统错误",
		zap.Error(err),
		zap.String("path", c.Request.URL.Path),
	)
	Error(c, errors.CodeServerError, "系统错误")
}

// Unauthorized 未认证响应
func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code: errors.CodeUnauthorized,
		Msg:  msg,
		Data: nil,
	})
	c.Abort()
}

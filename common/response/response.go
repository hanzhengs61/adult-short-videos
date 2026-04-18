package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 响应结构体
type Response struct {
	Code int         `json:"code"`           // 状态码
	Msg  string      `json:"msg"`            // 提示信息
	Data interface{} `json:"data,omitempty"` // 数据（可选）
}

// 业务状态码定义
// 自己定义的业务码，不是 HTTP 状态码
const (
	SUCCESS       = 200 // 成功
	ERROR         = 500 // 服务器错误
	INVALID_PARAM = 400 // 参数错误
	UNAUTHORIZED  = 401 // 未认证
	FORBIDDEN     = 403 // 无权限
	NOT_FOUND     = 404 // 资源不存在
)

// Success 响应成功
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: SUCCESS,
		Msg:  "success",
		Data: data,
	})
}

// SuccessWithMsg 自定义成功响应
// 用法： response.SuccessWithMsg(c, "注册成功", nil)
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: SUCCESS,
		Msg:  msg,
		Data: data,
	})
}

// Error 响应错误
// 用法： response.Error(c, response.INVALID_PARAM, "用户名已存在")
func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

// Unauthorized 未认证响应
// 这个会返回 HTTP 401，并终止请求处理
func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code: UNAUTHORIZED,
		Msg:  msg,
		Data: nil,
	})
	c.Abort() // 终止后续处理
}

// ServerError 服务器错误
// 用于处理意外错误
func ServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code: ERROR,
		Msg:  msg,
		Data: nil,
	})
}

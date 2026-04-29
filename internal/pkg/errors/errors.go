package errors

import (
	"errors"
	"fmt"
)

// 错误码定义
const (
	CodeSuccess         = 200
	CodeInvalidParam    = 400
	CodeUnauthorized    = 401
	CodeForbidden       = 403
	CodeNotFound        = 404
	CodeTooManyRequests = 429
	CodeServerError     = 500

	// CodeDatabaseError 数据库错误 1xxx
	CodeDatabaseError = 1001
	CodeCacheError    = 1002

	// CodeUserExists 用户错误 2xxx
	CodeUserExists    = 2001
	CodeUserNotFound  = 2002
	CodePasswordError = 2003
	CodeTokenInvalid  = 2004
	CodeTokenExpired  = 2005
	CodeEmailExists   = 2006
	CodeEncryptError  = 2007

	// CodeVideoNotFound 视频错误 3xxx
	CodeVideoNotFound = 3001
	CodeVideoOffline  = 3002

	// CodeAlreadyFavorite 收藏错误 4xxx
	CodeAlreadyFavorite = 4001
	CodeNotFavorite     = 4002

	// CodeCommentTooLong 评论错误 5xxx
	CodeCommentTooLong = 5001
	CodeAlreadyLiked   = 5002
)

// BizError 业务错误
type BizError struct {
	Code int
	Msg  string
}

func (e *BizError) Error() string {
	return e.Msg
}

// New 创建业务错误
func New(code int, msg string) *BizError {
	return &BizError{
		Code: code,
		Msg:  msg,
	}
}

// Newf 创建格式化业务错误
func Newf(code int, format string, args ...interface{}) *BizError {
	return &BizError{
		Code: code,
		Msg:  fmt.Sprintf(format, args...),
	}
}

// 预定义错误
var (
	// ErrInvalidParam 通用错误
	ErrInvalidParam    = New(CodeInvalidParam, "参数错误")
	ErrUnauthorized    = New(CodeUnauthorized, "未登录")
	ErrForbidden       = New(CodeForbidden, "无权限")
	ErrNotFound        = New(CodeNotFound, "资源不存在")
	ErrTooManyRequests = New(CodeTooManyRequests, "请求过于频繁")
	ErrServerError     = New(CodeServerError, "服务器错误")

	// ErrUserExists 用户错误
	ErrUserExists    = New(CodeUserExists, "用户已存在")
	ErrUserNotFound  = New(CodeUserNotFound, "用户不存在")
	ErrPasswordError = New(CodePasswordError, "密码错误")
	ErrTokenInvalid  = New(CodeTokenInvalid, "Token无效")
	ErrTokenExpired  = New(CodeTokenExpired, "Token已过期")
	ErrEmailExists   = New(CodeEmailExists, "邮箱已被注册")
	ErrEncryptError  = New(CodeEncryptError, "密码加密失败")
	ErrUserDisabled  = New(CodeUserExists, "账号已被禁用")

	// ErrVideoNotFound 视频错误
	ErrVideoNotFound = New(CodeVideoNotFound, "视频不存在")
	ErrVideoOffline  = New(CodeVideoOffline, "视频已下架")

	// ErrAlreadyFavorite 收藏错误
	ErrAlreadyFavorite = New(CodeAlreadyFavorite, "已收藏")
	ErrNotFavorite     = New(CodeNotFavorite, "未收藏")

	// ErrCommentTooLong 评论错误
	ErrCommentTooLong = New(CodeCommentTooLong, "评论内容过长")
	ErrAlreadyLiked   = New(CodeAlreadyLiked, "已点赞")
)

// IsBizError 判断是否为业务错误
func IsBizError(err error) bool {
	var bizError *BizError
	ok := errors.As(err, &bizError)
	return ok
}

// GetBizError 获取业务错误
func GetBizError(err error) *BizError {
	var bizErr *BizError
	if errors.As(err, &bizErr) {
		return bizErr
	}
	return nil
}

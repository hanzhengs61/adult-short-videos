package errors

import (
	"errors"
	"fmt"
)

// 错误码定义
/**
通用错误 200...
数据库错误: 1xxx
用户错误: 2xxx
视频错误: 3xxx
收藏错误: 4xxx
评论错误: 5xxx
搜索错误: 6xxx
关注错误: 7xxx
**/
const (
	CodeSuccess         = 200 // 成功
	CodeInvalidParam    = 400 // 参数错误
	CodeUnauthorized    = 401 // 未授权
	CodeForbidden       = 403 // 拒绝访问
	CodeNotFound        = 404 // 资源不存在
	CodeTooManyRequests = 429 // 请求过于频繁
	CodeServerError     = 500 // 服务器错误

	CodeDatabaseError = 1001 // 数据库错误
	CodeCacheError    = 1002 // 缓存错误

	// CodeUserExists 用户错误 2xxx
	CodeUserExists    = 2001 // 用户已存在
	CodeUserNotFound  = 2002 // 用户不存在
	CodePasswordError = 2003 // 密码错误
	CodeTokenInvalid  = 2004 // token 无效
	CodeTokenExpired  = 2005 // token 过期
	CodeEmailExists   = 2006 // 邮箱已存在
	CodeEncryptError  = 2007 // 加密错误
	CodeUserDisabled  = 2008 // 用户被禁用

	CodeVideoNotFound = 3001 // 视频不存在
	CodeVideoOffline  = 3002 // 视频已下架

	CodeAlreadyFavorite = 4001 // 已经收藏了
	CodeNotFavorite     = 4002 // 未收藏

	CodeCommentTooLong      = 5001 // 评论内容过长
	CodeAlreadyLiked        = 5002 // 已经点赞评论了
	CodeCommentEmpty        = 5003 // 评论为空
	CodeCommentFailed       = 5004 // 评论失败
	CodeCommentListFailed   = 5005 // 查询评论列表失败
	CodeLikeCommentFailed   = 5007 // 评论点赞失败
	CodeNotLiked            = 5008 // 未点赞该评论
	CodeUnlikeCommentFailed = 5009 // 取消点赞失败

	CodeAlreadyFollowed = 7001 // 已关注该作者
	CodeNotFollowed     = 7002 // 未关注该作者
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

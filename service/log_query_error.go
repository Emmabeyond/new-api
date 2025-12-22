package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// QueryErrorType 查询错误类型
type QueryErrorType int

const (
	QueryErrorTypeUnknown QueryErrorType = iota
	QueryErrorTypeTimeout      // 用户可操作：缩小时间范围
	QueryErrorTypeRateLimit    // 用户可操作：稍后重试
	QueryErrorTypeConnection   // 系统错误：重试
	QueryErrorTypeInvalidParam // 用户可操作：修正参数
)

// QueryError 查询错误
type QueryError struct {
	Type       QueryErrorType `json:"type"`
	Code       string         `json:"code"`
	Message    string         `json:"message"`
	Suggestion string         `json:"suggestion"`
	RetryAfter time.Duration  `json:"-"`
	HTTPStatus int            `json:"-"`
}

// Error 实现 error 接口
func (e *QueryError) Error() string {
	return e.Message
}

// IsUserActionable 是否用户可操作
func (e *QueryError) IsUserActionable() bool {
	switch e.Type {
	case QueryErrorTypeTimeout, QueryErrorTypeRateLimit, QueryErrorTypeInvalidParam:
		return true
	default:
		return false
	}
}

// NewTimeoutError 创建超时错误
func NewTimeoutError() *QueryError {
	return &QueryError{
		Type:       QueryErrorTypeTimeout,
		Code:       "QUERY_TIMEOUT",
		Message:    "查询超时",
		Suggestion: "请缩小时间范围或添加过滤条件",
		HTTPStatus: http.StatusGatewayTimeout,
	}
}

// NewRateLimitError 创建速率限制错误
func NewRateLimitError(retryAfter time.Duration) *QueryError {
	return &QueryError{
		Type:       QueryErrorTypeRateLimit,
		Code:       "RATE_LIMIT_EXCEEDED",
		Message:    "请求过于频繁",
		Suggestion: fmt.Sprintf("请等待 %d 秒后重试", int(retryAfter.Seconds())),
		RetryAfter: retryAfter,
		HTTPStatus: http.StatusTooManyRequests,
	}
}

// NewConnectionError 创建连接错误
func NewConnectionError() *QueryError {
	return &QueryError{
		Type:       QueryErrorTypeConnection,
		Code:       "CONNECTION_ERROR",
		Message:    "数据库连接异常",
		Suggestion: "系统正在恢复，请稍后重试",
		HTTPStatus: http.StatusServiceUnavailable,
	}
}

// NewInvalidParamError 创建参数错误
func NewInvalidParamError(message string) *QueryError {
	return &QueryError{
		Type:       QueryErrorTypeInvalidParam,
		Code:       "INVALID_PARAM",
		Message:    message,
		Suggestion: "请检查并修正请求参数",
		HTTPStatus: http.StatusBadRequest,
	}
}

// NewUnknownError 创建未知错误
func NewUnknownError(err error) *QueryError {
	return &QueryError{
		Type:       QueryErrorTypeUnknown,
		Code:       "UNKNOWN_ERROR",
		Message:    err.Error(),
		Suggestion: "请稍后重试，如问题持续请联系管理员",
		HTTPStatus: http.StatusInternalServerError,
	}
}

// ClassifyError 对错误进行分类
func ClassifyError(err error) *QueryError {
	if err == nil {
		return nil
	}

	// 检查超时错误
	if IsQueryTimeout(err) {
		return NewTimeoutError()
	}

	// 检查上下文取消
	if errors.Is(err, context.Canceled) {
		return NewTimeoutError()
	}

	// 检查连接错误
	if isConnectionError(err) {
		return NewConnectionError()
	}

	// 检查 GORM 记录未找到错误
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NewInvalidParamError("未找到相关记录")
	}

	// 默认返回未知错误
	return NewUnknownError(err)
}

// isConnectionError 检查是否为连接错误
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	connectionKeywords := []string{
		"connection refused",
		"connection reset",
		"connection timed out",
		"no such host",
		"network is unreachable",
		"dial tcp",
		"broken pipe",
		"eof",
		"driver: bad connection",
		"invalid connection",
	}

	for _, keyword := range connectionKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}

// HandleQueryError 处理查询错误并返回 HTTP 响应
func HandleQueryError(c *gin.Context, err error) {
	queryErr := ClassifyError(err)
	if queryErr == nil {
		return
	}

	// 设置 Retry-After 头（如果适用）
	if queryErr.RetryAfter > 0 {
		c.Header("Retry-After", fmt.Sprintf("%d", int(queryErr.RetryAfter.Seconds())))
	}

	c.JSON(queryErr.HTTPStatus, gin.H{
		"success":    false,
		"message":    queryErr.Message,
		"code":       queryErr.Code,
		"suggestion": queryErr.Suggestion,
	})
}

// HandleQueryErrorWithStatus 处理查询错误并返回指定状态码的 HTTP 响应
func HandleQueryErrorWithStatus(c *gin.Context, err error, statusCode int) {
	queryErr := ClassifyError(err)
	if queryErr == nil {
		return
	}

	// 覆盖状态码
	queryErr.HTTPStatus = statusCode

	// 设置 Retry-After 头（如果适用）
	if queryErr.RetryAfter > 0 {
		c.Header("Retry-After", fmt.Sprintf("%d", int(queryErr.RetryAfter.Seconds())))
	}

	c.JSON(queryErr.HTTPStatus, gin.H{
		"success":    false,
		"message":    queryErr.Message,
		"code":       queryErr.Code,
		"suggestion": queryErr.Suggestion,
	})
}

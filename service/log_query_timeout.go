package service

import (
	"context"
	"errors"
	"time"

	"github.com/QuantumNous/new-api/common"
)

// 默认超时配置
const (
	DefaultQueryTimeout = 10 * time.Second // 默认超时 10 秒
	MaxQueryTimeout     = 30 * time.Second // 最大超时 30 秒
)

// QueryTimeoutConfig 查询超时配置
type QueryTimeoutConfig struct {
	DefaultTimeout time.Duration // 默认超时时间
	MaxTimeout     time.Duration // 最大允许超时时间
}

// DefaultQueryTimeoutConfig 返回默认的超时配置
func DefaultQueryTimeoutConfig() *QueryTimeoutConfig {
	return &QueryTimeoutConfig{
		DefaultTimeout: DefaultQueryTimeout,
		MaxTimeout:     MaxQueryTimeout,
	}
}

// globalQueryTimeoutConfig 全局超时配置
var globalQueryTimeoutConfig = DefaultQueryTimeoutConfig()

// SetQueryTimeoutConfig 设置全局超时配置
func SetQueryTimeoutConfig(config *QueryTimeoutConfig) {
	if config != nil {
		globalQueryTimeoutConfig = config
	}
}

// GetQueryTimeoutConfig 获取全局超时配置
func GetQueryTimeoutConfig() *QueryTimeoutConfig {
	return globalQueryTimeoutConfig
}

// WithQueryTimeout 为查询添加超时上下文
// 如果 timeout <= 0，使用默认超时
// 如果 timeout > MaxTimeout，使用 MaxTimeout
func WithQueryTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	config := globalQueryTimeoutConfig

	// 使用默认超时
	if timeout <= 0 {
		timeout = config.DefaultTimeout
	}

	// 限制最大超时
	if timeout > config.MaxTimeout {
		timeout = config.MaxTimeout
	}

	return context.WithTimeout(ctx, timeout)
}

// WithDefaultQueryTimeout 使用默认超时创建上下文
func WithDefaultQueryTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return WithQueryTimeout(ctx, 0)
}

// ErrQueryTimeout 查询超时错误
var ErrQueryTimeout = errors.New("query timeout")

// IsQueryTimeout 检查错误是否为查询超时
func IsQueryTimeout(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, ErrQueryTimeout)
}

// HandleQueryTimeout 处理查询超时，记录日志
func HandleQueryTimeout(err error, queryDesc string) {
	if IsQueryTimeout(err) {
		common.SysLog("query timeout: " + queryDesc)
	}
}

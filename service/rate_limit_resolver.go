package service

import (
	"github.com/QuantumNous/new-api/model"
)

// GetRateLimitForGroup 获取指定渠道分组的限流配置
// 优先级：group_rate_limits > rate_limit > nil（无限制）
//
// 参数：
//   - benefits: 等级权益配置，可以为 nil
//   - channelGroup: 渠道分组名称（如 "default", "vip" 等）
//
// 返回值：
//   - *RateLimitConfig: 限流配置，nil 表示无限制
func GetRateLimitForGroup(benefits *model.LevelBenefits, channelGroup string) *model.RateLimitConfig {
	if benefits == nil {
		return nil // 无限制
	}

	// 1. 优先检查渠道分组限流
	if benefits.GroupRateLimits != nil && channelGroup != "" {
		if limit, exists := benefits.GroupRateLimits[channelGroup]; exists && limit != nil {
			return limit
		}
	}

	// 2. 回退到全局限流
	if benefits.RateLimit != nil {
		return benefits.RateLimit
	}

	// 3. 无限制
	return nil
}

// GetRateLimitValues 获取限流配置的具体值
// 返回 totalCount, successCount, found
func GetRateLimitValues(config *model.RateLimitConfig) (int, int, bool) {
	if config == nil {
		return 0, 0, false
	}

	if config.TotalCount > 0 || config.SuccessCount > 0 {
		return config.TotalCount, config.SuccessCount, true
	}

	return 0, 0, false
}

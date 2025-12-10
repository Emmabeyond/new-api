package common

import (
	"errors"
	"time"
)

// MaxQueryDays 最大查询天数限制
const MaxQueryDays = 90

// 时间验证相关错误
var (
	ErrInvalidTimeRange  = errors.New("invalid time range: start time must be before end time")
	ErrInvalidTimestamp  = errors.New("invalid timestamp: timestamp must be positive")
	ErrTimeRangeTooLarge = errors.New("time range too large: automatically truncated to max allowed days")
)

// ValidateTimeRange 验证并规范化时间范围
// 如果范围超过 MaxQueryDays，自动截断为最近 MaxQueryDays 天（保持 endTimestamp 不变）
// 返回: 规范化后的 startTimestamp, endTimestamp, error
func ValidateTimeRange(startTimestamp, endTimestamp int64) (int64, int64, error) {
	// 验证时间戳是否为正数
	if startTimestamp < 0 || endTimestamp < 0 {
		return 0, 0, ErrInvalidTimestamp
	}

	// 验证开始时间是否小于结束时间
	if startTimestamp > endTimestamp {
		return 0, 0, ErrInvalidTimeRange
	}

	// 计算时间范围（秒）
	maxRangeSeconds := int64(MaxQueryDays * 24 * 60 * 60)
	timeRange := endTimestamp - startTimestamp

	// 如果范围超过最大限制，截断开始时间
	if timeRange > maxRangeSeconds {
		startTimestamp = endTimestamp - maxRangeSeconds
	}

	return startTimestamp, endTimestamp, nil
}

// ValidateTimeRangeWithDefault 验证时间范围，如果未提供则使用默认值
// defaultDays: 默认查询天数
func ValidateTimeRangeWithDefault(startTimestamp, endTimestamp int64, defaultDays int) (int64, int64, error) {
	now := time.Now().Unix()

	// 如果结束时间为 0，使用当前时间
	if endTimestamp == 0 {
		endTimestamp = now
	}

	// 如果开始时间为 0，使用默认天数
	if startTimestamp == 0 {
		startTimestamp = endTimestamp - int64(defaultDays*24*60*60)
	}

	return ValidateTimeRange(startTimestamp, endTimestamp)
}

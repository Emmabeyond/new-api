package operation_setting

import (
	"github.com/QuantumNous/new-api/setting/config"
)

// EmptyResponseSetting 空回复处理配置
type EmptyResponseSetting struct {
	// Enabled 是否启用空回复检测和重试
	Enabled bool `json:"enabled"`

	// MaxRetryCount 最大重试次数（针对空回复），默认为 2
	MaxRetryCount int `json:"max_retry_count"`

	// ExcludedModels 排除检测的模型列表（支持前缀匹配，逗号分隔）
	ExcludedModels string `json:"excluded_models"`

	// AlertThreshold 告警阈值（空回复率百分比），默认为 10
	AlertThreshold float64 `json:"alert_threshold"`

	// NonEmptyFinishReasons 非空 finish_reason 列表（逗号分隔）
	// 当 finish_reason 在此列表中时，即使内容为空也不判定为空回复
	// 默认值: tool_calls,tool_use,function_call
	NonEmptyFinishReasons string `json:"non_empty_finish_reasons"`
}

// 默认配置
var emptyResponseSetting = EmptyResponseSetting{
	Enabled:               true,
	MaxRetryCount:         2,
	ExcludedModels:        "",
	AlertThreshold:        10.0,
	NonEmptyFinishReasons: "tool_calls,tool_use,function_call",
}

// EmptyResponseConfigSetter 配置设置回调函数类型
type EmptyResponseConfigSetter func(enabled bool, maxRetryCount int, excludedModels []string, alertThreshold float64, nonEmptyFinishReasons []string)

// emptyResponseConfigSetter 配置设置回调
var emptyResponseConfigSetter EmptyResponseConfigSetter

// RegisterEmptyResponseConfigSetter 注册配置设置回调
func RegisterEmptyResponseConfigSetter(setter EmptyResponseConfigSetter) {
	emptyResponseConfigSetter = setter
}

func init() {
	// 注册到全局配置管理器
	config.GlobalConfig.Register("empty_response_setting", &emptyResponseSetting)
}

// GetEmptyResponseSetting 获取空回复配置
func GetEmptyResponseSetting() *EmptyResponseSetting {
	return &emptyResponseSetting
}

// IsEmptyResponseRetryEnabled 是否启用空回复重试
func IsEmptyResponseRetryEnabled() bool {
	return emptyResponseSetting.Enabled
}

// GetEmptyResponseMaxRetryCount 获取空回复最大重试次数
func GetEmptyResponseMaxRetryCount() int {
	if emptyResponseSetting.MaxRetryCount <= 0 {
		return 2
	}
	return emptyResponseSetting.MaxRetryCount
}

// GetEmptyResponseAlertThreshold 获取空回复告警阈值
func GetEmptyResponseAlertThreshold() float64 {
	if emptyResponseSetting.AlertThreshold <= 0 {
		return 10.0
	}
	return emptyResponseSetting.AlertThreshold
}

// SyncToService 同步配置到 service 层
func (s *EmptyResponseSetting) SyncToService() {
	if emptyResponseConfigSetter == nil {
		return
	}

	excludedModels := []string{}
	if s.ExcludedModels != "" {
		// 解析逗号分隔的模型列表
		for _, model := range splitAndTrim(s.ExcludedModels, ",") {
			if model != "" {
				excludedModels = append(excludedModels, model)
			}
		}
	}

	nonEmptyFinishReasons := []string{}
	if s.NonEmptyFinishReasons != "" {
		// 解析逗号分隔的 finish_reason 列表
		for _, reason := range splitAndTrim(s.NonEmptyFinishReasons, ",") {
			if reason != "" {
				nonEmptyFinishReasons = append(nonEmptyFinishReasons, reason)
			}
		}
	}

	emptyResponseConfigSetter(s.Enabled, s.MaxRetryCount, excludedModels, s.AlertThreshold, nonEmptyFinishReasons)
}

// splitAndTrim 分割字符串并去除空白
func splitAndTrim(s string, sep string) []string {
	parts := []string{}
	for _, part := range splitString(s, sep) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

// splitString 简单的字符串分割
func splitString(s string, sep string) []string {
	if s == "" {
		return []string{}
	}
	result := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

// trimSpace 去除字符串首尾空白
func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

// InitEmptyResponseSetting 初始化空回复配置（在配置加载后调用）
func InitEmptyResponseSetting() {
	emptyResponseSetting.SyncToService()
}

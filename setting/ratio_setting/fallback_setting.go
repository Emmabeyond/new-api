package ratio_setting

import (
	"math"

	"github.com/QuantumNous/new-api/setting/config"
)

// 内置默认倍率
const BuiltinDefaultRatio = 37.5

// RatioFallbackSetting 模型倍率兜底配置
type RatioFallbackSetting struct {
	// 是否启用兜底倍率（非自用模式下也生效）
	EnableFallbackRatio bool `json:"enable_fallback_ratio"`
	// 全局默认倍率
	DefaultRatio float64 `json:"default_ratio"`
}

// 默认配置
var ratioFallbackSetting = RatioFallbackSetting{
	EnableFallbackRatio: false,
	DefaultRatio:        BuiltinDefaultRatio,
}

func init() {
	// 注册到全局配置管理器
	config.GlobalConfig.Register("ratio_fallback_setting", &ratioFallbackSetting)
}

// GetRatioFallbackSetting 获取兜底配置
func GetRatioFallbackSetting() *RatioFallbackSetting {
	return &ratioFallbackSetting
}

// GetEffectiveDefaultRatio 获取有效的默认倍率
// 如果配置的默认倍率无效（<=0 或 NaN），返回内置默认值
func GetEffectiveDefaultRatio() float64 {
	if ratioFallbackSetting.DefaultRatio <= 0 || math.IsNaN(ratioFallbackSetting.DefaultRatio) {
		return BuiltinDefaultRatio
	}
	return ratioFallbackSetting.DefaultRatio
}

// IsFallbackRatioEnabled 检查是否启用兜底倍率
func IsFallbackRatioEnabled() bool {
	return ratioFallbackSetting.EnableFallbackRatio
}

// SetRatioFallbackSetting 设置兜底配置（用于测试或动态更新）
func SetRatioFallbackSetting(setting RatioFallbackSetting) {
	ratioFallbackSetting = setting
}

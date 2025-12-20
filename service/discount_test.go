package service

import (
	"testing"

	"github.com/QuantumNous/new-api/model"
)

// TestBackwardCompatibility_OnlyDiscountRatio 测试向后兼容：只有全局折扣的旧配置
func TestBackwardCompatibility_OnlyDiscountRatio(t *testing.T) {
	tests := []struct {
		name          string
		benefits      *model.LevelBenefits
		channelGroup  string
		expectedRatio float64
	}{
		{
			name: "旧配置：只有discount_ratio，无group_discount_ratios",
			benefits: &model.LevelBenefits{
				DiscountRatio:       0.8,
				GroupDiscountRatios: nil, // 旧配置没有这个字段
			},
			channelGroup:  "OpenAI",
			expectedRatio: 0.8, // 应该返回全局折扣
		},
		{
			name: "旧配置：只有discount_ratio，group_discount_ratios为空map",
			benefits: &model.LevelBenefits{
				DiscountRatio:       0.9,
				GroupDiscountRatios: map[string]float64{}, // 空map
			},
			channelGroup:  "Claude",
			expectedRatio: 0.9, // 应该返回全局折扣
		},
		{
			name: "旧配置：discount_ratio为1.0（无折扣）",
			benefits: &model.LevelBenefits{
				DiscountRatio:       1.0,
				GroupDiscountRatios: nil,
			},
			channelGroup:  "Gemini",
			expectedRatio: 1.0, // 应该返回1.0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ratio := GetDiscountForChannel(tt.benefits, tt.channelGroup)
			if ratio != tt.expectedRatio {
				t.Errorf("GetDiscountForChannel() = %v, want %v", ratio, tt.expectedRatio)
			}
		})
	}
}

// TestBackwardCompatibility_MixedConfiguration 测试混合配置（同时有新旧字段）
func TestBackwardCompatibility_MixedConfiguration(t *testing.T) {
	tests := []struct {
		name          string
		benefits      *model.LevelBenefits
		channelGroup  string
		expectedRatio float64
	}{
		{
			name: "混合配置：查询有专属折扣的渠道",
			benefits: &model.LevelBenefits{
				DiscountRatio: 0.9, // 全局折扣
				GroupDiscountRatios: map[string]float64{
					"OpenAI": 0.85, // 专属折扣
				},
			},
			channelGroup:  "OpenAI",
			expectedRatio: 0.85, // 应该返回专属折扣
		},
		{
			name: "混合配置：查询无专属折扣的渠道",
			benefits: &model.LevelBenefits{
				DiscountRatio: 0.9, // 全局折扣
				GroupDiscountRatios: map[string]float64{
					"OpenAI": 0.85, // 只配置了OpenAI
				},
			},
			channelGroup:  "Claude",
			expectedRatio: 0.9, // 应该回退到全局折扣
		},
		{
			name: "混合配置：多个渠道专属折扣",
			benefits: &model.LevelBenefits{
				DiscountRatio: 1.0, // 全局无折扣
				GroupDiscountRatios: map[string]float64{
					"OpenAI": 0.85,
					"Claude": 0.80,
					"Gemini": 0.90,
				},
			},
			channelGroup:  "Claude",
			expectedRatio: 0.80, // 应该返回Claude的专属折扣
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ratio := GetDiscountForChannel(tt.benefits, tt.channelGroup)
			if ratio != tt.expectedRatio {
				t.Errorf("GetDiscountForChannel() = %v, want %v", ratio, tt.expectedRatio)
			}
		})
	}
}

// TestBackwardCompatibility_EdgeCases 测试边界情况
func TestBackwardCompatibility_EdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		benefits      *model.LevelBenefits
		channelGroup  string
		expectedRatio float64
	}{
		{
			name:          "空配置：benefits为nil",
			benefits:      nil,
			channelGroup:  "OpenAI",
			expectedRatio: 1.0, // 应该返回默认值1.0
		},
		{
			name: "空配置：所有字段都是零值",
			benefits: &model.LevelBenefits{
				DiscountRatio:       0,
				GroupDiscountRatios: nil,
			},
			channelGroup:  "OpenAI",
			expectedRatio: 1.0, // 应该返回默认值1.0
		},
		{
			name: "空字符串渠道",
			benefits: &model.LevelBenefits{
				DiscountRatio: 0.8,
			},
			channelGroup:  "",
			expectedRatio: 0.8, // 应该返回全局折扣
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ratio := GetDiscountForChannel(tt.benefits, tt.channelGroup)
			if ratio != tt.expectedRatio {
				t.Errorf("GetDiscountForChannel() = %v, want %v", ratio, tt.expectedRatio)
			}
		})
	}
}

// TestBackwardCompatibility_NewConfiguration 测试新配置格式
func TestBackwardCompatibility_NewConfiguration(t *testing.T) {
	tests := []struct {
		name          string
		benefits      *model.LevelBenefits
		channelGroup  string
		expectedRatio float64
	}{
		{
			name: "新配置：只有group_discount_ratios，无discount_ratio",
			benefits: &model.LevelBenefits{
				DiscountRatio: 0, // 未设置全局折扣
				GroupDiscountRatios: map[string]float64{
					"OpenAI": 0.85,
					"Claude": 0.80,
				},
			},
			channelGroup:  "OpenAI",
			expectedRatio: 0.85, // 应该返回专属折扣
		},
		{
			name: "新配置：查询未配置的渠道，无全局折扣",
			benefits: &model.LevelBenefits{
				DiscountRatio: 0, // 未设置全局折扣
				GroupDiscountRatios: map[string]float64{
					"OpenAI": 0.85,
				},
			},
			channelGroup:  "Claude",
			expectedRatio: 1.0, // 应该返回默认值1.0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ratio := GetDiscountForChannel(tt.benefits, tt.channelGroup)
			if ratio != tt.expectedRatio {
				t.Errorf("GetDiscountForChannel() = %v, want %v", ratio, tt.expectedRatio)
			}
		})
	}
}

// TestCheckLegacyConfiguration 测试旧配置检测
func TestCheckLegacyConfiguration(t *testing.T) {
	tests := []struct {
		name           string
		levelId        string
		benefits       *model.LevelBenefits
		expectedLegacy bool
	}{
		{
			name:    "旧配置：只有discount_ratio",
			levelId: "tier_1",
			benefits: &model.LevelBenefits{
				DiscountRatio:       0.8,
				GroupDiscountRatios: nil,
			},
			expectedLegacy: true,
		},
		{
			name:    "旧配置：discount_ratio和空的group_discount_ratios",
			levelId: "tier_2",
			benefits: &model.LevelBenefits{
				DiscountRatio:       0.9,
				GroupDiscountRatios: map[string]float64{},
			},
			expectedLegacy: true,
		},
		{
			name:    "新配置：有group_discount_ratios",
			levelId: "tier_3",
			benefits: &model.LevelBenefits{
				DiscountRatio: 0.9,
				GroupDiscountRatios: map[string]float64{
					"OpenAI": 0.85,
				},
			},
			expectedLegacy: false,
		},
		{
			name:    "新配置：只有group_discount_ratios",
			levelId: "tier_4",
			benefits: &model.LevelBenefits{
				DiscountRatio: 0,
				GroupDiscountRatios: map[string]float64{
					"OpenAI": 0.85,
					"Claude": 0.80,
				},
			},
			expectedLegacy: false,
		},
		{
			name:    "无折扣配置：discount_ratio为1.0",
			levelId: "tier_5",
			benefits: &model.LevelBenefits{
				DiscountRatio:       1.0,
				GroupDiscountRatios: nil,
			},
			expectedLegacy: false, // 1.0表示无折扣，不算旧配置
		},
		{
			name:           "空配置：benefits为nil",
			levelId:        "tier_6",
			benefits:       nil,
			expectedLegacy: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isLegacy, message := CheckLegacyConfiguration(tt.levelId, tt.benefits)
			if isLegacy != tt.expectedLegacy {
				t.Errorf("CheckLegacyConfiguration() isLegacy = %v, want %v", isLegacy, tt.expectedLegacy)
			}
			if isLegacy && message == "" {
				t.Error("CheckLegacyConfiguration() should return a message for legacy configuration")
			}
			if !isLegacy && message != "" {
				t.Error("CheckLegacyConfiguration() should not return a message for non-legacy configuration")
			}
		})
	}
}

package service

import (
	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
)

// GetDiscountForChannel 获取指定等级和渠道的折扣倍率
// 优先返回渠道专属折扣，如果不存在则返回全局折扣
//
// 向后兼容性说明：
// - 旧配置：只有 discount_ratio 字段（全局折扣），无 group_discount_ratios
//   行为：所有渠道都使用 discount_ratio 的值
// - 新配置：有 group_discount_ratios 字段（渠道专属折扣）
//   行为：优先使用渠道专属折扣，未配置的渠道回退到 discount_ratio
// - 混合配置：同时存在两个字段
//   行为：已配置渠道使用专属折扣，未配置渠道使用全局折扣
//
// 参数：
//   - benefits: 等级权益配置，可以为 nil
//   - channelGroup: 渠道分组名称（如 "OpenAI", "Claude" 等）
//
// 返回值：
//   - 折扣倍率（0-1之间），1.0 表示无折扣
func GetDiscountForChannel(benefits *model.LevelBenefits, channelGroup string) float64 {
	if benefits == nil {
		return 1.0 // 默认无折扣
	}

	// 优先使用渠道专属折扣（新配置）
	if benefits.GroupDiscountRatios != nil {
		if discount, exists := benefits.GroupDiscountRatios[channelGroup]; exists {
			return discount
		}
	}

	// 回退到全局折扣（旧配置兼容）
	if benefits.DiscountRatio > 0 {
		return benefits.DiscountRatio
	}

	return 1.0 // 默认无折扣
}

// CheckLegacyConfiguration 检查是否使用旧配置格式并记录提示
// 这是一个可选的辅助函数，用于在系统启动或配置加载时检测旧配置
//
// 参数：
//   - levelId: 等级ID
//   - benefits: 等级权益配置
//
// 返回值：
//   - isLegacy: 是否为旧配置格式
//   - message: 提示信息
func CheckLegacyConfiguration(levelId string, benefits *model.LevelBenefits) (isLegacy bool, message string) {
	if benefits == nil {
		return false, ""
	}

	// 检查是否为旧配置：有 discount_ratio 但没有 group_discount_ratios
	hasGlobalDiscount := benefits.DiscountRatio > 0 && benefits.DiscountRatio < 1
	hasGroupDiscounts := benefits.GroupDiscountRatios != nil && len(benefits.GroupDiscountRatios) > 0

	if hasGlobalDiscount && !hasGroupDiscounts {
		message = "等级 " + levelId + " 使用旧配置格式（仅全局折扣）。" +
			"系统将继续正常工作，但建议升级到新的渠道级别折扣配置以获得更精细的控制。"
		return true, message
	}

	return false, ""
}

// LogLegacyConfigurationHints 在系统启动时检查并记录所有旧配置的提示
// 这是一个可选的辅助函数，可以在系统初始化时调用
func LogLegacyConfigurationHints() {
	levels, err := model.GetAllLevelConfigs()
	if err != nil {
		return
	}

	legacyCount := 0
	for _, level := range levels {
		benefits, err := level.GetBenefits()
		if err != nil {
			continue
		}

		isLegacy, message := CheckLegacyConfiguration(level.Id, benefits)
		if isLegacy {
			legacyCount++
			common.SysLog(message)
		}
	}

	if legacyCount > 0 {
		common.SysLog("检测到 " + string(rune(legacyCount)) + " 个等级使用旧配置格式。" +
			"所有功能将正常工作，但建议在管理界面升级到新的渠道级别折扣配置。")
	}
}

// CalculateDiscountedFee 计算折扣后的费用
func CalculateDiscountedFee(originalFee float64, discountRatio float64) float64 {
	if originalFee <= 0 {
		return 0
	}
	if discountRatio <= 0 || discountRatio > 1 {
		return originalFee // 无效折扣，返回原价
	}
	return originalFee * discountRatio
}

// ConvertToUSD 将充值金额转换为美元额度
// amount: 充值金额（人民币）
// rate: 兑换比率（例如：20元=100美元，rate=5）
func ConvertToUSD(amount float64, rate float64) float64 {
	if amount <= 0 || rate <= 0 {
		return 0
	}
	return amount * rate
}

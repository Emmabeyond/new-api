package service

import (
	"math/rand"
	"time"

	"github.com/QuantumNous/new-api/setting"
)

// 惊喜奖励和补签配置（不可通过系统设置覆盖的常量）
var (
	CheckinBonusProbability = 0.10 // 10% 概率触发惊喜奖励
	CheckinMakeupMaxDays    = 3    // 最多可补签天数
	CheckinMakeupRewardRate = 0.5  // 补签奖励比例（50%）
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// CalculateBaseReward 根据连续签到天数计算基础奖励
// 奖励层级：
// - 第1-6天: 1000 quota/天
// - 第7天: 5000 quota (周奖励)
// - 第8-13天: 1500 quota/天
// - 第14-29天: 2000 quota/天
// - 第30天及30的倍数: 20000 quota (月奖励)
// - 第31天+: 2500 quota/天 (封顶)
func CalculateBaseReward(consecutiveDays int) int {
	if consecutiveDays <= 0 {
		consecutiveDays = 1
	}

	// 检查是否是30的倍数（月奖励）
	if consecutiveDays >= 30 && consecutiveDays%30 == 0 {
		return setting.CheckinRewardDay30
	}

	// 检查是否是第7天或7的倍数（周奖励，但不与月奖励重叠）
	if consecutiveDays >= 7 && consecutiveDays%7 == 0 {
		return setting.CheckinRewardDay7
	}

	// 根据天数范围计算奖励
	switch {
	case consecutiveDays <= 6:
		return setting.CheckinRewardDay1To6
	case consecutiveDays <= 13:
		return setting.CheckinRewardDay8To13
	case consecutiveDays <= 29:
		return setting.CheckinRewardDay14To29
	default:
		return setting.CheckinRewardDay31Plus
	}
}

// CalculateMakeupReward 计算补签奖励（基础奖励的50%）
func CalculateMakeupReward(consecutiveDays int) int {
	baseReward := CalculateBaseReward(consecutiveDays)
	return int(float64(baseReward) * CheckinMakeupRewardRate)
}

// GenerateBonusReward 生成随机惊喜奖励
// 返回 (bonus, triggered)
// - bonus: 惊喜奖励金额（如果触发）
// - triggered: 是否触发了惊喜奖励
func GenerateBonusReward() (int, bool) {
	// 判断是否触发惊喜奖励
	if rand.Float64() >= CheckinBonusProbability {
		return 0, false
	}

	// 生成随机奖励金额 [CheckinBonusMin, CheckinBonusMax]
	bonus := setting.CheckinBonusMin + rand.Intn(setting.CheckinBonusMax-setting.CheckinBonusMin+1)
	return bonus, true
}

// GetMakeupCost 获取补签消耗额度
func GetMakeupCost() int {
	return setting.CheckinMakeupCost
}

// GetMakeupMaxDays 获取最大可补签天数
func GetMakeupMaxDays() int {
	return CheckinMakeupMaxDays
}

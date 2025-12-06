package setting

import (
	"strconv"

	"github.com/QuantumNous/new-api/common"
)

// 签到设置选项名称
const (
	OptionCheckinEnabled         = "CheckinEnabled"
	OptionCheckinRewardDay1To6   = "CheckinRewardDay1To6"
	OptionCheckinRewardDay7      = "CheckinRewardDay7"
	OptionCheckinRewardDay8To13  = "CheckinRewardDay8To13"
	OptionCheckinRewardDay14To29 = "CheckinRewardDay14To29"
	OptionCheckinRewardDay30     = "CheckinRewardDay30"
	OptionCheckinRewardDay31Plus = "CheckinRewardDay31Plus"
	OptionCheckinBonusMin        = "CheckinBonusMin"
	OptionCheckinBonusMax        = "CheckinBonusMax"
	OptionCheckinMakeupCost      = "CheckinMakeupCost"
)

// InitCheckinSettings 初始化签到设置
func InitCheckinSettings() {
	// 从数据库加载设置
	loadCheckinSettings()
}

// getOptionString 从 OptionMap 获取字符串值
func getOptionString(key string) string {
	common.OptionMapRWMutex.RLock()
	defer common.OptionMapRWMutex.RUnlock()
	return common.OptionMap[key]
}

// getOptionInt 从 OptionMap 获取整数值
func getOptionInt(key string) (int, bool) {
	val := getOptionString(key)
	if val == "" {
		return 0, false
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, false
	}
	return intVal, true
}

// 签到配置变量（由 loadCheckinSettings 更新）
var (
	CheckinEnabled         = true
	CheckinRewardDay1To6   = 1000
	CheckinRewardDay7      = 5000
	CheckinRewardDay8To13  = 1500
	CheckinRewardDay14To29 = 2000
	CheckinRewardDay30     = 20000
	CheckinRewardDay31Plus = 2500
	CheckinBonusMin        = 500
	CheckinBonusMax        = 5000
	CheckinMakeupCost      = 2000
)

// loadCheckinSettings 从 OptionMap 加载签到设置
func loadCheckinSettings() {
	// 签到功能开关
	if val := getOptionString(OptionCheckinEnabled); val != "" {
		CheckinEnabled = val == "true"
	}

	// 奖励配置
	if val, ok := getOptionInt(OptionCheckinRewardDay1To6); ok && val > 0 {
		CheckinRewardDay1To6 = val
	}
	if val, ok := getOptionInt(OptionCheckinRewardDay7); ok && val > 0 {
		CheckinRewardDay7 = val
	}
	if val, ok := getOptionInt(OptionCheckinRewardDay8To13); ok && val > 0 {
		CheckinRewardDay8To13 = val
	}
	if val, ok := getOptionInt(OptionCheckinRewardDay14To29); ok && val > 0 {
		CheckinRewardDay14To29 = val
	}
	if val, ok := getOptionInt(OptionCheckinRewardDay30); ok && val > 0 {
		CheckinRewardDay30 = val
	}
	if val, ok := getOptionInt(OptionCheckinRewardDay31Plus); ok && val > 0 {
		CheckinRewardDay31Plus = val
	}

	// 惊喜奖励配置
	if val, ok := getOptionInt(OptionCheckinBonusMin); ok && val > 0 {
		CheckinBonusMin = val
	}
	if val, ok := getOptionInt(OptionCheckinBonusMax); ok && val > 0 {
		CheckinBonusMax = val
	}

	// 补签配置
	if val, ok := getOptionInt(OptionCheckinMakeupCost); ok && val > 0 {
		CheckinMakeupCost = val
	}
}

// GetCheckinSettings 获取签到设置
func GetCheckinSettings() map[string]interface{} {
	return map[string]interface{}{
		"enabled":            CheckinEnabled,
		"reward_day_1_6":     CheckinRewardDay1To6,
		"reward_day_7":       CheckinRewardDay7,
		"reward_day_8_13":    CheckinRewardDay8To13,
		"reward_day_14_29":   CheckinRewardDay14To29,
		"reward_day_30":      CheckinRewardDay30,
		"reward_day_31_plus": CheckinRewardDay31Plus,
		"bonus_min":          CheckinBonusMin,
		"bonus_max":          CheckinBonusMax,
		"bonus_probability":  0.10,
		"makeup_cost":        CheckinMakeupCost,
		"makeup_max_days":    3,
	}
}

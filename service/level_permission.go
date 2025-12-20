package service

import (
	"github.com/QuantumNous/new-api/model"
)

// CheckUserLevelGroupPermission 检查用户等级是否有权访问指定渠道分组
func CheckUserLevelGroupPermission(userId int, groupName string) (bool, error) {
	groups, err := GetUserLevelAvailableGroups(userId)
	if err != nil {
		return false, err
	}

	for _, g := range groups {
		if g == groupName {
			return true, nil
		}
	}
	return false, nil
}

// GetUserLevelAvailableGroups 获取用户等级可用的渠道分组列表
func GetUserLevelAvailableGroups(userId int) ([]string, error) {
	user, err := model.GetUserById(userId, false)
	if err != nil {
		return nil, err
	}

	return GetLevelAvailableGroups(user.Level)
}

// GetLevelAvailableGroups 获取等级可用的渠道分组列表
func GetLevelAvailableGroups(levelId string) ([]string, error) {
	level, err := GetLevelConfigById(levelId)
	if err != nil {
		// 如果等级不存在，返回默认等级的分组
		level, err = GetDefaultLevelConfig()
		if err != nil {
			return []string{"default"}, nil
		}
	}

	benefits, err := level.GetBenefits()
	if err != nil {
		return []string{"default"}, nil
	}

	if len(benefits.AvailableChannelGroups) == 0 {
		return []string{"default"}, nil
	}

	return benefits.AvailableChannelGroups, nil
}

// GetUserLevelDiscountRatio 获取用户等级的优惠倍率
func GetUserLevelDiscountRatio(userId int, groupName string) (float64, error) {
	user, err := model.GetUserById(userId, false)
	if err != nil {
		return 1.0, err
	}

	return GetLevelDiscountRatio(user.Level, groupName)
}

// GetLevelDiscountRatio 获取等级对指定分组的优惠倍率
func GetLevelDiscountRatio(levelId string, groupName string) (float64, error) {
	level, err := GetLevelConfigById(levelId)
	if err != nil {
		return 1.0, nil
	}

	benefits, err := level.GetBenefits()
	if err != nil {
		return 1.0, nil
	}

	// 优先使用特定分组的倍率
	if benefits.GroupDiscountRatios != nil {
		if ratio, ok := benefits.GroupDiscountRatios[groupName]; ok {
			return ratio, nil
		}
	}

	// 使用全局优惠倍率
	if benefits.DiscountRatio > 0 {
		return benefits.DiscountRatio, nil
	}

	return 1.0, nil
}

// GetUserLevelRateLimit 获取用户等级的速率限制配置
func GetUserLevelRateLimit(userId int) (*model.RateLimitConfig, error) {
	user, err := model.GetUserById(userId, false)
	if err != nil {
		return nil, err
	}

	return GetLevelRateLimit(user.Level)
}

// GetLevelRateLimit 获取等级的速率限制配置
func GetLevelRateLimit(levelId string) (*model.RateLimitConfig, error) {
	level, err := GetLevelConfigById(levelId)
	if err != nil {
		return nil, nil
	}

	benefits, err := level.GetBenefits()
	if err != nil {
		return nil, nil
	}

	return benefits.RateLimit, nil
}

// GetUserLevelModelRateLimit 获取用户等级对特定模型的速率限制
func GetUserLevelModelRateLimit(userId int, modelName string) (*model.RateLimitConfig, error) {
	user, err := model.GetUserById(userId, false)
	if err != nil {
		return nil, err
	}

	return GetLevelModelRateLimit(user.Level, modelName)
}

// GetLevelModelRateLimit 获取等级对特定模型的速率限制
func GetLevelModelRateLimit(levelId string, modelName string) (*model.RateLimitConfig, error) {
	level, err := GetLevelConfigById(levelId)
	if err != nil {
		return nil, nil
	}

	benefits, err := level.GetBenefits()
	if err != nil {
		return nil, nil
	}

	// 优先使用模型级别的限制
	if benefits.ModelRateLimits != nil {
		if limit, ok := benefits.ModelRateLimits[modelName]; ok {
			return limit, nil
		}
	}

	// 返回全局限制
	return benefits.RateLimit, nil
}

// GetUserLevelInfo 获取用户等级完整信息
type UserLevelInfo struct {
	Level           *model.LevelConfig `json:"level"`
	Benefits        *model.LevelBenefits `json:"benefits"`
	UpgradeProgress *UpgradeProgress   `json:"upgrade_progress"`
}

func GetUserLevelInfo(userId int) (*UserLevelInfo, error) {
	user, err := model.GetUserById(userId, false)
	if err != nil {
		return nil, err
	}

	level, err := GetLevelConfigById(user.Level)
	if err != nil {
		level, err = GetDefaultLevelConfig()
		if err != nil {
			return nil, err
		}
	}

	benefits, err := level.GetBenefits()
	if err != nil {
		benefits = &model.LevelBenefits{DiscountRatio: 1.0}
	}

	progress, err := GetUserUpgradeProgress(userId)
	if err != nil {
		progress = nil
	}

	return &UserLevelInfo{
		Level:           level,
		Benefits:        benefits,
		UpgradeProgress: progress,
	}, nil
}

package service

import (
	"fmt"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
)

// CheckLevelUpgradeByUSD 基于美元额度检查用户是否满足升级条件
// 返回用户可以升级到的等级，如果没有可升级的等级则返回 nil
func CheckLevelUpgradeByUSD(userId int) (*model.LevelConfig, error) {
	// 获取用户信息
	user, err := model.GetUserById(userId, true)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 获取用户当前等级
	currentLevel, err := model.GetLevelConfigById(user.Level)
	if err != nil {
		return nil, fmt.Errorf("获取当前等级配置失败: %w", err)
	}

	// 查找更高等级
	higherLevel, err := model.GetHigherLevelConfig(user.CumulativeRecharge, currentLevel.Priority)
	if err != nil {
		return nil, fmt.Errorf("查询更高等级失败: %w", err)
	}

	return higherLevel, nil
}

// UpgradeUserLevel 升级用户等级
func UpgradeUserLevel(userId int, newLevelId string) error {
	user, err := model.GetUserById(userId, true)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 验证新等级是否存在
	newLevel, err := model.GetLevelConfigById(newLevelId)
	if err != nil {
		return fmt.Errorf("新等级不存在: %w", err)
	}

	// 验证用户是否满足升级条件
	if user.CumulativeRecharge < newLevel.MinCumulativeRecharge {
		return fmt.Errorf("用户累计充值不足，无法升级到该等级")
	}

	// 更新用户等级
	user.Level = newLevelId
	user.LevelUpdatedAt = common.GetTimestamp()
	err = user.Update(false)
	if err != nil {
		return fmt.Errorf("更新用户等级失败: %w", err)
	}

	// 记录日志
	model.RecordLog(userId, model.LogTypeSystem, fmt.Sprintf("等级升级: %s", newLevel.Name))

	return nil
}

// AutoUpgradeUserLevel 自动升级用户等级（充值后调用）
func AutoUpgradeUserLevel(userId int, amount float64) error {
	// 检查是否有可升级的等级
	higherLevel, err := CheckLevelUpgradeByUSD(userId)
	if err != nil {
		return err
	}

	// 如果没有可升级的等级，直接返回
	if higherLevel == nil {
		return nil
	}

	// 执行升级
	err = UpgradeUserLevel(userId, higherLevel.Id)
	if err != nil {
		return err
	}

	return nil
}

// 注册充值成功回调
func init() {
	model.RegisterTopUpSuccessCallback(AutoUpgradeUserLevel)
}

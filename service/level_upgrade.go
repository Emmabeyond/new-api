package service

import (
	"errors"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
)

// LevelChangeResult 等级变更结果
type LevelChangeResult struct {
	Changed      bool   `json:"changed"`
	OldLevel     string `json:"old_level"`
	NewLevel     string `json:"new_level"`
	ChangeReason string `json:"change_reason"`
}

// UpgradeProgress 升级进度
type UpgradeProgress struct {
	CurrentLevel     *model.LevelConfig `json:"current_level"`
	NextLevel        *model.LevelConfig `json:"next_level"`
	CurrentRecharge  float64            `json:"current_recharge"`
	RequiredRecharge float64            `json:"required_recharge"`
	ProgressPercent  float64            `json:"progress_percent"`
}

// CheckAndUpgradeUserLevel 检查并升级用户等级
func CheckAndUpgradeUserLevel(userId int) (*LevelChangeResult, error) {
	// 获取用户信息
	user, err := model.GetUserById(userId, false)
	if err != nil {
		return nil, err
	}

	// 获取当前等级配置
	currentLevel, err := GetLevelConfigById(user.Level)
	if err != nil {
		// 如果当前等级不存在，使用默认等级
		currentLevel, err = GetDefaultLevelConfig()
		if err != nil {
			return nil, err
		}
	}

	// 查找可升级的更高等级
	targetLevel, err := model.GetHigherLevelConfig(user.CumulativeRecharge, currentLevel.Priority)
	if err != nil {
		return nil, err
	}

	// 没有可升级的等级
	if targetLevel == nil {
		return &LevelChangeResult{
			Changed:  false,
			OldLevel: user.Level,
			NewLevel: user.Level,
		}, nil
	}

	// 执行升级
	oldLevel := user.Level
	newLevel := targetLevel.Id

	// 更新用户等级
	if err := UpdateUserLevel(userId, newLevel); err != nil {
		return nil, err
	}

	// 记录变更日志
	if err := model.CreateLevelChangeLog(userId, oldLevel, newLevel, model.LevelChangeReasonAutoUpgrade, 0); err != nil {
		common.SysLog("failed to create level change log: " + err.Error())
	}

	return &LevelChangeResult{
		Changed:      true,
		OldLevel:     oldLevel,
		NewLevel:     newLevel,
		ChangeReason: model.LevelChangeReasonAutoUpgrade,
	}, nil
}

// ManualSetUserLevel 手动调整用户等级
func ManualSetUserLevel(userId int, levelId string, reason string, operatorId int) error {
	// 验证目标等级是否存在
	_, err := GetLevelConfigById(levelId)
	if err != nil {
		return errors.New("目标等级不存在")
	}

	// 获取用户当前等级
	user, err := model.GetUserById(userId, false)
	if err != nil {
		return err
	}

	oldLevel := user.Level
	if oldLevel == levelId {
		return nil // 等级相同，无需变更
	}

	// 更新用户等级
	if err := UpdateUserLevel(userId, levelId); err != nil {
		return err
	}

	// 记录变更日志
	changeReason := model.LevelChangeReasonManualAdjust
	if reason != "" {
		changeReason = reason
	}
	return model.CreateLevelChangeLog(userId, oldLevel, levelId, changeReason, operatorId)
}


// UpdateUserLevel 更新用户等级
func UpdateUserLevel(userId int, levelId string) error {
	return model.DB.Model(&model.User{}).Where("id = ?", userId).Updates(map[string]interface{}{
		"level":            levelId,
		"level_updated_at": common.GetTimestamp(),
	}).Error
}

// UpdateUserCumulativeRecharge 更新用户累计充值金额
func UpdateUserCumulativeRecharge(userId int, amount float64) error {
	return model.DB.Model(&model.User{}).Where("id = ?", userId).
		Update("cumulative_recharge", model.DB.Raw("cumulative_recharge + ?", amount)).Error
}

// SyncUserCumulativeRecharge 同步用户累计充值金额（从充值记录计算）
func SyncUserCumulativeRecharge(userId int) (float64, error) {
	// 计算用户所有成功充值的总金额
	var totalMoney float64
	err := model.DB.Model(&model.TopUp{}).
		Where("user_id = ? AND status = ?", userId, common.TopUpStatusSuccess).
		Select("COALESCE(SUM(money), 0)").
		Scan(&totalMoney).Error
	if err != nil {
		return 0, err
	}

	// 更新用户累计充值金额
	err = model.DB.Model(&model.User{}).Where("id = ?", userId).
		Update("cumulative_recharge", totalMoney).Error
	if err != nil {
		return 0, err
	}

	return totalMoney, nil
}

// GetUserUpgradeProgress 获取用户升级进度
func GetUserUpgradeProgress(userId int) (*UpgradeProgress, error) {
	// 获取用户信息
	user, err := model.GetUserById(userId, false)
	if err != nil {
		return nil, err
	}

	// 获取当前等级配置
	currentLevel, err := GetLevelConfigById(user.Level)
	if err != nil {
		currentLevel, err = GetDefaultLevelConfig()
		if err != nil {
			return nil, err
		}
	}

	// 获取所有等级配置
	allLevels, err := GetAllLevelConfigs()
	if err != nil {
		return nil, err
	}

	// 查找下一个等级
	var nextLevel *model.LevelConfig
	for _, level := range allLevels {
		if level.Priority > currentLevel.Priority {
			if nextLevel == nil || level.Priority < nextLevel.Priority {
				nextLevel = level
			}
		}
	}

	progress := &UpgradeProgress{
		CurrentLevel:    currentLevel,
		NextLevel:       nextLevel,
		CurrentRecharge: user.CumulativeRecharge,
	}

	if nextLevel != nil {
		progress.RequiredRecharge = nextLevel.MinCumulativeRecharge
		if progress.RequiredRecharge > 0 {
			progress.ProgressPercent = (user.CumulativeRecharge / progress.RequiredRecharge) * 100
			if progress.ProgressPercent > 100 {
				progress.ProgressPercent = 100
			}
		}
	}

	return progress, nil
}

// OnUserTopUpSuccess 用户充值成功后的回调
func OnUserTopUpSuccess(userId int, amount float64) error {
	// 更新累计充值金额
	if err := UpdateUserCumulativeRecharge(userId, amount); err != nil {
		return err
	}

	// 检查并升级等级
	_, err := CheckAndUpgradeUserLevel(userId)
	return err
}


// InitLevelUpgradeCallback 初始化等级升级回调
func InitLevelUpgradeCallback() {
	model.RegisterTopUpSuccessCallback(OnUserTopUpSuccess)
}

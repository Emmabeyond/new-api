package model

import (
	"github.com/QuantumNous/new-api/common"
)

// LevelChangeLog 等级变更日志
type LevelChangeLog struct {
	Id           int    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId       int    `json:"user_id" gorm:"index"`
	OldLevel     string `json:"old_level" gorm:"type:varchar(32)"`
	NewLevel     string `json:"new_level" gorm:"type:varchar(32)"`
	ChangeReason string `json:"change_reason" gorm:"type:varchar(255)"` // auto_upgrade / manual_adjust
	OperatorId   int    `json:"operator_id" gorm:"default:0"`           // 操作者ID，0表示系统自动
	CreatedAt    int64  `json:"created_at" gorm:"autoCreateTime"`
}

const (
	LevelChangeReasonAutoUpgrade  = "auto_upgrade"   // 自动升级
	LevelChangeReasonManualAdjust = "manual_adjust"  // 手动调整
	LevelChangeReasonMigration    = "migration"      // 数据迁移
)

func (LevelChangeLog) TableName() string {
	return "level_change_logs"
}

// CreateLevelChangeLog 创建等级变更日志
func CreateLevelChangeLog(userId int, oldLevel, newLevel, reason string, operatorId int) error {
	log := &LevelChangeLog{
		UserId:       userId,
		OldLevel:     oldLevel,
		NewLevel:     newLevel,
		ChangeReason: reason,
		OperatorId:   operatorId,
	}
	return DB.Create(log).Error
}

// GetUserLevelChangeLogs 获取用户等级变更日志
func GetUserLevelChangeLogs(userId int, pageInfo *common.PageInfo) ([]*LevelChangeLog, int64, error) {
	var logs []*LevelChangeLog
	var total int64

	tx := DB.Begin()
	if tx.Error != nil {
		return nil, 0, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&LevelChangeLog{}).Where("user_id = ?", userId).Count(&total).Error; err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	if err := tx.Where("user_id = ?", userId).Order("id desc").
		Limit(pageInfo.GetPageSize()).Offset(pageInfo.GetStartIdx()).Find(&logs).Error; err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetAllLevelChangeLogs 获取所有等级变更日志（管理员）
func GetAllLevelChangeLogs(pageInfo *common.PageInfo) ([]*LevelChangeLog, int64, error) {
	var logs []*LevelChangeLog
	var total int64

	tx := DB.Begin()
	if tx.Error != nil {
		return nil, 0, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&LevelChangeLog{}).Count(&total).Error; err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	if err := tx.Order("id desc").Limit(pageInfo.GetPageSize()).Offset(pageInfo.GetStartIdx()).Find(&logs).Error; err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

package model

import (
	"time"
)

// CheckinAudit 签到审计日志表
type CheckinAudit struct {
	Id        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	AdminId   int       `json:"admin_id" gorm:"index;not null"`
	UserId    int       `json:"user_id" gorm:"index;not null"`
	Action    string    `json:"action" gorm:"type:varchar(50);not null"`
	OldValue  string    `json:"old_value" gorm:"type:varchar(255)"`
	NewValue  string    `json:"new_value" gorm:"type:varchar(255)"`
	CreatedAt time.Time `json:"created_at"`
}

// 审计操作类型常量
const (
	CheckinAuditActionAdjustConsecutive = "adjust_consecutive_days"
	CheckinAuditActionEnableFeature     = "enable_checkin"
	CheckinAuditActionDisableFeature    = "disable_checkin"
)

// CreateCheckinAudit 创建审计日志
func CreateCheckinAudit(audit *CheckinAudit) error {
	return DB.Create(audit).Error
}

// GetCheckinAuditsByUserId 获取指定用户的审计日志
func GetCheckinAuditsByUserId(userId int) ([]*CheckinAudit, error) {
	var audits []*CheckinAudit
	err := DB.Where("user_id = ?", userId).Order("created_at DESC").Find(&audits).Error
	return audits, err
}

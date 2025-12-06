package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Checkin 用户签到状态表
type Checkin struct {
	Id              int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId          int       `json:"user_id" gorm:"uniqueIndex;not null"`
	ConsecutiveDays int       `json:"consecutive_days" gorm:"default:0"`
	TotalCheckins   int       `json:"total_checkins" gorm:"default:0"`
	TotalQuota      int       `json:"total_quota" gorm:"default:0"`
	LastCheckinDate time.Time `json:"last_checkin_date" gorm:"index"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// GetCheckinByUserId 根据用户ID获取签到状态
func GetCheckinByUserId(userId int) (*Checkin, error) {
	var checkin Checkin
	err := DB.Where("user_id = ?", userId).First(&checkin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &checkin, nil
}

// GetOrCreateCheckin 获取或创建用户签到状态
func GetOrCreateCheckin(userId int) (*Checkin, error) {
	checkin, err := GetCheckinByUserId(userId)
	if err != nil {
		return nil, err
	}
	if checkin == nil {
		checkin = &Checkin{
			UserId:          userId,
			ConsecutiveDays: 0,
			TotalCheckins:   0,
			TotalQuota:      0,
		}
		if err := DB.Create(checkin).Error; err != nil {
			return nil, err
		}
	}
	return checkin, nil
}

// Update 更新签到状态
func (c *Checkin) Update() error {
	return DB.Save(c).Error
}

// UpdateConsecutiveDays 更新连续签到天数
func UpdateConsecutiveDays(userId int, days int) error {
	return DB.Model(&Checkin{}).Where("user_id = ?", userId).Update("consecutive_days", days).Error
}

// GetCheckinStats 获取签到统计数据
func GetCheckinStats() (todayCheckins int64, activeUsers int64, avgConsecutive float64, err error) {
	today := time.Now().Truncate(24 * time.Hour)
	
	// 今日签到数
	err = DB.Model(&CheckinRecord{}).Where("checkin_date >= ?", today).Count(&todayCheckins).Error
	if err != nil {
		return
	}
	
	// 活跃用户数（有签到记录的用户）
	err = DB.Model(&Checkin{}).Where("total_checkins > 0").Count(&activeUsers).Error
	if err != nil {
		return
	}
	
	// 平均连续天数
	var result struct {
		Avg float64
	}
	err = DB.Model(&Checkin{}).Select("AVG(consecutive_days) as avg").Where("consecutive_days > 0").Scan(&result).Error
	if err != nil {
		return
	}
	avgConsecutive = result.Avg
	
	return
}

// GetTodayQuotaDistributed 获取今日发放的总额度
func GetTodayQuotaDistributed() (int64, error) {
	today := time.Now().Truncate(24 * time.Hour)
	var result struct {
		Total int64
	}
	err := DB.Model(&CheckinRecord{}).
		Select("COALESCE(SUM(base_reward + bonus_reward), 0) as total").
		Where("checkin_date >= ?", today).
		Scan(&result).Error
	return result.Total, err
}

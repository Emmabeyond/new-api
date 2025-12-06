package model

import (
	"time"

	"github.com/QuantumNous/new-api/common"
)

// 签到类型常量
const (
	CheckinTypeNormal = 1 // 正常签到
	CheckinTypeMakeup = 2 // 补签
)

// CheckinRecord 签到记录表
type CheckinRecord struct {
	Id          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId      int       `json:"user_id" gorm:"index;not null"`
	CheckinDate time.Time `json:"checkin_date" gorm:"index;not null"`
	CheckinType int       `json:"checkin_type" gorm:"default:1"` // 1: normal, 2: makeup
	BaseReward  int       `json:"base_reward" gorm:"default:0"`
	BonusReward int       `json:"bonus_reward" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateCheckinRecord 创建签到记录
func CreateCheckinRecord(record *CheckinRecord) error {
	return DB.Create(record).Error
}

// GetCheckinRecordByDate 获取指定日期的签到记录
func GetCheckinRecordByDate(userId int, date time.Time) (*CheckinRecord, error) {
	var record CheckinRecord
	startOfDay := date.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	err := DB.Where("user_id = ? AND checkin_date >= ? AND checkin_date < ?", userId, startOfDay, endOfDay).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// HasCheckinRecord 检查指定日期是否有签到记录
func HasCheckinRecord(userId int, date time.Time) (bool, error) {
	startOfDay := date.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	var count int64
	err := DB.Model(&CheckinRecord{}).Where("user_id = ? AND checkin_date >= ? AND checkin_date < ?", userId, startOfDay, endOfDay).Count(&count).Error
	return count > 0, err
}

// GetCheckinHistory 获取签到历史记录（分页）
func GetCheckinHistory(userId int, pageInfo *common.PageInfo) ([]*CheckinRecord, int64, error) {
	var records []*CheckinRecord
	var total int64
	
	query := DB.Model(&CheckinRecord{}).Where("user_id = ?", userId)
	
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	err = query.Order("checkin_date DESC").
		Limit(pageInfo.GetPageSize()).
		Offset(pageInfo.GetStartIdx()).
		Find(&records).Error
	
	return records, total, err
}

// GetMonthlyCheckinDates 获取指定月份的签到日期列表
func GetMonthlyCheckinDates(userId int, year int, month int) ([]int, error) {
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)
	
	var records []CheckinRecord
	err := DB.Where("user_id = ? AND checkin_date >= ? AND checkin_date < ?", userId, startOfMonth, endOfMonth).
		Select("checkin_date").
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	
	days := make([]int, 0, len(records))
	for _, record := range records {
		days = append(days, record.CheckinDate.Day())
	}
	return days, nil
}

// GetAllCheckinRecords 获取所有签到记录（管理员用，支持筛选）
func GetAllCheckinRecords(pageInfo *common.PageInfo, userId int, startDate, endDate *time.Time, checkinType int) ([]*CheckinRecord, int64, error) {
	var records []*CheckinRecord
	var total int64
	
	query := DB.Model(&CheckinRecord{})
	
	if userId > 0 {
		query = query.Where("user_id = ?", userId)
	}
	if startDate != nil {
		query = query.Where("checkin_date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("checkin_date < ?", *endDate)
	}
	if checkinType > 0 {
		query = query.Where("checkin_type = ?", checkinType)
	}
	
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	err = query.Order("created_at DESC").
		Limit(pageInfo.GetPageSize()).
		Offset(pageInfo.GetStartIdx()).
		Find(&records).Error
	
	return records, total, err
}

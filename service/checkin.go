package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting"
)

var (
	ErrAlreadyCheckedIn    = errors.New("今日已签到")
	ErrCheckinDisabled     = errors.New("签到功能已关闭")
	ErrMakeupDateInvalid   = errors.New("补签日期无效，只能补签最近3天内的日期")
	ErrMakeupAlreadyExists = errors.New("该日期已有签到记录")
	ErrInsufficientQuota   = errors.New("额度不足，无法补签")
)

// HasCheckedInToday 检查用户今日是否已签到
func HasCheckedInToday(userId int) (bool, error) {
	today := time.Now().Truncate(24 * time.Hour)
	return model.HasCheckinRecord(userId, today)
}

// DoCheckin 执行签到
func DoCheckin(userId int) (*dto.CheckinResult, error) {
	// 检查签到功能是否启用
	if !setting.CheckinEnabled {
		return nil, ErrCheckinDisabled
	}

	// 检查今日是否已签到
	checked, err := HasCheckedInToday(userId)
	if err != nil {
		return nil, err
	}
	if checked {
		return nil, ErrAlreadyCheckedIn
	}

	// 获取或创建用户签到状态
	checkin, err := model.GetOrCreateCheckin(userId)
	if err != nil {
		return nil, err
	}

	// 计算连续签到天数
	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)
	
	var newConsecutiveDays int
	if checkin.LastCheckinDate.IsZero() {
		// 首次签到
		newConsecutiveDays = 1
	} else {
		lastCheckinDay := checkin.LastCheckinDate.Truncate(24 * time.Hour)
		if lastCheckinDay.Equal(yesterday) {
			// 连续签到
			newConsecutiveDays = checkin.ConsecutiveDays + 1
		} else {
			// 断签，重置为1
			newConsecutiveDays = 1
		}
	}

	// 计算奖励
	baseReward := CalculateBaseReward(newConsecutiveDays)
	bonusReward, bonusTriggered := GenerateBonusReward()
	totalReward := baseReward + bonusReward

	// 开始事务
	tx := model.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建签到记录
	record := &model.CheckinRecord{
		UserId:      userId,
		CheckinDate: time.Now(),
		CheckinType: model.CheckinTypeNormal,
		BaseReward:  baseReward,
		BonusReward: bonusReward,
	}
	if err := tx.Create(record).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新用户签到状态
	checkin.ConsecutiveDays = newConsecutiveDays
	checkin.TotalCheckins++
	checkin.TotalQuota += totalReward
	checkin.LastCheckinDate = time.Now()
	if err := tx.Save(checkin).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 增加用户额度
	if err := tx.Model(&model.User{}).Where("id = ?", userId).
		Update("quota", model.DB.Raw("quota + ?", totalReward)).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 记录日志
	model.RecordLog(userId, model.LogTypeSystem, fmt.Sprintf("签到成功，获得 %s（基础 %s%s）",
		logger.LogQuota(totalReward),
		logger.LogQuota(baseReward),
		func() string {
			if bonusTriggered {
				return fmt.Sprintf("，惊喜奖励 %s", logger.LogQuota(bonusReward))
			}
			return ""
		}()))

	return &dto.CheckinResult{
		Success:         true,
		BaseReward:      baseReward,
		BonusReward:     bonusReward,
		TotalReward:     totalReward,
		ConsecutiveDays: newConsecutiveDays,
		BonusTriggered:  bonusTriggered,
	}, nil
}

// GetUserCheckinStats 获取用户签到统计
func GetUserCheckinStats(userId int) (*dto.CheckinStats, error) {
	checkin, err := model.GetCheckinByUserId(userId)
	if err != nil {
		return nil, err
	}

	checkedToday, err := HasCheckedInToday(userId)
	if err != nil {
		return nil, err
	}

	if checkin == nil {
		return &dto.CheckinStats{
			ConsecutiveDays: 0,
			TotalCheckins:   0,
			TotalQuota:      0,
			CheckedInToday:  false,
		}, nil
	}

	return &dto.CheckinStats{
		ConsecutiveDays: checkin.ConsecutiveDays,
		TotalCheckins:   checkin.TotalCheckins,
		TotalQuota:      checkin.TotalQuota,
		LastCheckinDate: checkin.LastCheckinDate,
		CheckedInToday:  checkedToday,
	}, nil
}

// GetCheckinHistory 获取签到历史
func GetCheckinHistory(userId int, pageInfo *common.PageInfo) ([]*model.CheckinRecord, int64, error) {
	return model.GetCheckinHistory(userId, pageInfo)
}

// GetMonthlyCalendar 获取月度签到日历
func GetMonthlyCalendar(userId int, year, month int) ([]int, error) {
	return model.GetMonthlyCheckinDates(userId, year, month)
}


// DoMakeupCheckin 执行补签
func DoMakeupCheckin(userId int, targetDate time.Time) (*dto.CheckinResult, error) {
	// 检查签到功能是否启用
	if !setting.CheckinEnabled {
		return nil, ErrCheckinDisabled
	}

	// 标准化目标日期
	targetDay := targetDate.Truncate(24 * time.Hour)
	today := time.Now().Truncate(24 * time.Hour)

	// 检查目标日期是否在允许范围内（最近3天）
	maxDays := GetMakeupMaxDays()
	earliestAllowed := today.AddDate(0, 0, -maxDays)
	if targetDay.Before(earliestAllowed) || !targetDay.Before(today) {
		return nil, ErrMakeupDateInvalid
	}

	// 检查目标日期是否已有签到记录
	exists, err := model.HasCheckinRecord(userId, targetDay)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrMakeupAlreadyExists
	}

	// 获取用户当前额度
	userQuota, err := model.GetUserQuota(userId, true)
	if err != nil {
		return nil, err
	}

	// 检查额度是否足够
	makeupCost := GetMakeupCost()
	if userQuota < makeupCost {
		return nil, ErrInsufficientQuota
	}

	// 获取或创建用户签到状态
	checkin, err := model.GetOrCreateCheckin(userId)
	if err != nil {
		return nil, err
	}

	// 计算补签后的连续天数
	// 需要重新计算从目标日期开始的连续天数链
	newConsecutiveDays := calculateConsecutiveDaysWithMakeup(userId, targetDay, checkin)

	// 计算补签奖励（50%）
	baseReward := CalculateMakeupReward(newConsecutiveDays)

	// 开始事务
	tx := model.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 扣除补签费用
	if err := tx.Model(&model.User{}).Where("id = ? AND quota >= ?", userId, makeupCost).
		Update("quota", model.DB.Raw("quota - ?", makeupCost)).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建补签记录
	record := &model.CheckinRecord{
		UserId:      userId,
		CheckinDate: targetDay,
		CheckinType: model.CheckinTypeMakeup,
		BaseReward:  baseReward,
		BonusReward: 0, // 补签不触发惊喜奖励
	}
	if err := tx.Create(record).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新用户签到状态
	checkin.ConsecutiveDays = newConsecutiveDays
	checkin.TotalCheckins++
	checkin.TotalQuota += baseReward
	// 如果补签日期比当前最后签到日期更新，则更新最后签到日期
	if targetDay.After(checkin.LastCheckinDate) {
		checkin.LastCheckinDate = targetDay
	}
	if err := tx.Save(checkin).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 增加用户额度（补签奖励）
	if err := tx.Model(&model.User{}).Where("id = ?", userId).
		Update("quota", model.DB.Raw("quota + ?", baseReward)).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 记录日志
	model.RecordLog(userId, model.LogTypeSystem, fmt.Sprintf("补签成功（%s），消耗 %s，获得 %s",
		targetDay.Format("2006-01-02"),
		logger.LogQuota(makeupCost),
		logger.LogQuota(baseReward)))

	return &dto.CheckinResult{
		Success:         true,
		BaseReward:      baseReward,
		BonusReward:     0,
		TotalReward:     baseReward,
		ConsecutiveDays: newConsecutiveDays,
		BonusTriggered:  false,
	}, nil
}

// calculateConsecutiveDaysWithMakeup 计算补签后的连续天数
func calculateConsecutiveDaysWithMakeup(userId int, targetDay time.Time, checkin *model.Checkin) int {
	// 简化实现：从目标日期向前查找连续签到记录
	consecutiveDays := 1
	checkDay := targetDay.AddDate(0, 0, -1)

	for {
		exists, err := model.HasCheckinRecord(userId, checkDay)
		if err != nil || !exists {
			break
		}
		consecutiveDays++
		checkDay = checkDay.AddDate(0, 0, -1)
	}

	// 向后检查是否有连续的签到记录
	checkDay = targetDay.AddDate(0, 0, 1)
	today := time.Now().Truncate(24 * time.Hour)
	for !checkDay.After(today) {
		exists, err := model.HasCheckinRecord(userId, checkDay)
		if err != nil || !exists {
			break
		}
		consecutiveDays++
		checkDay = checkDay.AddDate(0, 0, 1)
	}

	return consecutiveDays
}

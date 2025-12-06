package dto

import "time"

// CheckinResult 签到结果
type CheckinResult struct {
	Success         bool `json:"success"`
	BaseReward      int  `json:"base_reward"`
	BonusReward     int  `json:"bonus_reward"`
	TotalReward     int  `json:"total_reward"`
	ConsecutiveDays int  `json:"consecutive_days"`
	BonusTriggered  bool `json:"bonus_triggered"`
}

// CheckinStats 签到统计
type CheckinStats struct {
	ConsecutiveDays int       `json:"consecutive_days"`
	TotalCheckins   int       `json:"total_checkins"`
	TotalQuota      int       `json:"total_quota"`
	LastCheckinDate time.Time `json:"last_checkin_date"`
	CheckedInToday  bool      `json:"checked_in_today"`
}

// MakeupCheckinRequest 补签请求
type MakeupCheckinRequest struct {
	TargetDate string `json:"target_date" binding:"required"` // YYYY-MM-DD
}

// AdminCheckinDashboard 管理员签到仪表盘
type AdminCheckinDashboard struct {
	TodayCheckins         int64   `json:"today_checkins"`
	TodayQuotaDistributed int64   `json:"today_quota_distributed"`
	ActiveUsers           int64   `json:"active_users"`
	AvgConsecutiveDays    float64 `json:"avg_consecutive_days"`
}

// AdminAdjustConsecutiveRequest 管理员调整连续天数请求
type AdminAdjustConsecutiveRequest struct {
	ConsecutiveDays int `json:"consecutive_days" binding:"required,min=0"`
}

// CheckinCalendarRequest 签到日历请求
type CheckinCalendarRequest struct {
	Year  int `form:"year"`
	Month int `form:"month"`
}

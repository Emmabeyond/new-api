package model

import (
	"time"

	"github.com/QuantumNous/new-api/common"
)

// LogStatsAggregation 日志统计聚合表
// 用于存储按小时预聚合的统计数据，加速统计查询
type LogStatsAggregation struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	HourStart    time.Time `json:"hour_start" gorm:"index:idx_hour_user;index:idx_hour_model;index:idx_hour_channel;uniqueIndex:idx_unique_agg,priority:1"`
	UserID       int       `json:"user_id" gorm:"index:idx_hour_user;uniqueIndex:idx_unique_agg,priority:2;default:0"`
	Username     string    `json:"username" gorm:"index;default:''"`
	ModelName    string    `json:"model_name" gorm:"index:idx_hour_model;uniqueIndex:idx_unique_agg,priority:3;default:''"`
	ChannelID    int       `json:"channel_id" gorm:"index:idx_hour_channel;uniqueIndex:idx_unique_agg,priority:4;default:0"`
	GroupName    string    `json:"group_name" gorm:"uniqueIndex:idx_unique_agg,priority:5;default:''"`
	TotalQuota   int64     `json:"total_quota" gorm:"default:0"`
	TotalTokens  int64     `json:"total_tokens" gorm:"default:0"` // prompt_tokens + completion_tokens
	RequestCount int64     `json:"request_count" gorm:"default:0"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (LogStatsAggregation) TableName() string {
	return "log_stats_aggregations"
}

// CreateLogStatsAggregation 创建或更新聚合记录
func CreateLogStatsAggregation(agg *LogStatsAggregation) error {
	return LOG_DB.Create(agg).Error
}

// UpsertLogStatsAggregation 插入或更新聚合记录（使用 ON CONFLICT）
func UpsertLogStatsAggregation(agg *LogStatsAggregation) error {
	// 使用 GORM 的 Clauses 实现 upsert
	if common.UsingPostgreSQL || common.LogSqlType == common.DatabaseTypePostgreSQL {
		return LOG_DB.Exec(`
			INSERT INTO log_stats_aggregations 
				(hour_start, user_id, username, model_name, channel_id, group_name, total_quota, total_tokens, request_count, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (hour_start, user_id, model_name, channel_id, group_name) 
			DO UPDATE SET 
				total_quota = EXCLUDED.total_quota,
				total_tokens = EXCLUDED.total_tokens,
				request_count = EXCLUDED.request_count,
				updated_at = EXCLUDED.updated_at
		`, agg.HourStart, agg.UserID, agg.Username, agg.ModelName, agg.ChannelID, agg.GroupName,
			agg.TotalQuota, agg.TotalTokens, agg.RequestCount, agg.UpdatedAt).Error
	}

	// MySQL 使用 ON DUPLICATE KEY UPDATE
	if common.UsingMySQL || common.LogSqlType == common.DatabaseTypeMySQL {
		return LOG_DB.Exec(`
			INSERT INTO log_stats_aggregations 
				(hour_start, user_id, username, model_name, channel_id, group_name, total_quota, total_tokens, request_count, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE 
				total_quota = VALUES(total_quota),
				total_tokens = VALUES(total_tokens),
				request_count = VALUES(request_count),
				updated_at = VALUES(updated_at)
		`, agg.HourStart, agg.UserID, agg.Username, agg.ModelName, agg.ChannelID, agg.GroupName,
			agg.TotalQuota, agg.TotalTokens, agg.RequestCount, agg.UpdatedAt).Error
	}

	// SQLite 使用 INSERT OR REPLACE
	return LOG_DB.Exec(`
		INSERT OR REPLACE INTO log_stats_aggregations 
			(hour_start, user_id, username, model_name, channel_id, group_name, total_quota, total_tokens, request_count, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, agg.HourStart, agg.UserID, agg.Username, agg.ModelName, agg.ChannelID, agg.GroupName,
		agg.TotalQuota, agg.TotalTokens, agg.RequestCount, agg.UpdatedAt).Error
}

// GetAggregatedStatsByTimeRange 获取指定时间范围内的聚合统计数据
func GetAggregatedStatsByTimeRange(startHour, endHour time.Time, username string, modelName string, channelID int, groupName string) ([]LogStatsAggregation, error) {
	var results []LogStatsAggregation
	query := LOG_DB.Model(&LogStatsAggregation{}).
		Where("hour_start >= ? AND hour_start < ?", startHour, endHour)

	if username != "" {
		query = query.Where("username = ?", username)
	}
	if modelName != "" {
		query = query.Where("model_name LIKE ?", modelName)
	}
	if channelID > 0 {
		query = query.Where("channel_id = ?", channelID)
	}
	if groupName != "" {
		query = query.Where("group_name = ?", groupName)
	}

	err := query.Find(&results).Error
	return results, err
}

// SumAggregatedStats 汇总指定时间范围内的聚合统计数据
func SumAggregatedStats(startHour, endHour time.Time, username string, modelName string, channelID int, groupName string) (*LogStat, error) {
	var result struct {
		TotalQuota  int64
		TotalTokens int64
		TotalCount  int64
	}

	query := LOG_DB.Model(&LogStatsAggregation{}).
		Select("COALESCE(SUM(total_quota), 0) as total_quota, COALESCE(SUM(total_tokens), 0) as total_tokens, COALESCE(SUM(request_count), 0) as total_count").
		Where("hour_start >= ? AND hour_start < ?", startHour, endHour)

	if username != "" {
		query = query.Where("username = ?", username)
	}
	if modelName != "" {
		query = query.Where("model_name LIKE ?", modelName)
	}
	if channelID > 0 {
		query = query.Where("channel_id = ?", channelID)
	}
	if groupName != "" {
		query = query.Where("group_name = ?", groupName)
	}

	err := query.Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return &LogStat{
		Quota: result.TotalQuota,
		Tpm:   result.TotalTokens,
		Rpm:   result.TotalCount,
	}, nil
}

// DeleteOldAggregations 删除指定时间之前的聚合数据
func DeleteOldAggregations(before time.Time, limit int) (int64, error) {
	result := LOG_DB.Where("hour_start < ?", before).Limit(limit).Delete(&LogStatsAggregation{})
	return result.RowsAffected, result.Error
}

// GetLatestAggregationHour 获取最新的聚合小时
func GetLatestAggregationHour() (time.Time, error) {
	var agg LogStatsAggregation
	err := LOG_DB.Model(&LogStatsAggregation{}).Order("hour_start DESC").First(&agg).Error
	if err != nil {
		return time.Time{}, err
	}
	return agg.HourStart, nil
}

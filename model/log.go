package model

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/types"

	"github.com/gin-gonic/gin"

	"github.com/bytedance/gopkg/util/gopool"
	"gorm.io/gorm"
)

type Log struct {
	Id               int    `json:"id" gorm:"index:idx_created_at_id,priority:1"`
	UserId           int    `json:"user_id" gorm:"index;index:idx_user_created,priority:1"`
	CreatedAt        int64  `json:"created_at" gorm:"bigint;index:idx_created_at_id,priority:2;index:idx_created_at_type,priority:2;index:idx_user_created,priority:2;index:idx_username_created,priority:2;index:idx_channel_created,priority:2"`
	Type             int    `json:"type" gorm:"index:idx_created_at_type,priority:1"`
	Content          string `json:"content"`
	Username         string `json:"username" gorm:"index;index:index_username_model_name,priority:2;index:idx_username_created,priority:1;default:''"`
	TokenName        string `json:"token_name" gorm:"index;default:''"`
	ModelName        string `json:"model_name" gorm:"index;index:index_username_model_name,priority:1;default:''"`
	Quota            int    `json:"quota" gorm:"default:0"`
	PromptTokens     int    `json:"prompt_tokens" gorm:"default:0"`
	CompletionTokens int    `json:"completion_tokens" gorm:"default:0"`
	UseTime          int    `json:"use_time" gorm:"default:0"`
	IsStream         bool   `json:"is_stream"`
	ChannelId        int    `json:"channel" gorm:"index;index:idx_channel_created,priority:1"`
	ChannelName      string `json:"channel_name" gorm:"->"`
	TokenId          int    `json:"token_id" gorm:"default:0;index"`
	Group            string `json:"group" gorm:"index"`
	Ip               string `json:"ip" gorm:"index;default:''"`
	Other            string `json:"other"`
}

// don't use iota, avoid change log type value
const (
	LogTypeUnknown = 0
	LogTypeTopup   = 1
	LogTypeConsume = 2
	LogTypeManage  = 3
	LogTypeSystem  = 4
	LogTypeError   = 5
	LogTypeRefund  = 6
)

func formatUserLogs(logs []*Log) {
	for i := range logs {
		logs[i].ChannelName = ""
		var otherMap map[string]interface{}
		otherMap, _ = common.StrToMap(logs[i].Other)
		if otherMap != nil {
			// delete admin
			delete(otherMap, "admin_info")
		}
		logs[i].Other = common.MapToJsonStr(otherMap)
		logs[i].Id = logs[i].Id % 1024
	}
}

// FillChannelNames fills channel names for a list of logs by batch lookup
func FillChannelNames(logs []*Log) {
	channelIds := types.NewSet[int]()
	for _, log := range logs {
		if log.ChannelId != 0 {
			channelIds.Add(log.ChannelId)
		}
	}

	if channelIds.Len() > 0 {
		var channels []struct {
			Id   int    `gorm:"column:id"`
			Name string `gorm:"column:name"`
		}
		if err := DB.Table("channels").Select("id, name").Where("id IN ?", channelIds.Items()).Find(&channels).Error; err != nil {
			return
		}
		channelMap := make(map[int]string, len(channels))
		for _, channel := range channels {
			channelMap[channel.Id] = channel.Name
		}
		for i := range logs {
			logs[i].ChannelName = channelMap[logs[i].ChannelId]
		}
	}
}

func GetLogByKey(key string) (logs []*Log, err error) {
	if os.Getenv("LOG_SQL_DSN") != "" {
		var tk Token
		if err = DB.Model(&Token{}).Where(logKeyCol+"=?", strings.TrimPrefix(key, "sk-")).First(&tk).Error; err != nil {
			return nil, err
		}
		err = LOG_DB.Model(&Log{}).Where("token_id=?", tk.Id).Find(&logs).Error
	} else {
		err = LOG_DB.Joins("left join tokens on tokens.id = logs.token_id").Where("tokens.key = ?", strings.TrimPrefix(key, "sk-")).Find(&logs).Error
	}
	formatUserLogs(logs)
	return logs, err
}

func RecordLog(userId int, logType int, content string) {
	if logType == LogTypeConsume && !common.LogConsumeEnabled {
		return
	}
	username, _ := GetUsernameById(userId, false)
	log := &Log{
		UserId:    userId,
		Username:  username,
		CreatedAt: common.GetTimestamp(),
		Type:      logType,
		Content:   content,
	}
	err := LOG_DB.Create(log).Error
	if err != nil {
		common.SysLog("failed to record log: " + err.Error())
	}
}

func RecordErrorLog(c *gin.Context, userId int, channelId int, modelName string, tokenName string, content string, tokenId int, useTimeSeconds int,
	isStream bool, group string, other map[string]interface{}) {
	logger.LogInfo(c, fmt.Sprintf("record error log: userId=%d, channelId=%d, modelName=%s, tokenName=%s, content=%s", userId, channelId, modelName, tokenName, content))
	username := c.GetString("username")
	otherStr := common.MapToJsonStr(other)
	// 判断是否需要记录 IP
	needRecordIp := false
	if settingMap, err := GetUserSetting(userId, false); err == nil {
		if settingMap.RecordIpLog {
			needRecordIp = true
		}
	}
	log := &Log{
		UserId:           userId,
		Username:         username,
		CreatedAt:        common.GetTimestamp(),
		Type:             LogTypeError,
		Content:          content,
		PromptTokens:     0,
		CompletionTokens: 0,
		TokenName:        tokenName,
		ModelName:        modelName,
		Quota:            0,
		ChannelId:        channelId,
		TokenId:          tokenId,
		UseTime:          useTimeSeconds,
		IsStream:         isStream,
		Group:            group,
		Ip: func() string {
			if needRecordIp {
				return c.ClientIP()
			}
			return ""
		}(),
		Other: otherStr,
	}
	err := LOG_DB.Create(log).Error
	if err != nil {
		logger.LogError(c, "failed to record log: "+err.Error())
	}
}

type RecordConsumeLogParams struct {
	ChannelId        int                    `json:"channel_id"`
	PromptTokens     int                    `json:"prompt_tokens"`
	CompletionTokens int                    `json:"completion_tokens"`
	ModelName        string                 `json:"model_name"`
	TokenName        string                 `json:"token_name"`
	Quota            int                    `json:"quota"`
	Content          string                 `json:"content"`
	TokenId          int                    `json:"token_id"`
	UseTimeSeconds   int                    `json:"use_time_seconds"`
	IsStream         bool                   `json:"is_stream"`
	Group            string                 `json:"group"`
	Other            map[string]interface{} `json:"other"`
}

func RecordConsumeLog(c *gin.Context, userId int, params RecordConsumeLogParams) {
	if !common.LogConsumeEnabled {
		return
	}
	logger.LogInfo(c, fmt.Sprintf("record consume log: userId=%d, params=%s", userId, common.GetJsonString(params)))
	username := c.GetString("username")
	otherStr := common.MapToJsonStr(params.Other)
	// 判断是否需要记录 IP - 使用缓存的用户设置
	needRecordIp := false
	if settingMap, err := GetUserSettingCached(userId); err == nil {
		if settingMap.RecordIpLog {
			needRecordIp = true
		}
	}
	clientIP := ""
	if needRecordIp {
		clientIP = c.ClientIP()
	}
	createdAt := common.GetTimestamp()
	log := &Log{
		UserId:           userId,
		Username:         username,
		CreatedAt:        createdAt,
		Type:             LogTypeConsume,
		Content:          params.Content,
		PromptTokens:     params.PromptTokens,
		CompletionTokens: params.CompletionTokens,
		TokenName:        params.TokenName,
		ModelName:        params.ModelName,
		Quota:            params.Quota,
		ChannelId:        params.ChannelId,
		TokenId:          params.TokenId,
		UseTime:          params.UseTimeSeconds,
		IsStream:         params.IsStream,
		Group:            params.Group,
		Ip:               clientIP,
		Other:            otherStr,
	}
	// 异步写入日志，避免阻塞请求
	gopool.Go(func() {
		err := LOG_DB.Create(log).Error
		if err != nil {
			common.SysLog("failed to record consume log: " + err.Error())
		}
	})
	if common.DataExportEnabled {
		gopool.Go(func() {
			LogQuotaData(userId, username, params.ModelName, params.Quota, createdAt, params.PromptTokens+params.CompletionTokens)
		})
	}
}

func GetAllLogs(logType int, startTimestamp int64, endTimestamp int64, modelName string, username string, tokenName string, startIdx int, num int, channel int, group string) (logs []*Log, total int64, err error) {
	var tx *gorm.DB
	if logType == LogTypeUnknown {
		tx = LOG_DB
	} else {
		tx = LOG_DB.Where("logs.type = ?", logType)
	}

	if modelName != "" {
		tx = tx.Where("logs.model_name like ?", modelName)
	}
	if username != "" {
		tx = tx.Where("logs.username = ?", username)
	}
	if tokenName != "" {
		tx = tx.Where("logs.token_name = ?", tokenName)
	}
	if startTimestamp != 0 {
		tx = tx.Where("logs.created_at >= ?", startTimestamp)
	}
	if endTimestamp != 0 {
		tx = tx.Where("logs.created_at <= ?", endTimestamp)
	}
	if channel != 0 {
		tx = tx.Where("logs.channel_id = ?", channel)
	}
	if group != "" {
		tx = tx.Where("logs."+logGroupCol+" = ?", group)
	}
	err = tx.Model(&Log{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = tx.Order("logs.id desc").Limit(num).Offset(startIdx).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	// Fill channel names using the shared helper function
	FillChannelNames(logs)

	return logs, total, err
}

func GetUserLogs(userId int, logType int, startTimestamp int64, endTimestamp int64, modelName string, tokenName string, startIdx int, num int, group string) (logs []*Log, total int64, err error) {
	var tx *gorm.DB
	if logType == LogTypeUnknown {
		tx = LOG_DB.Where("logs.user_id = ?", userId)
	} else {
		tx = LOG_DB.Where("logs.user_id = ? and logs.type = ?", userId, logType)
	}

	if modelName != "" {
		tx = tx.Where("logs.model_name like ?", modelName)
	}
	if tokenName != "" {
		tx = tx.Where("logs.token_name = ?", tokenName)
	}
	if startTimestamp != 0 {
		tx = tx.Where("logs.created_at >= ?", startTimestamp)
	}
	if endTimestamp != 0 {
		tx = tx.Where("logs.created_at <= ?", endTimestamp)
	}
	if group != "" {
		tx = tx.Where("logs."+logGroupCol+" = ?", group)
	}
	err = tx.Model(&Log{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = tx.Order("logs.id desc").Limit(num).Offset(startIdx).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	formatUserLogs(logs)
	return logs, total, err
}

func SearchAllLogs(keyword string) (logs []*Log, err error) {
	err = LOG_DB.Where("type = ? or content LIKE ?", keyword, keyword+"%").Order("id desc").Limit(common.MaxRecentItems).Find(&logs).Error
	return logs, err
}

func SearchUserLogs(userId int, keyword string) (logs []*Log, err error) {
	err = LOG_DB.Where("user_id = ? and type = ?", userId, keyword).Order("id desc").Limit(common.MaxRecentItems).Find(&logs).Error
	formatUserLogs(logs)
	return logs, err
}

type Stat struct {
	Quota int `json:"quota"`
	Rpm   int `json:"rpm"`
	Tpm   int `json:"tpm"`
}

// LogQueryOptions defines the options for cursor-based log queries
type LogQueryOptions struct {
	LogType        int
	StartTimestamp int64
	EndTimestamp   int64
	ModelName      string
	Username       string
	TokenName      string
	ChannelId      int
	Group          string
	PageSize       int
	Cursor         int64 // ID of the last record from previous page, 0 for first page
}

// LogQueryResult represents the result of a cursor-based log query
type LogQueryResult struct {
	Items      []*Log `json:"items"`
	NextCursor int64  `json:"next_cursor"` // Cursor for next page, 0 if no more data
	HasMore    bool   `json:"has_more"`
}

// applyLogFilters applies common filter conditions to a query
func applyLogFilters(query *gorm.DB, opts LogQueryOptions) *gorm.DB {
	if opts.LogType != LogTypeUnknown {
		query = query.Where("type = ?", opts.LogType)
	}
	if opts.StartTimestamp > 0 {
		query = query.Where("created_at >= ?", opts.StartTimestamp)
	}
	if opts.EndTimestamp > 0 {
		query = query.Where("created_at <= ?", opts.EndTimestamp)
	}
	if opts.Username != "" {
		query = query.Where("username = ?", opts.Username)
	}
	if opts.ModelName != "" {
		query = query.Where("model_name LIKE ?", opts.ModelName)
	}
	if opts.TokenName != "" {
		query = query.Where("token_name = ?", opts.TokenName)
	}
	if opts.ChannelId > 0 {
		query = query.Where("channel_id = ?", opts.ChannelId)
	}
	if opts.Group != "" {
		query = query.Where(logGroupCol+" = ?", opts.Group)
	}
	return query
}

// GetLogsByCursor retrieves logs using cursor-based pagination
func GetLogsByCursor(opts LogQueryOptions) (*LogQueryResult, error) {
	query := LOG_DB.Model(&Log{})

	// Apply cursor condition (for pagination)
	if opts.Cursor > 0 {
		query = query.Where("id < ?", opts.Cursor)
	}

	// Apply filter conditions
	query = applyLogFilters(query, opts)

	// Query pageSize + 1 records to determine if there are more
	pageSize := opts.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	var logs []*Log
	err := query.Order("id DESC").Limit(pageSize + 1).Find(&logs).Error
	if err != nil {
		return nil, err
	}

	result := &LogQueryResult{
		Items:   logs,
		HasMore: len(logs) > pageSize,
	}

	if result.HasMore {
		result.Items = logs[:pageSize]
		result.NextCursor = int64(logs[pageSize-1].Id)
	}

	return result, nil
}

// GetLogsCount returns the total count of logs matching the filter conditions
func GetLogsCount(opts LogQueryOptions) (int64, error) {
	query := LOG_DB.Model(&Log{})

	// Apply filter conditions (reuse the same filter logic)
	query = applyLogFilters(query, opts)

	var count int64
	err := query.Count(&count).Error
	return count, err
}

// LogQueryOptionsWithUser extends LogQueryOptions with user-specific filtering
type LogQueryOptionsWithUser struct {
	LogQueryOptions
	UserId int
}

// GetUserLogsByCursor retrieves user-specific logs using cursor-based pagination
func GetUserLogsByCursor(opts LogQueryOptionsWithUser) (*LogQueryResult, error) {
	query := LOG_DB.Model(&Log{}).Where("user_id = ?", opts.UserId)

	// Apply cursor condition (for pagination)
	if opts.Cursor > 0 {
		query = query.Where("id < ?", opts.Cursor)
	}

	// Apply filter conditions
	query = applyLogFilters(query, opts.LogQueryOptions)

	// Query pageSize + 1 records to determine if there are more
	pageSize := opts.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	var logs []*Log
	err := query.Order("id DESC").Limit(pageSize + 1).Find(&logs).Error
	if err != nil {
		return nil, err
	}

	result := &LogQueryResult{
		Items:   logs,
		HasMore: len(logs) > pageSize,
	}

	if result.HasMore {
		result.Items = logs[:pageSize]
		result.NextCursor = int64(logs[pageSize-1].Id)
	}

	// Format user logs (hide sensitive info)
	formatUserLogs(result.Items)

	return result, nil
}

// LogStat represents aggregated statistics for logs
type LogStat struct {
	Quota int64 `json:"quota"`
	Rpm   int64 `json:"rpm"`
	Tpm   int64 `json:"tpm"`
}

// GetLogStats retrieves quota, rpm, and tpm statistics in a single optimized query
func GetLogStats(opts LogQueryOptions) (*LogStat, error) {
	now := time.Now().Unix()
	recentStart := now - 60 // Last 60 seconds for rpm/tpm

	// Single query to get all statistics using CASE WHEN
	var result struct {
		Quota     int64
		RecentRpm int64
		RecentTpm int64
	}

	query := LOG_DB.Table("logs").
		Select(`
			SUM(CASE WHEN created_at >= ? AND created_at <= ? THEN quota ELSE 0 END) as quota,
			SUM(CASE WHEN created_at >= ? THEN 1 ELSE 0 END) as recent_rpm,
			SUM(CASE WHEN created_at >= ? THEN prompt_tokens + completion_tokens ELSE 0 END) as recent_tpm
		`, opts.StartTimestamp, opts.EndTimestamp, recentStart, recentStart).
		Where("type = ?", LogTypeConsume)

	// Apply filter conditions
	if opts.Username != "" {
		query = query.Where("username = ?", opts.Username)
	}
	if opts.ModelName != "" {
		query = query.Where("model_name LIKE ?", opts.ModelName)
	}
	if opts.TokenName != "" {
		query = query.Where("token_name = ?", opts.TokenName)
	}
	if opts.ChannelId > 0 {
		query = query.Where("channel_id = ?", opts.ChannelId)
	}
	if opts.Group != "" {
		query = query.Where(logGroupCol+" = ?", opts.Group)
	}

	err := query.Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return &LogStat{
		Quota: result.Quota,
		Rpm:   result.RecentRpm,
		Tpm:   result.RecentTpm,
	}, nil
}

func SumUsedQuota(logType int, startTimestamp int64, endTimestamp int64, modelName string, username string, tokenName string, channel int, group string) (stat Stat) {
	tx := LOG_DB.Table("logs").Select("sum(quota) quota")

	// 为rpm和tpm创建单独的查询
	rpmTpmQuery := LOG_DB.Table("logs").Select("count(*) rpm, sum(prompt_tokens) + sum(completion_tokens) tpm")

	if username != "" {
		tx = tx.Where("username = ?", username)
		rpmTpmQuery = rpmTpmQuery.Where("username = ?", username)
	}
	if tokenName != "" {
		tx = tx.Where("token_name = ?", tokenName)
		rpmTpmQuery = rpmTpmQuery.Where("token_name = ?", tokenName)
	}
	if startTimestamp != 0 {
		tx = tx.Where("created_at >= ?", startTimestamp)
	}
	if endTimestamp != 0 {
		tx = tx.Where("created_at <= ?", endTimestamp)
	}
	if modelName != "" {
		tx = tx.Where("model_name like ?", modelName)
		rpmTpmQuery = rpmTpmQuery.Where("model_name like ?", modelName)
	}
	if channel != 0 {
		tx = tx.Where("channel_id = ?", channel)
		rpmTpmQuery = rpmTpmQuery.Where("channel_id = ?", channel)
	}
	if group != "" {
		tx = tx.Where(logGroupCol+" = ?", group)
		rpmTpmQuery = rpmTpmQuery.Where(logGroupCol+" = ?", group)
	}

	tx = tx.Where("type = ?", LogTypeConsume)
	rpmTpmQuery = rpmTpmQuery.Where("type = ?", LogTypeConsume)

	// 只统计最近60秒的rpm和tpm
	rpmTpmQuery = rpmTpmQuery.Where("created_at >= ?", time.Now().Add(-60*time.Second).Unix())

	// 执行查询
	tx.Scan(&stat)
	rpmTpmQuery.Scan(&stat)

	return stat
}

func SumUsedToken(logType int, startTimestamp int64, endTimestamp int64, modelName string, username string, tokenName string) (token int) {
	tx := LOG_DB.Table("logs").Select("ifnull(sum(prompt_tokens),0) + ifnull(sum(completion_tokens),0)")
	if username != "" {
		tx = tx.Where("username = ?", username)
	}
	if tokenName != "" {
		tx = tx.Where("token_name = ?", tokenName)
	}
	if startTimestamp != 0 {
		tx = tx.Where("created_at >= ?", startTimestamp)
	}
	if endTimestamp != 0 {
		tx = tx.Where("created_at <= ?", endTimestamp)
	}
	if modelName != "" {
		tx = tx.Where("model_name = ?", modelName)
	}
	tx.Where("type = ?", LogTypeConsume).Scan(&token)
	return token
}

func DeleteOldLog(ctx context.Context, targetTimestamp int64, limit int) (int64, error) {
	var total int64 = 0

	for {
		if nil != ctx.Err() {
			return total, ctx.Err()
		}

		result := LOG_DB.Where("created_at < ?", targetTimestamp).Limit(limit).Delete(&Log{})
		if nil != result.Error {
			return total, result.Error
		}

		total += result.RowsAffected

		if result.RowsAffected < int64(limit) {
			break
		}
	}

	return total, nil
}

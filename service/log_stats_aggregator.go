package service

import (
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"

	"github.com/bytedance/gopkg/util/gopool"
)

// LogStatsAggregator 日志统计聚合服务
type LogStatsAggregator struct {
	mu              sync.RWMutex
	enabled         bool
	aggregationChan chan aggregationTask
	stopChan        chan struct{}
}

type aggregationTask struct {
	hourStart time.Time
}

var (
	statsAggregator     *LogStatsAggregator
	statsAggregatorOnce sync.Once
)

// GetLogStatsAggregator 获取聚合服务单例
func GetLogStatsAggregator() *LogStatsAggregator {
	statsAggregatorOnce.Do(func() {
		statsAggregator = &LogStatsAggregator{
			enabled:         true,
			aggregationChan: make(chan aggregationTask, 100),
			stopChan:        make(chan struct{}),
		}
		// 启动后台聚合处理器
		go statsAggregator.processAggregations()
	})
	return statsAggregator
}

// processAggregations 后台处理聚合任务
func (a *LogStatsAggregator) processAggregations() {
	for {
		select {
		case task := <-a.aggregationChan:
			a.doAggregateHourlyStats(task.hourStart)
		case <-a.stopChan:
			return
		}
	}
}

// Stop 停止聚合服务
func (a *LogStatsAggregator) Stop() {
	close(a.stopChan)
}

// IsEnabled 检查聚合服务是否启用
func (a *LogStatsAggregator) IsEnabled() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.enabled
}

// SetEnabled 设置聚合服务启用状态
func (a *LogStatsAggregator) SetEnabled(enabled bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.enabled = enabled
}

// TriggerAggregationAsync 异步触发指定小时的聚合
// 不阻塞调用方，符合 Requirements 5.4
func (a *LogStatsAggregator) TriggerAggregationAsync(hourStart time.Time) {
	if !a.IsEnabled() {
		return
	}

	// 使用 gopool 异步执行，避免阻塞
	gopool.Go(func() {
		select {
		case a.aggregationChan <- aggregationTask{hourStart: hourStart}:
		default:
			// 队列满时跳过，避免阻塞
			common.SysLog("log stats aggregation queue full, skipping hour: " + hourStart.Format(time.RFC3339))
		}
	})
}

// AggregateHourlyStats 聚合指定小时的统计数据
// 同步执行，用于手动触发或定时任务
func (a *LogStatsAggregator) AggregateHourlyStats(hourStart time.Time) error {
	return a.doAggregateHourlyStats(hourStart)
}

// doAggregateHourlyStats 执行实际的聚合逻辑
func (a *LogStatsAggregator) doAggregateHourlyStats(hourStart time.Time) error {
	// 确保 hourStart 是整点
	hourStart = hourStart.Truncate(time.Hour)
	hourEnd := hourStart.Add(time.Hour)

	// 从原始日志表查询该小时的聚合数据
	// 按 user_id, model_name, channel_id, group 分组
	type AggResult struct {
		UserID       int
		Username     string
		ModelName    string
		ChannelID    int
		GroupName    string
		TotalQuota   int64
		TotalTokens  int64
		RequestCount int64
	}

	var results []AggResult
	err := model.LOG_DB.Table("logs").
		Select(`
			user_id,
			username,
			model_name,
			channel_id,
			`+"`group`"+` as group_name,
			COALESCE(SUM(quota), 0) as total_quota,
			COALESCE(SUM(prompt_tokens + completion_tokens), 0) as total_tokens,
			COUNT(*) as request_count
		`).
		Where("type = ?", model.LogTypeConsume).
		Where("created_at >= ? AND created_at < ?", hourStart.Unix(), hourEnd.Unix()).
		Group("user_id, username, model_name, channel_id, `group`").
		Scan(&results).Error

	if err != nil {
		common.SysLog("failed to aggregate hourly stats: " + err.Error())
		return err
	}

	// 批量写入聚合表
	now := time.Now()
	for _, r := range results {
		agg := &model.LogStatsAggregation{
			HourStart:    hourStart,
			UserID:       r.UserID,
			Username:     r.Username,
			ModelName:    r.ModelName,
			ChannelID:    r.ChannelID,
			GroupName:    r.GroupName,
			TotalQuota:   r.TotalQuota,
			TotalTokens:  r.TotalTokens,
			RequestCount: r.RequestCount,
			UpdatedAt:    now,
		}

		if err := model.UpsertLogStatsAggregation(agg); err != nil {
			common.SysLog("failed to upsert aggregation: " + err.Error())
			// 继续处理其他记录
		}
	}

	return nil
}

// AggregateRecentHours 聚合最近 N 小时的数据
// 用于初始化或补充聚合数据
func (a *LogStatsAggregator) AggregateRecentHours(hours int) error {
	now := time.Now().Truncate(time.Hour)

	for i := hours; i > 0; i-- {
		hourStart := now.Add(-time.Duration(i) * time.Hour)
		if err := a.AggregateHourlyStats(hourStart); err != nil {
			common.SysLog("failed to aggregate hour " + hourStart.Format(time.RFC3339) + ": " + err.Error())
			// 继续处理其他小时
		}
	}

	return nil
}

// GetCurrentHourStart 获取当前小时的开始时间
func GetCurrentHourStart() time.Time {
	return time.Now().Truncate(time.Hour)
}

// GetPreviousHourStart 获取上一个小时的开始时间
func GetPreviousHourStart() time.Time {
	return time.Now().Truncate(time.Hour).Add(-time.Hour)
}

// TruncateToHour 将时间戳截断到小时
func TruncateToHour(timestamp int64) time.Time {
	return time.Unix(timestamp, 0).Truncate(time.Hour)
}

// IsCompleteHour 检查时间范围是否包含完整的小时
// 返回完整小时的开始和结束时间，以及是否有部分小时
func IsCompleteHour(startTimestamp, endTimestamp int64) (completeStart, completeEnd time.Time, hasPartial bool) {
	start := time.Unix(startTimestamp, 0)
	end := time.Unix(endTimestamp, 0)

	// 计算第一个完整小时的开始
	completeStart = start.Truncate(time.Hour)
	if completeStart.Before(start) {
		completeStart = completeStart.Add(time.Hour)
	}

	// 计算最后一个完整小时的结束
	completeEnd = end.Truncate(time.Hour)

	// 检查是否有部分小时
	hasPartial = !start.Equal(start.Truncate(time.Hour)) || !end.Equal(end.Truncate(time.Hour))

	return completeStart, completeEnd, hasPartial
}

// GetPartialHourRanges 获取部分小时的时间范围
// 返回开始部分和结束部分的时间范围
func GetPartialHourRanges(startTimestamp, endTimestamp int64) (startPartial, endPartial *TimeRange) {
	start := time.Unix(startTimestamp, 0)
	end := time.Unix(endTimestamp, 0)

	startHour := start.Truncate(time.Hour)
	endHour := end.Truncate(time.Hour)

	// 开始部分：从 start 到 startHour + 1h
	if !start.Equal(startHour) {
		nextHour := startHour.Add(time.Hour)
		if nextHour.Before(end) || nextHour.Equal(end) {
			startPartial = &TimeRange{
				Start: startTimestamp,
				End:   nextHour.Unix(),
			}
		} else {
			// 整个范围都在一个小时内
			startPartial = &TimeRange{
				Start: startTimestamp,
				End:   endTimestamp,
			}
			return startPartial, nil
		}
	}

	// 结束部分：从 endHour 到 end
	if !end.Equal(endHour) && (startPartial == nil || endHour.Unix() > startPartial.End) {
		endPartial = &TimeRange{
			Start: endHour.Unix(),
			End:   endTimestamp,
		}
	}

	return startPartial, endPartial
}

// TimeRange 时间范围
type TimeRange struct {
	Start int64
	End   int64
}

// InitLogStatsAggregator 初始化日志统计聚合服务
// 注册聚合数据查询回调到 model 层
func InitLogStatsAggregator() {
	// 获取聚合服务实例
	aggregator := GetLogStatsAggregator()

	// 注册聚合数据查询回调
	model.RegisterLogStatsAggregationCallback(func(startHour, endHour time.Time, username, modelName string, channelID int, groupName string) (*model.LogStat, error) {
		return model.SumAggregatedStats(startHour, endHour, username, modelName, channelID, groupName)
	})

	// 启动时异步聚合最近 24 小时的历史数据
	gopool.Go(func() {
		// 等待数据库连接稳定
		time.Sleep(5 * time.Second)
		common.SysLog("Starting initial log stats aggregation for last 24 hours...")
		if err := aggregator.AggregateRecentHours(24); err != nil {
			common.SysLog("Initial aggregation completed with some errors: " + err.Error())
		} else {
			common.SysLog("Initial log stats aggregation completed successfully")
		}
	})

	// 启动定时聚合任务（每小时执行一次）
	go aggregator.startHourlyAggregation()

	common.SysLog("Log stats aggregator initialized")
}

// startHourlyAggregation 启动每小时定时聚合任务
func (a *LogStatsAggregator) startHourlyAggregation() {
	// 计算到下一个整点的时间
	now := time.Now()
	nextHour := now.Truncate(time.Hour).Add(time.Hour)
	waitDuration := nextHour.Sub(now)

	// 等待到下一个整点
	time.Sleep(waitDuration)

	// 创建每小时触发的 ticker
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !a.IsEnabled() {
				continue
			}
			// 聚合上一个完整小时的数据
			prevHour := GetPreviousHourStart()
			gopool.Go(func() {
				if err := a.AggregateHourlyStats(prevHour); err != nil {
					common.SysLog("Hourly aggregation failed for " + prevHour.Format(time.RFC3339) + ": " + err.Error())
				} else {
					common.SysLog("Hourly aggregation completed for " + prevHour.Format(time.RFC3339))
				}
			})
		case <-a.stopChan:
			return
		}
	}
}

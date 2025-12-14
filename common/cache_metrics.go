package common

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// CacheMetrics 缓存指标收集器
// 用于收集和统计多级缓存的性能指标
type CacheMetrics struct {
	stats sync.Map // cacheName -> *CacheMetricsData
}

// CacheMetricsData 单个缓存的指标数据
type CacheMetricsData struct {
	Name     string
	L1Hits   int64
	L1Misses int64
	L2Hits   int64
	L2Misses int64

	// 渠道选择耗时统计
	SelectCount      int64
	SelectTotalNanos int64
	SelectMaxNanos   int64
}

// globalCacheMetrics 全局缓存指标收集器
var globalCacheMetrics = &CacheMetrics{}

// GetCacheMetrics 获取全局缓存指标收集器
func GetCacheMetrics() *CacheMetrics {
	return globalCacheMetrics
}

// Record 记录缓存访问
// cacheName: 缓存名称
// hit: 是否命中
// level: 缓存层级 (1=L1, 2=L2)
func (m *CacheMetrics) Record(cacheName string, hit bool, level int) {
	data := m.getOrCreate(cacheName)

	if level == 1 {
		if hit {
			atomic.AddInt64(&data.L1Hits, 1)
		} else {
			atomic.AddInt64(&data.L1Misses, 1)
		}
	} else if level == 2 {
		if hit {
			atomic.AddInt64(&data.L2Hits, 1)
		} else {
			atomic.AddInt64(&data.L2Misses, 1)
		}
	}
}


// RecordSelectTime 记录渠道选择耗时
func (m *CacheMetrics) RecordSelectTime(cacheName string, duration time.Duration) {
	data := m.getOrCreate(cacheName)

	nanos := duration.Nanoseconds()
	atomic.AddInt64(&data.SelectCount, 1)
	atomic.AddInt64(&data.SelectTotalNanos, nanos)

	// 更新最大耗时（使用 CAS 操作）
	for {
		current := atomic.LoadInt64(&data.SelectMaxNanos)
		if nanos <= current {
			break
		}
		if atomic.CompareAndSwapInt64(&data.SelectMaxNanos, current, nanos) {
			break
		}
	}
}

// GetStats 获取指定缓存的统计信息
func (m *CacheMetrics) GetStats(cacheName string) *CacheMetricsData {
	if val, ok := m.stats.Load(cacheName); ok {
		data := val.(*CacheMetricsData)
		// 返回副本，避免并发问题
		return &CacheMetricsData{
			Name:             data.Name,
			L1Hits:           atomic.LoadInt64(&data.L1Hits),
			L1Misses:         atomic.LoadInt64(&data.L1Misses),
			L2Hits:           atomic.LoadInt64(&data.L2Hits),
			L2Misses:         atomic.LoadInt64(&data.L2Misses),
			SelectCount:      atomic.LoadInt64(&data.SelectCount),
			SelectTotalNanos: atomic.LoadInt64(&data.SelectTotalNanos),
			SelectMaxNanos:   atomic.LoadInt64(&data.SelectMaxNanos),
		}
	}
	return nil
}

// GetAllStats 获取所有缓存的统计信息
func (m *CacheMetrics) GetAllStats() []*CacheMetricsData {
	var result []*CacheMetricsData
	m.stats.Range(func(key, value interface{}) bool {
		data := value.(*CacheMetricsData)
		result = append(result, &CacheMetricsData{
			Name:             data.Name,
			L1Hits:           atomic.LoadInt64(&data.L1Hits),
			L1Misses:         atomic.LoadInt64(&data.L1Misses),
			L2Hits:           atomic.LoadInt64(&data.L2Hits),
			L2Misses:         atomic.LoadInt64(&data.L2Misses),
			SelectCount:      atomic.LoadInt64(&data.SelectCount),
			SelectTotalNanos: atomic.LoadInt64(&data.SelectTotalNanos),
			SelectMaxNanos:   atomic.LoadInt64(&data.SelectMaxNanos),
		})
		return true
	})
	return result
}

// HitRate 计算命中率
func (d *CacheMetricsData) HitRate() float64 {
	totalHits := d.L1Hits + d.L2Hits
	totalMisses := d.L1Misses // L1 miss 才会查 L2
	total := totalHits + totalMisses
	if total == 0 {
		return 0
	}
	return float64(totalHits) / float64(total)
}

// AvgSelectTime 计算平均选择耗时
func (d *CacheMetricsData) AvgSelectTime() time.Duration {
	count := d.SelectCount
	if count == 0 {
		return 0
	}
	return time.Duration(d.SelectTotalNanos / count)
}

// MaxSelectTime 获取最大选择耗时
func (d *CacheMetricsData) MaxSelectTime() time.Duration {
	return time.Duration(d.SelectMaxNanos)
}

// Reset 重置统计数据
func (m *CacheMetrics) Reset(cacheName string) {
	if val, ok := m.stats.Load(cacheName); ok {
		data := val.(*CacheMetricsData)
		atomic.StoreInt64(&data.L1Hits, 0)
		atomic.StoreInt64(&data.L1Misses, 0)
		atomic.StoreInt64(&data.L2Hits, 0)
		atomic.StoreInt64(&data.L2Misses, 0)
		atomic.StoreInt64(&data.SelectCount, 0)
		atomic.StoreInt64(&data.SelectTotalNanos, 0)
		atomic.StoreInt64(&data.SelectMaxNanos, 0)
	}
}

// getOrCreate 获取或创建缓存指标数据
func (m *CacheMetrics) getOrCreate(cacheName string) *CacheMetricsData {
	if val, ok := m.stats.Load(cacheName); ok {
		return val.(*CacheMetricsData)
	}

	data := &CacheMetricsData{Name: cacheName}
	actual, _ := m.stats.LoadOrStore(cacheName, data)
	return actual.(*CacheMetricsData)
}


// CheckAndWarn 检查指标并记录告警
// 命中率低于 70% 时记录警告
// 渠道选择耗时超过 5ms 时记录警告
func (m *CacheMetrics) CheckAndWarn() {
	m.stats.Range(func(key, value interface{}) bool {
		cacheName := key.(string)
		data := value.(*CacheMetricsData)

		// 检查命中率
		hitRate := data.HitRate()
		totalRequests := atomic.LoadInt64(&data.L1Hits) + atomic.LoadInt64(&data.L1Misses)
		if totalRequests > 100 && hitRate < 0.7 {
			SysLog("Warning: cache " + cacheName + " hit rate " + 
				formatPercent(hitRate) + " is below 70%")
		}

		// 检查渠道选择耗时
		if cacheName == "channel_selector" {
			avgTime := data.AvgSelectTime()
			maxTime := data.MaxSelectTime()
			if avgTime > 5*time.Millisecond {
				SysLog("Warning: channel selector avg time " + 
					avgTime.String() + " exceeds 5ms")
			}
			if maxTime > 10*time.Millisecond {
				SysLog("Warning: channel selector max time " + 
					maxTime.String() + " exceeds 10ms")
			}
		}

		return true
	})
}

// formatPercent 格式化百分比
func formatPercent(rate float64) string {
	return fmt.Sprintf("%.2f%%", rate*100)
}

// CheckCacheUsage 检查缓存使用率
// 使用率超过 80% 时记录警告
func CheckCacheUsage(cacheName string, size, capacity int) {
	if capacity <= 0 {
		return
	}
	usageRate := float64(size) / float64(capacity)
	if usageRate > 0.8 {
		SysLog(fmt.Sprintf("Warning: cache %s usage %.2f%% exceeds 80%%", 
			cacheName, usageRate*100))
	}
}

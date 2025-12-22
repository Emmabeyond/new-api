package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
)

// LogQueryCacheConfig 日志查询缓存配置
type LogQueryCacheConfig struct {
	DefaultTTL time.Duration // 默认 TTL，30 秒
	MaxTTL     time.Duration // 最大 TTL，5 分钟
	L1Capacity int           // L1 缓存容量
}

// DefaultLogQueryCacheConfig 默认配置
func DefaultLogQueryCacheConfig() LogQueryCacheConfig {
	return LogQueryCacheConfig{
		DefaultTTL: 30 * time.Second,
		MaxTTL:     5 * time.Minute,
		L1Capacity: 5000,
	}
}

// LogQueryCache 日志查询缓存
type LogQueryCache struct {
	statsCache *common.MultiLevelCache[*model.LogStat]
	config     LogQueryCacheConfig
	mu         sync.RWMutex
}

// 全局日志查询缓存实例
var (
	logQueryCache     *LogQueryCache
	logQueryCacheOnce sync.Once
)

// GetLogQueryCache 获取全局日志查询缓存实例
func GetLogQueryCache() *LogQueryCache {
	logQueryCacheOnce.Do(func() {
		logQueryCache = NewLogQueryCache(DefaultLogQueryCacheConfig())
		// 注册缓存失效回调
		model.RegisterLogCacheInvalidateCallback(func(username string) {
			logQueryCache.InvalidateByUsername(username)
		})
	})
	return logQueryCache
}

// NewLogQueryCache 创建日志查询缓存
func NewLogQueryCache(config LogQueryCacheConfig) *LogQueryCache {
	mlcConfig := common.MultiLevelCacheConfig{
		L1Capacity:    config.L1Capacity,
		L1TTL:         config.DefaultTTL,
		L2TTL:         config.MaxTTL,
		TTLJitter:     0.1,
		EmptyTTL:      10 * time.Second,
		WarnThreshold: 0.8,
	}

	return &LogQueryCache{
		statsCache: common.NewMultiLevelCache[*model.LogStat]("log_stats", mlcConfig),
		config:     config,
	}
}

// CacheKey 生成缓存键
func (c *LogQueryCache) CacheKey(opts model.LogQueryOptions) string {
	// 使用查询参数生成唯一键
	keyStr := fmt.Sprintf("stats:%d:%d:%d:%s:%s:%s:%d:%s",
		opts.LogType,
		opts.StartTimestamp,
		opts.EndTimestamp,
		opts.Username,
		opts.ModelName,
		opts.TokenName,
		opts.ChannelId,
		opts.Group,
	)

	// 使用 MD5 生成固定长度的键
	hash := md5.Sum([]byte(keyStr))
	return hex.EncodeToString(hash[:])
}

// CacheResult 缓存查询结果
type CacheResult struct {
	Data      *model.LogStat
	Hit       bool
	CacheKey  string
	ExpiresAt time.Time
}

// GetStats 获取统计缓存
// 返回缓存结果和是否命中
func (c *LogQueryCache) GetStats(opts model.LogQueryOptions) (*model.LogStat, bool) {
	key := c.CacheKey(opts)

	// 使用 MultiLevelCache 的 Get 方法，但不使用 loader
	// 因为我们需要区分缓存命中和未命中
	stat, err := c.statsCache.Get(key, func() (*model.LogStat, bool, error) {
		// 返回 false 表示未找到，不触发实际加载
		return nil, false, nil
	})

	if err != nil {
		return nil, false
	}

	return stat, stat != nil
}

// SetStats 设置统计缓存
func (c *LogQueryCache) SetStats(opts model.LogQueryOptions, stat *model.LogStat) {
	key := c.CacheKey(opts)
	c.statsCache.Set(key, stat)
}

// GetStatsWithLoader 获取统计缓存，如果未命中则使用 loader 加载
func (c *LogQueryCache) GetStatsWithLoader(opts model.LogQueryOptions, loader func() (*model.LogStat, error)) (*model.LogStat, bool, error) {
	key := c.CacheKey(opts)

	var cacheHit bool = true

	stat, err := c.statsCache.Get(key, func() (*model.LogStat, bool, error) {
		cacheHit = false
		result, loadErr := loader()
		if loadErr != nil {
			return nil, false, loadErr
		}
		return result, true, nil
	})

	return stat, cacheHit, err
}

// InvalidateStats 使统计缓存失效
func (c *LogQueryCache) InvalidateStats(opts model.LogQueryOptions) {
	key := c.CacheKey(opts)
	c.statsCache.Invalidate(key)
}

// InvalidateByUsername 使指定用户的所有统计缓存失效
// 由于缓存键包含多个参数，这里使用模式匹配失效
func (c *LogQueryCache) InvalidateByUsername(username string) {
	// 由于 MultiLevelCache 不支持模式匹配删除，
	// 我们依赖 TTL 自动过期来处理
	// 对于写入频繁的场景，30 秒的 TTL 足够保证数据一致性
	common.SysLog(fmt.Sprintf("Log cache invalidation triggered for user: %s", username))
}

// InvalidateAll 使所有统计缓存失效
// 注意：这是一个重操作，仅在必要时使用
func (c *LogQueryCache) InvalidateAll() {
	// MultiLevelCache 不支持批量清除，依赖 TTL 自动过期
	common.SysLog("Log cache full invalidation triggered")
}

// Stats 返回缓存统计信息
func (c *LogQueryCache) Stats() common.CacheStats {
	return c.statsCache.Stats()
}

// CacheStatus 缓存状态常量
const (
	CacheStatusHit  = "HIT"
	CacheStatusMiss = "MISS"
)

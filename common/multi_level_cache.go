package common

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// MultiLevelCacheConfig 多级缓存配置
type MultiLevelCacheConfig struct {
	L1Capacity    int           // L1 缓存容量，默认 10000
	L1TTL         time.Duration // L1 TTL，默认 60s
	L2TTL         time.Duration // L2 TTL，默认 300s
	TTLJitter     float64       // TTL 抖动比例，默认 0.1 (±10%)
	EmptyTTL      time.Duration // 空结果 TTL，默认 30s
	WarnThreshold float64       // 使用率警告阈值，默认 0.8
}

// DefaultMultiLevelCacheConfig 默认配置
func DefaultMultiLevelCacheConfig() MultiLevelCacheConfig {
	return MultiLevelCacheConfig{
		L1Capacity:    10000,
		L1TTL:         60 * time.Second,
		L2TTL:         300 * time.Second,
		TTLJitter:     0.1,
		EmptyTTL:      30 * time.Second,
		WarnThreshold: 0.8,
	}
}

// MultiLevelCache 多级缓存管理器
// L1: 本地 LRU 缓存（进程内）
// L2: Redis 缓存（分布式）
type MultiLevelCache[T any] struct {
	name      string
	l1        *LRUCache[T]
	l2Enabled bool
	config    MultiLevelCacheConfig

	// 单飞模式，防止缓存击穿
	singleFlight sync.Map // key -> *singleFlightCall[T]

	// 统计指标
	l1Hits   int64
	l1Misses int64
	l2Hits   int64
	l2Misses int64
}

// singleFlightCall 单飞调用
type singleFlightCall[T any] struct {
	wg     sync.WaitGroup
	val    T
	err    error
	done   bool
	isEmpty bool
}

// NewMultiLevelCache 创建多级缓存
func NewMultiLevelCache[T any](name string, config MultiLevelCacheConfig) *MultiLevelCache[T] {
	return &MultiLevelCache[T]{
		name:      name,
		l1:        NewLRUCache[T](config.L1Capacity),
		l2Enabled: RedisEnabled,
		config:    config,
	}
}

// addJitter 添加 TTL 抖动，防止缓存雪崩
func (c *MultiLevelCache[T]) addJitter(ttl time.Duration) time.Duration {
	if c.config.TTLJitter <= 0 {
		return ttl
	}
	// 生成 [-jitter, +jitter] 范围内的随机抖动
	jitterRange := float64(ttl) * c.config.TTLJitter
	jitter := (rand.Float64()*2 - 1) * jitterRange
	return ttl + time.Duration(jitter)
}

// Get 获取缓存，按 L1 -> L2 -> loader 顺序查询
// loader: 当缓存未命中时的数据加载函数
func (c *MultiLevelCache[T]) Get(key string, loader func() (T, bool, error)) (T, error) {
	var zero T
	metrics := GetCacheMetrics()

	// 1. 查询 L1 缓存
	if val, hit, isEmpty := c.l1.Get(key); hit {
		atomic.AddInt64(&c.l1Hits, 1)
		metrics.Record(c.name, true, 1)
		if isEmpty {
			return zero, fmt.Errorf("key %s not found (cached empty)", key)
		}
		return val, nil
	}
	atomic.AddInt64(&c.l1Misses, 1)
	metrics.Record(c.name, false, 1)

	// 2. 查询 L2 缓存 (Redis)
	if c.l2Enabled {
		if val, err := c.getFromL2(key); err == nil {
			atomic.AddInt64(&c.l2Hits, 1)
			metrics.Record(c.name, true, 2)
			// 回填 L1
			c.l1.Set(key, val, c.addJitter(c.config.L1TTL))
			return val, nil
		}
		atomic.AddInt64(&c.l2Misses, 1)
		metrics.Record(c.name, false, 2)
	}

	// 3. 使用单飞模式加载数据
	return c.loadWithSingleFlight(key, loader)
}

// loadWithSingleFlight 使用单飞模式加载数据
// 相同 key 的并发请求只会执行一次 loader
func (c *MultiLevelCache[T]) loadWithSingleFlight(key string, loader func() (T, bool, error)) (T, error) {
	var zero T

	// 尝试获取或创建单飞调用
	call := &singleFlightCall[T]{}
	call.wg.Add(1)

	actual, loaded := c.singleFlight.LoadOrStore(key, call)
	if loaded {
		// 已有进行中的调用，等待结果
		existingCall := actual.(*singleFlightCall[T])
		existingCall.wg.Wait()
		if existingCall.err != nil {
			return zero, existingCall.err
		}
		if existingCall.isEmpty {
			return zero, fmt.Errorf("key %s not found", key)
		}
		return existingCall.val, nil
	}

	// 执行 loader
	defer func() {
		call.done = true
		call.wg.Done()
		// 延迟删除，给其他等待的 goroutine 一点时间
		go func() {
			time.Sleep(100 * time.Millisecond)
			c.singleFlight.Delete(key)
		}()
	}()

	val, exists, err := loader()
	if err != nil {
		call.err = err
		return zero, err
	}

	if !exists {
		// 缓存空结果，防止穿透
		call.isEmpty = true
		c.l1.SetEmpty(key, c.config.EmptyTTL)
		if c.l2Enabled {
			c.setEmptyToL2(key)
		}
		return zero, fmt.Errorf("key %s not found", key)
	}

	call.val = val

	// 回填两级缓存
	c.l1.Set(key, val, c.addJitter(c.config.L1TTL))
	if c.l2Enabled {
		c.setToL2(key, val)
	}

	return val, nil
}

// Set 设置缓存，同时更新 L1 和 L2
func (c *MultiLevelCache[T]) Set(key string, value T) {
	c.l1.Set(key, value, c.addJitter(c.config.L1TTL))
	if c.l2Enabled {
		c.setToL2(key, value)
	}
}

// Invalidate 使缓存失效
func (c *MultiLevelCache[T]) Invalidate(key string) {
	c.l1.Delete(key)
	if c.l2Enabled {
		c.deleteFromL2(key)
	}
}

// Stats 返回缓存统计信息
func (c *MultiLevelCache[T]) Stats() CacheStats {
	l1Hits, l1Misses, size, capacity := c.l1.Stats()
	l2Hits := atomic.LoadInt64(&c.l2Hits)
	l2Misses := atomic.LoadInt64(&c.l2Misses)

	totalHits := l1Hits + l2Hits
	totalMisses := l1Misses // L1 miss 才会查 L2

	var hitRate float64
	if totalHits+totalMisses > 0 {
		hitRate = float64(totalHits) / float64(totalHits+totalMisses)
	}

	return CacheStats{
		Name:     c.name,
		L1Hits:   l1Hits,
		L1Misses: l1Misses,
		L2Hits:   l2Hits,
		L2Misses: l2Misses,
		HitRate:  hitRate,
		Size:     size,
		Capacity: capacity,
	}
}

// CheckAndWarn 检查缓存状态并记录警告
func (c *MultiLevelCache[T]) CheckAndWarn() {
	stats := c.Stats()

	// 检查命中率
	if stats.HitRate < 0.7 && (stats.L1Hits+stats.L1Misses) > 100 {
		SysLog(fmt.Sprintf("Warning: cache %s hit rate %.2f%% is below 70%%", c.name, stats.HitRate*100))
	}

	// 检查使用率
	if stats.Capacity > 0 {
		usageRate := float64(stats.Size) / float64(stats.Capacity)
		if usageRate > c.config.WarnThreshold {
			SysLog(fmt.Sprintf("Warning: cache %s usage %.2f%% exceeds %.0f%%", c.name, usageRate*100, c.config.WarnThreshold*100))
		}
	}
}

// L2 Redis 操作

func (c *MultiLevelCache[T]) l2Key(key string) string {
	return fmt.Sprintf("mlc:%s:%s", c.name, key)
}

func (c *MultiLevelCache[T]) getFromL2(key string) (T, error) {
	var zero T
	if !RedisEnabled {
		return zero, fmt.Errorf("redis not enabled")
	}

	data, err := RedisGet(c.l2Key(key))
	if err != nil {
		return zero, err
	}

	// 检查是否为空结果标记
	if data == "__EMPTY__" {
		return zero, fmt.Errorf("empty marker")
	}

	var val T
	if err := json.Unmarshal([]byte(data), &val); err != nil {
		return zero, err
	}

	return val, nil
}

func (c *MultiLevelCache[T]) setToL2(key string, value T) {
	if !RedisEnabled {
		return
	}

	data, err := json.Marshal(value)
	if err != nil {
		SysLog(fmt.Sprintf("failed to marshal cache value: %v", err))
		return
	}

	ttl := c.addJitter(c.config.L2TTL)
	if err := RedisSet(c.l2Key(key), string(data), ttl); err != nil {
		SysLog(fmt.Sprintf("failed to set L2 cache: %v", err))
	}
}

func (c *MultiLevelCache[T]) setEmptyToL2(key string) {
	if !RedisEnabled {
		return
	}

	if err := RedisSet(c.l2Key(key), "__EMPTY__", c.config.EmptyTTL); err != nil {
		SysLog(fmt.Sprintf("failed to set L2 empty cache: %v", err))
	}
}

func (c *MultiLevelCache[T]) deleteFromL2(key string) {
	if !RedisEnabled {
		return
	}

	if err := RedisDel(c.l2Key(key)); err != nil {
		SysLog(fmt.Sprintf("failed to delete L2 cache: %v", err))
	}
}

// CacheStats 缓存统计
type CacheStats struct {
	Name     string  `json:"name"`
	L1Hits   int64   `json:"l1_hits"`
	L1Misses int64   `json:"l1_misses"`
	L2Hits   int64   `json:"l2_hits"`
	L2Misses int64   `json:"l2_misses"`
	HitRate  float64 `json:"hit_rate"`
	Size     int     `json:"size"`
	Capacity int     `json:"capacity"`
}

// StartCacheCleanup 启动后台缓存清理任务
func (c *MultiLevelCache[T]) StartCacheCleanup(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				cleaned := c.l1.CleanExpired()
				if cleaned > 0 && DebugEnabled {
					SysLog(fmt.Sprintf("cache %s cleaned %d expired entries", c.name, cleaned))
				}
				c.CheckAndWarn()
			}
		}
	}()
}

package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// LogQueryRateLimiter 日志查询速率限制器
// 使用滑动窗口算法，默认 30 次/分钟
type LogQueryRateLimiter struct {
	maxPerMin    int           // 每分钟最大请求数，默认30
	windowSize   time.Duration // 窗口大小，默认1分钟
	redisEnabled bool
}

var (
	logQueryRateLimiter     *LogQueryRateLimiter
	logQueryRateLimiterOnce sync.Once
	logQueryInMemoryLimiter common.InMemoryRateLimiter
)

// DefaultLogQueryMaxPerMin 默认每分钟最大查询次数
const DefaultLogQueryMaxPerMin = 30

// DefaultLogQueryWindowSize 默认窗口大小（秒）
const DefaultLogQueryWindowSize = 60

// NewLogQueryRateLimiter 创建日志查询速率限制器
func NewLogQueryRateLimiter() *LogQueryRateLimiter {
	logQueryRateLimiterOnce.Do(func() {
		logQueryRateLimiter = &LogQueryRateLimiter{
			maxPerMin:    DefaultLogQueryMaxPerMin,
			windowSize:   time.Duration(DefaultLogQueryWindowSize) * time.Second,
			redisEnabled: common.RedisEnabled,
		}
		if !common.RedisEnabled {
			logQueryInMemoryLimiter.Init(common.RateLimitKeyExpirationDuration)
		}
	})
	return logQueryRateLimiter
}

// Allow 检查是否允许查询
// 返回是否允许和需要等待的时间
func (r *LogQueryRateLimiter) Allow(userID int) (bool, time.Duration) {
	key := r.buildKey(userID)

	if r.redisEnabled {
		return r.redisAllow(key)
	}
	return r.memoryAllow(key)
}

// buildKey 构建限流键
func (r *LogQueryRateLimiter) buildKey(userID int) string {
	return fmt.Sprintf("log_query_rate:%d", userID)
}

// redisAllow 使用 Redis 实现滑动窗口限流
func (r *LogQueryRateLimiter) redisAllow(key string) (bool, time.Duration) {
	ctx := context.Background()
	rdb := common.RDB
	now := time.Now()
	windowStart := now.Add(-r.windowSize).UnixMilli()
	nowMilli := now.UnixMilli()

	// 使用 Redis ZSET 实现滑动窗口
	// 1. 移除窗口外的旧记录
	rdb.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))

	// 2. 获取当前窗口内的请求数
	count, err := rdb.ZCard(ctx, key).Result()
	if err != nil {
		// Redis 错误时放行
		common.SysLog(fmt.Sprintf("LogQueryRateLimiter redis error: %v", err))
		return true, 0
	}

	// 3. 检查是否超过限制
	if count >= int64(r.maxPerMin) {
		// 获取最早的请求时间，计算需要等待的时间
		oldest, err := rdb.ZRange(ctx, key, 0, 0).Result()
		if err == nil && len(oldest) > 0 {
			oldestTime, _ := strconv.ParseInt(oldest[0], 10, 64)
			retryAfter := time.Duration(oldestTime+r.windowSize.Milliseconds()-nowMilli) * time.Millisecond
			if retryAfter < 0 {
				retryAfter = time.Second
			}
			return false, retryAfter
		}
		return false, time.Second
	}

	// 4. 添加当前请求记录
	rdb.ZAdd(ctx, key, &redis.Z{
		Score:  float64(nowMilli),
		Member: strconv.FormatInt(nowMilli, 10),
	})
	rdb.Expire(ctx, key, r.windowSize+time.Minute) // 设置过期时间

	return true, 0
}

// memoryAllow 使用内存实现滑动窗口限流
func (r *LogQueryRateLimiter) memoryAllow(key string) (bool, time.Duration) {
	allowed := logQueryInMemoryLimiter.Request(key, r.maxPerMin, int64(r.windowSize.Seconds()))
	if !allowed {
		// 简单估算重试时间
		return false, time.Second * 5
	}
	return true, 0
}

// GetRetryAfter 获取重试等待时间
func (r *LogQueryRateLimiter) GetRetryAfter(userID int) time.Duration {
	_, retryAfter := r.Allow(userID)
	return retryAfter
}

// LogQueryRateLimit 日志查询速率限制中间件
// 返回 HTTP 429 状态码和 Retry-After 响应头
func LogQueryRateLimit() gin.HandlerFunc {
	limiter := NewLogQueryRateLimiter()

	return func(c *gin.Context) {
		// 获取用户 ID
		userID := c.GetInt("id")
		if userID == 0 {
			// 未认证用户使用 IP 作为标识
			userID = hashIP(c.ClientIP())
		}

		allowed, retryAfter := limiter.Allow(userID)
		if !allowed {
			// 设置 Retry-After 响应头（秒）
			retryAfterSec := int(retryAfter.Seconds())
			if retryAfterSec < 1 {
				retryAfterSec = 1
			}
			c.Header("Retry-After", strconv.Itoa(retryAfterSec))
			c.Header("X-RateLimit-Limit", strconv.Itoa(limiter.maxPerMin))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(retryAfter).Unix(), 10))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": fmt.Sprintf("请求过于频繁，请 %d 秒后重试", retryAfterSec),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// hashIP 将 IP 地址转换为整数用于限流
func hashIP(ip string) int {
	hash := 0
	for _, c := range ip {
		hash = hash*31 + int(c)
	}
	if hash < 0 {
		hash = -hash
	}
	return hash
}

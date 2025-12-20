package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/common/limiter"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const (
	ModelRequestRateLimitCountMark        = "MRRL"
	ModelRequestRateLimitSuccessCountMark = "MRRLS"
)

// 检查Redis中的请求限制
func checkRedisRateLimit(ctx context.Context, rdb *redis.Client, key string, maxCount int, duration int64) (bool, error) {
	// 如果maxCount为0，表示不限制
	if maxCount == 0 {
		return true, nil
	}

	// 获取当前计数
	length, err := rdb.LLen(ctx, key).Result()
	if err != nil {
		return false, err
	}

	// 如果未达到限制，允许请求
	if length < int64(maxCount) {
		return true, nil
	}

	// 检查时间窗口
	oldTimeStr, _ := rdb.LIndex(ctx, key, -1).Result()
	oldTime, err := time.Parse(timeFormat, oldTimeStr)
	if err != nil {
		return false, err
	}

	nowTimeStr := time.Now().Format(timeFormat)
	nowTime, err := time.Parse(timeFormat, nowTimeStr)
	if err != nil {
		return false, err
	}
	// 如果在时间窗口内已达到限制，拒绝请求
	subTime := nowTime.Sub(oldTime).Seconds()
	if int64(subTime) < duration {
		rdb.Expire(ctx, key, time.Duration(setting.ModelRequestRateLimitDurationMinutes)*time.Minute)
		return false, nil
	}

	return true, nil
}

// 记录Redis请求
func recordRedisRequest(ctx context.Context, rdb *redis.Client, key string, maxCount int) {
	// 如果maxCount为0，不记录请求
	if maxCount == 0 {
		return
	}

	now := time.Now().Format(timeFormat)
	rdb.LPush(ctx, key, now)
	rdb.LTrim(ctx, key, 0, int64(maxCount-1))
	rdb.Expire(ctx, key, time.Duration(setting.ModelRequestRateLimitDurationMinutes)*time.Minute)
}

// Redis限流处理器
func redisRateLimitHandler(duration int64, totalMaxCount, successMaxCount int) gin.HandlerFunc {
	return redisRateLimitHandlerWithGroup(duration, totalMaxCount, successMaxCount, "")
}

// Redis限流处理器（支持渠道分组）
func redisRateLimitHandlerWithGroup(duration int64, totalMaxCount, successMaxCount int, channelGroup string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := strconv.Itoa(c.GetInt("id"))
		ctx := context.Background()
		rdb := common.RDB

		// 构建 key，包含渠道分组信息
		var successKey, totalKey string
		if channelGroup != "" {
			successKey = fmt.Sprintf("rateLimit:%s:%s:%s", ModelRequestRateLimitSuccessCountMark, userId, channelGroup)
			totalKey = fmt.Sprintf("rateLimit:%s:%s", userId, channelGroup)
		} else {
			successKey = fmt.Sprintf("rateLimit:%s:%s", ModelRequestRateLimitSuccessCountMark, userId)
			totalKey = fmt.Sprintf("rateLimit:%s", userId)
		}

		// 1. 检查成功请求数限制
		allowed, err := checkRedisRateLimit(ctx, rdb, successKey, successMaxCount, duration)
		if err != nil {
			fmt.Println("检查成功请求数限制失败:", err.Error())
			abortWithOpenAiMessage(c, http.StatusInternalServerError, "rate_limit_check_failed")
			return
		}
		if !allowed {
			msg := fmt.Sprintf("您已达到请求数限制：%d分钟内最多请求%d次", setting.ModelRequestRateLimitDurationMinutes, successMaxCount)
			if channelGroup != "" {
				msg = fmt.Sprintf("您已达到请求数限制：%d分钟内最多请求%d次（渠道分组：%s）", setting.ModelRequestRateLimitDurationMinutes, successMaxCount, channelGroup)
			}
			abortWithOpenAiMessage(c, http.StatusTooManyRequests, msg)
			return
		}

		// 2. 检查总请求数限制并记录总请求（当totalMaxCount为0时会自动跳过，使用令牌桶限流器）
		if totalMaxCount > 0 {
			// 初始化
			tb := limiter.New(ctx, rdb)
			allowed, err = tb.Allow(
				ctx,
				totalKey,
				limiter.WithCapacity(int64(totalMaxCount)*duration),
				limiter.WithRate(int64(totalMaxCount)),
				limiter.WithRequested(duration),
			)

			if err != nil {
				fmt.Println("检查总请求数限制失败:", err.Error())
				abortWithOpenAiMessage(c, http.StatusInternalServerError, "rate_limit_check_failed")
				return
			}

			if !allowed {
				msg := fmt.Sprintf("您已达到总请求数限制：%d分钟内最多请求%d次，包括失败次数，请检查您的请求是否正确", setting.ModelRequestRateLimitDurationMinutes, totalMaxCount)
				if channelGroup != "" {
					msg = fmt.Sprintf("您已达到总请求数限制：%d分钟内最多请求%d次（渠道分组：%s），包括失败次数，请检查您的请求是否正确", setting.ModelRequestRateLimitDurationMinutes, totalMaxCount, channelGroup)
				}
				abortWithOpenAiMessage(c, http.StatusTooManyRequests, msg)
			}
		}

		// 4. 处理请求
		c.Next()

		// 5. 如果请求成功，记录成功请求
		if c.Writer.Status() < 400 {
			recordRedisRequest(ctx, rdb, successKey, successMaxCount)
		}
	}
}

// 内存限流处理器
func memoryRateLimitHandler(duration int64, totalMaxCount, successMaxCount int) gin.HandlerFunc {
	return memoryRateLimitHandlerWithGroup(duration, totalMaxCount, successMaxCount, "")
}

// 内存限流处理器（支持渠道分组）
func memoryRateLimitHandlerWithGroup(duration int64, totalMaxCount, successMaxCount int, channelGroup string) gin.HandlerFunc {
	inMemoryRateLimiter.Init(time.Duration(setting.ModelRequestRateLimitDurationMinutes) * time.Minute)

	return func(c *gin.Context) {
		userId := strconv.Itoa(c.GetInt("id"))
		
		// 构建 key，包含渠道分组信息
		var totalKey, successKey string
		if channelGroup != "" {
			totalKey = ModelRequestRateLimitCountMark + userId + ":" + channelGroup
			successKey = ModelRequestRateLimitSuccessCountMark + userId + ":" + channelGroup
		} else {
			totalKey = ModelRequestRateLimitCountMark + userId
			successKey = ModelRequestRateLimitSuccessCountMark + userId
		}

		// 1. 检查总请求数限制（当totalMaxCount为0时跳过）
		if totalMaxCount > 0 && !inMemoryRateLimiter.Request(totalKey, totalMaxCount, duration) {
			c.Status(http.StatusTooManyRequests)
			c.Abort()
			return
		}

		// 2. 检查成功请求数限制
		// 使用一个临时key来检查限制，这样可以避免实际记录
		checkKey := successKey + "_check"
		if !inMemoryRateLimiter.Request(checkKey, successMaxCount, duration) {
			c.Status(http.StatusTooManyRequests)
			c.Abort()
			return
		}

		// 3. 处理请求
		c.Next()

		// 4. 如果请求成功，记录到实际的成功请求计数中
		if c.Writer.Status() < 400 {
			inMemoryRateLimiter.Request(successKey, successMaxCount, duration)
		}
	}
}

// ModelRequestRateLimit 模型请求限流中间件
func ModelRequestRateLimit() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 在每个请求时检查是否启用限流
		if !setting.ModelRequestRateLimitEnabled {
			c.Next()
			return
		}

		// 计算限流参数
		duration := int64(setting.ModelRequestRateLimitDurationMinutes * 60)
		totalMaxCount := setting.ModelRequestRateLimitCount
		successMaxCount := setting.ModelRequestRateLimitSuccessCount

		// 获取渠道分组
		channelGroup := common.GetContextKeyString(c, constant.ContextKeyTokenGroup)
		if channelGroup == "" {
			channelGroup = common.GetContextKeyString(c, constant.ContextKeyUserGroup)
		}

		// 获取分组的限流配置
		groupTotalCount, groupSuccessCount, found := setting.GetGroupRateLimit(channelGroup)
		if found {
			totalMaxCount = groupTotalCount
			successMaxCount = groupSuccessCount
		}

		// 获取用户等级的限流配置（优先级高于分组配置）
		// 优先使用渠道分组级别的限流配置
		userId := c.GetInt("id")
		if userId > 0 {
			levelTotalCount, levelSuccessCount, levelFound := GetUserLevelRateLimitForGroup(userId, channelGroup)
			if levelFound {
				totalMaxCount = levelTotalCount
				successMaxCount = levelSuccessCount
			}
		}

		// 根据存储类型选择并执行限流处理器
		if common.RedisEnabled {
			redisRateLimitHandlerWithGroup(duration, totalMaxCount, successMaxCount, channelGroup)(c)
		} else {
			memoryRateLimitHandlerWithGroup(duration, totalMaxCount, successMaxCount, channelGroup)(c)
		}
	}
}

// GetUserLevelRateLimit 获取用户等级的速率限制配置（全局限流）
func GetUserLevelRateLimit(userId int) (totalCount int, successCount int, found bool) {
	// 使用 model 包获取用户等级
	user, err := model.GetUserById(userId, false)
	if err != nil || user == nil {
		return 0, 0, false
	}

	levelConfig, err := model.GetLevelConfigById(user.Level)
	if err != nil || levelConfig == nil {
		return 0, 0, false
	}

	benefits, err := levelConfig.GetBenefits()
	if err != nil || benefits == nil || benefits.RateLimit == nil {
		return 0, 0, false
	}

	// 如果配置了速率限制
	if benefits.RateLimit.TotalCount > 0 || benefits.RateLimit.SuccessCount > 0 {
		return benefits.RateLimit.TotalCount, benefits.RateLimit.SuccessCount, true
	}

	return 0, 0, false
}

// GetUserLevelRateLimitForGroup 获取用户等级在指定渠道分组的速率限制配置
// 优先级：group_rate_limits > rate_limit > 无限制
func GetUserLevelRateLimitForGroup(userId int, channelGroup string) (totalCount int, successCount int, found bool) {
	// 使用 model 包获取用户等级
	user, err := model.GetUserById(userId, false)
	if err != nil || user == nil {
		return 0, 0, false
	}

	levelConfig, err := model.GetLevelConfigById(user.Level)
	if err != nil || levelConfig == nil {
		return 0, 0, false
	}

	benefits, err := levelConfig.GetBenefits()
	if err != nil || benefits == nil {
		return 0, 0, false
	}

	// 1. 优先检查渠道分组限流
	if benefits.GroupRateLimits != nil && channelGroup != "" {
		if limit, exists := benefits.GroupRateLimits[channelGroup]; exists && limit != nil {
			if limit.TotalCount > 0 || limit.SuccessCount > 0 {
				return limit.TotalCount, limit.SuccessCount, true
			}
		}
	}

	// 2. 回退到全局限流
	if benefits.RateLimit != nil {
		if benefits.RateLimit.TotalCount > 0 || benefits.RateLimit.SuccessCount > 0 {
			return benefits.RateLimit.TotalCount, benefits.RateLimit.SuccessCount, true
		}
	}

	// 3. 无限制
	return 0, 0, false
}

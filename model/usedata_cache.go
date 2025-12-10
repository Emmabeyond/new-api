package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/QuantumNous/new-api/common"
)

// QuotaDataCacheTTL 缓存过期时间
const QuotaDataCacheTTL = 5 * time.Minute

// QuotaDataCacheKeyFmt 缓存 key 格式
// type: admin (管理员全局查询) 或 user:{userId} (用户查询)
const QuotaDataCacheKeyFmt = "quota_data:%s:%d:%d:%s"

// QuotaDataCacheKey 生成缓存 key
// userId: 用户 ID，0 表示管理员全局查询
// startTime, endTime: 时间范围（Unix 时间戳）
// timeUnit: 时间单位 (hour/day/week)
func QuotaDataCacheKey(userId int, startTime, endTime int64, timeUnit string) string {
	var typeStr string
	if userId == 0 {
		typeStr = "admin"
	} else {
		typeStr = fmt.Sprintf("user:%d", userId)
	}
	return fmt.Sprintf(QuotaDataCacheKeyFmt, typeStr, startTime, endTime, timeUnit)
}

// GetQuotaDataWithCache 带缓存的查询
// userId: 用户 ID，0 表示管理员全局查询
// username: 用户名（仅管理员查询时使用）
// startTime, endTime: 时间范围
// timeUnit: 时间单位
func GetQuotaDataWithCache(userId int, username string, startTime, endTime int64, timeUnit string) ([]*QuotaData, error) {
	// 如果 Redis 未启用，直接查询数据库
	if !common.RedisEnabled {
		return getQuotaDataFromDB(userId, username, startTime, endTime)
	}

	// 生成缓存 key
	cacheKey := QuotaDataCacheKey(userId, startTime, endTime, timeUnit)

	// 尝试从缓存获取
	cachedData, err := common.RedisGet(cacheKey)
	if err == nil && cachedData != "" {
		var result []*QuotaData
		if err := json.Unmarshal([]byte(cachedData), &result); err == nil {
			return result, nil
		}
		// 反序列化失败，删除无效缓存
		common.RedisDel(cacheKey)
	}

	// 缓存未命中，查询数据库
	result, err := getQuotaDataFromDB(userId, username, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// 将结果写入缓存
	if data, err := json.Marshal(result); err == nil {
		common.RedisSet(cacheKey, string(data), QuotaDataCacheTTL)
	}

	return result, nil
}

// getQuotaDataFromDB 从数据库查询数据
func getQuotaDataFromDB(userId int, username string, startTime, endTime int64) ([]*QuotaData, error) {
	if userId == 0 {
		// 管理员查询
		return GetAllQuotaDates(startTime, endTime, username)
	}
	// 用户查询
	return GetQuotaDataByUserId(userId, startTime, endTime)
}

// InvalidateQuotaDataCache 使用户相关的缓存失效
// 由于缓存 key 包含时间范围，无法精确删除所有相关缓存
// 这里使用模式匹配删除（需要 Redis SCAN 支持）
// 简化实现：不主动删除，依赖 TTL 自动过期
func InvalidateQuotaDataCache(userId int) {
	// 简化实现：缓存 TTL 较短（5分钟），依赖自动过期
	// 如果需要立即失效，可以使用 Redis SCAN 命令扫描并删除匹配的 key
	// 但这会增加复杂度和 Redis 负载
	//
	// 未来优化方向：
	// 1. 使用 Redis 的 key 前缀 + SCAN 命令批量删除
	// 2. 使用缓存版本号机制
}

// InvalidateAllQuotaDataCache 使所有 quota_data 缓存失效
// 用于管理员操作或数据批量更新场景
func InvalidateAllQuotaDataCache() {
	// 同上，依赖 TTL 自动过期
}

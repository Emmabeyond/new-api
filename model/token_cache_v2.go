package model

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/bytedance/gopkg/util/gopool"
)

// TokenCacheService Token 多级缓存服务
// 实现 L1 本地缓存 + L2 Redis 缓存的两级缓存架构
type TokenCacheService struct {
	cache *common.MultiLevelCache[*Token]
}

// tokenCacheServiceInstance 全局单例
var tokenCacheServiceInstance *TokenCacheService

// GetTokenCacheService 获取 Token 缓存服务单例
func GetTokenCacheService() *TokenCacheService {
	return tokenCacheServiceInstance
}

// InitTokenCacheService 初始化 Token 缓存服务
func InitTokenCacheService() {
	config := common.DefaultMultiLevelCacheConfig()

	// 从环境变量读取配置
	if capacity := os.Getenv("L1_CACHE_CAPACITY"); capacity != "" {
		if val, err := strconv.Atoi(capacity); err == nil && val > 0 {
			config.L1Capacity = val
		}
	}

	if ttl := os.Getenv("L1_CACHE_TTL"); ttl != "" {
		if val, err := strconv.Atoi(ttl); err == nil && val > 0 {
			config.L1TTL = time.Duration(val) * time.Second
		}
	}

	if jitter := os.Getenv("CACHE_TTL_JITTER"); jitter != "" {
		if val, err := strconv.ParseFloat(jitter, 64); err == nil && val >= 0 && val <= 1 {
			config.TTLJitter = val
		}
	}

	tokenCacheServiceInstance = &TokenCacheService{
		cache: common.NewMultiLevelCache[*Token]("token", config),
	}

	// 启动后台清理任务
	tokenCacheServiceInstance.cache.StartCacheCleanup(context.Background(), 5*time.Minute)

	common.SysLog("Token cache service initialized with L1 capacity: " + strconv.Itoa(config.L1Capacity))
}

// generateCacheKey 生成缓存 key
// 使用 HMAC 对原始 key 进行哈希，保护敏感信息
func (s *TokenCacheService) generateCacheKey(key string) string {
	return common.GenerateHMAC(key)
}

// GetByKey 通过 key 获取 Token
// 使用多级缓存：L1 -> L2 -> DB
// Requirements: 3.1, 3.2, 3.3, 5.1
func (s *TokenCacheService) GetByKey(key string) (*Token, error) {
	if key == "" {
		return nil, fmt.Errorf("token key is empty")
	}

	cacheKey := s.generateCacheKey(key)

	token, err := s.cache.Get(cacheKey, func() (*Token, bool, error) {
		// 从数据库加载
		var token Token
		err := DB.Where(commonKeyCol+" = ?", key).First(&token).Error
		if err != nil {
			// Token 不存在
			return nil, false, nil
		}
		// 清理敏感信息后缓存
		token.Clean()
		return &token, true, nil
	})

	if err != nil {
		return nil, err
	}

	// 恢复原始 key（因为缓存中的 token 已经 Clean 过）
	token.Key = key
	return token, nil
}

// UpdateQuota 更新配额并同步缓存
// Requirements: 6.1
func (s *TokenCacheService) UpdateQuota(id int, key string, delta int64) error {
	if id == 0 {
		return fmt.Errorf("token id is empty")
	}

	// 1. 更新数据库
	var err error
	if delta > 0 {
		err = DB.Model(&Token{}).Where("id = ?", id).Updates(
			map[string]interface{}{
				"remain_quota":  DB.Raw("remain_quota + ?", delta),
				"used_quota":    DB.Raw("used_quota - ?", delta),
				"accessed_time": common.GetTimestamp(),
			},
		).Error
	} else {
		err = DB.Model(&Token{}).Where("id = ?", id).Updates(
			map[string]interface{}{
				"remain_quota":  DB.Raw("remain_quota - ?", -delta),
				"used_quota":    DB.Raw("used_quota + ?", -delta),
				"accessed_time": common.GetTimestamp(),
			},
		).Error
	}

	if err != nil {
		return err
	}

	// 2. 使缓存失效，下次访问时重新加载
	if key != "" {
		cacheKey := s.generateCacheKey(key)
		s.cache.Invalidate(cacheKey)

		// 同时更新旧的 Redis 缓存（兼容性）
		if common.RedisEnabled {
			gopool.Go(func() {
				_ = cacheIncrTokenQuota(key, delta)
			})
		}
	}

	return nil
}

// Disable 禁用 Token 并使缓存失效
// Requirements: 6.2
func (s *TokenCacheService) Disable(key string) error {
	if key == "" {
		return fmt.Errorf("token key is empty")
	}

	// 1. 更新数据库状态
	err := DB.Model(&Token{}).Where(commonKeyCol+" = ?", key).Update("status", common.TokenStatusDisabled).Error
	if err != nil {
		return err
	}

	// 2. 立即使缓存失效
	cacheKey := s.generateCacheKey(key)
	s.cache.Invalidate(cacheKey)

	// 同时删除旧的 Redis 缓存（兼容性）
	if common.RedisEnabled {
		gopool.Go(func() {
			_ = cacheDeleteToken(key)
		})
	}

	return nil
}

// Set 直接设置缓存（用于更新后回填）
func (s *TokenCacheService) Set(token *Token) {
	if token == nil || token.Key == "" {
		return
	}

	cacheKey := s.generateCacheKey(token.Key)
	// 复制一份并清理敏感信息
	cachedToken := *token
	cachedToken.Clean()
	s.cache.Set(cacheKey, &cachedToken)
}

// Invalidate 使指定 key 的缓存失效
func (s *TokenCacheService) Invalidate(key string) {
	if key == "" {
		return
	}
	cacheKey := s.generateCacheKey(key)
	s.cache.Invalidate(cacheKey)
}

// Stats 返回缓存统计信息
func (s *TokenCacheService) Stats() common.CacheStats {
	return s.cache.Stats()
}

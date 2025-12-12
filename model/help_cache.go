package model

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
)

const (
	HelpCategoriesCacheKey = "help:categories"
	HelpDocumentsCacheKey  = "help:documents"
	HelpDocumentCacheKey   = "help:document:%d"
	HelpCacheTTL           = 3600 // 1 hour in seconds
)

// 内存缓存（Redis 不可用时的降级方案）
var (
	memoryCacheMutex      sync.RWMutex
	memoryCacheCategories []HelpCategoryWithDocuments
	memoryCacheExpireAt   time.Time
)

// GetCachedHelpData 获取缓存的帮助中心数据（分类+文档）
func GetCachedHelpData() ([]HelpCategoryWithDocuments, error) {
	// 尝试从 Redis 获取
	if common.RedisEnabled {
		data, err := getCachedHelpDataFromRedis()
		if err == nil && data != nil {
			return data, nil
		}
		// Redis 获取失败，尝试从数据库获取并缓存
	}

	// 尝试从内存缓存获取
	memoryCacheMutex.RLock()
	if memoryCacheCategories != nil && time.Now().Before(memoryCacheExpireAt) {
		data := memoryCacheCategories
		memoryCacheMutex.RUnlock()
		return data, nil
	}
	memoryCacheMutex.RUnlock()

	// 从数据库获取
	data, err := GetHelpDocumentsGroupedByCategory()
	if err != nil {
		return nil, err
	}

	// 缓存数据
	cacheHelpData(data)

	return data, nil
}

// getCachedHelpDataFromRedis 从 Redis 获取缓存数据
func getCachedHelpDataFromRedis() ([]HelpCategoryWithDocuments, error) {
	jsonStr, err := common.RedisGet(HelpDocumentsCacheKey)
	if err != nil {
		return nil, err
	}

	var data []HelpCategoryWithDocuments
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// cacheHelpData 缓存帮助中心数据
func cacheHelpData(data []HelpCategoryWithDocuments) {
	// 缓存到 Redis
	if common.RedisEnabled {
		jsonBytes, err := json.Marshal(data)
		if err == nil {
			_ = common.RedisSet(HelpDocumentsCacheKey, string(jsonBytes), time.Duration(HelpCacheTTL)*time.Second)
		}
	}

	// 同时缓存到内存（作为降级方案）
	memoryCacheMutex.Lock()
	memoryCacheCategories = data
	memoryCacheExpireAt = time.Now().Add(time.Duration(HelpCacheTTL) * time.Second)
	memoryCacheMutex.Unlock()
}

// GetCachedHelpDocument 获取缓存的单个文档
func GetCachedHelpDocument(id int) (*HelpDocument, error) {
	cacheKey := fmt.Sprintf(HelpDocumentCacheKey, id)

	// 尝试从 Redis 获取
	if common.RedisEnabled {
		jsonStr, err := common.RedisGet(cacheKey)
		if err == nil {
			var doc HelpDocument
			if json.Unmarshal([]byte(jsonStr), &doc) == nil {
				return &doc, nil
			}
		}
	}

	// 从数据库获取
	doc, err := GetHelpDocumentByIdPublic(id)
	if err != nil {
		return nil, err
	}

	// 缓存到 Redis
	if common.RedisEnabled {
		jsonBytes, err := json.Marshal(doc)
		if err == nil {
			_ = common.RedisSet(cacheKey, string(jsonBytes), time.Duration(HelpCacheTTL)*time.Second)
		}
	}

	return doc, nil
}

// InvalidateHelpCache 清除帮助中心缓存
func InvalidateHelpCache() {
	// 清除 Redis 缓存
	if common.RedisEnabled {
		_ = common.RedisDel(HelpDocumentsCacheKey)
		_ = common.RedisDel(HelpCategoriesCacheKey)
		// 注意：单个文档缓存会在 TTL 后自动过期
	}

	// 清除内存缓存
	memoryCacheMutex.Lock()
	memoryCacheCategories = nil
	memoryCacheExpireAt = time.Time{}
	memoryCacheMutex.Unlock()
}

// InvalidateHelpDocumentCache 清除单个文档缓存
func InvalidateHelpDocumentCache(id int) {
	if common.RedisEnabled {
		cacheKey := fmt.Sprintf(HelpDocumentCacheKey, id)
		_ = common.RedisDel(cacheKey)
	}
	// 同时清除列表缓存，因为文档信息可能在列表中
	InvalidateHelpCache()
}

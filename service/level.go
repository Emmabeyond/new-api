package service

import (
	"errors"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/model"
)

var (
	levelConfigCache     []*model.LevelConfig
	levelConfigCacheMap  map[string]*model.LevelConfig
	levelConfigCacheLock sync.RWMutex
	levelConfigCacheTime time.Time
	levelCacheTTL        = 1 * time.Hour
)

// RefreshLevelConfigCache 刷新等级配置缓存
func RefreshLevelConfigCache() error {
	levels, err := model.GetAllLevelConfigs()
	if err != nil {
		return err
	}

	levelConfigCacheLock.Lock()
	defer levelConfigCacheLock.Unlock()

	levelConfigCache = levels
	levelConfigCacheMap = make(map[string]*model.LevelConfig)
	for _, level := range levels {
		levelConfigCacheMap[level.Id] = level
	}
	levelConfigCacheTime = time.Now()
	return nil
}

// GetAllLevelConfigs 获取所有等级配置（带缓存）
func GetAllLevelConfigs() ([]*model.LevelConfig, error) {
	levelConfigCacheLock.RLock()
	if levelConfigCache != nil && time.Since(levelConfigCacheTime) < levelCacheTTL {
		defer levelConfigCacheLock.RUnlock()
		return levelConfigCache, nil
	}
	levelConfigCacheLock.RUnlock()

	// 缓存过期，重新加载
	if err := RefreshLevelConfigCache(); err != nil {
		return nil, err
	}

	levelConfigCacheLock.RLock()
	defer levelConfigCacheLock.RUnlock()
	return levelConfigCache, nil
}

// GetLevelConfigById 根据ID获取等级配置（带缓存）
func GetLevelConfigById(levelId string) (*model.LevelConfig, error) {
	levelConfigCacheLock.RLock()
	if levelConfigCacheMap != nil && time.Since(levelConfigCacheTime) < levelCacheTTL {
		if level, ok := levelConfigCacheMap[levelId]; ok {
			levelConfigCacheLock.RUnlock()
			return level, nil
		}
	}
	levelConfigCacheLock.RUnlock()

	// 缓存未命中，从数据库获取
	return model.GetLevelConfigById(levelId)
}

// GetDefaultLevelConfig 获取默认等级配置
func GetDefaultLevelConfig() (*model.LevelConfig, error) {
	levels, err := GetAllLevelConfigs()
	if err != nil {
		return nil, err
	}
	for _, level := range levels {
		if level.IsDefault {
			return level, nil
		}
	}
	return nil, errors.New("no default level config found")
}


// CreateLevelConfig 创建等级配置
func CreateLevelConfig(level *model.LevelConfig) error {
	// 验证配置
	if err := model.ValidateLevelConfig(level); err != nil {
		return err
	}

	// 检查ID是否已存在
	exists, err := model.CheckLevelIdExists(level.Id)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("等级ID已存在")
	}

	// 检查优先级是否已存在
	exists, err = model.CheckLevelPriorityExists(level.Priority, "")
	if err != nil {
		return err
	}
	if exists {
		return errors.New("等级优先级已存在")
	}

	// 创建等级
	if err := model.CreateLevelConfig(level); err != nil {
		return err
	}

	// 刷新缓存
	return RefreshLevelConfigCache()
}

// UpdateLevelConfig 更新等级配置
func UpdateLevelConfig(level *model.LevelConfig) error {
	// 验证配置
	if err := model.ValidateLevelConfig(level); err != nil {
		return err
	}

	// 检查优先级是否已存在（排除自身）
	exists, err := model.CheckLevelPriorityExists(level.Priority, level.Id)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("等级优先级已存在")
	}

	// 更新等级
	if err := model.UpdateLevelConfig(level); err != nil {
		return err
	}

	// 刷新缓存
	return RefreshLevelConfigCache()
}

// DeleteLevelConfig 删除等级配置
func DeleteLevelConfig(levelId string) error {
	// 检查是否为默认等级
	level, err := model.GetLevelConfigById(levelId)
	if err != nil {
		return err
	}
	if level.IsDefault {
		return errors.New("不能删除默认等级")
	}

	// 检查是否有用户使用该等级
	count, err := model.GetLevelUserCount(levelId)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该等级下有用户，无法删除")
	}

	// 删除等级
	if err := model.DeleteLevelConfig(levelId); err != nil {
		return err
	}

	// 刷新缓存
	return RefreshLevelConfigCache()
}

// GetLevelUserStats 获取等级用户统计
func GetLevelUserStats() (map[string]int64, error) {
	return model.GetAllLevelUserStats()
}

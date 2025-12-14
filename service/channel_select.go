package service

import (
	"errors"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting"
	"github.com/gin-gonic/gin"
)

// useOptimizedSelector 是否使用优化选择器
// 可通过环境变量 USE_OPTIMIZED_CHANNEL_SELECTOR 控制
var useOptimizedSelector = true

func init() {
	// 注册优化选择器重建函数到 model 包
	model.SetOptimizedSelectorRebuildFunc(func(channels []*model.Channel, abilities []*model.Ability) {
		GetOptimizedSelector().Rebuild(channels, abilities)
	})
}

// getChannelWithOptimizedSelector 使用优化选择器获取渠道
func getChannelWithOptimizedSelector(group string, modelName string, retry int) (*model.Channel, error) {
	selector := GetOptimizedSelector()
	if !selector.IsInitialized() {
		// 优化选择器未初始化，降级到原有逻辑
		return model.GetRandomSatisfiedChannel(group, modelName, retry)
	}
	return selector.Select(group, modelName, retry)
}

func CacheGetRandomSatisfiedChannel(c *gin.Context, group string, modelName string, retry int) (*model.Channel, string, error) {
	var channel *model.Channel
	var err error
	selectGroup := group
	userGroup := common.GetContextKeyString(c, constant.ContextKeyUserGroup)
	if group == "auto" {
		if len(setting.GetAutoGroups()) == 0 {
			return nil, selectGroup, errors.New("auto groups is not enabled")
		}
		for _, autoGroup := range GetUserAutoGroup(userGroup) {
			logger.LogDebug(c, "Auto selecting group:", autoGroup)
			// 优先使用优化选择器
			if useOptimizedSelector {
				channel, _ = getChannelWithOptimizedSelector(autoGroup, modelName, retry)
			} else {
				channel, _ = model.GetRandomSatisfiedChannel(autoGroup, modelName, retry)
			}
			if channel == nil {
				continue
			} else {
				c.Set("auto_group", autoGroup)
				selectGroup = autoGroup
				logger.LogDebug(c, "Auto selected group:", autoGroup)
				break
			}
		}
	} else {
		// 优先使用优化选择器
		if useOptimizedSelector {
			channel, err = getChannelWithOptimizedSelector(group, modelName, retry)
		} else {
			channel, err = model.GetRandomSatisfiedChannel(group, modelName, retry)
		}
		if err != nil {
			return nil, group, err
		}
	}
	return channel, selectGroup, nil
}

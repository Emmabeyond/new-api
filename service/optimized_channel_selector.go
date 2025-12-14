package service

import (
	"errors"
	"math/rand"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting/ratio_setting"
)

// ChannelGroup 表示同一优先级的渠道组
// 预计算权重累积数组，支持 O(log n) 的二分查找选择
type ChannelGroup struct {
	Priority         int64            // 优先级
	Channels         []*model.Channel // 渠道列表
	CumulativeWeight []int            // 预计算的权重累积数组 [w1, w1+w2, w1+w2+w3, ...]
	TotalWeight      int              // 总权重
}

// SelectorData 选择器数据结构
// group -> model -> []ChannelGroup (按优先级降序排列)
type SelectorData struct {
	GroupModelChannels map[string]map[string][]ChannelGroup
	ChannelsById       map[int]*model.Channel
}

// OptimizedSelector 优化的渠道选择器
// 使用原子指针实现无锁读取，支持热更新
type OptimizedSelector struct {
	data atomic.Pointer[SelectorData]
}

// 全局优化选择器实例
var globalOptimizedSelector = &OptimizedSelector{}

// GetOptimizedSelector 获取全局优化选择器实例
func GetOptimizedSelector() *OptimizedSelector {
	return globalOptimizedSelector
}

// buildCumulativeWeights 构建权重累积数组
// 输入: 渠道列表
// 输出: 累积权重数组和总权重
// 示例: 权重 [10, 20, 30] -> 累积 [10, 30, 60], 总权重 60
func buildCumulativeWeights(channels []*model.Channel) ([]int, int) {
	if len(channels) == 0 {
		return nil, 0
	}

	cumulative := make([]int, len(channels))
	sum := 0

	for i, ch := range channels {
		weight := ch.GetWeight()
		// 权重为 0 时，给予最小权重 1，确保有被选中的机会
		if weight <= 0 {
			weight = 1
		}
		sum += weight
		cumulative[i] = sum
	}

	return cumulative, sum
}

// binarySearchChannel 使用二分查找选择渠道
// 时间复杂度: O(log n)
// randomWeight: 0 到 TotalWeight-1 之间的随机数
func binarySearchChannel(group *ChannelGroup, randomWeight int) *model.Channel {
	if group == nil || len(group.Channels) == 0 {
		return nil
	}

	// 边界情况：单个渠道
	if len(group.Channels) == 1 {
		return group.Channels[0]
	}

	// 使用 sort.Search 进行二分查找
	// 找到第一个累积权重 > randomWeight 的索引
	idx := sort.Search(len(group.CumulativeWeight), func(i int) bool {
		return group.CumulativeWeight[i] > randomWeight
	})

	// 确保索引在有效范围内
	if idx >= len(group.Channels) {
		idx = len(group.Channels) - 1
	}

	return group.Channels[idx]
}

// Select 使用优化算法选择渠道
// group: 用户分组
// modelName: 模型名称
// retry: 重试次数，用于降级到次优先级
// 返回: 选中的渠道，错误信息
func (s *OptimizedSelector) Select(group, modelName string, retry int) (*model.Channel, error) {
	startTime := time.Now()
	defer func() {
		// 记录选择耗时
		common.GetCacheMetrics().RecordSelectTime("channel_selector", time.Since(startTime))
	}()

	data := s.data.Load()
	if data == nil {
		return nil, errors.New("selector not initialized")
	}

	// 查找 group -> model 对应的渠道组列表
	modelChannels, ok := data.GroupModelChannels[group]
	if !ok {
		return nil, nil
	}

	channelGroups, ok := modelChannels[modelName]
	if !ok || len(channelGroups) == 0 {
		// 尝试使用规范化的模型名称
		normalizedModel := ratio_setting.FormatMatchingModelName(modelName)
		channelGroups, ok = modelChannels[normalizedModel]
		if !ok || len(channelGroups) == 0 {
			return nil, nil
		}
	}

	// 根据 retry 次数选择优先级组
	// retry=0 选择最高优先级，retry=1 选择次高优先级，以此类推
	groupIndex := retry
	if groupIndex >= len(channelGroups) {
		groupIndex = len(channelGroups) - 1
	}

	targetGroup := &channelGroups[groupIndex]

	// 如果总权重为 0，返回第一个渠道
	if targetGroup.TotalWeight <= 0 {
		if len(targetGroup.Channels) > 0 {
			return targetGroup.Channels[0], nil
		}
		return nil, nil
	}

	// 生成随机数并使用二分查找选择渠道
	randomWeight := rand.Intn(targetGroup.TotalWeight)
	return binarySearchChannel(targetGroup, randomWeight), nil
}

// Rebuild 重建选择数据结构
// 在渠道配置变化时调用
// 使用原子操作更新，保证并发安全
func (s *OptimizedSelector) Rebuild(channels []*model.Channel, abilities []*model.Ability) {
	if len(channels) == 0 {
		s.data.Store(nil)
		return
	}

	// 构建 channelId -> channel 映射
	channelsById := make(map[int]*model.Channel, len(channels))
	for _, ch := range channels {
		channelsById[ch.Id] = ch
	}

	// 收集所有 group
	groups := make(map[string]bool)
	for _, ability := range abilities {
		groups[ability.Group] = true
	}

	// 构建 group -> model -> []channelId 映射
	groupModelChannelIds := make(map[string]map[string][]int)
	for group := range groups {
		groupModelChannelIds[group] = make(map[string][]int)
	}

	// 遍历启用的渠道，构建映射
	for _, ch := range channels {
		if ch.Status != common.ChannelStatusEnabled {
			continue
		}

		channelGroups := strings.Split(ch.Group, ",")
		models := strings.Split(ch.Models, ",")

		for _, g := range channelGroups {
			g = strings.TrimSpace(g)
			if _, ok := groupModelChannelIds[g]; !ok {
				groupModelChannelIds[g] = make(map[string][]int)
			}

			for _, m := range models {
				m = strings.TrimSpace(m)
				if m == "" {
					continue
				}
				groupModelChannelIds[g][m] = append(groupModelChannelIds[g][m], ch.Id)
			}
		}
	}

	// 构建最终的数据结构
	result := make(map[string]map[string][]ChannelGroup)

	for group, modelChannels := range groupModelChannelIds {
		result[group] = make(map[string][]ChannelGroup)

		for modelName, channelIds := range modelChannels {
			if len(channelIds) == 0 {
				continue
			}

			// 按优先级分组
			priorityMap := make(map[int64][]*model.Channel)
			for _, chId := range channelIds {
				ch, ok := channelsById[chId]
				if !ok {
					continue
				}
				priority := ch.GetPriority()
				priorityMap[priority] = append(priorityMap[priority], ch)
			}

			// 提取并排序优先级（降序）
			priorities := make([]int64, 0, len(priorityMap))
			for p := range priorityMap {
				priorities = append(priorities, p)
			}
			sort.Slice(priorities, func(i, j int) bool {
				return priorities[i] > priorities[j]
			})

			// 为每个优先级构建 ChannelGroup
			channelGroups := make([]ChannelGroup, 0, len(priorities))
			for _, priority := range priorities {
				chs := priorityMap[priority]
				cumulative, total := buildCumulativeWeights(chs)
				channelGroups = append(channelGroups, ChannelGroup{
					Priority:         priority,
					Channels:         chs,
					CumulativeWeight: cumulative,
					TotalWeight:      total,
				})
			}

			result[group][modelName] = channelGroups
		}
	}

	// 原子更新数据
	newData := &SelectorData{
		GroupModelChannels: result,
		ChannelsById:       channelsById,
	}
	s.data.Store(newData)

	common.SysLog("optimized channel selector rebuilt")
}

// IsInitialized 检查选择器是否已初始化
func (s *OptimizedSelector) IsInitialized() bool {
	return s.data.Load() != nil
}

// GetChannelById 从缓存中获取渠道
func (s *OptimizedSelector) GetChannelById(id int) (*model.Channel, bool) {
	data := s.data.Load()
	if data == nil {
		return nil, false
	}
	ch, ok := data.ChannelsById[id]
	return ch, ok
}

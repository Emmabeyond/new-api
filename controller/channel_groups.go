package controller

import (
	"net/http"

	"github.com/QuantumNous/new-api/setting"
	"github.com/gin-gonic/gin"
)

// ChannelGroupInfo 渠道分组信息
type ChannelGroupInfo struct {
	Key   string `json:"key"`   // 分组键名
	Label string `json:"label"` // 分组显示名称
}

// GetChannelGroups 获取所有可用的渠道分组列表
func GetChannelGroups(c *gin.Context) {
	// 从系统设置中获取用户可用分组
	userUsableGroups := setting.GetUserUsableGroupsCopy()

	// 构建分组信息列表（包含键和显示名称）
	groups := make([]ChannelGroupInfo, 0, len(userUsableGroups))
	for key, label := range userUsableGroups {
		groups = append(groups, ChannelGroupInfo{
			Key:   key,
			Label: label,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    groups,
	})
}

package controller

import (
	"net/http"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"

	"github.com/gin-gonic/gin"
)

// MigrateLinuxDOSetting 迁移旧的 LinuxDO 配置到 linuxdo.* 配置系统
func MigrateLinuxDOSetting(c *gin.Context) {
	// 读取全部 option
	opts, err := model.AllOption()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// 建立 map
	valMap := map[string]string{}
	for _, o := range opts {
		valMap[o.Key] = o.Value
	}

	// 迁移 LinuxDO 配置
	migrated := false
	if clientId := valMap["LinuxDOClientId"]; clientId != "" {
		model.UpdateOption("linuxdo.client_id", clientId)
		migrated = true
	}
	if clientSecret := valMap["LinuxDOClientSecret"]; clientSecret != "" {
		model.UpdateOption("linuxdo.client_secret", clientSecret)
		migrated = true
	}
	if trustLevel := valMap["LinuxDOMinimumTrustLevel"]; trustLevel != "" {
		model.UpdateOption("linuxdo.minimum_trust_level", trustLevel)
		migrated = true
	}
	if enabled := valMap["LinuxDOOAuthEnabled"]; enabled != "" {
		model.UpdateOption("linuxdo.enabled", enabled)
		migrated = true
	}

	if migrated {
		// 重新加载 OptionMap
		model.InitOptionMap()
		common.SysLog("LinuxDO setting migrated to new config system")
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "LinuxDO 配置已迁移到新系统"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "无需迁移，未找到旧的 LinuxDO 配置"})
	}
}

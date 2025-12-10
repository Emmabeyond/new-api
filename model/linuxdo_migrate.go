package model

import (
	"github.com/QuantumNous/new-api/common"
)

// MigrateLinuxDOSettingOnStartup 在应用启动时自动迁移 LinuxDO 配置
func MigrateLinuxDOSettingOnStartup() {
	// 读取全部 option
	opts, err := AllOption()
	if err != nil {
		common.SysError("failed to read options for LinuxDO migration: " + err.Error())
		return
	}

	// 建立 map
	valMap := map[string]string{}
	for _, o := range opts {
		valMap[o.Key] = o.Value
	}

	// 检查是否已经迁移（如果新配置已存在，则跳过）
	if _, exists := valMap["linuxdo.client_id"]; exists {
		// 已经迁移过了
		return
	}

	// 检查是否有旧配置需要迁移
	hasOldConfig := false
	if _, exists := valMap["LinuxDOClientId"]; exists {
		hasOldConfig = true
	}
	if _, exists := valMap["LinuxDOClientSecret"]; exists {
		hasOldConfig = true
	}

	if !hasOldConfig {
		// 没有旧配置，无需迁移
		return
	}

	common.SysLog("Migrating LinuxDO settings to new config system...")

	// 迁移 LinuxDO 配置
	if clientId := valMap["LinuxDOClientId"]; clientId != "" {
		UpdateOption("linuxdo.client_id", clientId)
	}
	if clientSecret := valMap["LinuxDOClientSecret"]; clientSecret != "" {
		UpdateOption("linuxdo.client_secret", clientSecret)
	}
	if trustLevel := valMap["LinuxDOMinimumTrustLevel"]; trustLevel != "" {
		UpdateOption("linuxdo.minimum_trust_level", trustLevel)
	}
	if enabled := valMap["LinuxDOOAuthEnabled"]; enabled != "" {
		UpdateOption("linuxdo.enabled", enabled)
	}

	common.SysLog("LinuxDO settings migrated successfully")
}

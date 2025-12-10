package system_setting

import "github.com/QuantumNous/new-api/setting/config"

type LinuxDOSettings struct {
	Enabled            bool   `json:"enabled"`
	ClientId           string `json:"client_id"`
	ClientSecret       string `json:"client_secret"`
	MinimumTrustLevel  int    `json:"minimum_trust_level"`
}

// 默认配置
var defaultLinuxDOSettings = LinuxDOSettings{
	Enabled:           false,
	ClientId:          "",
	ClientSecret:      "",
	MinimumTrustLevel: 1,
}

func init() {
	// 注册到全局配置管理器
	config.GlobalConfig.Register("linuxdo", &defaultLinuxDOSettings)
}

func GetLinuxDOSettings() *LinuxDOSettings {
	return &defaultLinuxDOSettings
}

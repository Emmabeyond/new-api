package system_setting

import (
	"errors"
	"strings"

	"github.com/QuantumNous/new-api/setting/config"
)

// WeChatThirdPartySettings 微信第三方登录配置
// 用于集成第三方微信认证平台（如 abu117.cn）
type WeChatThirdPartySettings struct {
	Enabled      bool   `json:"enabled"`       // 是否启用微信第三方登录
	ClientKey    string `json:"client_key"`    // 平台 ClientKey，格式: wac_xxxxxxxxxxxx
	ClientSecret string `json:"client_secret"` // 平台 ClientSecret，格式: wacs_xxxxxxxxxxxxxxxxxxxxxxxx
	BaseURL      string `json:"base_url"`      // 平台 API 地址，默认 https://www.abu117.cn
}

// 默认配置
var defaultWeChatThirdPartySettings = WeChatThirdPartySettings{
	Enabled:      false,
	ClientKey:    "",
	ClientSecret: "",
	BaseURL:      "https://www.abu117.cn",
}

func init() {
	// 注册到全局配置管理器
	config.GlobalConfig.Register("wechat_third_party", &defaultWeChatThirdPartySettings)
}

// GetWeChatThirdPartySettings 获取微信第三方登录配置
func GetWeChatThirdPartySettings() *WeChatThirdPartySettings {
	return &defaultWeChatThirdPartySettings
}

// ValidateWeChatThirdPartySettings 验证微信第三方登录配置
// 当启用功能时，验证必填字段是否已配置
func ValidateWeChatThirdPartySettings(settings *WeChatThirdPartySettings) error {
	if !settings.Enabled {
		return nil // 未启用时不需要验证
	}

	// 验证 ClientKey
	if strings.TrimSpace(settings.ClientKey) == "" {
		return errors.New("ClientKey 不能为空")
	}

	// 验证 ClientSecret
	if strings.TrimSpace(settings.ClientSecret) == "" {
		return errors.New("ClientSecret 不能为空")
	}

	// 验证 BaseURL
	if strings.TrimSpace(settings.BaseURL) == "" {
		return errors.New("BaseURL 不能为空")
	}

	// 验证 BaseURL 格式（必须是 HTTPS）
	if !strings.HasPrefix(settings.BaseURL, "https://") {
		return errors.New("BaseURL 必须使用 HTTPS 协议")
	}

	return nil
}

// UpdateWeChatThirdPartySettings 更新微信第三方登录配置
// 更新前会进行配置验证
func UpdateWeChatThirdPartySettings(settings *WeChatThirdPartySettings) error {
	// 验证配置
	if err := ValidateWeChatThirdPartySettings(settings); err != nil {
		return err
	}

	// 更新配置
	defaultWeChatThirdPartySettings.Enabled = settings.Enabled
	defaultWeChatThirdPartySettings.ClientKey = settings.ClientKey
	defaultWeChatThirdPartySettings.ClientSecret = settings.ClientSecret
	defaultWeChatThirdPartySettings.BaseURL = settings.BaseURL

	return nil
}

// IsWeChatThirdPartyEnabled 检查微信第三方登录是否启用
func IsWeChatThirdPartyEnabled() bool {
	return defaultWeChatThirdPartySettings.Enabled
}

// GetWeChatThirdPartyBaseURL 获取平台 API 地址
func GetWeChatThirdPartyBaseURL() string {
	if defaultWeChatThirdPartySettings.BaseURL == "" {
		return "https://www.abu117.cn"
	}
	return defaultWeChatThirdPartySettings.BaseURL
}

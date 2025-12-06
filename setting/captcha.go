package setting

import (
	"strconv"

	"github.com/QuantumNous/new-api/common"
)

// 验证码设置选项名称
const (
	OptionCaptchaEnabled           = "CaptchaEnabled"
	OptionCaptchaToleranceRange    = "CaptchaToleranceRange"
	OptionCaptchaRequireOnLogin    = "CaptchaRequireOnLogin"
	OptionCaptchaRequireOnRegister = "CaptchaRequireOnRegister"
	OptionCaptchaRequireOnCheckin  = "CaptchaRequireOnCheckin"
	OptionCaptchaMaxAttempts       = "CaptchaMaxAttempts"
	OptionCaptchaBlockDuration     = "CaptchaBlockDuration"
)

// 验证码配置变量
var (
	CaptchaEnabled           = false
	CaptchaToleranceRange    = 5
	CaptchaRequireOnLogin    = true
	CaptchaRequireOnRegister = true
	CaptchaRequireOnCheckin  = true
	CaptchaMaxAttempts       = 5
	CaptchaBlockDuration     = 300 // 秒
)

// InitCaptchaSettings 初始化验证码设置
func InitCaptchaSettings() {
	loadCaptchaSettings()
}

// loadCaptchaSettings 从 OptionMap 加载验证码设置
func loadCaptchaSettings() {
	// 验证码功能开关
	if val := getOptionString(OptionCaptchaEnabled); val != "" {
		CaptchaEnabled = val == "true"
	}

	// 容差范围
	if val, ok := getOptionInt(OptionCaptchaToleranceRange); ok && val > 0 {
		CaptchaToleranceRange = val
	}

	// 登录时需要验证
	if val := getOptionString(OptionCaptchaRequireOnLogin); val != "" {
		CaptchaRequireOnLogin = val == "true"
	}

	// 注册时需要验证
	if val := getOptionString(OptionCaptchaRequireOnRegister); val != "" {
		CaptchaRequireOnRegister = val == "true"
	}

	// 签到时需要验证
	if val := getOptionString(OptionCaptchaRequireOnCheckin); val != "" {
		CaptchaRequireOnCheckin = val == "true"
	}

	// 最大尝试次数
	if val, ok := getOptionInt(OptionCaptchaMaxAttempts); ok && val > 0 {
		CaptchaMaxAttempts = val
	}

	// 封禁时长
	if val, ok := getOptionInt(OptionCaptchaBlockDuration); ok && val > 0 {
		CaptchaBlockDuration = val
	}
}

// GetCaptchaSettings 获取验证码设置
func GetCaptchaSettings() map[string]interface{} {
	return map[string]interface{}{
		"enabled":             CaptchaEnabled,
		"tolerance_range":     CaptchaToleranceRange,
		"require_on_login":    CaptchaRequireOnLogin,
		"require_on_register": CaptchaRequireOnRegister,
		"require_on_checkin":  CaptchaRequireOnCheckin,
		"max_attempts":        CaptchaMaxAttempts,
		"block_duration":      CaptchaBlockDuration,
	}
}

// UpdateCaptchaSettings 更新验证码设置
func UpdateCaptchaSettings(settings map[string]interface{}) {
	if enabled, ok := settings["enabled"].(bool); ok {
		CaptchaEnabled = enabled
		common.OptionMapRWMutex.Lock()
		common.OptionMap[OptionCaptchaEnabled] = strconv.FormatBool(enabled)
		common.OptionMapRWMutex.Unlock()
	}

	if tolerance, ok := settings["tolerance_range"].(float64); ok {
		CaptchaToleranceRange = int(tolerance)
		common.OptionMapRWMutex.Lock()
		common.OptionMap[OptionCaptchaToleranceRange] = strconv.Itoa(int(tolerance))
		common.OptionMapRWMutex.Unlock()
	}

	if requireLogin, ok := settings["require_on_login"].(bool); ok {
		CaptchaRequireOnLogin = requireLogin
		common.OptionMapRWMutex.Lock()
		common.OptionMap[OptionCaptchaRequireOnLogin] = strconv.FormatBool(requireLogin)
		common.OptionMapRWMutex.Unlock()
	}

	if requireRegister, ok := settings["require_on_register"].(bool); ok {
		CaptchaRequireOnRegister = requireRegister
		common.OptionMapRWMutex.Lock()
		common.OptionMap[OptionCaptchaRequireOnRegister] = strconv.FormatBool(requireRegister)
		common.OptionMapRWMutex.Unlock()
	}

	if requireCheckin, ok := settings["require_on_checkin"].(bool); ok {
		CaptchaRequireOnCheckin = requireCheckin
		common.OptionMapRWMutex.Lock()
		common.OptionMap[OptionCaptchaRequireOnCheckin] = strconv.FormatBool(requireCheckin)
		common.OptionMapRWMutex.Unlock()
	}

	if maxAttempts, ok := settings["max_attempts"].(float64); ok {
		CaptchaMaxAttempts = int(maxAttempts)
		common.OptionMapRWMutex.Lock()
		common.OptionMap[OptionCaptchaMaxAttempts] = strconv.Itoa(int(maxAttempts))
		common.OptionMapRWMutex.Unlock()
	}

	if blockDuration, ok := settings["block_duration"].(float64); ok {
		CaptchaBlockDuration = int(blockDuration)
		common.OptionMapRWMutex.Lock()
		common.OptionMap[OptionCaptchaBlockDuration] = strconv.Itoa(int(blockDuration))
		common.OptionMapRWMutex.Unlock()
	}
}

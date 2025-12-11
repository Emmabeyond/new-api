package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting"
	"github.com/QuantumNous/new-api/setting/console_setting"
	"github.com/QuantumNous/new-api/setting/ratio_setting"
	"github.com/QuantumNous/new-api/setting/system_setting"

	"github.com/gin-gonic/gin"
)

func GetOptions(c *gin.Context) {
	var options []*model.Option
	common.OptionMapRWMutex.Lock()
	for k, v := range common.OptionMap {
		if strings.HasSuffix(k, "Token") || strings.HasSuffix(k, "Secret") || strings.HasSuffix(k, "Key") {
			continue
		}
		options = append(options, &model.Option{
			Key:   k,
			Value: common.Interface2String(v),
		})
	}
	common.OptionMapRWMutex.Unlock()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    options,
	})
	return
}

// BatchGetOptions 批量获取配置项
// GET /api/option/batch?keys=key1,key2,key3
func BatchGetOptions(c *gin.Context) {
	keysParam := c.Query("keys")
	if keysParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "缺少 keys 参数",
		})
		return
	}

	keys := strings.Split(keysParam, ",")
	data := make(map[string]any)

	common.OptionMapRWMutex.RLock()
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		// 过滤敏感字段
		if strings.HasSuffix(key, "Token") || strings.HasSuffix(key, "Secret") || strings.HasSuffix(key, "Key") {
			continue
		}
		if v, ok := common.OptionMap[key]; ok {
			data[key] = common.Interface2String(v)
		}
	}
	common.OptionMapRWMutex.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

type OptionUpdateRequest struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func UpdateOption(c *gin.Context) {
	var option OptionUpdateRequest
	err := json.NewDecoder(c.Request.Body).Decode(&option)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	switch option.Value.(type) {
	case bool:
		option.Value = common.Interface2String(option.Value.(bool))
	case float64:
		option.Value = common.Interface2String(option.Value.(float64))
	case int:
		option.Value = common.Interface2String(option.Value.(int))
	default:
		option.Value = fmt.Sprintf("%v", option.Value)
	}
	switch option.Key {
	case "GitHubOAuthEnabled":
		if option.Value == "true" && common.GitHubClientId == "" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用 GitHub OAuth，请先填入 GitHub Client Id 以及 GitHub Client Secret！",
			})
			return
		}
	case "discord.enabled":
		if option.Value == "true" && system_setting.GetDiscordSettings().ClientId == "" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用 Discord OAuth，请先填入 Discord Client Id 以及 Discord Client Secret！",
			})
			return
		}
	case "oidc.enabled":
		if option.Value == "true" && system_setting.GetOIDCSettings().ClientId == "" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用 OIDC 登录，请先填入 OIDC Client Id 以及 OIDC Client Secret！",
			})
			return
		}
	case "linuxdo.enabled":
		linuxdoSettings := system_setting.GetLinuxDOSettings()
		if option.Value == "true" && linuxdoSettings.ClientId == "" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用 LinuxDO OAuth，请先填入 LinuxDO Client Id 以及 LinuxDO Client Secret！",
			})
			return
		}
	case "EmailDomainRestrictionEnabled":
		if option.Value == "true" && len(common.EmailDomainWhitelist) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用邮箱域名限制，请先填入限制的邮箱域名！",
			})
			return
		}
	case "WeChatAuthEnabled":
		if option.Value == "true" && common.WeChatServerAddress == "" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用微信登录，请先填入微信登录相关配置信息！",
			})
			return
		}
	case "TurnstileCheckEnabled":
		if option.Value == "true" && common.TurnstileSiteKey == "" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用 Turnstile 校验，请先填入 Turnstile 校验相关配置信息！",
			})

			return
		}
	case "TelegramOAuthEnabled":
		if option.Value == "true" && common.TelegramBotToken == "" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用 Telegram OAuth，请先填入 Telegram Bot Token！",
			})
			return
		}
	case "GroupRatio":
		err = ratio_setting.CheckGroupRatio(option.Value.(string))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
	case "ImageRatio":
		err = ratio_setting.UpdateImageRatioByJSONString(option.Value.(string))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "图片倍率设置失败: " + err.Error(),
			})
			return
		}
	case "AudioRatio":
		err = ratio_setting.UpdateAudioRatioByJSONString(option.Value.(string))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "音频倍率设置失败: " + err.Error(),
			})
			return
		}
	case "AudioCompletionRatio":
		err = ratio_setting.UpdateAudioCompletionRatioByJSONString(option.Value.(string))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "音频补全倍率设置失败: " + err.Error(),
			})
			return
		}
	case "ModelRequestRateLimitGroup":
		err = setting.CheckModelRequestRateLimitGroup(option.Value.(string))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
	case "console_setting.api_info":
		err = console_setting.ValidateConsoleSettings(option.Value.(string), "ApiInfo")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
	case "console_setting.announcements":
		err = console_setting.ValidateConsoleSettings(option.Value.(string), "Announcements")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
	case "console_setting.faq":
		err = console_setting.ValidateConsoleSettings(option.Value.(string), "FAQ")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
	case "console_setting.uptime_kuma_groups":
		err = console_setting.ValidateConsoleSettings(option.Value.(string), "UptimeKumaGroups")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
	}
	err = model.UpdateOption(option.Key, option.Value.(string))
	if err != nil {
		common.ApiError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
	return
}

// BatchUpdateOptions 批量更新配置项
// PUT /api/option/batch
func BatchUpdateOptions(c *gin.Context) {
	var req dto.BatchOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的请求参数",
		})
		return
	}

	if len(req.Options) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "options 不能为空",
		})
		return
	}

	results := make([]dto.OptionResult, 0, len(req.Options))
	successCount := 0
	failureCount := 0

	for _, item := range req.Options {
		result := dto.OptionResult{
			Key:     item.Key,
			Success: true,
		}

		// 转换 value 为字符串
		var valueStr string
		switch v := item.Value.(type) {
		case bool:
			valueStr = common.Interface2String(v)
		case float64:
			valueStr = common.Interface2String(v)
		case int:
			valueStr = common.Interface2String(v)
		default:
			valueStr = fmt.Sprintf("%v", item.Value)
		}

		// 验证特定配置项
		validationErr := validateOption(item.Key, valueStr)
		if validationErr != nil {
			result.Success = false
			result.Error = validationErr.Error()
			failureCount++
			results = append(results, result)
			continue
		}

		// 保存配置项
		if err := model.UpdateOption(item.Key, valueStr); err != nil {
			result.Success = false
			result.Error = err.Error()
			failureCount++
		} else {
			successCount++
		}
		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": dto.BatchOptionResponse{
			Results:      results,
			SuccessCount: successCount,
			FailureCount: failureCount,
		},
	})
}

// validateOption 验证单个配置项
func validateOption(key, value string) error {
	switch key {
	case "GitHubOAuthEnabled":
		if value == "true" && common.GitHubClientId == "" {
			return fmt.Errorf("无法启用 GitHub OAuth，请先填入 GitHub Client Id 以及 GitHub Client Secret！")
		}
	case "discord.enabled":
		if value == "true" && system_setting.GetDiscordSettings().ClientId == "" {
			return fmt.Errorf("无法启用 Discord OAuth，请先填入 Discord Client Id 以及 Discord Client Secret！")
		}
	case "oidc.enabled":
		if value == "true" && system_setting.GetOIDCSettings().ClientId == "" {
			return fmt.Errorf("无法启用 OIDC 登录，请先填入 OIDC Client Id 以及 OIDC Client Secret！")
		}
	case "linuxdo.enabled":
		linuxdoSettings := system_setting.GetLinuxDOSettings()
		if value == "true" && linuxdoSettings.ClientId == "" {
			return fmt.Errorf("无法启用 LinuxDO OAuth，请先填入 LinuxDO Client Id 以及 LinuxDO Client Secret！")
		}
	case "EmailDomainRestrictionEnabled":
		if value == "true" && len(common.EmailDomainWhitelist) == 0 {
			return fmt.Errorf("无法启用邮箱域名限制，请先填入限制的邮箱域名！")
		}
	case "WeChatAuthEnabled":
		if value == "true" && common.WeChatServerAddress == "" {
			return fmt.Errorf("无法启用微信登录，请先填入微信登录相关配置信息！")
		}
	case "TurnstileCheckEnabled":
		if value == "true" && common.TurnstileSiteKey == "" {
			return fmt.Errorf("无法启用 Turnstile 校验，请先填入 Turnstile 校验相关配置信息！")
		}
	case "TelegramOAuthEnabled":
		if value == "true" && common.TelegramBotToken == "" {
			return fmt.Errorf("无法启用 Telegram OAuth，请先填入 Telegram Bot Token！")
		}
	case "GroupRatio":
		return ratio_setting.CheckGroupRatio(value)
	case "ImageRatio":
		return ratio_setting.UpdateImageRatioByJSONString(value)
	case "AudioRatio":
		return ratio_setting.UpdateAudioRatioByJSONString(value)
	case "AudioCompletionRatio":
		return ratio_setting.UpdateAudioCompletionRatioByJSONString(value)
	case "ModelRequestRateLimitGroup":
		return setting.CheckModelRequestRateLimitGroup(value)
	case "console_setting.api_info":
		return console_setting.ValidateConsoleSettings(value, "ApiInfo")
	case "console_setting.announcements":
		return console_setting.ValidateConsoleSettings(value, "Announcements")
	case "console_setting.faq":
		return console_setting.ValidateConsoleSettings(value, "FAQ")
	case "console_setting.uptime_kuma_groups":
		return console_setting.ValidateConsoleSettings(value, "UptimeKumaGroups")
	}
	return nil
}

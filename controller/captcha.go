package controller

import (
	"fmt"
	"net/http"

	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting"

	"github.com/gin-gonic/gin"
)

// ChallengeRequest 获取挑战请求
type ChallengeRequest struct{}

// VerifyRequest 验证请求
type VerifyRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	X         int    `json:"x" binding:"required"`
}

// CaptchaSettingsRequest 验证码设置请求
type CaptchaSettingsRequest struct {
	Enabled           bool `json:"enabled"`
	ToleranceRange    int  `json:"tolerance_range"`
	RequireOnLogin    bool `json:"require_on_login"`
	RequireOnRegister bool `json:"require_on_register"`
	RequireOnCheckin  bool `json:"require_on_checkin"`
	MaxAttempts       int  `json:"max_attempts"`
	BlockDuration     int  `json:"block_duration"`
}

// GetCaptchaChallenge 获取验证码挑战
// GET /api/captcha/challenge
func GetCaptchaChallenge(c *gin.Context) {
	clientIP := c.ClientIP()

	challenge, err := service.GenerateChallenge(clientIP)
	if err != nil {
		if err == service.ErrCaptchaDisabled {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "验证码功能已关闭",
			})
			return
		}
		if err == service.ErrRateLimited {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "RATE_LIMITED",
					"message": "操作过于频繁，请稍后再试",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "生成验证码失败",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    challenge,
	})
}

// VerifyCaptcha 验证滑块位置
// POST /api/captcha/verify
func VerifyCaptcha(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "无效的请求参数",
			},
		})
		return
	}

	clientIP := c.ClientIP()

	result, err := service.VerifyChallenge(req.SessionID, req.X, clientIP)
	if err != nil {
		if err == service.ErrCaptchaDisabled {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "验证码功能已关闭",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "验证失败",
			},
		})
		return
	}

	if !result.Success {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VERIFY_FAILED",
				"message": result.Message,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"token": result.Token,
		},
	})
}

// GetCaptchaSettings 获取验证码设置（管理员）
// GET /api/captcha/settings
func GetCaptchaSettings(c *gin.Context) {
	settings := setting.GetCaptchaSettings()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    settings,
	})
}

// UpdateCaptchaSettings 更新验证码设置（管理员）
// PUT /api/captcha/settings
func UpdateCaptchaSettings(c *gin.Context) {
	var req CaptchaSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的请求参数",
		})
		return
	}

	// 验证参数
	if req.ToleranceRange < 1 || req.ToleranceRange > 20 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "容差范围必须在 1-20 之间",
		})
		return
	}

	if req.MaxAttempts < 1 || req.MaxAttempts > 20 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "最大尝试次数必须在 1-20 之间",
		})
		return
	}

	if req.BlockDuration < 60 || req.BlockDuration > 3600 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "封禁时长必须在 60-3600 秒之间",
		})
		return
	}

	// 更新设置
	settingsMap := map[string]interface{}{
		"enabled":             req.Enabled,
		"tolerance_range":     float64(req.ToleranceRange),
		"require_on_login":    req.RequireOnLogin,
		"require_on_register": req.RequireOnRegister,
		"require_on_checkin":  req.RequireOnCheckin,
		"max_attempts":        float64(req.MaxAttempts),
		"block_duration":      float64(req.BlockDuration),
	}

	setting.UpdateCaptchaSettings(settingsMap)

	// 保存到数据库
	optionUpdates := map[string]string{
		setting.OptionCaptchaEnabled:           boolToString(req.Enabled),
		setting.OptionCaptchaToleranceRange:    intToString(req.ToleranceRange),
		setting.OptionCaptchaRequireOnLogin:    boolToString(req.RequireOnLogin),
		setting.OptionCaptchaRequireOnRegister: boolToString(req.RequireOnRegister),
		setting.OptionCaptchaRequireOnCheckin:  boolToString(req.RequireOnCheckin),
		setting.OptionCaptchaMaxAttempts:       intToString(req.MaxAttempts),
		setting.OptionCaptchaBlockDuration:     intToString(req.BlockDuration),
	}

	for key, value := range optionUpdates {
		if err := model.UpdateOption(key, value); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "保存设置失败: " + err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "设置已更新",
	})
}

// GetCaptchaStatus 获取验证码状态（公开接口）
// GET /api/captcha/status
func GetCaptchaStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"enabled":             setting.CaptchaEnabled,
			"require_on_login":    setting.CaptchaRequireOnLogin,
			"require_on_register": setting.CaptchaRequireOnRegister,
			"require_on_checkin":  setting.CaptchaRequireOnCheckin,
		},
	})
}

// intToString 将整数转换为字符串
func intToString(i int) string {
	return fmt.Sprintf("%d", i)
}

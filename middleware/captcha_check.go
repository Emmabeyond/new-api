package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting"

	"github.com/gin-gonic/gin"
)

// CaptchaProtectedRequest 包含验证码令牌的请求
type CaptchaProtectedRequest struct {
	CaptchaToken string `json:"captcha_token"`
}

// CaptchaRequired 验证码验证中间件
// operation: "login", "register", "checkin"
func CaptchaRequired(operation string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查验证码是否启用
		if !setting.CaptchaEnabled {
			c.Next()
			return
		}

		// 检查该操作是否需要验证码
		if !service.IsCaptchaRequiredForOperation(operation) {
			c.Next()
			return
		}

		// 尝试从多个来源获取验证码令牌
		captchaToken := getCaptchaToken(c)

		if captchaToken == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "CAPTCHA_REQUIRED",
					"message": "请完成人机验证",
				},
			})
			c.Abort()
			return
		}

		// 验证令牌
		clientIP := c.ClientIP()
		err := service.ValidateCaptchaToken(captchaToken, clientIP)
		if err != nil {
			var errorCode, errorMessage string
			switch err {
			case service.ErrTokenInvalid:
				errorCode = "TOKEN_INVALID"
				errorMessage = "验证令牌无效"
			case service.ErrTokenExpired:
				errorCode = "TOKEN_EXPIRED"
				errorMessage = "验证令牌已过期"
			case service.ErrTokenUsed:
				errorCode = "TOKEN_USED"
				errorMessage = "验证令牌已使用"
			default:
				errorCode = "CAPTCHA_ERROR"
				errorMessage = "验证失败"
			}

			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    errorCode,
					"message": errorMessage,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getCaptchaToken 从多个来源获取验证码令牌
// 优先级: Header > Query > Body
func getCaptchaToken(c *gin.Context) string {
	// 1. 尝试从 Header 获取
	token := c.GetHeader("X-Captcha-Token")
	if token != "" {
		return token
	}

	// 2. 尝试从查询参数获取
	token = c.Query("captcha_token")
	if token != "" {
		return token
	}

	// 3. 尝试从请求体获取（需要保留请求体供后续处理）
	if c.Request.Body != nil && c.Request.ContentLength > 0 {
		// 读取请求体
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return ""
		}
		// 恢复请求体供后续处理
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// 尝试解析 JSON
		var bodyMap map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &bodyMap); err == nil {
			if tokenVal, ok := bodyMap["captcha_token"]; ok {
				if tokenStr, ok := tokenVal.(string); ok {
					return tokenStr
				}
			}
		}
	}

	return ""
}

// CaptchaRequiredForLogin 登录验证码中间件
func CaptchaRequiredForLogin() gin.HandlerFunc {
	return CaptchaRequired("login")
}

// CaptchaRequiredForRegister 注册验证码中间件
func CaptchaRequiredForRegister() gin.HandlerFunc {
	return CaptchaRequired("register")
}

// CaptchaRequiredForCheckin 签到验证码中间件
func CaptchaRequiredForCheckin() gin.HandlerFunc {
	return CaptchaRequired("checkin")
}

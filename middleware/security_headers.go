package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityHeadersConfig 安全响应头配置
type SecurityHeadersConfig struct {
	EnableCSP                bool   `json:"enable_csp"`
	CSPPolicy                string `json:"csp_policy"`
	EnableXSSProtection      bool   `json:"enable_xss_protection"`
	EnableFrameOptions       bool   `json:"enable_frame_options"`
	EnableHSTS               bool   `json:"enable_hsts"`
	HSTSMaxAge               int    `json:"hsts_max_age"`
	EnableContentTypeNosniff bool   `json:"enable_content_type_nosniff"`
}

// DefaultCSPPolicy 默认 CSP 策略
// 注意：'unsafe-inline' 和 'unsafe-eval' 是为了兼容现有前端代码
// 生产环境建议使用更严格的策略
const DefaultCSPPolicy = "default-src 'self'; " +
	"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
	"style-src 'self' 'unsafe-inline'; " +
	"img-src 'self' data: https:; " +
	"font-src 'self' data:; " +
	"connect-src 'self'; " +
	"frame-src 'self' https:; " +
	"frame-ancestors 'self'; " +
	"base-uri 'self'; " +
	"form-action 'self'"

// DefaultSecurityHeadersConfig 返回默认配置
func DefaultSecurityHeadersConfig() *SecurityHeadersConfig {
	config := &SecurityHeadersConfig{
		EnableCSP:                true,
		CSPPolicy:                DefaultCSPPolicy,
		EnableXSSProtection:      true,
		EnableFrameOptions:       true,
		EnableHSTS:               true,
		HSTSMaxAge:               31536000, // 1 year in seconds
		EnableContentTypeNosniff: true,
	}

	// 从环境变量读取 CSP 策略
	if cspPolicy := os.Getenv("CSP_POLICY"); cspPolicy != "" {
		config.CSPPolicy = cspPolicy
	}

	return config
}

// SecurityHeaders 返回安全响应头中间件
// 该中间件添加以下安全响应头：
// - Content-Security-Policy: 限制资源加载来源
// - X-XSS-Protection: 启用浏览器 XSS 过滤器
// - X-Content-Type-Options: 防止 MIME 类型嗅探
// - X-Frame-Options: 防止点击劫持
// - Strict-Transport-Security: 强制 HTTPS（仅 HTTPS 连接时添加）
func SecurityHeaders(config *SecurityHeadersConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultSecurityHeadersConfig()
	}

	return func(c *gin.Context) {
		// Content-Security-Policy (Requirements 2.1)
		if config.EnableCSP && config.CSPPolicy != "" {
			c.Header("Content-Security-Policy", config.CSPPolicy)
		}

		// X-XSS-Protection (Requirements 2.2)
		if config.EnableXSSProtection {
			c.Header("X-XSS-Protection", "1; mode=block")
		}

		// X-Content-Type-Options (Requirements 2.3)
		if config.EnableContentTypeNosniff {
			c.Header("X-Content-Type-Options", "nosniff")
		}

		// X-Frame-Options (Requirements 2.4)
		if config.EnableFrameOptions {
			c.Header("X-Frame-Options", "SAMEORIGIN")
		}

		// Strict-Transport-Security (Requirements 2.5)
		// 仅在 HTTPS 连接时添加
		if config.EnableHSTS {
			isHTTPS := isHTTPSRequest(c)
			if isHTTPS {
				hstsValue := fmt.Sprintf("max-age=%d; includeSubDomains", config.HSTSMaxAge)
				c.Header("Strict-Transport-Security", hstsValue)
			}
		}

		c.Next()
	}
}

// isHTTPSRequest 检查请求是否为 HTTPS
// 支持直接 TLS 连接和反向代理场景
func isHTTPSRequest(c *gin.Context) bool {
	// 直接 TLS 连接
	if c.Request.TLS != nil {
		return true
	}

	// 反向代理场景：检查 X-Forwarded-Proto 头
	proto := c.GetHeader("X-Forwarded-Proto")
	if strings.ToLower(proto) == "https" {
		return true
	}

	// 检查 X-Forwarded-Ssl 头（某些代理使用）
	ssl := c.GetHeader("X-Forwarded-Ssl")
	if strings.ToLower(ssl) == "on" {
		return true
	}

	return false
}

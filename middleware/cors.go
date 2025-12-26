package middleware

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-gonic/gin"
)

// CORSConfig CORS 配置结构体
type CORSConfig struct {
	AllowedOrigins   []string `json:"allowed_origins"`
	AllowCredentials bool     `json:"allow_credentials"`
	AllowedMethods   []string `json:"allowed_methods"`
	AllowedHeaders   []string `json:"allowed_headers"`
	MaxAge           int      `json:"max_age"`
}

// DefaultCORSConfig 返回默认的 CORS 配置（仅允许本地开发）
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		MaxAge:           3600,
	}
}

// LoadCORSConfigFromEnv 从环境变量加载 CORS 配置
func LoadCORSConfigFromEnv() *CORSConfig {
	config := DefaultCORSConfig()

	// 加载 ALLOWED_ORIGINS
	originsStr := os.Getenv("ALLOWED_ORIGINS")
	if originsStr == "" {
		common.SysLog("ALLOWED_ORIGINS not set, using default (localhost only)")
	} else {
		origins := strings.Split(originsStr, ",")
		validOrigins := []string{}

		for _, origin := range origins {
			origin = strings.TrimSpace(origin)
			if origin == "" {
				continue
			}

			// 验证 URL 格式
			if isValidOrigin(origin) {
				validOrigins = append(validOrigins, origin)
			} else {
				common.SysError(fmt.Sprintf("Invalid origin in ALLOWED_ORIGINS: %s", origin))
			}
		}

		if len(validOrigins) > 0 {
			config.AllowedOrigins = validOrigins
			common.SysLog(fmt.Sprintf("CORS allowed origins: %v", validOrigins))
		} else {
			common.SysError("No valid origins found in ALLOWED_ORIGINS, using default")
		}
	}

	// 加载 CORS_MAX_AGE
	maxAgeStr := os.Getenv("CORS_MAX_AGE")
	if maxAgeStr != "" {
		maxAge, err := strconv.Atoi(maxAgeStr)
		if err != nil {
			common.SysError(fmt.Sprintf("Invalid CORS_MAX_AGE value: %s, using default: %d", maxAgeStr, config.MaxAge))
		} else if maxAge > 0 {
			config.MaxAge = maxAge
		}
	}

	// 加载 CORS_ALLOW_CREDENTIALS
	allowCredStr := os.Getenv("CORS_ALLOW_CREDENTIALS")
	if allowCredStr != "" {
		allowCred, err := strconv.ParseBool(allowCredStr)
		if err != nil {
			common.SysError(fmt.Sprintf("Invalid CORS_ALLOW_CREDENTIALS value: %s", allowCredStr))
		} else {
			config.AllowCredentials = allowCred
		}
	}

	// 加载 CORS_ALLOWED_HEADERS
	headersStr := os.Getenv("CORS_ALLOWED_HEADERS")
	if headersStr != "" {
		headers := strings.Split(headersStr, ",")
		validHeaders := []string{}
		for _, header := range headers {
			header = strings.TrimSpace(header)
			if header != "" {
				validHeaders = append(validHeaders, header)
			}
		}
		if len(validHeaders) > 0 {
			config.AllowedHeaders = validHeaders
		}
	}

	return config
}

// isValidOrigin 验证来源 URL 格式是否有效
func isValidOrigin(origin string) bool {
	// 允许通配符
	if origin == "*" {
		return true
	}

	// 解析 URL
	parsedURL, err := url.Parse(origin)
	if err != nil {
		return false
	}

	// 必须有 scheme（http 或 https）
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	// 必须有 host
	if parsedURL.Host == "" {
		return false
	}

	return true
}

// isOriginAllowed 检查来源是否在白名单中
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}

	for _, allowed := range allowedOrigins {
		// 支持通配符
		if allowed == "*" {
			return true
		}

		// 精确匹配
		if origin == allowed {
			return true
		}

		// 支持子域名通配符（如 *.example.com）
		if strings.HasPrefix(allowed, "*.") {
			domain := strings.TrimPrefix(allowed, "*")
			if strings.HasSuffix(origin, domain) {
				// 验证是否是有效的子域名匹配
				parsedOrigin, err := url.Parse(origin)
				if err == nil && strings.HasSuffix(parsedOrigin.Host, strings.TrimPrefix(domain, ".")) {
					return true
				}
			}
		}
	}

	return false
}

// isHeaderAllowed 检查请求头是否在白名单中
func isHeaderAllowed(header string, allowedHeaders []string) bool {
	header = strings.ToLower(strings.TrimSpace(header))
	if header == "" {
		return true
	}

	for _, allowed := range allowedHeaders {
		if strings.ToLower(allowed) == header {
			return true
		}
	}

	return false
}

// CORS 返回 CORS 中间件（向后兼容：无参数时从环境变量加载配置）
func CORS(configs ...*CORSConfig) gin.HandlerFunc {
	var config *CORSConfig
	if len(configs) > 0 && configs[0] != nil {
		config = configs[0]
	} else {
		config = LoadCORSConfigFromEnv()
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// 如果没有 Origin 头，说明不是跨域请求
		if origin == "" {
			c.Next()
			return
		}

		// 验证来源是否在白名单中
		if !isOriginAllowed(origin, config.AllowedOrigins) {
			// 记录被拒绝的跨域请求
			common.SysLog(fmt.Sprintf("CORS: Rejected request from origin: %s, path: %s", origin, c.Request.URL.Path))

			// 对于预检请求，返回 403
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(403)
				return
			}

			// 对于其他请求，不设置 CORS 头，让浏览器阻止
			c.Next()
			return
		}

		// 设置 CORS 响应头
		c.Header("Access-Control-Allow-Origin", origin)

		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			// 验证请求头
			requestHeaders := c.GetHeader("Access-Control-Request-Headers")
			if requestHeaders != "" {
				headers := strings.Split(requestHeaders, ",")
				for _, header := range headers {
					if !isHeaderAllowed(header, config.AllowedHeaders) {
						common.SysLog(fmt.Sprintf("CORS: Rejected header: %s from origin: %s", header, origin))
						c.AbortWithStatus(403)
						return
					}
				}
			}

			c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
			c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
			c.Header("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))

			c.AbortWithStatus(204)
			return
		}

		// 对于非预检请求，也设置暴露的响应头
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")

		c.Next()
	}
}

// CORSWithConfig 使用指定配置的 CORS 中间件
func CORSWithConfig(config *CORSConfig) gin.HandlerFunc {
	return CORS(config)
}

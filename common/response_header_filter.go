package common

import (
	"net/http"
	"strings"
	"sync"
)

// DefaultHeaderBlacklist 默认需要过滤的响应头模式
// 这些头可能泄露上游服务商或代理层级信息
var DefaultHeaderBlacklist = []string{
	"X-Request-Id",
	"X-Ratelimit-*",
	"X-RateLimit-*",
	"CF-*",
	"Via",
	"Server",
	"X-Powered-By",
	"X-OpenAI-*",
	"X-Claude-*",
	"Anthropic-*",
	"X-Amz-*",
	"X-Goog-*",
	"X-Cloud-*",
	"X-Azure-*",
	"X-Ms-*",
	"Alt-Svc",
	"Report-To",
	"NEL",
}

// EssentialHeaders 必须保留的响应头（不受过滤影响）
var EssentialHeaders = []string{
	"Content-Type",
	"Content-Length",
	"Transfer-Encoding",
	"Cache-Control",
	"Content-Encoding",
	"Content-Disposition",
	"Accept-Ranges",
}

// HeaderFilterMode 过滤模式
type HeaderFilterMode string

const (
	HeaderFilterModeBlacklist HeaderFilterMode = "blacklist"
	HeaderFilterModeWhitelist HeaderFilterMode = "whitelist"
)

// HeaderFilterConfig 响应头过滤配置
type HeaderFilterConfig struct {
	// 是否启用响应头过滤
	EnableHeaderFilter bool `json:"enable_header_filter"`
	// 过滤模式: "blacklist" 或 "whitelist"
	FilterMode HeaderFilterMode `json:"filter_mode"`
	// 黑名单模式下要过滤的头模式
	HeaderBlacklist []string `json:"header_blacklist"`
	// 白名单模式下允许的头
	HeaderWhitelist []string `json:"header_whitelist"`
	// 是否标准化请求ID
	StandardizeRequestId bool `json:"standardize_request_id"`
	// 是否记录过滤日志
	LogFilteredHeaders bool `json:"log_filtered_headers"`
}

// DefaultHeaderFilterConfig 默认配置
var DefaultHeaderFilterConfig = HeaderFilterConfig{
	EnableHeaderFilter:   true,
	FilterMode:           HeaderFilterModeBlacklist,
	HeaderBlacklist:      DefaultHeaderBlacklist,
	HeaderWhitelist:      EssentialHeaders,
	StandardizeRequestId: true,
	LogFilteredHeaders:   false,
}

// ResponseHeaderFilter 响应头过滤器
type ResponseHeaderFilter struct {
	config  *HeaderFilterConfig
	matcher *HeaderPatternMatcher
	mu      sync.RWMutex
}

// headerFilterConfigProvider 配置提供者函数
var headerFilterConfigProvider func() *HeaderFilterConfig

// SetHeaderFilterConfigProvider 设置配置提供者
func SetHeaderFilterConfigProvider(provider func() *HeaderFilterConfig) {
	headerFilterConfigProvider = provider
}

// getHeaderFilterConfig 获取当前配置
func getHeaderFilterConfig() *HeaderFilterConfig {
	if headerFilterConfigProvider != nil {
		return headerFilterConfigProvider()
	}
	return &DefaultHeaderFilterConfig
}

var (
	globalHeaderFilter     *ResponseHeaderFilter
	headerFilterOnce       sync.Once
)

// GetResponseHeaderFilter 获取全局响应头过滤器实例
func GetResponseHeaderFilter() *ResponseHeaderFilter {
	headerFilterOnce.Do(func() {
		globalHeaderFilter = &ResponseHeaderFilter{
			config:  getHeaderFilterConfig(),
			matcher: GetHeaderPatternMatcher(),
		}
	})
	return globalHeaderFilter
}

// FilterHeaders 过滤响应头
// 返回过滤后的响应头副本
func (f *ResponseHeaderFilter) FilterHeaders(headers http.Header) http.Header {
	config := getHeaderFilterConfig()
	if !config.EnableHeaderFilter {
		return headers.Clone()
	}

	result := make(http.Header)

	for name, values := range headers {
		if f.shouldKeep(name, config) {
			result[name] = values
		} else if config.LogFilteredHeaders && DebugEnabled {
			SysLog("Filtered response header: " + name)
		}
	}

	return result
}

// shouldKeep 判断头是否应该保留
func (f *ResponseHeaderFilter) shouldKeep(headerName string, config *HeaderFilterConfig) bool {
	// 必要头始终保留
	if f.isEssentialHeader(headerName) {
		return true
	}

	switch config.FilterMode {
	case HeaderFilterModeWhitelist:
		// 白名单模式：只保留白名单中的头
		return f.matcher.MatchAny(headerName, config.HeaderWhitelist)
	case HeaderFilterModeBlacklist:
		fallthrough
	default:
		// 黑名单模式：过滤黑名单中的头
		blacklist := config.HeaderBlacklist
		if len(blacklist) == 0 {
			blacklist = DefaultHeaderBlacklist
		}
		return !f.matcher.MatchAny(headerName, blacklist)
	}
}

// isEssentialHeader 检查是否为必要头
func (f *ResponseHeaderFilter) isEssentialHeader(headerName string) bool {
	headerLower := strings.ToLower(headerName)
	for _, essential := range EssentialHeaders {
		if strings.ToLower(essential) == headerLower {
			return true
		}
	}
	return false
}

// ShouldFilter 判断单个头是否应该被过滤
func (f *ResponseHeaderFilter) ShouldFilter(headerName string) bool {
	config := getHeaderFilterConfig()
	if !config.EnableHeaderFilter {
		return false
	}
	return !f.shouldKeep(headerName, config)
}

// GetStandardizedRequestId 获取标准化的请求ID
// 如果启用了请求ID标准化，返回本地生成的ID
// 否则返回原始ID
func (f *ResponseHeaderFilter) GetStandardizedRequestId(localId string, upstreamId string) string {
	config := getHeaderFilterConfig()
	if config.StandardizeRequestId {
		return localId
	}
	if upstreamId != "" {
		return upstreamId
	}
	return localId
}

// FilterAndCopyHeaders 过滤并复制响应头到目标
func FilterAndCopyHeaders(dst http.Header, src http.Header, localRequestId string) {
	filter := GetResponseHeaderFilter()
	config := getHeaderFilterConfig()

	for name, values := range src {
		// 跳过请求ID头，稍后统一处理
		if strings.EqualFold(name, RequestIdKey) || strings.EqualFold(name, "X-Request-Id") {
			continue
		}

		if !filter.ShouldFilter(name) {
			for _, v := range values {
				dst.Add(name, v)
			}
		}
	}

	// 设置标准化的请求ID
	if config.StandardizeRequestId && localRequestId != "" {
		dst.Set(RequestIdKey, localRequestId)
	}
}

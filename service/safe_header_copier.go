package service

import (
	"net/http"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-gonic/gin"
)

// SafeCopyResponseHeaders 安全复制响应头到客户端响应
// 自动过滤敏感头信息，防止泄露上游服务商或代理层级信息
func SafeCopyResponseHeaders(c *gin.Context, resp *http.Response) {
	if resp == nil {
		return
	}

	localRequestId := c.GetString(common.RequestIdKey)
	common.FilterAndCopyHeaders(c.Writer.Header(), resp.Header, localRequestId)
}

// SafeCopyResponseHeadersWithStatus 安全复制响应头并设置状态码
func SafeCopyResponseHeadersWithStatus(c *gin.Context, resp *http.Response) {
	SafeCopyResponseHeaders(c, resp)
	c.Writer.WriteHeader(resp.StatusCode)
}

// SafeSetResponseHeader 安全设置单个响应头
// 如果头名称在黑名单中，则不设置
func SafeSetResponseHeader(c *gin.Context, name string, value string) {
	filter := common.GetResponseHeaderFilter()
	if !filter.ShouldFilter(name) {
		c.Writer.Header().Set(name, value)
	}
}

// SafeAddResponseHeader 安全添加响应头
// 如果头名称在黑名单中，则不添加
func SafeAddResponseHeader(c *gin.Context, name string, value string) {
	filter := common.GetResponseHeaderFilter()
	if !filter.ShouldFilter(name) {
		c.Writer.Header().Add(name, value)
	}
}

// CopyFilteredHeaders 复制过滤后的响应头（不使用gin.Context）
// 用于需要直接操作http.Header的场景
func CopyFilteredHeaders(dst http.Header, src http.Header, localRequestId string) {
	common.FilterAndCopyHeaders(dst, src, localRequestId)
}

package middleware

import (
	"context"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-gonic/gin"
)

// RequestId 生成标准化的请求ID中间件
// 使用UUID格式，不包含时间戳信息，防止通过ID格式探测代理层级
func RequestId() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 使用UUID生成标准化请求ID，不暴露部署信息
		id := common.GenerateStandardRequestId()
		c.Set(common.RequestIdKey, id)
		ctx := context.WithValue(c.Request.Context(), common.RequestIdKey, id)
		c.Request = c.Request.WithContext(ctx)
		c.Header(common.RequestIdKey, id)
		c.Next()
	}
}

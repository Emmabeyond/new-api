package middleware

import (
	"fmt"
	"net/http"

	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/security_setting"
	"github.com/gin-gonic/gin"
)

// AntiAbuseCheck middleware checks for abusive behavior patterns
// It should be placed after TokenAuth middleware
func AntiAbuseCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if anti-abuse is enabled
		if !security_setting.IsAntiAbuseEnabled() {
			c.Next()
			return
		}

		// Get token info from context (set by TokenAuth middleware)
		tokenID, exists := c.Get("token_id")
		if !exists {
			// No token info, skip abuse check
			c.Next()
			return
		}

		userID, _ := c.Get("id")
		group, _ := c.Get("token_group")

		// Get model name from request (will be set by distributor or parsed from request)
		modelName := c.GetString("model")
		if modelName == "" {
			// Try to get from request body later, for now skip
			c.Next()
			return
		}

		// Get content from request (simplified - in production would parse request body)
		content := "" // Content extraction would be done in the relay handler

		// Perform abuse check
		detector := service.GetAbuseDetector()
		result, err := detector.CheckRequest(
			tokenID.(int),
			userID.(int),
			groupToString(group),
			modelName,
			content,
		)

		if err != nil {
			// On error, log and allow the request
			c.Next()
			return
		}

		if !result.Allowed {
			// Request is blocked due to abuse
			c.Header("Retry-After", fmt.Sprintf("%d", result.RetryAfter))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"message": fmt.Sprintf("Request blocked: %s. Retry after %d seconds.", result.Reason, result.RetryAfter),
					"type":    "abuse_detection",
					"code":    result.PenaltyType,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// groupToString safely converts group interface to string
func groupToString(group interface{}) string {
	if group == nil {
		return ""
	}
	if s, ok := group.(string); ok {
		return s
	}
	return ""
}

// RecordRequestForAbuse records a request for abuse detection metrics
// This should be called from relay handlers where we have access to the full request
func RecordRequestForAbuse(c *gin.Context, modelName string, content string) {
	if !security_setting.IsAntiAbuseEnabled() {
		return
	}

	tokenID, exists := c.Get("token_id")
	if !exists {
		return
	}

	detector := service.GetAbuseDetector()
	_ = detector.RecordRequest(tokenID.(int), modelName, content)
}

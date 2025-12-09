package middleware

import (
	"fmt"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"
	"github.com/gin-gonic/gin"
)

func abortWithOpenAiMessage(c *gin.Context, statusCode int, message string, code ...string) {
	codeStr := ""
	if len(code) > 0 {
		codeStr = code[0]
	}
	userId := c.GetInt("id")

	// Apply channel masking to remove sensitive channel information from user-facing message
	maskedMessage := common.MaskChannelInfo(message)

	c.JSON(statusCode, gin.H{
		"error": gin.H{
			"message": common.MessageWithRequestId(maskedMessage, c.GetString(common.RequestIdKey)),
			"type":    "new_api_error",
			"code":    codeStr,
		},
	})
	c.Abort()

	// Log original message with full details for admin debugging
	logger.LogError(c.Request.Context(), fmt.Sprintf("user %d | %s", userId, message))
}

func abortWithMidjourneyMessage(c *gin.Context, statusCode int, code int, description string) {
	// Apply channel masking to remove sensitive channel information from user-facing message
	maskedDescription := common.MaskChannelInfo(description)

	c.JSON(statusCode, gin.H{
		"description": maskedDescription,
		"type":        "new_api_error",
		"code":        code,
	})
	c.Abort()

	// Log original message with full details for admin debugging
	logger.LogError(c.Request.Context(), description)
}

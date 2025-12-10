package controller

import (
	"net/http"

	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/config"
	"github.com/QuantumNous/new-api/setting/security_setting"
	"github.com/gin-gonic/gin"
)

// GetSecuritySettings returns the current security settings
func GetSecuritySettings(c *gin.Context) {
	settings := security_setting.GetSecuritySetting()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    settings,
	})
}

// UpdateSecuritySettings updates the security settings
func UpdateSecuritySettings(c *gin.Context) {
	var newSettings security_setting.SecuritySetting
	if err := c.ShouldBindJSON(&newSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的参数: " + err.Error(),
		})
		return
	}

	// Validate settings
	if newSettings.ModelSwitchWindowMinutes < 1 {
		newSettings.ModelSwitchWindowMinutes = 5
	}
	if newSettings.ModelSwitchThreshold < 1 {
		newSettings.ModelSwitchThreshold = 10
	}
	if newSettings.MinContentLength < 1 {
		newSettings.MinContentLength = 10
	}
	if newSettings.TestContentThreshold < 1 {
		newSettings.TestContentThreshold = 20
	}
	if newSettings.TempBanDurationMinutes < 1 {
		newSettings.TempBanDurationMinutes = 30
	}
	if newSettings.RateLimitRequests < 1 {
		newSettings.RateLimitRequests = 5
	}

	// Convert to map for saving
	configMap, err := config.ConfigToMap(&newSettings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "配置转换失败: " + err.Error(),
		})
		return
	}

	// Save to database
	for key, value := range configMap {
		dbKey := "security_setting." + key
		err := model.UpdateOption(dbKey, value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "保存配置失败: " + err.Error(),
			})
			return
		}
	}

	// Update in-memory settings
	currentSettings := security_setting.GetSecuritySetting()
	*currentSettings = newSettings

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "安全设置已更新",
	})
}


// GetActivePenalties returns all active penalties
func GetActivePenalties(c *gin.Context) {
	page := parsePageNum(c)
	pageSize := parsePageSize(c)

	penalties, total, err := model.GetActivePenalties(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取处罚列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"penalties": penalties,
			"total":     total,
		},
	})
}

// LiftPenalty lifts a penalty for a token
func LiftPenalty(c *gin.Context) {
	tokenIdStr := c.Param("token_id")
	var tokenId int
	if _, err := parseIntParam(tokenIdStr, &tokenId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的 token_id",
		})
		return
	}

	// Get current user ID for audit
	userId := c.GetInt("id")

	// Lift penalty in database
	if err := model.LiftTokenPenalty(tokenId, userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "解除处罚失败: " + err.Error(),
		})
		return
	}

	// Lift penalty in cache (Redis/memory)
	penaltyManager := service.GetPenaltyManager()
	if err := penaltyManager.LiftPenalty(tokenId); err != nil {
		// Log error but don't fail the request
		// The database record is already updated
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "处罚已解除",
	})
}

// parseIntParam parses an integer parameter
func parseIntParam(s string, result *int) (bool, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return false, nil
		}
		n = n*10 + int(c-'0')
	}
	*result = n
	return true, nil
}

// parsePageNum parses page number from query
func parsePageNum(c *gin.Context) int {
	pageStr := c.DefaultQuery("page", "1")
	var page int
	parseIntParam(pageStr, &page)
	if page < 1 {
		page = 1
	}
	return page
}

// parsePageSize parses page size from query
func parsePageSize(c *gin.Context) int {
	sizeStr := c.DefaultQuery("page_size", "10")
	var size int
	parseIntParam(sizeStr, &size)
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return size
}


// GetTokenAbuseInfo returns abuse information for a specific token
func GetTokenAbuseInfo(c *gin.Context) {
	tokenIdStr := c.Param("token_id")
	var tokenId int
	if _, err := parseIntParam(tokenIdStr, &tokenId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的 token_id",
		})
		return
	}

	// Get abuse score
	detector := service.GetAbuseDetector()
	scoreResult, err := detector.GetAbuseScore(tokenId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取滥用评分失败: " + err.Error(),
		})
		return
	}

	// Get penalty status
	penaltyManager := service.GetPenaltyManager()
	penaltyStatus, err := penaltyManager.CheckPenalty(tokenId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取处罚状态失败: " + err.Error(),
		})
		return
	}

	// Get penalty history
	history, total, err := model.GetPenaltyHistory(tokenId, 1, 10)
	if err != nil {
		// Non-fatal error, continue with empty history
		history = []*model.TokenPenalty{}
		total = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"abuse_score":    scoreResult,
			"penalty_status": penaltyStatus,
			"penalty_history": gin.H{
				"items": history,
				"total": total,
			},
		},
	})
}

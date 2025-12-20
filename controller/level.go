package controller

import (
	"net/http"
	"strconv"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"

	"github.com/gin-gonic/gin"
)

// GetAllLevels 获取所有等级配置（管理员）
func GetAllLevels(c *gin.Context) {
	levels, err := service.GetAllLevelConfigs()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    levels,
	})
}

// CreateLevel 创建等级配置（管理员）
func CreateLevel(c *gin.Context) {
	var level model.LevelConfig
	if err := c.ShouldBindJSON(&level); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数: " + err.Error(),
		})
		return
	}

	if err := service.CreateLevelConfig(&level); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建成功",
	})
}

// UpdateLevel 更新等级配置（管理员）
func UpdateLevel(c *gin.Context) {
	levelId := c.Param("id")
	if levelId == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "等级ID不能为空",
		})
		return
	}

	var level model.LevelConfig
	if err := c.ShouldBindJSON(&level); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数: " + err.Error(),
		})
		return
	}
	level.Id = levelId

	if err := service.UpdateLevelConfig(&level); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
	})
}

// DeleteLevel 删除等级配置（管理员）
func DeleteLevel(c *gin.Context) {
	levelId := c.Param("id")
	if levelId == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "等级ID不能为空",
		})
		return
	}

	if err := service.DeleteLevelConfig(levelId); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}

// GetLevelStats 获取等级用户统计（管理员）
func GetLevelStats(c *gin.Context) {
	stats, err := service.GetLevelUserStats()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}


// GetUserLevel 获取当前用户等级信息
func GetUserLevel(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户未登录",
		})
		return
	}

	info, err := service.GetUserLevelInfo(userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    info,
	})
}

// GetUserLevelProgress 获取用户升级进度
func GetUserLevelProgress(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户未登录",
		})
		return
	}

	progress, err := service.GetUserUpgradeProgress(userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    progress,
	})
}

// GetUserLevelBenefits 获取用户等级权益详情
func GetUserLevelBenefits(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户未登录",
		})
		return
	}

	user, err := model.GetUserById(userId, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	level, err := service.GetLevelConfigById(user.Level)
	if err != nil {
		level, err = service.GetDefaultLevelConfig()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
	}

	benefits, err := level.GetBenefits()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"level":    level,
			"benefits": benefits,
		},
	})
}

// GetLevelComparison 获取等级对比信息
func GetLevelComparison(c *gin.Context) {
	levels, err := service.GetAllLevelConfigs()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 构建对比数据
	type LevelCompare struct {
		*model.LevelConfig
		Benefits *model.LevelBenefits `json:"benefits_parsed"`
	}

	var comparisons []LevelCompare
	for _, level := range levels {
		benefits, _ := level.GetBenefits()
		comparisons = append(comparisons, LevelCompare{
			LevelConfig: level,
			Benefits:    benefits,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    comparisons,
	})
}

// SetUserLevel 管理员调整用户等级
func SetUserLevel(c *gin.Context) {
	userIdStr := c.Param("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil || userId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的用户ID",
		})
		return
	}

	var req struct {
		LevelId string `json:"level_id"`
		Reason  string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}

	operatorId := c.GetInt("id")
	if err := service.ManualSetUserLevel(userId, req.LevelId, req.Reason, operatorId); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "等级调整成功",
	})
}

// GetUserLevelHistory 获取用户等级变更历史
func GetUserLevelHistory(c *gin.Context) {
	userIdStr := c.Param("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil || userId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的用户ID",
		})
		return
	}

	pageInfo := common.GetPageQuery(c)
	logs, total, err := model.GetUserLevelChangeLogs(userId, pageInfo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    logs,
		"total":   total,
	})
}

// SyncUserRecharge 同步用户累计充值金额
func SyncUserRecharge(c *gin.Context) {
	userIdStr := c.Param("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil || userId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的用户ID",
		})
		return
	}

	totalMoney, err := service.SyncUserCumulativeRecharge(userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 同步后检查是否需要升级
	result, err := service.CheckAndUpgradeUserLevel(userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "同步成功但升级检查失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "同步成功",
		"data": gin.H{
			"cumulative_recharge": totalMoney,
			"level_changed":       result.Changed,
			"new_level":           result.NewLevel,
		},
	})
}

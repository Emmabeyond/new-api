package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting"

	"github.com/gin-gonic/gin"
)

// GetAllCheckinRecords 获取所有签到记录（管理员）
// GET /api/admin/checkin/records
func GetAllCheckinRecords(c *gin.Context) {
	pageInfo := common.GetPageQuery(c)

	// 解析筛选参数
	userIdStr := c.Query("user_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	checkinTypeStr := c.Query("checkin_type")

	var userId, checkinType int
	var startDate, endDate *time.Time

	if userIdStr != "" {
		userId, _ = strconv.Atoi(userIdStr)
	}
	if checkinTypeStr != "" {
		checkinType, _ = strconv.Atoi(checkinTypeStr)
	}
	if startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endOfDay := t.Add(24 * time.Hour)
			endDate = &endOfDay
		}
	}

	records, total, err := model.GetAllCheckinRecords(pageInfo, userId, startDate, endDate, checkinType)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(records)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    pageInfo,
	})
}

// GetUserCheckinInfo 获取指定用户签到信息（管理员）
// GET /api/admin/checkin/user/:id
func GetUserCheckinInfo(c *gin.Context) {
	userIdStr := c.Param("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的用户ID",
		})
		return
	}

	// 获取用户签到状态
	checkin, err := model.GetCheckinByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 获取最近签到记录
	pageInfo := &common.PageInfo{
		Page:     1,
		PageSize: 10,
	}
	records, _, err := model.GetCheckinHistory(userId, pageInfo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"checkin":        checkin,
			"recent_records": records,
		},
	})
}

// GetCheckinDashboard 获取签到统计仪表盘（管理员）
// GET /api/admin/checkin/dashboard
func GetCheckinDashboard(c *gin.Context) {
	todayCheckins, activeUsers, avgConsecutive, err := model.GetCheckinStats()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	todayQuota, err := model.GetTodayQuotaDistributed()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	dashboard := dto.AdminCheckinDashboard{
		TodayCheckins:         todayCheckins,
		TodayQuotaDistributed: todayQuota,
		ActiveUsers:           activeUsers,
		AvgConsecutiveDays:    avgConsecutive,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    dashboard,
	})
}

// AdjustUserConsecutiveDays 调整用户连续签到天数（管理员）
// PUT /api/admin/checkin/user/:id/consecutive
func AdjustUserConsecutiveDays(c *gin.Context) {
	adminId := c.GetInt("id")
	userIdStr := c.Param("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的用户ID",
		})
		return
	}

	var req dto.AdminAdjustConsecutiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}

	// 获取当前签到状态
	checkin, err := model.GetCheckinByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if checkin == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户没有签到记录",
		})
		return
	}

	oldValue := checkin.ConsecutiveDays

	// 更新连续天数
	err = model.UpdateConsecutiveDays(userId, req.ConsecutiveDays)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 记录审计日志
	audit := &model.CheckinAudit{
		AdminId:   adminId,
		UserId:    userId,
		Action:    model.CheckinAuditActionAdjustConsecutive,
		OldValue:  fmt.Sprintf("%d", oldValue),
		NewValue:  fmt.Sprintf("%d", req.ConsecutiveDays),
		CreatedAt: time.Now(),
	}
	_ = model.CreateCheckinAudit(audit)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "调整成功",
	})
}

// GetCheckinSettings 获取签到设置（管理员）
// GET /api/admin/checkin/settings
func GetCheckinSettingsAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    setting.GetCheckinSettings(),
	})
}


// UpdateCheckinSettingsRequest 更新签到设置请求
type UpdateCheckinSettingsRequest struct {
	Enabled         *bool `json:"enabled"`
	RewardDay1To6   *int  `json:"reward_day_1_6"`
	RewardDay7      *int  `json:"reward_day_7"`
	RewardDay8To13  *int  `json:"reward_day_8_13"`
	RewardDay14To29 *int  `json:"reward_day_14_29"`
	RewardDay30     *int  `json:"reward_day_30"`
	RewardDay31Plus *int  `json:"reward_day_31_plus"`
	BonusMin        *int  `json:"bonus_min"`
	BonusMax        *int  `json:"bonus_max"`
	MakeupCost      *int  `json:"makeup_cost"`
}

// UpdateCheckinSettings 更新签到设置（管理员）
// PUT /api/admin/checkin/settings
func UpdateCheckinSettingsAPI(c *gin.Context) {
	var req UpdateCheckinSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}

	// 更新各项设置
	if req.Enabled != nil {
		value := "false"
		if *req.Enabled {
			value = "true"
		}
		if err := model.UpdateOption(setting.OptionCheckinEnabled, value); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
		setting.CheckinEnabled = *req.Enabled
	}

	if req.RewardDay1To6 != nil && *req.RewardDay1To6 > 0 {
		if err := model.UpdateOption(setting.OptionCheckinRewardDay1To6, strconv.Itoa(*req.RewardDay1To6)); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
		setting.CheckinRewardDay1To6 = *req.RewardDay1To6
	}

	if req.RewardDay7 != nil && *req.RewardDay7 > 0 {
		if err := model.UpdateOption(setting.OptionCheckinRewardDay7, strconv.Itoa(*req.RewardDay7)); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
		setting.CheckinRewardDay7 = *req.RewardDay7
	}

	if req.RewardDay8To13 != nil && *req.RewardDay8To13 > 0 {
		if err := model.UpdateOption(setting.OptionCheckinRewardDay8To13, strconv.Itoa(*req.RewardDay8To13)); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
		setting.CheckinRewardDay8To13 = *req.RewardDay8To13
	}

	if req.RewardDay14To29 != nil && *req.RewardDay14To29 > 0 {
		if err := model.UpdateOption(setting.OptionCheckinRewardDay14To29, strconv.Itoa(*req.RewardDay14To29)); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
		setting.CheckinRewardDay14To29 = *req.RewardDay14To29
	}

	if req.RewardDay30 != nil && *req.RewardDay30 > 0 {
		if err := model.UpdateOption(setting.OptionCheckinRewardDay30, strconv.Itoa(*req.RewardDay30)); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
		setting.CheckinRewardDay30 = *req.RewardDay30
	}

	if req.RewardDay31Plus != nil && *req.RewardDay31Plus > 0 {
		if err := model.UpdateOption(setting.OptionCheckinRewardDay31Plus, strconv.Itoa(*req.RewardDay31Plus)); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
		setting.CheckinRewardDay31Plus = *req.RewardDay31Plus
	}

	if req.BonusMin != nil && *req.BonusMin >= 0 {
		if err := model.UpdateOption(setting.OptionCheckinBonusMin, strconv.Itoa(*req.BonusMin)); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
		setting.CheckinBonusMin = *req.BonusMin
	}

	if req.BonusMax != nil && *req.BonusMax > 0 {
		if err := model.UpdateOption(setting.OptionCheckinBonusMax, strconv.Itoa(*req.BonusMax)); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
		setting.CheckinBonusMax = *req.BonusMax
	}

	if req.MakeupCost != nil && *req.MakeupCost >= 0 {
		if err := model.UpdateOption(setting.OptionCheckinMakeupCost, strconv.Itoa(*req.MakeupCost)); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
		setting.CheckinMakeupCost = *req.MakeupCost
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "设置已更新",
		"data":    setting.GetCheckinSettings(),
	})
}

package controller

import (
	"net/http"
	"strconv"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"

	"github.com/gin-gonic/gin"
)

func GetAllQuotaDates(c *gin.Context) {
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	username := c.Query("username")
	timeUnit := c.Query("default_time")
	if timeUnit == "" {
		timeUnit = "hour"
	}

	// 验证并规范化时间范围
	startTimestamp, endTimestamp, err := common.ValidateTimeRange(startTimestamp, endTimestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 使用带缓存的查询
	dates, err := model.GetQuotaDataWithCache(0, username, startTimestamp, endTimestamp, timeUnit)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    dates,
	})
}

func GetUserQuotaDates(c *gin.Context) {
	userId := c.GetInt("id")
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	timeUnit := c.Query("default_time")
	if timeUnit == "" {
		timeUnit = "hour"
	}

	// 验证并规范化时间范围
	startTimestamp, endTimestamp, err := common.ValidateTimeRange(startTimestamp, endTimestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 使用带缓存的查询
	dates, err := model.GetQuotaDataWithCache(userId, "", startTimestamp, endTimestamp, timeUnit)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    dates,
	})
}

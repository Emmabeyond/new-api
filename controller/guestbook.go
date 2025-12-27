package controller

import (
	"net/http"
	"strconv"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
)

// ==================== 请求结构体 ====================

// CreateMessageRequest 创建留言请求
type CreateMessageRequest struct {
	Content string `json:"content" binding:"required,min=1,max=500"`
}

// ReviewMessageRequest 审核留言请求
type ReviewMessageRequest struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}

// FeatureMessageRequest 设置精选请求
type FeatureMessageRequest struct {
	IsFeatured bool `json:"is_featured"`
}

// AdminReplyRequest 管理员回复请求
type AdminReplyRequest struct {
	Reply string `json:"reply" binding:"required,min=1,max=300"`
}

// ==================== 公开接口 ====================

// GetApprovedMessages 获取已审核通过的留言列表
// GET /api/guestbook/messages
func GetApprovedMessages(c *gin.Context) {
	pageInfo := common.GetPageQuery(c)

	messages, total, err := service.GetApprovedMessages(pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "获取留言失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"messages": messages,
			"total":    total,
			"page":     pageInfo.Page,
		},
	})
}

// GetFeaturedMessages 获取精选留言
// GET /api/guestbook/featured
func GetFeaturedMessages(c *gin.Context) {
	// 默认获取3条精选留言用于 Dashboard 展示
	limitStr := c.DefaultQuery("limit", "3")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 3
	}
	if limit > 5 {
		limit = 5
	}

	messages, err := service.GetFeaturedMessages(limit)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "获取精选留言失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    messages,
	})
}


// CreateMessage 提交留言（需登录）
// POST /api/guestbook/messages
func CreateMessage(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "请先登录",
		})
		return
	}

	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的请求参数",
		})
		return
	}

	// XSS 检测
	if common.DetectAndLogXSSAttempt(req.Content, userId, c.ClientIP(), c.GetHeader("User-Agent"), c.Request.URL.Path) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "留言内容包含非法字符",
		})
		return
	}

	// 获取用户名
	user, err := model.GetUserById(userId, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "获取用户信息失败",
		})
		return
	}

	// 创建留言
	_, err = service.CreateGuestbookMessage(userId, user.Username, req.Content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "留言提交成功，等待审核",
	})
}

// GetMyMessages 获取我的留言（需登录）
// GET /api/guestbook/my-messages
func GetMyMessages(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "请先登录",
		})
		return
	}

	messages, err := service.GetUserMessages(userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "获取留言失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    messages,
	})
}

// DeleteMyMessage 删除我的留言（需登录）
// DELETE /api/guestbook/messages/:id
func DeleteMyMessage(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "请先登录",
		})
		return
	}

	idStr := c.Param("id")
	messageId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的留言ID",
		})
		return
	}

	err = service.DeleteUserGuestbookMessage(userId, messageId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "留言删除成功",
	})
}


// ==================== 管理员接口 ====================

// AdminGetAllMessages 管理员获取所有留言
// GET /api/guestbook/admin/messages
func AdminGetAllMessages(c *gin.Context) {
	pageInfo := common.GetPageQuery(c)
	status := c.Query("status")
	keyword := c.Query("keyword")

	messages, total, err := service.GetAllMessagesAdmin(pageInfo.Page, pageInfo.PageSize, status, keyword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "获取留言失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"messages": messages,
			"total":    total,
			"page":     pageInfo.Page,
		},
	})
}

// AdminReviewMessage 管理员审核留言
// PUT /api/guestbook/admin/messages/:id/review
func AdminReviewMessage(c *gin.Context) {
	adminId := c.GetInt("id")

	idStr := c.Param("id")
	messageId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的留言ID",
		})
		return
	}

	var req ReviewMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的请求参数，status 必须为 approved 或 rejected",
		})
		return
	}

	err = service.ReviewGuestbookMessage(messageId, req.Status, adminId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	statusText := "通过"
	if req.Status == model.GuestbookStatusRejected {
		statusText = "拒绝"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "留言审核" + statusText,
	})
}

// AdminFeatureMessage 管理员设置/取消精选
// PUT /api/guestbook/admin/messages/:id/feature
func AdminFeatureMessage(c *gin.Context) {
	idStr := c.Param("id")
	messageId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的留言ID",
		})
		return
	}

	var req FeatureMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的请求参数",
		})
		return
	}

	err = service.FeatureGuestbookMessage(messageId, req.IsFeatured)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	statusText := "取消精选"
	if req.IsFeatured {
		statusText = "设为精选"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "留言已" + statusText,
	})
}

// AdminDeleteMessage 管理员删除留言
// DELETE /api/guestbook/admin/messages/:id
func AdminDeleteMessage(c *gin.Context) {
	idStr := c.Param("id")
	messageId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的留言ID",
		})
		return
	}

	err = service.DeleteGuestbookMessage(messageId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "留言删除成功",
	})
}


// AdminReplyToMessage 管理员回复留言
// PUT /api/guestbook/admin/messages/:id/reply
func AdminReplyToMessage(c *gin.Context) {
	adminId := c.GetInt("id")

	idStr := c.Param("id")
	messageId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的留言ID",
		})
		return
	}

	var req AdminReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的请求参数，回复内容不能为空且不能超过300字",
		})
		return
	}

	// XSS 检测
	if common.DetectAndLogXSSAttempt(req.Reply, adminId, c.ClientIP(), c.GetHeader("User-Agent"), c.Request.URL.Path) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "回复内容包含非法字符",
		})
		return
	}

	err = service.AdminReplyMessage(messageId, req.Reply, adminId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "回复成功",
	})
}

// AdminDeleteReply 管理员删除回复
// DELETE /api/guestbook/admin/messages/:id/reply
func AdminDeleteReply(c *gin.Context) {
	idStr := c.Param("id")
	messageId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的留言ID",
		})
		return
	}

	err = service.AdminDeleteReply(messageId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "回复已删除",
	})
}

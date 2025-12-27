package service

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
)

// 留言板相关错误
var (
	ErrGuestbookContentEmpty     = errors.New("留言内容不能为空")
	ErrGuestbookContentTooLong   = errors.New("留言内容不能超过500字")
	ErrGuestbookDailyLimitExceed = errors.New("今日留言次数已达上限")
	ErrGuestbookMessageNotFound  = errors.New("留言不存在")
	ErrGuestbookNoPermission     = errors.New("无权操作此留言")
	ErrGuestbookFeaturedLimit    = errors.New("精选留言数量已达上限(5条)")
	ErrGuestbookNotApproved      = errors.New("只能精选已审核通过的留言")
	ErrGuestbookInvalidStatus    = errors.New("无效的审核状态")
	ErrGuestbookAlreadyReviewed  = errors.New("该留言已审核")
	ErrGuestbookReplyEmpty       = errors.New("回复内容不能为空")
	ErrGuestbookReplyTooLong     = errors.New("回复内容不能超过300字")
	ErrGuestbookReplyNotApproved = errors.New("只能回复已审核通过的留言")
	ErrGuestbookNoReply          = errors.New("该留言没有回复")
)

// 留言板配置常量
const (
	GuestbookMaxContentLength = 500 // 最大留言长度
	GuestbookMinContentLength = 1   // 最小留言长度
	GuestbookDailyLimit       = 3   // 每日留言限制
	GuestbookMaxFeatured      = 5   // 最大精选数量
	GuestbookMaxReplyLength   = 300 // 最大回复长度
)

// ValidateGuestbookContent 验证留言内容
// 检查内容是否为空、是否超长、是否只包含空白字符
func ValidateGuestbookContent(content string) (string, error) {
	// 去除首尾空白
	trimmed := strings.TrimSpace(content)

	// 检查是否为空或只包含空白字符
	if trimmed == "" {
		return "", ErrGuestbookContentEmpty
	}

	// 检查长度（使用 rune 计数以正确处理中文）
	length := utf8.RuneCountInString(trimmed)
	if length > GuestbookMaxContentLength {
		return "", ErrGuestbookContentTooLong
	}

	return trimmed, nil
}

// SanitizeGuestbookContent 净化留言内容，防止 XSS 攻击
func SanitizeGuestbookContent(content string) string {
	// 使用 common 包的 XSS 净化功能
	config := common.InputValidationConfig{
		MaxLength:     GuestbookMaxContentLength,
		AllowEmpty:    false,
		TrimSpace:     true,
		AllowNewlines: true, // 允许换行
	}
	return common.SanitizeUserInputWithConfig(content, config)
}

// CheckGuestbookDailyLimit 检查用户是否超过每日留言限制
// 返回 true 表示可以继续提交，false 表示已达上限
func CheckGuestbookDailyLimit(userId int) (bool, error) {
	count, err := model.GetUserDailyMessageCount(userId)
	if err != nil {
		return false, err
	}
	return count < GuestbookDailyLimit, nil
}

// CreateGuestbookMessage 创建留言
// 包含完整的验证、净化和限制检查
func CreateGuestbookMessage(userId int, username string, content string) (*model.GuestbookMessage, error) {
	// 1. 验证内容
	validatedContent, err := ValidateGuestbookContent(content)
	if err != nil {
		return nil, err
	}

	// 2. 检查每日限制
	canSubmit, err := CheckGuestbookDailyLimit(userId)
	if err != nil {
		return nil, err
	}
	if !canSubmit {
		return nil, ErrGuestbookDailyLimitExceed
	}

	// 3. 净化内容（XSS 防护）
	sanitizedContent := SanitizeGuestbookContent(validatedContent)

	// 4. 创建留言记录
	message := &model.GuestbookMessage{
		UserId:   userId,
		Username: username,
		Content:  sanitizedContent,
		Status:   model.GuestbookStatusPending,
	}

	err = model.CreateGuestbookMessage(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// ReviewGuestbookMessage 审核留言
// status: "approved" 或 "rejected"
func ReviewGuestbookMessage(messageId int, status string, adminId int) error {
	// 验证状态值
	if status != model.GuestbookStatusApproved && status != model.GuestbookStatusRejected {
		return ErrGuestbookInvalidStatus
	}

	// 获取留言
	message, err := model.GetGuestbookMessageById(messageId)
	if err != nil {
		return ErrGuestbookMessageNotFound
	}

	// 检查是否已审核（可选：允许重复审核以修改状态）
	// 如果需要严格限制，取消下面的注释
	// if message.ReviewedAt != nil {
	// 	return ErrGuestbookAlreadyReviewed
	// }

	// 如果从已通过变为拒绝，且是精选留言，需要取消精选
	if message.Status == model.GuestbookStatusApproved && status == model.GuestbookStatusRejected && message.IsFeatured {
		err = model.FeatureGuestbookMessage(messageId, false)
		if err != nil {
			return err
		}
	}

	// 执行审核
	return model.ReviewGuestbookMessage(messageId, status, adminId)
}

// FeatureGuestbookMessage 设置/取消精选留言
func FeatureGuestbookMessage(messageId int, isFeatured bool) error {
	// 如果要设置为精选，需要额外检查
	if isFeatured {
		// 检查留言是否存在且已审核通过
		message, err := model.GetGuestbookMessageById(messageId)
		if err != nil {
			return ErrGuestbookMessageNotFound
		}
		if message.Status != model.GuestbookStatusApproved {
			return ErrGuestbookNotApproved
		}

		// 检查精选数量限制
		count, err := model.GetFeaturedCount()
		if err != nil {
			return err
		}
		if count >= GuestbookMaxFeatured {
			return ErrGuestbookFeaturedLimit
		}
	}

	return model.FeatureGuestbookMessage(messageId, isFeatured)
}

// DeleteGuestbookMessage 管理员删除留言
func DeleteGuestbookMessage(messageId int) error {
	// 检查留言是否存在
	_, err := model.GetGuestbookMessageById(messageId)
	if err != nil {
		return ErrGuestbookMessageNotFound
	}

	return model.DeleteGuestbookMessage(messageId)
}

// DeleteUserGuestbookMessage 用户删除自己的留言
func DeleteUserGuestbookMessage(userId, messageId int) error {
	return model.DeleteUserGuestbookMessage(userId, messageId)
}

// GetApprovedMessages 获取已审核通过的留言（公开接口）
func GetApprovedMessages(page, pageSize int) ([]model.GuestbookMessage, int64, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return model.GetApprovedMessages(page, pageSize)
}

// GetFeaturedMessages 获取精选留言（Dashboard 展示）
func GetFeaturedMessages(limit int) ([]model.GuestbookMessage, error) {
	if limit < 1 {
		limit = 3
	}
	if limit > 5 {
		limit = 5
	}

	return model.GetFeaturedMessages(limit)
}

// GetUserMessages 获取用户的留言
func GetUserMessages(userId int) ([]model.GuestbookMessage, error) {
	return model.GetUserMessages(userId)
}

// GetAllMessagesAdmin 管理员获取所有留言
func GetAllMessagesAdmin(page, pageSize int, status, keyword string) ([]model.GuestbookMessage, int64, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return model.GetAllMessagesAdmin(page, pageSize, status, keyword)
}


// ValidateGuestbookReply 验证回复内容
func ValidateGuestbookReply(reply string) (string, error) {
	// 去除首尾空白
	trimmed := strings.TrimSpace(reply)

	// 检查是否为空
	if trimmed == "" {
		return "", ErrGuestbookReplyEmpty
	}

	// 检查长度（使用 rune 计数以正确处理中文）
	length := utf8.RuneCountInString(trimmed)
	if length > GuestbookMaxReplyLength {
		return "", ErrGuestbookReplyTooLong
	}

	return trimmed, nil
}

// SanitizeGuestbookReply 净化回复内容，防止 XSS 攻击
func SanitizeGuestbookReply(reply string) string {
	config := common.InputValidationConfig{
		MaxLength:     GuestbookMaxReplyLength,
		AllowEmpty:    false,
		TrimSpace:     true,
		AllowNewlines: true,
	}
	return common.SanitizeUserInputWithConfig(reply, config)
}

// AdminReplyMessage 管理员回复留言
func AdminReplyMessage(messageId int, reply string, adminId int) error {
	// 1. 验证回复内容
	validatedReply, err := ValidateGuestbookReply(reply)
	if err != nil {
		return err
	}

	// 2. 获取留言
	message, err := model.GetGuestbookMessageById(messageId)
	if err != nil {
		return ErrGuestbookMessageNotFound
	}

	// 3. 检查留言是否已审核通过
	if message.Status != model.GuestbookStatusApproved {
		return ErrGuestbookReplyNotApproved
	}

	// 4. 净化回复内容
	sanitizedReply := SanitizeGuestbookReply(validatedReply)

	// 5. 保存回复
	return model.AdminReplyGuestbookMessage(messageId, sanitizedReply, adminId)
}

// AdminDeleteReply 管理员删除回复
func AdminDeleteReply(messageId int) error {
	// 1. 获取留言
	message, err := model.GetGuestbookMessageById(messageId)
	if err != nil {
		return ErrGuestbookMessageNotFound
	}

	// 2. 检查是否有回复
	if message.AdminReply == nil || *message.AdminReply == "" {
		return ErrGuestbookNoReply
	}

	// 3. 删除回复
	return model.DeleteAdminReply(messageId)
}

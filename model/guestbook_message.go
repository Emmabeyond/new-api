package model

import (
	"errors"
	"time"

	"github.com/QuantumNous/new-api/common"
	"gorm.io/gorm"
)

// 留言状态常量
const (
	GuestbookStatusPending  = "pending"  // 待审核
	GuestbookStatusApproved = "approved" // 已通过
	GuestbookStatusRejected = "rejected" // 已拒绝
)

// GuestbookMessage 留言板消息模型
type GuestbookMessage struct {
	Id         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId     int    `json:"user_id" gorm:"index;not null"`
	Username   string `json:"username" gorm:"type:varchar(100);not null"`
	Content    string `json:"content" gorm:"type:varchar(500);not null"`
	Status     string `json:"status" gorm:"type:varchar(20);default:'pending';index"`
	IsFeatured bool   `json:"is_featured" gorm:"default:false;index"`
	CreatedAt  int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  int64  `json:"updated_at" gorm:"autoUpdateTime"`
	ReviewedAt *int64 `json:"reviewed_at"`
	ReviewedBy *int   `json:"reviewed_by"`
	DeletedAt  *int64 `json:"deleted_at" gorm:"index"`
	// 管理员回复字段
	AdminReply   *string `json:"admin_reply" gorm:"type:varchar(300)"`
	AdminReplyAt *int64  `json:"admin_reply_at"`
	AdminReplyBy *int    `json:"admin_reply_by"`
}

// TableName 返回表名
func (GuestbookMessage) TableName() string {
	return "guestbook_messages"
}

// CreateGuestbookMessage 创建留言
func CreateGuestbookMessage(message *GuestbookMessage) error {
	message.CreatedAt = time.Now().Unix()
	message.UpdatedAt = message.CreatedAt
	message.Status = GuestbookStatusPending
	return DB.Create(message).Error
}

// GetGuestbookMessageById 根据ID获取留言
func GetGuestbookMessageById(id int) (*GuestbookMessage, error) {
	var message GuestbookMessage
	err := DB.Where("id = ? AND deleted_at IS NULL", id).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// UpdateGuestbookMessage 更新留言
func UpdateGuestbookMessage(message *GuestbookMessage) error {
	message.UpdatedAt = time.Now().Unix()
	return DB.Model(message).Updates(map[string]interface{}{
		"content":     message.Content,
		"status":      message.Status,
		"is_featured": message.IsFeatured,
		"updated_at":  message.UpdatedAt,
		"reviewed_at": message.ReviewedAt,
		"reviewed_by": message.ReviewedBy,
	}).Error
}

// DeleteGuestbookMessage 软删除留言
func DeleteGuestbookMessage(id int) error {
	now := time.Now().Unix()
	return DB.Model(&GuestbookMessage{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// GetApprovedMessages 获取已审核通过的留言（带分页）
// 精选留言优先显示，然后按创建时间倒序
func GetApprovedMessages(page, pageSize int) ([]GuestbookMessage, int64, error) {
	var messages []GuestbookMessage
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	// 获取总数
	err := DB.Model(&GuestbookMessage{}).
		Where("status = ? AND deleted_at IS NULL", GuestbookStatusApproved).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据，精选优先，然后按创建时间倒序
	err = DB.Where("status = ? AND deleted_at IS NULL", GuestbookStatusApproved).
		Order("is_featured DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// GetFeaturedMessages 获取精选留言（带限制）
func GetFeaturedMessages(limit int) ([]GuestbookMessage, error) {
	var messages []GuestbookMessage
	err := DB.Where("status = ? AND is_featured = ? AND deleted_at IS NULL", GuestbookStatusApproved, true).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

// GetUserMessages 获取用户的留言
func GetUserMessages(userId int) ([]GuestbookMessage, error) {
	var messages []GuestbookMessage
	err := DB.Where("user_id = ? AND deleted_at IS NULL", userId).
		Order("created_at DESC").
		Find(&messages).Error
	return messages, err
}

// GetAllMessagesAdmin 管理员获取所有留言（带筛选和搜索）
func GetAllMessagesAdmin(page, pageSize int, status, keyword string) ([]GuestbookMessage, int64, error) {
	var messages []GuestbookMessage
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	// 构建查询
	query := DB.Model(&GuestbookMessage{}).Where("deleted_at IS NULL")

	// 状态筛选
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	// 关键词搜索（内容或用户名）
	if keyword != "" {
		query = query.Where("content LIKE ? OR username LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// GetFeaturedCount 获取当前精选留言数量
func GetFeaturedCount() (int64, error) {
	var count int64
	err := DB.Model(&GuestbookMessage{}).
		Where("is_featured = ? AND deleted_at IS NULL", true).
		Count(&count).Error
	return count, err
}

// GetUserDailyMessageCount 获取用户当日留言数量
func GetUserDailyMessageCount(userId int) (int64, error) {
	var count int64
	// 获取今天零点的时间戳
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()

	err := DB.Model(&GuestbookMessage{}).
		Where("user_id = ? AND created_at >= ? AND deleted_at IS NULL", userId, todayStart).
		Count(&count).Error
	return count, err
}

// ReviewGuestbookMessage 审核留言
func ReviewGuestbookMessage(id int, status string, adminId int) error {
	if status != GuestbookStatusApproved && status != GuestbookStatusRejected {
		return errors.New("无效的审核状态")
	}

	now := time.Now().Unix()
	return DB.Model(&GuestbookMessage{}).Where("id = ? AND deleted_at IS NULL", id).Updates(map[string]interface{}{
		"status":      status,
		"reviewed_at": now,
		"reviewed_by": adminId,
		"updated_at":  now,
	}).Error
}

// FeatureGuestbookMessage 设置/取消精选
func FeatureGuestbookMessage(id int, isFeatured bool) error {
	// 如果要设置为精选，先检查数量限制
	if isFeatured {
		count, err := GetFeaturedCount()
		if err != nil {
			return err
		}
		if count >= 5 {
			return errors.New("精选留言数量已达上限(5条)")
		}

		// 检查留言是否已审核通过
		message, err := GetGuestbookMessageById(id)
		if err != nil {
			return err
		}
		if message.Status != GuestbookStatusApproved {
			return errors.New("只能精选已审核通过的留言")
		}
	}

	now := time.Now().Unix()
	return DB.Model(&GuestbookMessage{}).Where("id = ? AND deleted_at IS NULL", id).Updates(map[string]interface{}{
		"is_featured": isFeatured,
		"updated_at":  now,
	}).Error
}

// DeleteUserGuestbookMessage 用户删除自己的留言
func DeleteUserGuestbookMessage(userId, messageId int) error {
	// 先检查留言是否属于该用户
	var message GuestbookMessage
	err := DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", messageId, userId).First(&message).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("留言不存在或无权删除")
		}
		return err
	}

	return DeleteGuestbookMessage(messageId)
}

// createGuestbookIndexes 创建留言板相关索引
func createGuestbookIndexes() {
	var err error

	// 复合索引：状态+创建时间（用于留言墙查询）
	if common.UsingPostgreSQL {
		err = DB.Exec("CREATE INDEX IF NOT EXISTS idx_guestbook_status_created ON guestbook_messages(status, created_at DESC)").Error
	} else if common.UsingMySQL {
		var count int64
		DB.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'guestbook_messages' AND index_name = 'idx_guestbook_status_created'").Scan(&count)
		if count == 0 {
			err = DB.Exec("CREATE INDEX idx_guestbook_status_created ON guestbook_messages(status, created_at DESC)").Error
		}
	} else {
		err = DB.Exec("CREATE INDEX IF NOT EXISTS idx_guestbook_status_created ON guestbook_messages(status, created_at DESC)").Error
	}
	if err != nil {
		common.SysLog("Warning: failed to create idx_guestbook_status_created: " + err.Error())
	}

	// 复合索引：用户ID+创建时间（用于用户留言查询）
	if common.UsingPostgreSQL {
		err = DB.Exec("CREATE INDEX IF NOT EXISTS idx_guestbook_user_created ON guestbook_messages(user_id, created_at DESC)").Error
	} else if common.UsingMySQL {
		var count int64
		DB.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'guestbook_messages' AND index_name = 'idx_guestbook_user_created'").Scan(&count)
		if count == 0 {
			err = DB.Exec("CREATE INDEX idx_guestbook_user_created ON guestbook_messages(user_id, created_at DESC)").Error
		}
	} else {
		err = DB.Exec("CREATE INDEX IF NOT EXISTS idx_guestbook_user_created ON guestbook_messages(user_id, created_at DESC)").Error
	}
	if err != nil {
		common.SysLog("Warning: failed to create idx_guestbook_user_created: " + err.Error())
	}

	// 复合索引：精选+状态+创建时间（用于精选留言查询）
	if common.UsingPostgreSQL {
		err = DB.Exec("CREATE INDEX IF NOT EXISTS idx_guestbook_featured ON guestbook_messages(is_featured, status, created_at DESC)").Error
	} else if common.UsingMySQL {
		var count int64
		DB.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'guestbook_messages' AND index_name = 'idx_guestbook_featured'").Scan(&count)
		if count == 0 {
			err = DB.Exec("CREATE INDEX idx_guestbook_featured ON guestbook_messages(is_featured, status, created_at DESC)").Error
		}
	} else {
		err = DB.Exec("CREATE INDEX IF NOT EXISTS idx_guestbook_featured ON guestbook_messages(is_featured, status, created_at DESC)").Error
	}
	if err != nil {
		common.SysLog("Warning: failed to create idx_guestbook_featured: " + err.Error())
	}

	common.SysLog("Guestbook indexes created/verified")
}


// InitDefaultGuestbookMessages 初始化默认留言数据
// 在系统首次部署时创建示例留言，展示留言板功能
func InitDefaultGuestbookMessages() error {
	// 检查是否已有留言数据
	var count int64
	if err := DB.Model(&GuestbookMessage{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil // 已有数据，不需要初始化
	}

	now := time.Now().Unix()
	reviewedAt := now

	// 创建示例留言数据
	// 使用系统用户 ID 1 (root) 作为示例用户和审核者
	defaultMessages := []GuestbookMessage{
		{
			UserId:     1,
			Username:   "AI爱好者",
			Content:    "用 Claude Sonnet 4.5 写代码真的很丝滑，Agent 能力太强了！",
			Status:     GuestbookStatusApproved,
			IsFeatured: true,
			CreatedAt:  now - 86400*2, // 2天前
			UpdatedAt:  now,
			ReviewedAt: &reviewedAt,
			ReviewedBy: func() *int { v := 1; return &v }(),
		},
		{
			UserId:     1,
			Username:   "开发者小王",
			Content:    "GPT-5.2 Thinking 模式的推理能力惊艳，400K 上下文处理文档很方便",
			Status:     GuestbookStatusApproved,
			IsFeatured: true,
			CreatedAt:  now - 86400, // 1天前
			UpdatedAt:  now,
			ReviewedAt: &reviewedAt,
			ReviewedBy: func() *int { v := 1; return &v }(),
		},
		{
			UserId:     1,
			Username:   "技术达人",
			Content:    "Claude Opus 4.5 降价后性价比无敌，复杂任务首选",
			Status:     GuestbookStatusApproved,
			IsFeatured: false,
			CreatedAt:  now - 43200, // 12小时前
			UpdatedAt:  now,
			ReviewedAt: &reviewedAt,
			ReviewedBy: func() *int { v := 1; return &v }(),
		},
		{
			UserId:     1,
			Username:   "研究员",
			Content:    "建议增加 GPT-5.2 Pro 的支持，研究级任务需要",
			Status:     GuestbookStatusApproved,
			IsFeatured: false,
			CreatedAt:  now - 21600, // 6小时前
			UpdatedAt:  now,
			ReviewedAt: &reviewedAt,
			ReviewedBy: func() *int { v := 1; return &v }(),
		},
		{
			UserId:     1,
			Username:   "新手开发者",
			Content:    "API 文档很清晰，接入很顺利，感谢团队！",
			Status:     GuestbookStatusApproved,
			IsFeatured: true,
			CreatedAt:  now - 3600, // 1小时前
			UpdatedAt:  now,
			ReviewedAt: &reviewedAt,
			ReviewedBy: func() *int { v := 1; return &v }(),
		},
	}

	for _, msg := range defaultMessages {
		if err := DB.Create(&msg).Error; err != nil {
			return err
		}
	}

	return nil
}


// AdminReplyGuestbookMessage 管理员回复留言
func AdminReplyGuestbookMessage(id int, reply string, adminId int) error {
	now := time.Now().Unix()
	return DB.Model(&GuestbookMessage{}).Where("id = ? AND deleted_at IS NULL", id).Updates(map[string]interface{}{
		"admin_reply":    reply,
		"admin_reply_at": now,
		"admin_reply_by": adminId,
		"updated_at":     now,
	}).Error
}

// DeleteAdminReply 删除管理员回复
func DeleteAdminReply(id int) error {
	now := time.Now().Unix()
	return DB.Model(&GuestbookMessage{}).Where("id = ? AND deleted_at IS NULL", id).Updates(map[string]interface{}{
		"admin_reply":    nil,
		"admin_reply_at": nil,
		"admin_reply_by": nil,
		"updated_at":     now,
	}).Error
}

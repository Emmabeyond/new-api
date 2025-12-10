package model

import (
	"time"

	"gorm.io/gorm"
)

// TokenPenalty represents a penalty record for audit logging
type TokenPenalty struct {
	Id          int            `json:"id" gorm:"primaryKey;autoIncrement"`
	TokenId     int            `json:"token_id" gorm:"index;not null"`
	UserId      int            `json:"user_id" gorm:"index;not null"`
	PenaltyType string         `json:"penalty_type" gorm:"type:varchar(20);not null"` // rate_limit, temp_ban, perm_ban
	Reason      string         `json:"reason" gorm:"type:text"`
	AbuseScore  int            `json:"abuse_score"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	ExpiresAt   *time.Time     `json:"expires_at" gorm:"index"`
	LiftedAt    *time.Time     `json:"lifted_at"`
	LiftedBy    *int           `json:"lifted_by"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for TokenPenalty
func (TokenPenalty) TableName() string {
	return "token_penalties"
}

// CreateTokenPenalty creates a new penalty record
func CreateTokenPenalty(penalty *TokenPenalty) error {
	return DB.Create(penalty).Error
}

// GetTokenPenaltyByTokenId gets the latest active penalty for a token
func GetTokenPenaltyByTokenId(tokenId int) (*TokenPenalty, error) {
	var penalty TokenPenalty
	err := DB.Where("token_id = ? AND (expires_at IS NULL OR expires_at > ?) AND lifted_at IS NULL",
		tokenId, time.Now()).
		Order("created_at DESC").
		First(&penalty).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &penalty, nil
}

// GetActivePenalties gets all active penalties
func GetActivePenalties(page, pageSize int) ([]*TokenPenalty, int64, error) {
	var penalties []*TokenPenalty
	var total int64

	query := DB.Model(&TokenPenalty{}).
		Where("(expires_at IS NULL OR expires_at > ?) AND lifted_at IS NULL", time.Now())

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&penalties).Error
	if err != nil {
		return nil, 0, err
	}

	return penalties, total, nil
}

// LiftTokenPenalty lifts a penalty by setting lifted_at and lifted_by
func LiftTokenPenalty(tokenId int, liftedBy int) error {
	return DB.Model(&TokenPenalty{}).
		Where("token_id = ? AND lifted_at IS NULL", tokenId).
		Updates(map[string]interface{}{
			"lifted_at": time.Now(),
			"lifted_by": liftedBy,
		}).Error
}

// GetPenaltyHistory gets penalty history for a token
func GetPenaltyHistory(tokenId int, page, pageSize int) ([]*TokenPenalty, int64, error) {
	var penalties []*TokenPenalty
	var total int64

	query := DB.Model(&TokenPenalty{}).Where("token_id = ?", tokenId)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&penalties).Error
	if err != nil {
		return nil, 0, err
	}

	return penalties, total, nil
}

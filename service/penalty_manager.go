package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/security_setting"
)

// PenaltyType constants
const (
	PenaltyTypeRateLimit = "rate_limit"
	PenaltyTypeTempBan   = "temp_ban"
	PenaltyTypePermBan   = "perm_ban"
)

// PenaltyStatus represents the current penalty status of a token
type PenaltyStatus struct {
	IsPenalized  bool      `json:"is_penalized"`
	PenaltyType  string    `json:"penalty_type"`
	Reason       string    `json:"reason"`
	AbuseScore   int       `json:"abuse_score"`
	ExpiresAt    time.Time `json:"expires_at"`
	RateLimitRPM int       `json:"rate_limit_rpm"` // Requests per minute when rate limited
}

// PenaltyRecord represents a penalty record for audit logging
type PenaltyRecord struct {
	TokenID     int       `json:"token_id"`
	TokenName   string    `json:"token_name"`
	UserID      int       `json:"user_id"`
	PenaltyType string    `json:"penalty_type"`
	Reason      string    `json:"reason"`
	AbuseScore  int       `json:"abuse_score"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// PenaltyManager manages penalties for tokens
type PenaltyManager interface {
	// ApplyPenalty applies a penalty to a token
	ApplyPenalty(tokenID int, penaltyType string, reason string, abuseScore int, durationMinutes int) error

	// CheckPenalty checks if a token is currently penalized
	CheckPenalty(tokenID int) (*PenaltyStatus, error)

	// LiftPenalty removes a penalty from a token
	LiftPenalty(tokenID int) error

	// GetActivePenalties returns all active penalties
	GetActivePenalties() ([]PenaltyRecord, error)
}


// redisPenaltyManager implements PenaltyManager using Redis
type redisPenaltyManager struct{}

// memoryPenaltyManager implements PenaltyManager using in-memory storage
type memoryPenaltyManager struct {
	mu        sync.RWMutex
	penalties map[int]*PenaltyStatus // tokenID -> penalty status
	records   []PenaltyRecord        // audit log
}

var (
	penaltyManager     PenaltyManager
	penaltyManagerOnce sync.Once
)

// GetPenaltyManager returns the singleton PenaltyManager instance
func GetPenaltyManager() PenaltyManager {
	penaltyManagerOnce.Do(func() {
		if common.RedisEnabled {
			penaltyManager = &redisPenaltyManager{}
		} else {
			penaltyManager = &memoryPenaltyManager{
				penalties: make(map[int]*PenaltyStatus),
				records:   make([]PenaltyRecord, 0),
			}
		}
	})
	return penaltyManager
}

// Redis key format for penalty tracking
const penaltyKeyPrefix = "abuse:penalty:"

func getPenaltyKey(tokenID int) string {
	return fmt.Sprintf("%s%d", penaltyKeyPrefix, tokenID)
}

// Redis implementation

func (r *redisPenaltyManager) ApplyPenalty(tokenID int, penaltyType string, reason string, abuseScore int, durationMinutes int) error {
	ctx := context.Background()
	key := getPenaltyKey(tokenID)

	settings := security_setting.GetSecuritySetting()

	status := &PenaltyStatus{
		IsPenalized:  true,
		PenaltyType:  penaltyType,
		Reason:       reason,
		AbuseScore:   abuseScore,
		RateLimitRPM: settings.RateLimitRequests,
	}

	// Set expiration based on penalty type
	var expiration time.Duration
	if penaltyType == PenaltyTypePermBan {
		// Permanent ban: set a very long expiration (1 year)
		expiration = 365 * 24 * time.Hour
		status.ExpiresAt = time.Now().Add(expiration)
	} else if durationMinutes > 0 {
		expiration = time.Duration(durationMinutes) * time.Minute
		status.ExpiresAt = time.Now().Add(expiration)
	} else {
		// Default to configured temp ban duration
		expiration = time.Duration(settings.TempBanDurationMinutes) * time.Minute
		status.ExpiresAt = time.Now().Add(expiration)
	}

	// Serialize and store
	data, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal penalty status: %w", err)
	}

	err = common.RDB.Set(ctx, key, string(data), expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to apply penalty: %w", err)
	}

	// Log the penalty (in production, this would go to database)
	common.SysLog(fmt.Sprintf("Penalty applied to token %d: type=%s, reason=%s, score=%d, expires=%v",
		tokenID, penaltyType, reason, abuseScore, status.ExpiresAt))

	return nil
}

func (r *redisPenaltyManager) CheckPenalty(tokenID int) (*PenaltyStatus, error) {
	ctx := context.Background()
	key := getPenaltyKey(tokenID)

	data, err := common.RDB.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return &PenaltyStatus{IsPenalized: false}, nil
		}
		return nil, fmt.Errorf("failed to check penalty: %w", err)
	}

	var status PenaltyStatus
	if err := json.Unmarshal([]byte(data), &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal penalty status: %w", err)
	}

	return &status, nil
}

func (r *redisPenaltyManager) LiftPenalty(tokenID int) error {
	ctx := context.Background()
	key := getPenaltyKey(tokenID)

	err := common.RDB.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to lift penalty: %w", err)
	}

	common.SysLog(fmt.Sprintf("Penalty lifted for token %d", tokenID))
	return nil
}

func (r *redisPenaltyManager) GetActivePenalties() ([]PenaltyRecord, error) {
	ctx := context.Background()

	// Scan for all penalty keys
	var cursor uint64
	var records []PenaltyRecord

	for {
		keys, nextCursor, err := common.RDB.Scan(ctx, cursor, penaltyKeyPrefix+"*", 100).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to scan penalties: %w", err)
		}

		for _, key := range keys {
			data, err := common.RDB.Get(ctx, key).Result()
			if err != nil {
				continue
			}

			var status PenaltyStatus
			if err := json.Unmarshal([]byte(data), &status); err != nil {
				continue
			}

			// Extract token ID from key
			var tokenID int
			fmt.Sscanf(key, penaltyKeyPrefix+"%d", &tokenID)

			records = append(records, PenaltyRecord{
				TokenID:     tokenID,
				PenaltyType: status.PenaltyType,
				Reason:      status.Reason,
				AbuseScore:  status.AbuseScore,
				CreatedAt:   status.ExpiresAt.Add(-time.Duration(security_setting.GetSecuritySetting().TempBanDurationMinutes) * time.Minute),
				ExpiresAt:   status.ExpiresAt,
			})
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return records, nil
}


// Memory implementation

func (m *memoryPenaltyManager) ApplyPenalty(tokenID int, penaltyType string, reason string, abuseScore int, durationMinutes int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	settings := security_setting.GetSecuritySetting()

	status := &PenaltyStatus{
		IsPenalized:  true,
		PenaltyType:  penaltyType,
		Reason:       reason,
		AbuseScore:   abuseScore,
		RateLimitRPM: settings.RateLimitRequests,
	}

	// Set expiration based on penalty type
	if penaltyType == PenaltyTypePermBan {
		status.ExpiresAt = time.Now().Add(365 * 24 * time.Hour)
	} else if durationMinutes > 0 {
		status.ExpiresAt = time.Now().Add(time.Duration(durationMinutes) * time.Minute)
	} else {
		status.ExpiresAt = time.Now().Add(time.Duration(settings.TempBanDurationMinutes) * time.Minute)
	}

	m.penalties[tokenID] = status

	// Add to audit log
	m.records = append(m.records, PenaltyRecord{
		TokenID:     tokenID,
		PenaltyType: penaltyType,
		Reason:      reason,
		AbuseScore:  abuseScore,
		CreatedAt:   time.Now(),
		ExpiresAt:   status.ExpiresAt,
	})

	common.SysLog(fmt.Sprintf("Penalty applied to token %d: type=%s, reason=%s, score=%d, expires=%v",
		tokenID, penaltyType, reason, abuseScore, status.ExpiresAt))

	return nil
}

func (m *memoryPenaltyManager) CheckPenalty(tokenID int) (*PenaltyStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status, exists := m.penalties[tokenID]
	if !exists {
		return &PenaltyStatus{IsPenalized: false}, nil
	}

	// Check if penalty has expired
	if time.Now().After(status.ExpiresAt) {
		// Penalty has expired, will be cleaned up later
		return &PenaltyStatus{IsPenalized: false}, nil
	}

	return status, nil
}

func (m *memoryPenaltyManager) LiftPenalty(tokenID int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.penalties, tokenID)
	common.SysLog(fmt.Sprintf("Penalty lifted for token %d", tokenID))
	return nil
}

func (m *memoryPenaltyManager) GetActivePenalties() ([]PenaltyRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var activeRecords []PenaltyRecord
	now := time.Now()

	for tokenID, status := range m.penalties {
		if status.IsPenalized && now.Before(status.ExpiresAt) {
			activeRecords = append(activeRecords, PenaltyRecord{
				TokenID:     tokenID,
				PenaltyType: status.PenaltyType,
				Reason:      status.Reason,
				AbuseScore:  status.AbuseScore,
				ExpiresAt:   status.ExpiresAt,
			})
		}
	}

	return activeRecords, nil
}

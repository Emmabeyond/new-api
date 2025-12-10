package service

import (
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting/security_setting"
)

// AbuseCheckResult contains the result of an abuse check
type AbuseCheckResult struct {
	Allowed    bool   `json:"allowed"`
	Reason     string `json:"reason"`
	AbuseScore int    `json:"abuse_score"`
	PenaltyType string `json:"penalty_type"`
	RetryAfter int    `json:"retry_after"` // Seconds until retry is allowed
}

// AbuseDetector is the main interface for abuse detection
type AbuseDetector interface {
	// CheckRequest checks if a request should be allowed
	CheckRequest(tokenID int, userID int, group string, modelName string, content string) (*AbuseCheckResult, error)

	// GetAbuseScore returns the current abuse score for a token
	GetAbuseScore(tokenID int) (*AbuseScoreResult, error)

	// IsWhitelisted checks if a token/user/group is whitelisted
	IsWhitelisted(tokenID int, userID int, group string) bool

	// RecordRequest records a request for metrics (even for whitelisted tokens)
	RecordRequest(tokenID int, modelName string, content string) error
}

// abuseDetector implements AbuseDetector
type abuseDetector struct {
	modelSwitchTracker  ModelSwitchTracker
	testContentDetector TestContentDetector
	scoreCalculator     *AbuseScoreCalculator
	penaltyManager      PenaltyManager
}

var defaultAbuseDetector *abuseDetector

// GetAbuseDetector returns the singleton AbuseDetector instance
func GetAbuseDetector() AbuseDetector {
	if defaultAbuseDetector == nil {
		defaultAbuseDetector = &abuseDetector{
			modelSwitchTracker:  GetModelSwitchTracker(),
			testContentDetector: GetTestContentDetector(),
			scoreCalculator:     NewAbuseScoreCalculator(),
			penaltyManager:      GetPenaltyManager(),
		}
	}
	return defaultAbuseDetector
}


// CheckRequest checks if a request should be allowed
func (d *abuseDetector) CheckRequest(tokenID int, userID int, group string, modelName string, content string) (*AbuseCheckResult, error) {
	settings := security_setting.GetSecuritySetting()

	// Check if anti-abuse is enabled
	if !settings.EnableAntiAbuse {
		return &AbuseCheckResult{Allowed: true}, nil
	}

	// Check whitelist (but still record metrics)
	isWhitelisted := d.IsWhitelisted(tokenID, userID, group)

	// Always record the request for metrics
	if err := d.RecordRequest(tokenID, modelName, content); err != nil {
		common.SysError("Failed to record request: " + err.Error())
	}

	// If whitelisted, allow the request
	if isWhitelisted {
		return &AbuseCheckResult{Allowed: true, Reason: "whitelisted"}, nil
	}

	// Check existing penalty
	penaltyStatus, err := d.penaltyManager.CheckPenalty(tokenID)
	if err != nil {
		common.SysError("Failed to check penalty: " + err.Error())
		// On error, allow the request but log the error
		return &AbuseCheckResult{Allowed: true}, nil
	}

	if penaltyStatus.IsPenalized {
		retryAfter := 0
		if !penaltyStatus.ExpiresAt.IsZero() {
			retryAfter = int(time.Until(penaltyStatus.ExpiresAt).Seconds())
			if retryAfter < 0 {
				retryAfter = 0
			}
		}

		return &AbuseCheckResult{
			Allowed:     false,
			Reason:      penaltyStatus.Reason,
			AbuseScore:  penaltyStatus.AbuseScore,
			PenaltyType: penaltyStatus.PenaltyType,
			RetryAfter:  retryAfter,
		}, nil
	}

	// Calculate abuse score
	scoreResult, err := d.scoreCalculator.CalculateScore(tokenID)
	if err != nil {
		common.SysError("Failed to calculate abuse score: " + err.Error())
		return &AbuseCheckResult{Allowed: true}, nil
	}

	// Check if action threshold is exceeded
	if scoreResult.ExceedsAction {
		// Apply penalty
		penaltyType := settings.PenaltyType
		durationMinutes := settings.TempBanDurationMinutes
		reason := "Abuse score exceeded action threshold"

		if err := d.penaltyManager.ApplyPenalty(tokenID, penaltyType, reason, scoreResult.TotalScore, durationMinutes); err != nil {
			common.SysError("Failed to apply penalty: " + err.Error())
		}

		// Record to database for audit
		penalty := &model.TokenPenalty{
			TokenId:     tokenID,
			UserId:      userID,
			PenaltyType: penaltyType,
			Reason:      reason,
			AbuseScore:  scoreResult.TotalScore,
		}
		if penaltyType != PenaltyTypePermBan {
			expiresAt := time.Now().Add(time.Duration(durationMinutes) * time.Minute)
			penalty.ExpiresAt = &expiresAt
		}
		if err := model.CreateTokenPenalty(penalty); err != nil {
			common.SysError("Failed to create penalty record: " + err.Error())
		}

		return &AbuseCheckResult{
			Allowed:     false,
			Reason:      reason,
			AbuseScore:  scoreResult.TotalScore,
			PenaltyType: penaltyType,
			RetryAfter:  durationMinutes * 60,
		}, nil
	}

	// Check if warning threshold is exceeded (log warning but allow)
	if scoreResult.ExceedsWarning {
		common.SysLog("Warning: Token " + string(rune(tokenID)) + " abuse score approaching threshold: " + string(rune(scoreResult.TotalScore)))
	}

	return &AbuseCheckResult{
		Allowed:    true,
		AbuseScore: scoreResult.TotalScore,
	}, nil
}

// GetAbuseScore returns the current abuse score for a token
func (d *abuseDetector) GetAbuseScore(tokenID int) (*AbuseScoreResult, error) {
	return d.scoreCalculator.CalculateScore(tokenID)
}

// IsWhitelisted checks if a token/user/group is whitelisted
func (d *abuseDetector) IsWhitelisted(tokenID int, userID int, group string) bool {
	// Check user whitelist
	if security_setting.IsUserWhitelisted(userID) {
		return true
	}

	// Check group whitelist
	if group != "" && security_setting.IsGroupWhitelisted(group) {
		return true
	}

	return false
}

// RecordRequest records a request for metrics
func (d *abuseDetector) RecordRequest(tokenID int, modelName string, content string) error {
	// Record model request
	if err := d.modelSwitchTracker.RecordModelRequest(tokenID, modelName); err != nil {
		return err
	}

	// Check and record test content
	if d.testContentDetector.IsTestContent(content) {
		if err := d.testContentDetector.RecordTestContent(tokenID); err != nil {
			return err
		}
	}

	return nil
}

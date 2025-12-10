package service

import (
	"github.com/QuantumNous/new-api/setting/security_setting"
)

// AbuseScoreResult contains the calculated abuse score and contributing factors
type AbuseScoreResult struct {
	TotalScore           int `json:"total_score"`
	ModelSwitchScore     int `json:"model_switch_score"`
	TestContentScore     int `json:"test_content_score"`
	ModelSwitchCount     int `json:"model_switch_count"`
	TestContentCount     int `json:"test_content_count"`
	ExceedsWarning       bool `json:"exceeds_warning"`
	ExceedsAction        bool `json:"exceeds_action"`
}

// AbuseScoreCalculator calculates abuse scores for tokens
type AbuseScoreCalculator struct {
	modelSwitchTracker  ModelSwitchTracker
	testContentDetector TestContentDetector
}

// NewAbuseScoreCalculator creates a new AbuseScoreCalculator
func NewAbuseScoreCalculator() *AbuseScoreCalculator {
	return &AbuseScoreCalculator{
		modelSwitchTracker:  GetModelSwitchTracker(),
		testContentDetector: GetTestContentDetector(),
	}
}

// CalculateScore calculates the abuse score for a token
func (c *AbuseScoreCalculator) CalculateScore(tokenID int) (*AbuseScoreResult, error) {
	settings := security_setting.GetSecuritySetting()

	result := &AbuseScoreResult{}

	// Get model switch count
	modelSwitchCount, err := c.modelSwitchTracker.GetDistinctModelCount(tokenID, settings.ModelSwitchWindowMinutes)
	if err != nil {
		return nil, err
	}
	result.ModelSwitchCount = modelSwitchCount

	// Get test content count
	testContentCount, err := c.testContentDetector.GetTestContentCount(tokenID, settings.TestContentWindowMinutes)
	if err != nil {
		return nil, err
	}
	result.TestContentCount = testContentCount

	// Calculate model switch score (0-50 points)
	// Score increases as count approaches and exceeds threshold
	if settings.ModelSwitchThreshold > 0 {
		ratio := float64(modelSwitchCount) / float64(settings.ModelSwitchThreshold)
		result.ModelSwitchScore = int(ratio * 50)
		if result.ModelSwitchScore > 50 {
			result.ModelSwitchScore = 50
		}
	}

	// Calculate test content score (0-50 points)
	// Score increases as count approaches and exceeds threshold
	if settings.TestContentThreshold > 0 {
		ratio := float64(testContentCount) / float64(settings.TestContentThreshold)
		result.TestContentScore = int(ratio * 50)
		if result.TestContentScore > 50 {
			result.TestContentScore = 50
		}
	}

	// Calculate total score (0-100)
	result.TotalScore = result.ModelSwitchScore + result.TestContentScore

	// Check thresholds
	result.ExceedsWarning = result.TotalScore >= settings.AbuseScoreWarningThreshold
	result.ExceedsAction = result.TotalScore >= settings.AbuseScoreActionThreshold

	return result, nil
}

// IsAbusive checks if a token's behavior is considered abusive
func (c *AbuseScoreCalculator) IsAbusive(tokenID int) (bool, *AbuseScoreResult, error) {
	result, err := c.CalculateScore(tokenID)
	if err != nil {
		return false, nil, err
	}
	return result.ExceedsAction, result, nil
}

// IsWarning checks if a token's behavior warrants a warning
func (c *AbuseScoreCalculator) IsWarning(tokenID int) (bool, *AbuseScoreResult, error) {
	result, err := c.CalculateScore(tokenID)
	if err != nil {
		return false, nil, err
	}
	return result.ExceedsWarning && !result.ExceedsAction, result, nil
}

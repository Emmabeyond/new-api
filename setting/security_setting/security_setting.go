package security_setting

import (
	"strings"
	"sync"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/config"
)

// Default test content patterns (one per line)
const DefaultTestContentPatterns = `hi
hello
test
ping
你好
测试`

// SecuritySetting manages all security-related configurations
// including channel masking and anti-abuse detection settings
type SecuritySetting struct {
	// Channel masking configuration
	EnableChannelMasking bool `json:"enable_channel_masking"`
	MaskChannelNames     bool `json:"mask_channel_names"`
	MaskChannelIDs       bool `json:"mask_channel_ids"`
	MaskChannelTypes     bool `json:"mask_channel_types"`

	// Anti-abuse detection configuration
	EnableAntiAbuse          bool   `json:"enable_anti_abuse"`
	ModelSwitchWindowMinutes int    `json:"model_switch_window_minutes"` // Time window in minutes
	ModelSwitchThreshold     int    `json:"model_switch_threshold"`      // Distinct model count threshold
	TestContentPatterns      string `json:"test_content_patterns"`       // Test content patterns (newline separated)
	MinContentLength         int    `json:"min_content_length"`          // Minimum content length
	TestContentThreshold     int    `json:"test_content_threshold"`      // Test content request threshold
	TestContentWindowMinutes int    `json:"test_content_window_minutes"` // Test content time window in minutes

	// Abuse score configuration
	AbuseScoreWarningThreshold int `json:"abuse_score_warning_threshold"` // Warning threshold
	AbuseScoreActionThreshold  int `json:"abuse_score_action_threshold"`  // Penalty threshold

	// Penalty configuration
	PenaltyType            string `json:"penalty_type"`              // rate_limit, temp_ban, perm_ban
	TempBanDurationMinutes int    `json:"temp_ban_duration_minutes"` // Temporary ban duration
	RateLimitRequests      int    `json:"rate_limit_requests"`       // Requests per minute when rate limited

	// Whitelist configuration
	WhitelistUserIDs string `json:"whitelist_user_ids"` // Comma separated user IDs
	WhitelistGroups  string `json:"whitelist_groups"`   // Comma separated group names
}

// Default configuration with sensible defaults
var securitySetting = SecuritySetting{
	// Channel masking defaults (enabled by default for security)
	EnableChannelMasking: true,
	MaskChannelNames:     true,
	MaskChannelIDs:       true,
	MaskChannelTypes:     true,

	// Anti-abuse defaults (disabled by default, admin must enable)
	EnableAntiAbuse:          false,
	ModelSwitchWindowMinutes: 5,
	ModelSwitchThreshold:     10,
	TestContentPatterns:      DefaultTestContentPatterns,
	MinContentLength:         10,
	TestContentThreshold:     20,
	TestContentWindowMinutes: 5,

	// Abuse score defaults
	AbuseScoreWarningThreshold: 50,
	AbuseScoreActionThreshold:  80,

	// Penalty defaults
	PenaltyType:            "rate_limit",
	TempBanDurationMinutes: 30,
	RateLimitRequests:      5,

	// Whitelist defaults (empty)
	WhitelistUserIDs: "",
	WhitelistGroups:  "",
}

var settingMutex sync.RWMutex

func init() {
	// Register with global config manager
	config.GlobalConfig.Register("security_setting", &securitySetting)

	// Register channel masking config provider to avoid circular imports
	common.SetChannelMaskingConfigProvider(func() *common.ChannelMaskingConfig {
		settingMutex.RLock()
		defer settingMutex.RUnlock()
		return &common.ChannelMaskingConfig{
			EnableChannelMasking: securitySetting.EnableChannelMasking,
			MaskChannelNames:     securitySetting.MaskChannelNames,
			MaskChannelIDs:       securitySetting.MaskChannelIDs,
			MaskChannelTypes:     securitySetting.MaskChannelTypes,
		}
	})
}

// GetSecuritySetting returns the current security setting configuration
func GetSecuritySetting() *SecuritySetting {
	settingMutex.RLock()
	defer settingMutex.RUnlock()
	return &securitySetting
}


// GetTestContentPatternList returns the test content patterns as a slice
func GetTestContentPatternList() []string {
	settingMutex.RLock()
	defer settingMutex.RUnlock()

	if securitySetting.TestContentPatterns == "" {
		return []string{}
	}

	patterns := strings.Split(securitySetting.TestContentPatterns, "\n")
	result := make([]string, 0, len(patterns))
	for _, p := range patterns {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, strings.ToLower(trimmed))
		}
	}
	return result
}

// GetWhitelistUserIDList returns the whitelist user IDs as a slice of integers
func GetWhitelistUserIDList() []int {
	settingMutex.RLock()
	defer settingMutex.RUnlock()

	if securitySetting.WhitelistUserIDs == "" {
		return []int{}
	}

	ids := strings.Split(securitySetting.WhitelistUserIDs, ",")
	result := make([]int, 0, len(ids))
	for _, id := range ids {
		trimmed := strings.TrimSpace(id)
		if trimmed != "" {
			var intID int
			if _, err := parseIntFromString(trimmed, &intID); err == nil {
				result = append(result, intID)
			}
		}
	}
	return result
}

// GetWhitelistGroupList returns the whitelist groups as a slice
func GetWhitelistGroupList() []string {
	settingMutex.RLock()
	defer settingMutex.RUnlock()

	if securitySetting.WhitelistGroups == "" {
		return []string{}
	}

	groups := strings.Split(securitySetting.WhitelistGroups, ",")
	result := make([]string, 0, len(groups))
	for _, g := range groups {
		trimmed := strings.TrimSpace(g)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// parseIntFromString is a helper to parse int from string
func parseIntFromString(s string, result *int) (bool, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return false, nil
		}
		n = n*10 + int(c-'0')
	}
	*result = n
	return true, nil
}

// IsChannelMaskingEnabled returns whether channel masking is enabled
func IsChannelMaskingEnabled() bool {
	settingMutex.RLock()
	defer settingMutex.RUnlock()
	return securitySetting.EnableChannelMasking
}

// IsMaskChannelNamesEnabled returns whether channel name masking is enabled
func IsMaskChannelNamesEnabled() bool {
	settingMutex.RLock()
	defer settingMutex.RUnlock()
	return securitySetting.MaskChannelNames
}

// IsMaskChannelIDsEnabled returns whether channel ID masking is enabled
func IsMaskChannelIDsEnabled() bool {
	settingMutex.RLock()
	defer settingMutex.RUnlock()
	return securitySetting.MaskChannelIDs
}

// IsMaskChannelTypesEnabled returns whether channel type masking is enabled
func IsMaskChannelTypesEnabled() bool {
	settingMutex.RLock()
	defer settingMutex.RUnlock()
	return securitySetting.MaskChannelTypes
}

// IsAntiAbuseEnabled returns whether anti-abuse detection is enabled
func IsAntiAbuseEnabled() bool {
	settingMutex.RLock()
	defer settingMutex.RUnlock()
	return securitySetting.EnableAntiAbuse
}

// IsUserWhitelisted checks if a user ID is in the whitelist
func IsUserWhitelisted(userID int) bool {
	for _, id := range GetWhitelistUserIDList() {
		if id == userID {
			return true
		}
	}
	return false
}

// IsGroupWhitelisted checks if a group is in the whitelist
func IsGroupWhitelisted(group string) bool {
	for _, g := range GetWhitelistGroupList() {
		if g == group {
			return true
		}
	}
	return false
}

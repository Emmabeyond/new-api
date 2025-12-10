package common

import (
	"regexp"
	"strings"
	"sync"
)

// ChannelMaskingConfig holds the configuration for channel masking
type ChannelMaskingConfig struct {
	EnableChannelMasking bool
	MaskChannelNames     bool
	MaskChannelIDs       bool
	MaskChannelTypes     bool
}

// channelMaskingConfigProvider is a function that returns the current channel masking configuration
// This is set by the security_setting package to avoid circular imports
var channelMaskingConfigProvider func() *ChannelMaskingConfig

// defaultChannelMaskingConfig returns the default configuration (all masking enabled)
func defaultChannelMaskingConfig() *ChannelMaskingConfig {
	return &ChannelMaskingConfig{
		EnableChannelMasking: true,
		MaskChannelNames:     true,
		MaskChannelIDs:       true,
		MaskChannelTypes:     true,
	}
}

// SetChannelMaskingConfigProvider sets the function that provides channel masking configuration
// This should be called by the security_setting package during initialization
func SetChannelMaskingConfigProvider(provider func() *ChannelMaskingConfig) {
	channelMaskingConfigProvider = provider
}

// getChannelMaskingConfig returns the current channel masking configuration
func getChannelMaskingConfig() *ChannelMaskingConfig {
	if channelMaskingConfigProvider != nil {
		return channelMaskingConfigProvider()
	}
	return defaultChannelMaskingConfig()
}

// Generic replacement messages for masked channel information
const (
	GenericChannelMessage = "upstream service"
	GenericServiceMessage = "An error occurred with the upstream service"
	SafeDefaultMessage    = "An error occurred while processing your request. Please try again later."
)

// channelMaskingPatterns holds compiled regex patterns for detecting channel information
var channelMaskingPatterns struct {
	once sync.Once

	// ChannelIDPattern matches channel IDs in various formats:
	// - #123, channel 456, 渠道 789, 渠道（#123）, 渠道(123)
	ChannelIDPattern *regexp.Regexp

	// ChannelNamePattern matches channel names in various formats:
	// - 通道「xxx」, channel 'xxx', channel "xxx", 渠道「xxx」
	ChannelNamePattern *regexp.Regexp

	// ChannelTypePattern matches channel type references:
	// - channel type 1, 渠道类型 2, type: 3
	ChannelTypePattern *regexp.Regexp

	// RetrySequencePattern matches retry sequences:
	// - 1->2->3, [1, 2, 3], 1 -> 2 -> 3
	RetrySequencePattern *regexp.Regexp

	// ChannelKeywordPattern matches common channel-related keywords with identifiers
	// - 渠道 xxx, channel xxx (followed by name or ID)
	ChannelKeywordPattern *regexp.Regexp

	// initErr stores any error during pattern compilation
	initErr error
}

// initPatterns compiles all regex patterns once
func initPatterns() {
	channelMaskingPatterns.once.Do(func() {
		var err error

		// Channel ID patterns - matches various formats of channel IDs
		// Examples: #123, channel 456, 渠道 789, 渠道（#123）, 渠道(123), 渠道 #123
		channelMaskingPatterns.ChannelIDPattern, err = regexp.Compile(
			`(?i)(?:` +
				`#\d+|` + // #123
				`channel\s+\d+|` + // channel 456
				`渠道\s*[（(]?\s*#?\d+\s*[）)]?|` + // 渠道 789, 渠道（#123）, 渠道(123)
				`通道\s*[（(]?\s*#?\d+\s*[）)]?|` + // 通道 789
				`(?:channel|渠道|通道)\s*(?:id|ID|Id)?\s*[:：]?\s*\d+` + // channel id: 123
				`)`)
		if err != nil {
			channelMaskingPatterns.initErr = err
			return
		}

		// Channel name patterns - matches channel names in quotes or brackets
		// Examples: 通道「OpenAI-Primary」, channel 'Claude-Backup', channel "test"
		channelMaskingPatterns.ChannelNamePattern, err = regexp.Compile(
			`(?i)(?:` +
				`(?:通道|渠道|channel)\s*[「『"']([^」』"']+)[」』"']|` + // 通道「xxx」, channel "xxx"
				`(?:通道|渠道|channel)\s*[:：]\s*([^\s,，。.]+)` + // 通道: xxx
				`)`)
		if err != nil {
			channelMaskingPatterns.initErr = err
			return
		}

		// Channel type patterns
		// Examples: channel type 1, 渠道类型 2, type: 3
		channelMaskingPatterns.ChannelTypePattern, err = regexp.Compile(
			`(?i)(?:` +
				`(?:channel\s+)?type\s*[:：]?\s*\d+|` + // channel type 1, type: 2
				`渠道类型\s*[:：]?\s*\d+|` + // 渠道类型 1
				`通道类型\s*[:：]?\s*\d+` + // 通道类型 1
				`)`)
		if err != nil {
			channelMaskingPatterns.initErr = err
			return
		}

		// Retry sequence patterns - matches channel retry sequences
		// Examples: 1->2->3, [1, 2, 3], 1 -> 2 -> 3
		channelMaskingPatterns.RetrySequencePattern, err = regexp.Compile(
			`(?:` +
				`\d+(?:\s*->\s*\d+)+|` + // 1->2->3, 1 -> 2 -> 3
				`\[\s*\d+(?:\s*,\s*\d+)+\s*\]` + // [1, 2, 3]
				`)`)
		if err != nil {
			channelMaskingPatterns.initErr = err
			return
		}

		// Channel keyword patterns - matches channel references with names
		// Examples: 渠道 xxx 已被禁用, channel xxx is disabled
		channelMaskingPatterns.ChannelKeywordPattern, err = regexp.Compile(
			`(?i)(?:` +
				`(?:渠道|通道|channel)\s+([^\s,，。.（(]+)(?:\s*(?:已被|is|was|has been))` + // 渠道 xxx 已被禁用
				`)`)
		if err != nil {
			channelMaskingPatterns.initErr = err
			return
		}
	})
}

// MaskChannelInfo masks channel-related sensitive information in error messages
// for end-user consumption while preserving the original error structure.
// It masks channel IDs, names, types, and retry sequences based on configuration.
// If masking fails, it returns a safe default message.
func MaskChannelInfo(message string) string {
	if message == "" {
		return message
	}

	// Get settings from configuration provider
	config := getChannelMaskingConfig()

	// Check if masking is enabled
	if !config.EnableChannelMasking {
		// Only apply general sensitive info masking when channel masking is disabled
		return MaskSensitiveInfo(message)
	}

	// Initialize patterns if not already done
	initPatterns()

	// If pattern compilation failed, return safe default message
	if channelMaskingPatterns.initErr != nil {
		SysError("Channel masking pattern compilation failed: " + channelMaskingPatterns.initErr.Error())
		return SafeDefaultMessage
	}

	// Use defer/recover to handle any unexpected panics during masking
	defer func() {
		if r := recover(); r != nil {
			SysError("Panic during channel masking, returning safe default")
		}
	}()

	result := message

	// Mask channel IDs (if enabled)
	if config.MaskChannelIDs && channelMaskingPatterns.ChannelIDPattern != nil {
		result = channelMaskingPatterns.ChannelIDPattern.ReplaceAllString(result, GenericChannelMessage)
	}

	// Mask channel names (if enabled)
	if config.MaskChannelNames && channelMaskingPatterns.ChannelNamePattern != nil {
		result = channelMaskingPatterns.ChannelNamePattern.ReplaceAllString(result, GenericChannelMessage)
	}

	// Mask channel types (if enabled)
	if config.MaskChannelTypes && channelMaskingPatterns.ChannelTypePattern != nil {
		result = channelMaskingPatterns.ChannelTypePattern.ReplaceAllString(result, GenericChannelMessage)
	}

	// Mask retry sequences (always mask these as they reveal channel selection)
	if channelMaskingPatterns.RetrySequencePattern != nil {
		result = channelMaskingPatterns.RetrySequencePattern.ReplaceAllString(result, "[masked]")
	}

	// Mask channel keyword patterns
	if channelMaskingPatterns.ChannelKeywordPattern != nil {
		result = channelMaskingPatterns.ChannelKeywordPattern.ReplaceAllStringFunc(result, func(match string) string {
			// Replace the channel name part while keeping the action part
			if strings.Contains(match, "已被") {
				return GenericChannelMessage + " 已被"
			}
			if strings.Contains(strings.ToLower(match), "is") {
				return GenericChannelMessage + " is"
			}
			if strings.Contains(strings.ToLower(match), "was") {
				return GenericChannelMessage + " was"
			}
			if strings.Contains(strings.ToLower(match), "has been") {
				return GenericChannelMessage + " has been"
			}
			return GenericChannelMessage
		})
	}

	// Apply general sensitive info masking (URLs, IPs, etc.)
	result = MaskSensitiveInfo(result)

	return result
}


// MaskChannelError creates a user-safe error message from a channel error.
// It takes the original error message and returns a generic error message
// without exposing channel details like ID, name, or type.
func MaskChannelError(channelID int, channelName string, channelType int, originalError string) string {
	if originalError == "" {
		return GenericServiceMessage
	}

	// First apply channel info masking
	maskedMessage := MaskChannelInfo(originalError)

	// If the message still looks like it might contain channel info after masking,
	// return a safe generic message
	if IsChannelInfoPresent(maskedMessage) {
		return GenericServiceMessage
	}

	return maskedMessage
}

// IsChannelInfoPresent detects if a string contains channel-identifying information.
// Returns true if channel ID, name, or type patterns are found.
// This is useful for validating that masking was successful.
func IsChannelInfoPresent(message string) bool {
	if message == "" {
		return false
	}

	// Initialize patterns if not already done
	initPatterns()

	// If pattern compilation failed, assume channel info might be present (safe default)
	if channelMaskingPatterns.initErr != nil {
		return true
	}

	// Check for channel ID patterns
	if channelMaskingPatterns.ChannelIDPattern != nil &&
		channelMaskingPatterns.ChannelIDPattern.MatchString(message) {
		return true
	}

	// Check for channel name patterns
	if channelMaskingPatterns.ChannelNamePattern != nil &&
		channelMaskingPatterns.ChannelNamePattern.MatchString(message) {
		return true
	}

	// Check for channel type patterns
	if channelMaskingPatterns.ChannelTypePattern != nil &&
		channelMaskingPatterns.ChannelTypePattern.MatchString(message) {
		return true
	}

	// Check for retry sequence patterns
	if channelMaskingPatterns.RetrySequencePattern != nil &&
		channelMaskingPatterns.RetrySequencePattern.MatchString(message) {
		return true
	}

	// Check for channel keyword patterns
	if channelMaskingPatterns.ChannelKeywordPattern != nil &&
		channelMaskingPatterns.ChannelKeywordPattern.MatchString(message) {
		return true
	}

	return false
}

// MaskChannelInfoInUpstreamError masks channel-identifying information in upstream error messages.
// This is specifically designed for errors returned from upstream services that might
// contain channel information embedded in the error response.
func MaskChannelInfoInUpstreamError(upstreamError string) string {
	if upstreamError == "" {
		return upstreamError
	}

	// Apply channel masking
	masked := MaskChannelInfo(upstreamError)

	// Double-check that no channel info remains
	if IsChannelInfoPresent(masked) {
		return GenericServiceMessage
	}

	return masked
}

// GetSafeErrorMessage returns a safe error message for user-facing responses.
// If the original message contains channel information, it returns a generic message.
// Otherwise, it returns the masked version of the original message.
func GetSafeErrorMessage(originalMessage string) string {
	if originalMessage == "" {
		return GenericServiceMessage
	}

	masked := MaskChannelInfo(originalMessage)

	// Final safety check
	if IsChannelInfoPresent(masked) {
		return GenericServiceMessage
	}

	return masked
}

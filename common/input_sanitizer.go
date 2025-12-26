package common

import (
	"html"
	"regexp"
	"strings"
	"unicode/utf8"
)

// InputValidationConfig 输入验证配置
type InputValidationConfig struct {
	MaxLength     int  // 最大长度限制
	AllowEmpty    bool // 是否允许空值
	TrimSpace     bool // 是否去除首尾空格
	AllowNewlines bool // 是否允许换行符
}

// DefaultInputValidationConfig 默认输入验证配置
var DefaultInputValidationConfig = InputValidationConfig{
	MaxLength:     10000,
	AllowEmpty:    true,
	TrimSpace:     true,
	AllowNewlines: true,
}

// XSS 攻击特征正则表达式
var xssPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)<script[^>]*>`),
	regexp.MustCompile(`(?i)</script>`),
	regexp.MustCompile(`(?i)javascript:`),
	regexp.MustCompile(`(?i)on\w+\s*=`),
	regexp.MustCompile(`(?i)<iframe[^>]*>`),
	regexp.MustCompile(`(?i)<object[^>]*>`),
	regexp.MustCompile(`(?i)<embed[^>]*>`),
	regexp.MustCompile(`(?i)<svg[^>]*onload`),
	regexp.MustCompile(`(?i)expression\s*\(`),
	regexp.MustCompile(`(?i)vbscript:`),
	regexp.MustCompile(`(?i)data:text/html`),
	regexp.MustCompile(`(?i)<img[^>]*onerror`),
	regexp.MustCompile(`(?i)<body[^>]*onload`),
	regexp.MustCompile(`(?i)<input[^>]*onfocus`),
	regexp.MustCompile(`(?i)<marquee[^>]*onstart`),
	regexp.MustCompile(`(?i)<video[^>]*onerror`),
	regexp.MustCompile(`(?i)<audio[^>]*onerror`),
	regexp.MustCompile(`(?i)<details[^>]*ontoggle`),
	regexp.MustCompile(`(?i)&#x?[0-9a-f]+;`), // HTML 实体编码
}

// XSS 关键字列表
var xssKeywords = []string{
	"<script",
	"javascript:",
	"onerror=",
	"onload=",
	"onclick=",
	"onmouseover=",
	"onfocus=",
	"onblur=",
	"onmouseout=",
	"onmouseenter=",
	"onmouseleave=",
	"onkeydown=",
	"onkeyup=",
	"onkeypress=",
	"onsubmit=",
	"onreset=",
	"onchange=",
	"oninput=",
	"ondblclick=",
	"oncontextmenu=",
	"ondrag=",
	"ondrop=",
	"oncopy=",
	"oncut=",
	"onpaste=",
	"<iframe",
	"<object",
	"<embed",
	"<svg",
	"<math",
	"<base",
	"<link",
	"<meta",
	"<style",
	"<form",
	"vbscript:",
	"data:text/html",
	"expression(",
}

// SanitizeUserInput 对用户输入进行 HTML 转义
// 使用默认配置进行验证和转义
func SanitizeUserInput(input string) string {
	return SanitizeUserInputWithConfig(input, DefaultInputValidationConfig)
}

// SanitizeUserInputWithConfig 使用指定配置对用户输入进行 HTML 转义
func SanitizeUserInputWithConfig(input string, config InputValidationConfig) string {
	if input == "" {
		return ""
	}

	// 去除首尾空格
	if config.TrimSpace {
		input = strings.TrimSpace(input)
	}

	// 检查空值
	if input == "" && !config.AllowEmpty {
		return ""
	}

	// 检查长度限制
	if config.MaxLength > 0 && utf8.RuneCountInString(input) > config.MaxLength {
		// 截断到最大长度
		runes := []rune(input)
		input = string(runes[:config.MaxLength])
	}

	// 移除换行符（如果不允许）
	if !config.AllowNewlines {
		input = strings.ReplaceAll(input, "\n", " ")
		input = strings.ReplaceAll(input, "\r", " ")
	}

	// HTML 转义
	return html.EscapeString(input)
}

// SanitizeErrorMessage 对错误信息进行转义
// 确保错误信息中的用户输入不会导致 XSS
func SanitizeErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	return html.EscapeString(err.Error())
}

// SanitizeErrorMessageString 对错误信息字符串进行转义
func SanitizeErrorMessageString(errMsg string) string {
	if errMsg == "" {
		return ""
	}
	return html.EscapeString(errMsg)
}

// DetectXSSAttempt 检测潜在的 XSS 攻击
// 返回 true 表示检测到 XSS 攻击特征
func DetectXSSAttempt(input string) bool {
	if input == "" {
		return false
	}

	// 检查所有 XSS 正则模式
	for _, pattern := range xssPatterns {
		if pattern.MatchString(input) {
			return true
		}
	}

	// 检查常见的 XSS 关键字（大小写不敏感）
	lowerInput := strings.ToLower(input)
	for _, keyword := range xssKeywords {
		if strings.Contains(lowerInput, keyword) {
			return true
		}
	}

	return false
}

// DetectAndLogXSSAttempt 检测 XSS 攻击并记录安全日志
// 返回 true 表示检测到 XSS 攻击特征
func DetectAndLogXSSAttempt(input string, userId int, ipAddress string, userAgent string, requestPath string) bool {
	if !DetectXSSAttempt(input) {
		return false
	}

	// 记录安全事件
	event := &SecurityEvent{
		UserId:      userId,
		EventType:   "xss_attempt",
		Details:     truncateString(input, 500), // 截断过长的输入
		IPAddress:   ipAddress,
		UserAgent:   truncateString(userAgent, 256),
		RequestPath: requestPath,
	}
	LogSecurityEvent(event)

	return true
}

// ValidateInputLength 验证输入长度
// 返回 true 表示长度有效
func ValidateInputLength(input string, minLength, maxLength int) bool {
	length := utf8.RuneCountInString(input)
	if minLength > 0 && length < minLength {
		return false
	}
	if maxLength > 0 && length > maxLength {
		return false
	}
	return true
}

// ValidateInputFormat 验证输入格式
// 使用正则表达式验证输入是否符合指定格式
func ValidateInputFormat(input string, pattern *regexp.Regexp) bool {
	if pattern == nil {
		return true
	}
	return pattern.MatchString(input)
}

// StripHTMLTags 移除所有 HTML 标签，只保留纯文本
func StripHTMLTags(input string) string {
	if input == "" {
		return ""
	}

	// 移除所有 HTML 标签
	re := regexp.MustCompile(`<[^>]*>`)
	result := re.ReplaceAllString(input, "")

	// 解码 HTML 实体
	result = html.UnescapeString(result)

	return strings.TrimSpace(result)
}

// truncateString 截断字符串到指定长度
func truncateString(s string, maxLen int) string {
	if maxLen <= 0 {
		return s
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

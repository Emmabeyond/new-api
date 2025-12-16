package common

import (
	"regexp"
	"strings"
	"sync"
)

// HeaderPatternMatcher 响应头模式匹配器
// 支持通配符模式匹配，用于过滤敏感响应头
type HeaderPatternMatcher struct {
	mu              sync.RWMutex
	compiledPatterns map[string]*regexp.Regexp
}

var (
	globalHeaderMatcher     *HeaderPatternMatcher
	headerMatcherOnce       sync.Once
)

// GetHeaderPatternMatcher 获取全局头模式匹配器实例
func GetHeaderPatternMatcher() *HeaderPatternMatcher {
	headerMatcherOnce.Do(func() {
		globalHeaderMatcher = &HeaderPatternMatcher{
			compiledPatterns: make(map[string]*regexp.Regexp),
		}
	})
	return globalHeaderMatcher
}

// Match 检查头名称是否匹配模式（大小写不敏感）
// 支持 * 通配符，如 "X-Ratelimit-*" 匹配 "X-Ratelimit-Limit"
func (m *HeaderPatternMatcher) Match(headerName string, pattern string) bool {
	if headerName == "" || pattern == "" {
		return false
	}

	// 转换为小写进行比较
	headerLower := strings.ToLower(headerName)
	patternLower := strings.ToLower(pattern)

	// 精确匹配（无通配符）
	if !strings.Contains(patternLower, "*") {
		return headerLower == patternLower
	}

	// 使用编译后的正则表达式匹配
	re, err := m.getOrCompilePattern(patternLower)
	if err != nil {
		return false
	}

	return re.MatchString(headerLower)
}

// MatchAny 检查头名称是否匹配任一模式
func (m *HeaderPatternMatcher) MatchAny(headerName string, patterns []string) bool {
	for _, pattern := range patterns {
		if m.Match(headerName, pattern) {
			return true
		}
	}
	return false
}

// getOrCompilePattern 获取或编译模式的正则表达式
func (m *HeaderPatternMatcher) getOrCompilePattern(pattern string) (*regexp.Regexp, error) {
	m.mu.RLock()
	if re, ok := m.compiledPatterns[pattern]; ok {
		m.mu.RUnlock()
		return re, nil
	}
	m.mu.RUnlock()

	// 编译新模式
	m.mu.Lock()
	defer m.mu.Unlock()

	// 双重检查
	if re, ok := m.compiledPatterns[pattern]; ok {
		return re, nil
	}

	re, err := compileWildcardPattern(pattern)
	if err != nil {
		return nil, err
	}

	m.compiledPatterns[pattern] = re
	return re, nil
}

// compileWildcardPattern 将通配符模式转换为正则表达式
// * 匹配任意字符（包括空字符串）
func compileWildcardPattern(pattern string) (*regexp.Regexp, error) {
	// 转义正则特殊字符，但保留 *
	escaped := regexp.QuoteMeta(pattern)
	// 将转义后的 \* 替换为 .*
	regexPattern := strings.ReplaceAll(escaped, `\*`, `.*`)
	// 添加锚点确保完整匹配
	regexPattern = "^" + regexPattern + "$"

	return regexp.Compile(regexPattern)
}

// CompilePatterns 预编译多个模式
func (m *HeaderPatternMatcher) CompilePatterns(patterns []string) error {
	for _, pattern := range patterns {
		patternLower := strings.ToLower(pattern)
		if strings.Contains(patternLower, "*") {
			_, err := m.getOrCompilePattern(patternLower)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ClearCache 清除编译缓存
func (m *HeaderPatternMatcher) ClearCache() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.compiledPatterns = make(map[string]*regexp.Regexp)
}

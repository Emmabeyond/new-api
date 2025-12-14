package ratio_setting

import (
	"strings"
)

// MatchWildcard 在 ratioMap 中查找最佳匹配的通配符配置
// 通配符规则：只支持末尾的 * 符号，如 "gpt-4-*"
// 返回值：匹配的 pattern, 对应的 ratio, 是否找到匹配
// 如果有多个匹配，返回最长前缀匹配（最具体的匹配）
func MatchWildcard(name string, ratioMap map[string]float64) (string, float64, bool) {
	var bestPattern string
	var bestRatio float64
	var bestPrefixLen int = -1

	for pattern, ratio := range ratioMap {
		// 检查是否是通配符模式（以 * 结尾）
		if !strings.HasSuffix(pattern, "*") {
			continue
		}

		// 获取通配符前缀（去掉末尾的 *）
		prefix := strings.TrimSuffix(pattern, "*")

		// 检查模型名是否以该前缀开头
		if strings.HasPrefix(name, prefix) {
			// 选择最长前缀匹配
			if len(prefix) > bestPrefixLen {
				bestPrefixLen = len(prefix)
				bestPattern = pattern
				bestRatio = ratio
			}
		}
	}

	if bestPrefixLen >= 0 {
		return bestPattern, bestRatio, true
	}

	return "", 0, false
}

// IsWildcardPattern 检查给定的模式是否是通配符模式
func IsWildcardPattern(pattern string) bool {
	return strings.HasSuffix(pattern, "*")
}

// GetWildcardPrefix 获取通配符模式的前缀部分
// 如果不是通配符模式，返回原字符串
func GetWildcardPrefix(pattern string) string {
	if strings.HasSuffix(pattern, "*") {
		return strings.TrimSuffix(pattern, "*")
	}
	return pattern
}

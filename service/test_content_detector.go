package service

import (
	"strings"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
)

// TestContentDetector 测试内容检测器接口
type TestContentDetector interface {
	// IsTestContent 检查内容是否为测试内容
	IsTestContent(content string) bool
	// RecordTestContent 记录测试内容请求
	RecordTestContent(tokenID int) error
	// GetTestContentCount 获取测试内容计数
	// windowMinutes: 时间窗口（分钟）
	GetTestContentCount(tokenID int, windowMinutes int) (int, error)
}

// testContentDetector 测试内容检测器实现
type testContentDetector struct {
	// tokenID -> count
	counts sync.Map
	// tokenID -> lastReset
	lastReset sync.Map
}

var defaultTestContentDetector *testContentDetector

// GetTestContentDetector 获取测试内容检测器单例
func GetTestContentDetector() TestContentDetector {
	if defaultTestContentDetector == nil {
		defaultTestContentDetector = &testContentDetector{}
	}
	return defaultTestContentDetector
}

// 常见的测试内容关键词
var testKeywords = []string{
	"hello", "hi", "test", "testing", "ping", "echo",
	"你好", "测试", "1234", "asdf",
}

// IsTestContent 检查内容是否为测试内容
func (d *testContentDetector) IsTestContent(content string) bool {
	if content == "" {
		return false
	}

	// 内容太短可能是测试
	if len(content) < 10 {
		return true
	}

	// 检查是否包含测试关键词
	lowerContent := strings.ToLower(content)
	for _, keyword := range testKeywords {
		if strings.Contains(lowerContent, keyword) && len(content) < 50 {
			return true
		}
	}

	return false
}


// RecordTestContent 记录测试内容请求
func (d *testContentDetector) RecordTestContent(tokenID int) error {
	// 检查是否需要重置计数（每小时重置）
	now := time.Now()
	if lastReset, ok := d.lastReset.Load(tokenID); ok {
		if now.Sub(lastReset.(time.Time)) > time.Hour {
			d.counts.Store(tokenID, 0)
			d.lastReset.Store(tokenID, now)
		}
	} else {
		d.lastReset.Store(tokenID, now)
	}

	// 增加计数
	count := 0
	if val, ok := d.counts.Load(tokenID); ok {
		count = val.(int)
	}
	d.counts.Store(tokenID, count+1)

	// 记录日志（如果计数较高）
	if count > 10 {
		common.SysLog("Warning: Token has high test content count: " + string(rune(tokenID)))
	}

	return nil
}

// GetTestContentCount 获取测试内容计数
// windowMinutes: 时间窗口（分钟），用于过滤过期的计数
func (d *testContentDetector) GetTestContentCount(tokenID int, windowMinutes int) (int, error) {
	// 检查是否在时间窗口内
	if lastReset, ok := d.lastReset.Load(tokenID); ok {
		if time.Since(lastReset.(time.Time)) > time.Duration(windowMinutes)*time.Minute {
			// 超出时间窗口，返回 0
			return 0, nil
		}
	}

	if val, ok := d.counts.Load(tokenID); ok {
		return val.(int), nil
	}
	return 0, nil
}

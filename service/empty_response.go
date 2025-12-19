package service

import (
	"strings"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/setting/operation_setting"
)

func init() {
	// 注册配置设置回调，避免循环导入
	operation_setting.RegisterEmptyResponseConfigSetter(func(enabled bool, maxRetryCount int, excludedModels []string, alertThreshold float64, nonEmptyFinishReasons []string) {
		// 如果未配置非空 finish_reason 列表，使用默认值
		if len(nonEmptyFinishReasons) == 0 {
			nonEmptyFinishReasons = defaultNonEmptyFinishReasons
		}
		SetEmptyResponseConfig(&EmptyResponseConfig{
			Enabled:               enabled,
			MaxRetryCount:         maxRetryCount,
			ExcludedModels:        excludedModels,
			AlertThreshold:        alertThreshold,
			NonEmptyFinishReasons: nonEmptyFinishReasons,
		})
	})
}

// EmptyResponseConfig 空回复处理配置
type EmptyResponseConfig struct {
	// Enabled 是否启用空回复检测和重试
	Enabled bool `json:"enabled"`

	// MaxRetryCount 最大重试次数（针对空回复），默认为 2
	MaxRetryCount int `json:"max_retry_count"`

	// ExcludedModels 排除检测的模型列表（支持前缀匹配）
	ExcludedModels []string `json:"excluded_models"`

	// AlertThreshold 告警阈值（空回复率百分比），默认为 10
	AlertThreshold float64 `json:"alert_threshold"`

	// NonEmptyFinishReasons 非空 finish_reason 列表
	// 当 finish_reason 在此列表中时，即使内容为空也不判定为空回复
	// 用于支持 Agentic 场景（如 Claude Code、Codex）中的工具调用响应
	NonEmptyFinishReasons []string `json:"non_empty_finish_reasons"`
}

// EmptyResponseEvent 空回复事件
type EmptyResponseEvent struct {
	ChannelID   int       `json:"channel_id"`
	ChannelName string    `json:"channel_name"`
	ModelName   string    `json:"model_name"`
	Timestamp   time.Time `json:"timestamp"`
	RetryCount  int       `json:"retry_count"`
	FinalResult string    `json:"final_result"` // "success", "failed"
	RequestID   string    `json:"request_id"`
}

// EmptyResponseStats 空回复统计数据
type EmptyResponseStats struct {
	ChannelID         int     `json:"channel_id"`
	ModelName         string  `json:"model_name"`
	TotalRequests     int64   `json:"total_requests"`
	EmptyResponses    int64   `json:"empty_responses"`
	EmptyRate         float64 `json:"empty_rate"`
	AvgRetryCount     float64 `json:"avg_retry_count"`
	SuccessAfterRetry int64   `json:"success_after_retry"`
}

// 默认非空 finish_reason 列表
// 这些 finish_reason 表示工具调用，即使内容为空也不应判定为空回复
var defaultNonEmptyFinishReasons = []string{
	"tool_calls",    // OpenAI 格式的工具调用
	"tool_use",      // Claude 格式的工具使用
	"function_call", // OpenAI 旧版函数调用
}

// 默认配置
var defaultEmptyResponseConfig = &EmptyResponseConfig{
	Enabled:               true,
	MaxRetryCount:         2,
	ExcludedModels:        []string{},
	AlertThreshold:        10.0,
	NonEmptyFinishReasons: defaultNonEmptyFinishReasons,
}

var (
	emptyResponseConfig     *EmptyResponseConfig
	emptyResponseConfigLock sync.RWMutex
)

// GetEmptyResponseConfig 获取空回复配置
func GetEmptyResponseConfig() *EmptyResponseConfig {
	emptyResponseConfigLock.RLock()
	defer emptyResponseConfigLock.RUnlock()

	if emptyResponseConfig == nil {
		return defaultEmptyResponseConfig
	}
	return emptyResponseConfig
}

// SetEmptyResponseConfig 设置空回复配置
func SetEmptyResponseConfig(config *EmptyResponseConfig) {
	emptyResponseConfigLock.Lock()
	defer emptyResponseConfigLock.Unlock()
	emptyResponseConfig = config
}

// IsModelExcluded 检查模型是否在排除列表中
func (c *EmptyResponseConfig) IsModelExcluded(modelName string) bool {
	if c == nil || len(c.ExcludedModels) == 0 {
		return false
	}
	modelLower := strings.ToLower(modelName)
	for _, excluded := range c.ExcludedModels {
		excludedLower := strings.ToLower(excluded)
		// 支持前缀匹配
		if strings.HasPrefix(modelLower, excludedLower) || modelLower == excludedLower {
			return true
		}
	}
	return false
}

// IsNonEmptyFinishReason 检查 finish_reason 是否表示非空响应
// 当 finish_reason 在非空列表中时（如 tool_calls、tool_use、function_call），返回 true
// 这用于支持 Agentic 场景，避免将工具调用响应误判为空回复
func IsNonEmptyFinishReason(finishReason string) bool {
	if finishReason == "" {
		return false
	}
	config := GetEmptyResponseConfig()
	if config == nil {
		return false
	}
	// 如果空回复检测被禁用，则不进行 finish_reason 判断
	// 这样可以保持行为一致性：禁用后完全跳过空回复相关逻辑
	if !config.Enabled {
		return false
	}
	finishReasonLower := strings.ToLower(finishReason)
	for _, nonEmpty := range config.NonEmptyFinishReasons {
		if strings.ToLower(nonEmpty) == finishReasonLower {
			return true
		}
	}
	return false
}

// IsEmptyStreamResponse 检测流式响应是否为空
// 当 completionTokens 为 0 且 responseText 为空或仅包含空白字符时，返回 true
func IsEmptyStreamResponse(usage *dto.Usage, responseText string) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	// 检查 usage 是否为空
	if usage == nil {
		return true
	}

	// 如果有 completion tokens，说明有内容
	if usage.CompletionTokens > 0 {
		return false
	}

	// 检查响应文本是否为空或仅包含空白字符
	trimmedText := strings.TrimSpace(responseText)
	return trimmedText == ""
}

// IsEmptyStreamResponseWithFinishReason 带 finish_reason 的流式响应空检测
// 先检查 finish_reason 是否表示非空响应（如工具调用），再进行内容检测
func IsEmptyStreamResponseWithFinishReason(usage *dto.Usage, responseText string, finishReason string) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	// 先检查 finish_reason 是否表示非空响应
	if IsNonEmptyFinishReason(finishReason) {
		return false
	}

	// 回退到内容检测
	return IsEmptyStreamResponse(usage, responseText)
}

// IsEmptyNonStreamResponse 检测非流式响应是否为空
// 当 choices 为空，或所有 choices 的内容为空且没有 tool calls 时，返回 true
func IsEmptyNonStreamResponse(response *dto.OpenAITextResponse) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	if response == nil {
		return true
	}

	// 检查 choices 是否为空
	if len(response.Choices) == 0 {
		return true
	}

	// 检查所有 choices 是否都为空内容且没有 tool calls
	for _, choice := range response.Choices {
		// 如果有 tool calls，不算空回复
		if len(choice.Message.ToolCalls) > 0 {
			return false
		}

		// 检查内容是否非空
		content := choice.Message.StringContent()
		if strings.TrimSpace(content) != "" {
			return false
		}

		// 检查 reasoning content
		if strings.TrimSpace(choice.Message.ReasoningContent) != "" {
			return false
		}
		if strings.TrimSpace(choice.Message.Reasoning) != "" {
			return false
		}
	}

	// 所有 choices 都是空内容且没有 tool calls
	return true
}

// IsEmptyNonStreamResponseWithFinishReason 带 finish_reason 的非流式响应空检测
// 先检查 finish_reason 是否表示非空响应（如工具调用），再进行内容检测
func IsEmptyNonStreamResponseWithFinishReason(response *dto.OpenAITextResponse, finishReason string) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	// 先检查 finish_reason 是否表示非空响应
	if IsNonEmptyFinishReason(finishReason) {
		return false
	}

	// 回退到内容检测
	return IsEmptyNonStreamResponse(response)
}

// IsEmptyStreamResponseWithModel 带模型检查的流式响应空检测
func IsEmptyStreamResponseWithModel(usage *dto.Usage, responseText string, modelName string) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	// 检查模型是否在排除列表
	if config.IsModelExcluded(modelName) {
		return false
	}

	return IsEmptyStreamResponse(usage, responseText)
}

// IsEmptyNonStreamResponseWithModel 带模型检查的非流式响应空检测
func IsEmptyNonStreamResponseWithModel(response *dto.OpenAITextResponse, modelName string) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	// 检查模型是否在排除列表
	if config.IsModelExcluded(modelName) {
		return false
	}

	return IsEmptyNonStreamResponse(response)
}


// EmptyResponseEventRecorder 空回复事件记录器
type EmptyResponseEventRecorder struct {
	mu     sync.Mutex
	events []EmptyResponseEvent
}

var (
	emptyResponseRecorder     *EmptyResponseEventRecorder
	emptyResponseRecorderOnce sync.Once
)

// GetEmptyResponseRecorder 获取空回复事件记录器单例
func GetEmptyResponseRecorder() *EmptyResponseEventRecorder {
	emptyResponseRecorderOnce.Do(func() {
		emptyResponseRecorder = &EmptyResponseEventRecorder{
			events: make([]EmptyResponseEvent, 0, 1000),
		}
	})
	return emptyResponseRecorder
}

// Record 记录空回复事件
func (r *EmptyResponseEventRecorder) Record(event *EmptyResponseEvent) {
	if event == nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// 保持最近 1000 条记录
	if len(r.events) >= 1000 {
		r.events = r.events[1:]
	}
	r.events = append(r.events, *event)
}

// GetRecentEvents 获取最近的事件
func (r *EmptyResponseEventRecorder) GetRecentEvents(count int) []EmptyResponseEvent {
	r.mu.Lock()
	defer r.mu.Unlock()

	if count <= 0 || count > len(r.events) {
		count = len(r.events)
	}

	start := len(r.events) - count
	if start < 0 {
		start = 0
	}

	result := make([]EmptyResponseEvent, count)
	copy(result, r.events[start:])
	return result
}

// GetStatsByChannel 按渠道获取统计数据
func (r *EmptyResponseEventRecorder) GetStatsByChannel(channelID int, duration time.Duration) *EmptyResponseStats {
	r.mu.Lock()
	defer r.mu.Unlock()

	cutoff := time.Now().Add(-duration)
	stats := &EmptyResponseStats{
		ChannelID: channelID,
	}

	var totalRetries int64
	var successCount int64

	for _, event := range r.events {
		if event.ChannelID != channelID || event.Timestamp.Before(cutoff) {
			continue
		}

		stats.EmptyResponses++
		totalRetries += int64(event.RetryCount)

		if event.FinalResult == "success" {
			successCount++
		}
	}

	if stats.EmptyResponses > 0 {
		stats.AvgRetryCount = float64(totalRetries) / float64(stats.EmptyResponses)
		stats.SuccessAfterRetry = successCount
	}

	return stats
}

// RecordEmptyResponse 记录空回复事件的便捷函数
func RecordEmptyResponse(channelID int, channelName, modelName, requestID string, retryCount int, success bool) {
	result := "failed"
	if success {
		result = "success"
	}

	event := &EmptyResponseEvent{
		ChannelID:   channelID,
		ChannelName: channelName,
		ModelName:   modelName,
		Timestamp:   time.Now(),
		RetryCount:  retryCount,
		FinalResult: result,
		RequestID:   requestID,
	}

	GetEmptyResponseRecorder().Record(event)
}

// CheckAndAlertEmptyResponseRate 检查空回复率并告警
func CheckAndAlertEmptyResponseRate(channelID int, channelName string, totalRequests, emptyResponses int64) {
	if totalRequests == 0 {
		return
	}

	config := GetEmptyResponseConfig()
	if config == nil {
		return
	}

	rate := float64(emptyResponses) / float64(totalRequests) * 100
	if rate > config.AlertThreshold {
		// 记录警告日志（这里使用 common.SysLog，实际项目中可能需要更完善的告警机制）
		// common.SysLog(fmt.Sprintf("警告：渠道 %s (#%d) 空回复率 %.2f%% 超过阈值 %.2f%%", channelName, channelID, rate, config.AlertThreshold))
	}
}

// IsEmptyClaudeResponse 检测 Claude 非流式响应是否为空
// 当 Content 为空，或所有 Content 的文本内容为空且没有 tool_use 时，返回 true
func IsEmptyClaudeResponse(response *dto.ClaudeResponse) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	if response == nil {
		return true
	}

	// 检查 Content 是否为空
	if len(response.Content) == 0 {
		// 对于 Completion 模式，检查 Completion 字段
		if strings.TrimSpace(response.Completion) != "" {
			return false
		}
		return true
	}

	// 检查所有 Content 是否都为空内容且没有 tool_use
	for _, content := range response.Content {
		// 如果有 tool_use，不算空回复
		if content.Type == "tool_use" {
			return false
		}

		// 检查文本内容是否非空
		if content.Type == "text" && strings.TrimSpace(content.GetText()) != "" {
			return false
		}

		// 检查 thinking 内容
		if content.Type == "thinking" && content.Thinking != nil && strings.TrimSpace(*content.Thinking) != "" {
			return false
		}
	}

	// 所有 Content 都是空内容且没有 tool_use
	return true
}

// IsEmptyClaudeResponseWithModel 带模型检查的 Claude 响应空检测
func IsEmptyClaudeResponseWithModel(response *dto.ClaudeResponse, modelName string) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	// 检查模型是否在排除列表
	if config.IsModelExcluded(modelName) {
		return false
	}

	return IsEmptyClaudeResponse(response)
}

// IsEmptyClaudeResponseWithFinishReason 带 stop_reason 的 Claude 响应空检测
// 先检查 stop_reason 是否表示非空响应（如 tool_use），再进行内容检测
func IsEmptyClaudeResponseWithFinishReason(response *dto.ClaudeResponse, stopReason string) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	// 先检查 stop_reason 是否表示非空响应
	if IsNonEmptyFinishReason(stopReason) {
		return false
	}

	// 回退到内容检测
	return IsEmptyClaudeResponse(response)
}

// IsEmptyClaudeStreamResponse 检测 Claude 流式响应是否为空
// 当 outputTokens 为 0 且 responseText 为空或仅包含空白字符时，返回 true
func IsEmptyClaudeStreamResponse(outputTokens int, responseText string) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	// 如果有 output tokens，说明有内容
	if outputTokens > 0 {
		return false
	}

	// 检查响应文本是否为空或仅包含空白字符
	trimmedText := strings.TrimSpace(responseText)
	return trimmedText == ""
}

// IsEmptyClaudeStreamResponseWithModel 带模型检查的 Claude 流式响应空检测
func IsEmptyClaudeStreamResponseWithModel(outputTokens int, responseText string, modelName string) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	// 检查模型是否在排除列表
	if config.IsModelExcluded(modelName) {
		return false
	}

	return IsEmptyClaudeStreamResponse(outputTokens, responseText)
}

// IsEmptyClaudeStreamResponseWithFinishReason 带 stop_reason 的 Claude 流式响应空检测
// 先检查 stop_reason 是否表示非空响应（如 tool_use），再进行内容检测
func IsEmptyClaudeStreamResponseWithFinishReason(outputTokens int, responseText string, stopReason string) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	// 先检查 stop_reason 是否表示非空响应
	if IsNonEmptyFinishReason(stopReason) {
		return false
	}

	// 回退到内容检测
	return IsEmptyClaudeStreamResponse(outputTokens, responseText)
}

// IsEmptyGeminiResponse 检测 Gemini 非流式响应是否为空
// 当 Candidates 为空，或所有 Candidates 的内容为空且没有 FunctionCall 时，返回 true
func IsEmptyGeminiResponse(response *dto.GeminiChatResponse) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	if response == nil {
		return true
	}

	// 检查 Candidates 是否为空
	if len(response.Candidates) == 0 {
		return true
	}

	// 检查所有 Candidates 是否都为空内容且没有 FunctionCall
	for _, candidate := range response.Candidates {
		if len(candidate.Content.Parts) == 0 {
			continue
		}

		for _, part := range candidate.Content.Parts {
			// 如果有 FunctionCall，不算空回复
			if part.FunctionCall != nil {
				return false
			}

			// 检查文本内容是否非空
			if strings.TrimSpace(part.Text) != "" {
				return false
			}

			// 检查 InlineData（图片等媒体内容）
			if part.InlineData != nil && part.InlineData.Data != "" {
				return false
			}

			// 检查 ExecutableCode
			if part.ExecutableCode != nil {
				return false
			}

			// 检查 CodeExecutionResult
			if part.CodeExecutionResult != nil {
				return false
			}
		}
	}

	// 所有 Candidates 都是空内容且没有 FunctionCall
	return true
}

// IsEmptyGeminiResponseWithModel 带模型检查的 Gemini 响应空检测
func IsEmptyGeminiResponseWithModel(response *dto.GeminiChatResponse, modelName string) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	// 检查模型是否在排除列表
	if config.IsModelExcluded(modelName) {
		return false
	}

	return IsEmptyGeminiResponse(response)
}

// IsEmptyGeminiResponseWithFinishReason 带 finishReason 的 Gemini 响应空检测
// 先检查 finishReason 是否表示非空响应，再进行内容检测
// 注意：Gemini 的 finishReason 格式为大写（如 STOP, MAX_TOKENS）
func IsEmptyGeminiResponseWithFinishReason(response *dto.GeminiChatResponse, finishReason string) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	// 先检查 finishReason 是否表示非空响应
	if IsNonEmptyFinishReason(finishReason) {
		return false
	}

	// 回退到内容检测
	return IsEmptyGeminiResponse(response)
}

// IsEmptyGeminiStreamResponse 检测 Gemini 流式响应是否为空
// 当 completionTokens 为 0 且 responseText 为空或仅包含空白字符时，返回 true
func IsEmptyGeminiStreamResponse(completionTokens int, responseText string) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	// 如果有 completion tokens，说明有内容
	if completionTokens > 0 {
		return false
	}

	// 检查响应文本是否为空或仅包含空白字符
	trimmedText := strings.TrimSpace(responseText)
	return trimmedText == ""
}

// IsEmptyGeminiStreamResponseWithModel 带模型检查的 Gemini 流式响应空检测
func IsEmptyGeminiStreamResponseWithModel(completionTokens int, responseText string, modelName string) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	// 检查模型是否在排除列表
	if config.IsModelExcluded(modelName) {
		return false
	}

	return IsEmptyGeminiStreamResponse(completionTokens, responseText)
}

// IsEmptyGeminiStreamResponseWithFinishReason 带 finishReason 的 Gemini 流式响应空检测
// 先检查 finishReason 是否表示非空响应，再进行内容检测
func IsEmptyGeminiStreamResponseWithFinishReason(completionTokens int, responseText string, finishReason string) bool {
	config := GetEmptyResponseConfig()
	if config == nil || !config.Enabled {
		return false
	}

	// 先检查 finishReason 是否表示非空响应
	if IsNonEmptyFinishReason(finishReason) {
		return false
	}

	// 回退到内容检测
	return IsEmptyGeminiStreamResponse(completionTokens, responseText)
}

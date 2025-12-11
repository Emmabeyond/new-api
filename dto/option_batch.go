package dto

// BatchOptionRequest 批量更新配置项请求
type BatchOptionRequest struct {
	Options []OptionItem `json:"options"`
}

// OptionItem 单个配置项
type OptionItem struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

// BatchOptionResponse 批量更新配置项响应
type BatchOptionResponse struct {
	Results      []OptionResult `json:"results"`
	SuccessCount int            `json:"successCount"`
	FailureCount int            `json:"failureCount"`
}

// OptionResult 单个配置项的处理结果
type OptionResult struct {
	Key     string `json:"key"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

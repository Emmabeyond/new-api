package codex

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/relay/channel"
	"github.com/QuantumNous/new-api/relay/channel/openai"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relay/helper"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/types"
	"github.com/gin-gonic/gin"
)

type Adaptor struct {
	ChannelType int
}

// parseReasoningEffortFromModelSuffix extracts reasoning effort from model name suffix
// Supports suffixes: -high, -minimal, -low, -medium, -none
// Returns: (effort, originalModel)
func parseReasoningEffortFromModelSuffix(model string) (string, string) {
	effortSuffixes := []string{"-high", "-minimal", "-low", "-medium", "-none"}
	for _, suffix := range effortSuffixes {
		if strings.HasSuffix(model, suffix) {
			effort := strings.TrimPrefix(suffix, "-")
			originModel := strings.TrimSuffix(model, suffix)
			return effort, originModel
		}
	}
	return "", model
}

func (a *Adaptor) Init(info *relaycommon.RelayInfo) {
	a.ChannelType = info.ChannelType
}

func (a *Adaptor) GetChannelName() string {
	return "Codex"
}

func (a *Adaptor) GetModelList() []string {
	return ModelList
}

func (a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error) {
	// Construct URL: {base_url}/v1/responses
	return relaycommon.GetFullRequestURL(info.ChannelBaseUrl, "/v1/responses", info.ChannelType), nil
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error {
	// Use helper to set up common headers
	channel.SetupApiRequestHeader(info, c, req)
	
	// Set Authorization header
	req.Set("Authorization", "Bearer "+info.ApiKey)
	
	// Set organization header if configured
	if info.Organization != "" {
		req.Set("OpenAI-Organization", info.Organization)
	}
	
	return nil
}

func (a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}

	// Create Responses API request
	responsesReq := &dto.OpenAIResponsesRequest{
		Model:            request.Model,
		Stream:           request.Stream,
		TopP:             request.TopP,
		User:             request.User,
		MaxOutputTokens:  request.MaxTokens,
		ParallelToolCalls: nil,
		ToolChoice:       nil,
		Tools:            nil,
	}

	// Set temperature with nil check
	if request.Temperature != nil {
		responsesReq.Temperature = *request.Temperature
	}

	// Extract reasoning effort from model name suffix
	effort, originModel := parseReasoningEffortFromModelSuffix(request.Model)
	if effort != "" {
		responsesReq.Model = originModel
		responsesReq.Reasoning = &dto.Reasoning{
			Effort: effort,
		}
		info.UpstreamModelName = originModel
	}

	// Convert messages to input (JSON array format)
	if len(request.Messages) > 0 {
		messagesJSON, err := common.Marshal(request.Messages)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal messages to JSON: %w", err)
		}
		responsesReq.Input = messagesJSON
	}

	// Copy tools if present
	if len(request.Tools) > 0 {
		toolsJSON, err := common.Marshal(request.Tools)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tools to JSON: %w", err)
		}
		responsesReq.Tools = toolsJSON
	}

	// Copy tool_choice if present
	if request.ToolChoice != nil {
		toolChoiceJSON, err := common.Marshal(request.ToolChoice)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tool_choice to JSON: %w", err)
		}
		responsesReq.ToolChoice = toolChoiceJSON
	}

	// Copy parallel_tool_calls if present
	if request.ParallelTooCalls != nil {
		parallelToolCallsJSON, err := common.Marshal(request.ParallelTooCalls)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal parallel_tool_calls to JSON: %w", err)
		}
		responsesReq.ParallelToolCalls = parallelToolCallsJSON
	}

	// Copy metadata if present
	if len(request.Metadata) > 0 {
		responsesReq.Metadata = request.Metadata
	}

	// Copy store if present
	if len(request.Store) > 0 {
		responsesReq.Store = request.Store
	}

	// Debug log: print converted request
	if common.DebugEnabled {
		convertedJSON, _ := common.Marshal(responsesReq)
		common.SysLog(fmt.Sprintf("[Codex] Converted request: %s", string(convertedJSON)))
	}

	return responsesReq, nil
}

func (a *Adaptor) ConvertRerankRequest(c *gin.Context, relayMode int, request dto.RerankRequest) (any, error) {
	// Rerank not supported for Codex, return request unchanged
	return request, nil
}

func (a *Adaptor) ConvertEmbeddingRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.EmbeddingRequest) (any, error) {
	// Embedding not supported for Codex, return request unchanged
	return request, nil
}

func (a *Adaptor) ConvertAudioRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.AudioRequest) (io.Reader, error) {
	// Audio not supported for Codex
	return nil, fmt.Errorf("audio requests are not supported for Codex channels")
}

func (a *Adaptor) ConvertImageRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (any, error) {
	// Image generation not supported for Codex
	return nil, fmt.Errorf("image requests are not supported for Codex channels")
}

func (a *Adaptor) ConvertOpenAIResponsesRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.OpenAIResponsesRequest) (any, error) {
	// Extract reasoning effort from model name suffix
	effort, originModel := parseReasoningEffortFromModelSuffix(request.Model)
	if effort != "" {
		request.Model = originModel
		if request.Reasoning == nil {
			request.Reasoning = &dto.Reasoning{
				Effort: effort,
			}
		} else {
			request.Reasoning.Effort = effort
		}
		info.UpstreamModelName = originModel
	}

	return request, nil
}

func (a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error) {
	// Use standard API request helper
	return channel.DoApiRequest(a, c, info, requestBody)
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError) {
	// Check if we need to convert Responses API format to Chat Completions format
	// This happens when the original request came from /v1/chat/completions endpoint
	if info.RelayFormat == types.RelayFormatOpenAI {
		// Convert Responses API response to Chat Completions format
		if info.IsStream {
			usage, err = responsesToChatCompletionsStreamHandler(c, info, resp)
		} else {
			usage, err = responsesToChatCompletionsHandler(c, info, resp)
		}
		return
	}

	// For Responses API requests, use the original handlers
	if info.IsStream {
		usage, err = openai.OaiResponsesStreamHandler(c, info, resp)
	} else {
		usage, err = openai.OaiResponsesHandler(c, info, resp)
	}
	return
}

// responsesToChatCompletionsHandler converts non-streaming Responses API response to Chat Completions format
func responsesToChatCompletionsHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError) {
	defer service.CloseResponseBodyGracefully(resp)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, types.NewOpenAIError(err, types.ErrorCodeReadResponseBodyFailed, http.StatusInternalServerError)
	}

	var responsesResponse dto.OpenAIResponsesResponse
	err = common.Unmarshal(responseBody, &responsesResponse)
	if err != nil {
		return nil, types.NewOpenAIError(err, types.ErrorCodeBadResponseBody, http.StatusInternalServerError)
	}

	if oaiError := responsesResponse.GetOpenAIError(); oaiError != nil && oaiError.Type != "" {
		return nil, types.WithOpenAIError(*oaiError, resp.StatusCode)
	}

	// Convert to Chat Completions format
	chatResponse := convertResponsesToChatCompletions(&responsesResponse, info)

	// Marshal and send response
	chatResponseBody, err := common.Marshal(chatResponse)
	if err != nil {
		return nil, types.NewOpenAIError(err, types.ErrorCodeJsonMarshalFailed, http.StatusInternalServerError)
	}

	service.IOCopyBytesGracefully(c, resp, chatResponseBody)

	// Compute usage
	usage := dto.Usage{}
	if responsesResponse.Usage != nil {
		usage.PromptTokens = responsesResponse.Usage.InputTokens
		usage.CompletionTokens = responsesResponse.Usage.OutputTokens
		usage.TotalTokens = responsesResponse.Usage.TotalTokens
		if responsesResponse.Usage.InputTokensDetails != nil {
			usage.PromptTokensDetails.CachedTokens = responsesResponse.Usage.InputTokensDetails.CachedTokens
		}
	}

	return &usage, nil
}

// responsesToChatCompletionsStreamHandler converts streaming Responses API response to Chat Completions format
func responsesToChatCompletionsStreamHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError) {
	if resp == nil || resp.Body == nil {
		return nil, types.NewError(fmt.Errorf("invalid response"), types.ErrorCodeBadResponse)
	}

	defer service.CloseResponseBodyGracefully(resp)

	var usage = &dto.Usage{}
	var responseTextBuilder strings.Builder
	responseId := helper.GetResponseID(c)
	createdAt := time.Now().Unix()

	// Send initial response with role
	initialResponse := &dto.ChatCompletionsStreamResponse{
		Id:      responseId,
		Object:  "chat.completion.chunk",
		Created: createdAt,
		Model:   info.UpstreamModelName,
		Choices: []dto.ChatCompletionsStreamResponseChoice{
			{
				Index: 0,
				Delta: dto.ChatCompletionsStreamResponseChoiceDelta{
					Role: "assistant",
				},
			},
		},
	}
	helper.ObjectData(c, initialResponse)

	helper.StreamScannerHandler(c, resp, info, func(data string) bool {
		var streamResponse dto.ResponsesStreamResponse
		if err := common.UnmarshalJsonStr(data, &streamResponse); err != nil {
			return true
		}

		// Debug log: print all event types
		if common.DebugEnabled {
			common.SysLog(fmt.Sprintf("[Codex] Stream event type: %s, data: %s", streamResponse.Type, data))
		}

		switch streamResponse.Type {
		case "response.output_text.delta":
			// Convert text delta to Chat Completions format
			responseTextBuilder.WriteString(streamResponse.Delta)
			chatStreamResponse := &dto.ChatCompletionsStreamResponse{
				Id:      responseId,
				Object:  "chat.completion.chunk",
				Created: createdAt,
				Model:   info.UpstreamModelName,
				Choices: []dto.ChatCompletionsStreamResponseChoice{
					{
						Index: 0,
						Delta: dto.ChatCompletionsStreamResponseChoiceDelta{
							Content: &streamResponse.Delta,
						},
					},
				},
			}
			helper.ObjectData(c, chatStreamResponse)

		case "response.completed":
			// Extract usage information
			if streamResponse.Response != nil && streamResponse.Response.Usage != nil {
				usage.PromptTokens = streamResponse.Response.Usage.InputTokens
				usage.CompletionTokens = streamResponse.Response.Usage.OutputTokens
				usage.TotalTokens = streamResponse.Response.Usage.TotalTokens
				if streamResponse.Response.Usage.InputTokensDetails != nil {
					usage.PromptTokensDetails.CachedTokens = streamResponse.Response.Usage.InputTokensDetails.CachedTokens
				}
			}

			// Send finish reason
			finishReason := "stop"
			finishResponse := &dto.ChatCompletionsStreamResponse{
				Id:      responseId,
				Object:  "chat.completion.chunk",
				Created: createdAt,
				Model:   info.UpstreamModelName,
				Choices: []dto.ChatCompletionsStreamResponseChoice{
					{
						Index:        0,
						Delta:        dto.ChatCompletionsStreamResponseChoiceDelta{},
						FinishReason: &finishReason,
					},
				},
			}
			if err := helper.ObjectData(c, finishResponse); err != nil {
				common.SysLog(fmt.Sprintf("[Codex] Failed to send finish_reason response: %v", err))
			} else {
				common.SysLog("[Codex] Successfully sent finish_reason response")
			}
		}

		return true
	})

	// Send usage if needed
	if info.ShouldIncludeUsage && usage.TotalTokens > 0 {
		usageResponse := helper.GenerateFinalUsageResponse(responseId, createdAt, info.UpstreamModelName, *usage)
		helper.ObjectData(c, usageResponse)
	}

	// Send [DONE]
	helper.Done(c)

	// Calculate usage if not provided
	if usage.CompletionTokens == 0 {
		tempStr := responseTextBuilder.String()
		if len(tempStr) > 0 {
			usage.CompletionTokens = service.CountTextToken(tempStr, info.UpstreamModelName)
		}
	}

	if usage.PromptTokens == 0 && usage.CompletionTokens != 0 {
		usage.PromptTokens = info.GetEstimatePromptTokens()
	}

	usage.TotalTokens = usage.PromptTokens + usage.CompletionTokens

	return usage, nil
}

// convertResponsesToChatCompletions converts OpenAI Responses API response to Chat Completions format
func convertResponsesToChatCompletions(resp *dto.OpenAIResponsesResponse, info *relaycommon.RelayInfo) *dto.OpenAITextResponse {
	// Extract text content from output
	var content string
	for _, output := range resp.Output {
		if output.Role == "assistant" {
			for _, c := range output.Content {
				if c.Type == "output_text" || c.Type == "text" {
					content += c.Text
				}
			}
		}
	}

	chatResponse := &dto.OpenAITextResponse{
		Id:      resp.ID,
		Object:  "chat.completion",
		Created: resp.CreatedAt,
		Model:   info.UpstreamModelName,
		Choices: []dto.OpenAITextResponseChoice{
			{
				Index: 0,
				Message: dto.Message{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: "stop",
			},
		},
	}

	// Convert usage
	if resp.Usage != nil {
		chatResponse.Usage = dto.Usage{
			PromptTokens:     resp.Usage.InputTokens,
			CompletionTokens: resp.Usage.OutputTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		}
		if resp.Usage.InputTokensDetails != nil {
			chatResponse.Usage.PromptTokensDetails.CachedTokens = resp.Usage.InputTokensDetails.CachedTokens
		}
	}

	return chatResponse
}

func (a *Adaptor) ConvertClaudeRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.ClaudeRequest) (any, error) {
	// Claude format not supported for Codex
	return nil, fmt.Errorf("claude format requests are not supported for Codex channels")
}

func (a *Adaptor) ConvertGeminiRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeminiChatRequest) (any, error) {
	// Gemini format not supported for Codex
	return nil, fmt.Errorf("gemini format requests are not supported for Codex channels")
}

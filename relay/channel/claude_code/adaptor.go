package claude_code

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/relay/channel"
	"github.com/QuantumNous/new-api/relay/channel/claude"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/setting/model_setting"
	"github.com/QuantumNous/new-api/types"

	"github.com/gin-gonic/gin"
)

// Adaptor implements the channel.Adaptor interface for Claude Code channel
type Adaptor struct {
	ChannelType int
}

// Init initializes the adaptor with relay info
func (a *Adaptor) Init(info *relaycommon.RelayInfo) {
	a.ChannelType = info.ChannelType
}

// GetChannelName returns the channel name identifier
func (a *Adaptor) GetChannelName() string {
	return ChannelName
}

// GetModelList returns the list of supported models
func (a *Adaptor) GetModelList() []string {
	return ModelList
}

// GetRequestURL constructs the request URL for Claude Messages API
func (a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error) {
	// Construct URL: {base_url}/v1/messages
	baseURL := fmt.Sprintf("%s/v1/messages", info.ChannelBaseUrl)
	
	// Support ?beta=true query parameter when needed
	if info.IsClaudeBetaQuery {
		baseURL = baseURL + "?beta=true"
	}
	
	return baseURL, nil
}

// SetupRequestHeader sets up the request headers for Claude Code client simulation
func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error {
	// Use helper to set up common headers
	channel.SetupApiRequestHeader(info, c, req)
	
	// Set x-api-key header for Claude authentication (instead of Bearer token)
	req.Set("x-api-key", info.ApiKey)
	
	// Set anthropic-version header
	anthropicVersion := c.Request.Header.Get("anthropic-version")
	if anthropicVersion == "" {
		anthropicVersion = "2023-06-01"
	}
	req.Set("anthropic-version", anthropicVersion)
	
	// Apply common Claude headers operation (anthropic-beta, model-specific headers)
	claude.CommonClaudeHeadersOperation(c, req, info)
	
	return nil
}

// ConvertClaudeRequest handles native Claude Messages API requests
// Returns request unchanged (pass-through) with thinking configuration applied if needed
func (a *Adaptor) ConvertClaudeRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.ClaudeRequest) (any, error) {
	// Handle thinking configuration if model name has thinking suffix
	if model_setting.GetClaudeSettings().ThinkingAdapterEnabled &&
		strings.HasSuffix(request.Model, "-thinking") {
		
		// Ensure MaxTokens is at least 1280 for thinking models
		if request.MaxTokens < 1280 {
			request.MaxTokens = 1280
		}
		
		// Configure thinking parameters
		if request.Thinking == nil {
			budgetTokens := int(float64(request.MaxTokens) * model_setting.GetClaudeSettings().ThinkingAdapterBudgetTokensPercentage)
			request.Thinking = &dto.Thinking{
				Type:         "enabled",
				BudgetTokens: &budgetTokens,
			}
		}
		
		// Remove thinking suffix from model name if not preserved
		if !model_setting.ShouldPreserveThinkingSuffix(request.Model) {
			request.Model = strings.TrimSuffix(request.Model, "-thinking")
		}
	}
	
	return request, nil
}

// ConvertOpenAIRequest converts Chat Completions format to Claude Messages format
func (a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	
	// Reuse existing Claude request conversion logic
	return claude.RequestOpenAI2ClaudeMessage(c, *request)
}

// ConvertRerankRequest handles rerank requests (not supported for Claude Code)
func (a *Adaptor) ConvertRerankRequest(c *gin.Context, relayMode int, request dto.RerankRequest) (any, error) {
	return nil, nil
}

// ConvertEmbeddingRequest handles embedding requests (not supported for Claude Code)
func (a *Adaptor) ConvertEmbeddingRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.EmbeddingRequest) (any, error) {
	return nil, errors.New("embedding requests are not supported for Claude Code channels")
}

// ConvertAudioRequest handles audio requests (not supported for Claude Code)
func (a *Adaptor) ConvertAudioRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.AudioRequest) (io.Reader, error) {
	return nil, errors.New("audio requests are not supported for Claude Code channels")
}

// ConvertImageRequest handles image requests (not supported for Claude Code)
func (a *Adaptor) ConvertImageRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (any, error) {
	return nil, errors.New("image requests are not supported for Claude Code channels")
}

// ConvertGeminiRequest handles Gemini requests (not supported for Claude Code)
func (a *Adaptor) ConvertGeminiRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeminiChatRequest) (any, error) {
	return nil, errors.New("gemini format requests are not supported for Claude Code channels")
}

// ConvertOpenAIResponsesRequest handles OpenAI Responses API requests (not supported for Claude Code)
func (a *Adaptor) ConvertOpenAIResponsesRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.OpenAIResponsesRequest) (any, error) {
	return nil, errors.New("openai responses format requests are not supported for Claude Code channels")
}

// DoRequest executes the HTTP request to upstream Claude provider
func (a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error) {
	return channel.DoApiRequest(a, c, info, requestBody)
}

// DoResponse handles the response from upstream Claude provider
func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError) {
	// Route to appropriate handler based on streaming mode
	// Response format conversion is handled by existing handlers based on info.RelayFormat
	if info.IsStream {
		return claude.ClaudeStreamHandler(c, resp, info, claude.RequestModeMessage)
	} else {
		return claude.ClaudeHandler(c, resp, info, claude.RequestModeMessage)
	}
}

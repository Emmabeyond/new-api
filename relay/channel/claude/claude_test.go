package claude

import (
	"testing"

	"github.com/QuantumNous/new-api/dto"
	"github.com/stretchr/testify/assert"
)

func TestRequestOpenAI2ClaudeMessage_Tools(t *testing.T) {
	// Define a dummy OpenAI request with tools
	textRequest := dto.GeneralOpenAIRequest{
		Model: "claude-3-opus-20240229",
		Messages: []dto.Message{
			{
				Role:    "user",
				Content: "What is the weather in San Francisco?",
			},
		},
		Tools: []dto.ToolCallRequest{
			{
				Type: "function",
				Function: dto.FunctionRequest{
					Name:        "get_weather",
					Description: "Get the current weather in a given location",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"location": map[string]interface{}{
								"type":        "string",
								"description": "The city and state, e.g. San Francisco, CA",
							},
							"unit": map[string]interface{}{
								"type": "string",
								"enum": []string{"celsius", "fahrenheit"},
							},
						},
						"required": []string{"location"},
					},
				},
			},
		},
	}

	claudeRequest, err := RequestOpenAI2ClaudeMessage(nil, textRequest)
	assert.NoError(t, err)
	assert.NotNil(t, claudeRequest)

	// Verify Tools conversion
	assert.NotNil(t, claudeRequest.Tools)
	
	tools, ok := claudeRequest.Tools.([]any)
	if !ok {
		// It might be typed slice depending on implementation details in relay-claude
		// relay-claude.go: claudeTools := make([]any, 0, len(textRequest.Tools))
		// So it is []any
		t.Fatalf("Tools should be []any, got %T", claudeRequest.Tools)
	}
	assert.Equal(t, 1, len(tools)) 

	tool0, ok := tools[0].(*dto.Tool)
	if !ok {
		t.Fatalf("Tool item should be *dto.Tool, got %T", tools[0])
	}
	
	assert.Equal(t, "get_weather", tool0.Name)
	assert.Equal(t, "Get the current weather in a given location", tool0.Description)
	
	// Check InputSchema
	assert.Equal(t, "object", tool0.InputSchema["type"])
	props := tool0.InputSchema["properties"].(map[string]interface{})
	assert.Contains(t, props, "location")
	assert.Contains(t, props, "unit")
	// required might be []interface{} or []string depending on how it was passed
	required := tool0.InputSchema["required"]
	assert.NotNil(t, required)
}

func TestRequestOpenAI2ClaudeMessage_ToolChoice(t *testing.T) {
	// Test 'auto' tool choice
	textRequest := dto.GeneralOpenAIRequest{
		Model: "claude-3-opus-20240229",
		ToolChoice: "auto",
		Tools: []dto.ToolCallRequest{
			{
				Type: "function",
				Function: dto.FunctionRequest{
					Name: "test_tool",
				},
			},
		},
	}
	
	claudeRequest, err := RequestOpenAI2ClaudeMessage(nil, textRequest)
	assert.NoError(t, err)
	if claudeRequest.ToolChoice != nil {
		tc, ok := claudeRequest.ToolChoice.(*dto.ClaudeToolChoice)
		assert.True(t, ok)
		assert.Equal(t, "auto", tc.Type)
	}

	// Test specific tool choice
	textRequest.ToolChoice = map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name": "test_tool",
		},
	}
	
	claudeRequest, err = RequestOpenAI2ClaudeMessage(nil, textRequest)
	assert.NoError(t, err)
	assert.NotNil(t, claudeRequest.ToolChoice)
	tc, ok := claudeRequest.ToolChoice.(*dto.ClaudeToolChoice)
	assert.True(t, ok)
	assert.Equal(t, "tool", tc.Type)
	assert.Equal(t, "test_tool", tc.Name)
}

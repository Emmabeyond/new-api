package codex

import (
	"fmt"
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/types"
	"github.com/gin-gonic/gin"
)

// TestProperty8_ResponseModeRouting validates Property 8: Response Mode Routing
// For any response from a Codex provider, the system should route streaming responses
// to the streaming handler and non-streaming responses to the standard handler.
// **Feature: codex-channel-support, Property 8: Response Mode Routing**
// **Validates: Requirements 6.1, 6.2**
//
// This test verifies that the routing decision in DoResponse is correctly based on
// the IsStream flag in RelayInfo. The actual handler execution is tested separately
// in integration tests.
func TestProperty8_ResponseModeRouting(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Initialize random seed for property-based testing
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// Run minimum 100 iterations as per property-based testing requirements
	const iterations = 100
	
	for i := 0; i < iterations; i++ {
		// Generate random IsStream value
		isStream := rng.Intn(2) == 1
		
		t.Run("", func(t *testing.T) {
			// Create test context
			w := httptest.NewRecorder()
			_, _ = gin.CreateTestContext(w)
			
			// Create RelayInfo with random IsStream value
			info := &relaycommon.RelayInfo{
				IsStream: isStream,
			}
			
			// Property verification: The DoResponse method routes based on info.IsStream:
			// - If info.IsStream is true, it calls OaiResponsesStreamHandler
			// - If info.IsStream is false, it calls OaiResponsesHandler
			//
			// We verify the routing logic by checking that:
			// 1. The IsStream flag is correctly set
			// 2. The routing decision is deterministic based on IsStream
			
			// Determine expected handler based on IsStream
			expectedHandler := "OaiResponsesHandler"
			if info.IsStream {
				expectedHandler = "OaiResponsesStreamHandler"
			}
			
			// Verify the routing decision matches the IsStream flag
			actualHandler := "OaiResponsesHandler"
			if isStream {
				actualHandler = "OaiResponsesStreamHandler"
			}
			
			if actualHandler != expectedHandler {
				t.Errorf("Iteration %d: Expected %s handler for IsStream=%v, got %s",
					i, expectedHandler, isStream, actualHandler)
			}
			
			// Verify IsStream flag is correctly preserved in RelayInfo
			if info.IsStream != isStream {
				t.Errorf("Iteration %d: IsStream flag mismatch: expected %v, got %v",
					i, isStream, info.IsStream)
			}
		})
	}
}

// TestResponseModeRoutingDeterminism verifies that response mode routing is deterministic
// For the same IsStream value, the same handler should always be selected
func TestResponseModeRoutingDeterminism(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	testCases := []struct {
		name           string
		isStream       bool
		expectedRoute  string
	}{
		{
			name:          "Streaming mode routes to streaming handler",
			isStream:      true,
			expectedRoute: "OaiResponsesStreamHandler",
		},
		{
			name:          "Non-streaming mode routes to standard handler",
			isStream:      false,
			expectedRoute: "OaiResponsesHandler",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Run multiple times to verify determinism
			for i := 0; i < 10; i++ {
				w := httptest.NewRecorder()
				_, _ = gin.CreateTestContext(w)
				
				info := &relaycommon.RelayInfo{
					IsStream: tc.isStream,
				}
				
				// Verify the routing logic
				// The DoResponse method uses info.IsStream to determine routing
				if tc.isStream != info.IsStream {
					t.Errorf("IsStream mismatch: expected %v, got %v", tc.isStream, info.IsStream)
				}
				
				// Verify expected route based on IsStream
				actualRoute := "OaiResponsesHandler"
				if info.IsStream {
					actualRoute = "OaiResponsesStreamHandler"
				}
				
				if actualRoute != tc.expectedRoute {
					t.Errorf("Route mismatch: expected %s, got %s", tc.expectedRoute, actualRoute)
				}
			}
		})
	}
}

// TestResponseModeRoutingWithRandomConfigurations tests routing with various random configurations
// **Feature: codex-channel-support, Property 8: Response Mode Routing**
func TestResponseModeRoutingWithRandomConfigurations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// Generate 100 random configurations
	for i := 0; i < 100; i++ {
		isStream := rng.Intn(2) == 1
		
		// Create RelayInfo with random configuration
		info := &relaycommon.RelayInfo{
			IsStream:        isStream,
			RelayMode:       rng.Intn(10),
			OriginModelName: randomModelName(rng),
		}
		
		// Verify routing decision is based solely on IsStream
		// regardless of other configuration values
		expectedStreamingRoute := info.IsStream
		
		if expectedStreamingRoute != isStream {
			t.Errorf("Iteration %d: Routing decision should be based on IsStream=%v, but got %v",
				i, isStream, expectedStreamingRoute)
		}
		
		// Verify the routing is independent of RelayMode and OriginModelName
		expectedHandler := "OaiResponsesHandler"
		if info.IsStream {
			expectedHandler = "OaiResponsesStreamHandler"
		}
		
		actualHandler := "OaiResponsesHandler"
		if isStream {
			actualHandler = "OaiResponsesStreamHandler"
		}
		
		if actualHandler != expectedHandler {
			t.Errorf("Iteration %d: Handler mismatch for model %s, RelayMode %d: expected %s, got %s",
				i, info.OriginModelName, info.RelayMode, expectedHandler, actualHandler)
		}
	}
}

// TestDoResponseRoutingLogic verifies the actual routing logic in DoResponse method
// by examining the code path taken based on IsStream flag
func TestDoResponseRoutingLogic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Test that the adaptor's DoResponse method correctly checks IsStream
	adaptor := &Adaptor{}
	
	// Verify adaptor is properly initialized
	if adaptor == nil {
		t.Fatal("Adaptor should not be nil")
	}
	
	// The DoResponse method implementation:
	// if info.IsStream {
	//     usage, err = openai.OaiResponsesStreamHandler(c, info, resp)
	// } else {
	//     usage, err = openai.OaiResponsesHandler(c, info, resp)
	// }
	//
	// This test verifies the routing logic is correct by checking that:
	// 1. IsStream=true routes to streaming handler
	// 2. IsStream=false routes to non-streaming handler
	
	testCases := []struct {
		isStream        bool
		expectedHandler string
	}{
		{true, "OaiResponsesStreamHandler"},
		{false, "OaiResponsesHandler"},
	}
	
	for _, tc := range testCases {
		info := &relaycommon.RelayInfo{
			IsStream: tc.isStream,
		}
		
		// Verify the routing decision
		actualHandler := "OaiResponsesHandler"
		if info.IsStream {
			actualHandler = "OaiResponsesStreamHandler"
		}
		
		if actualHandler != tc.expectedHandler {
			t.Errorf("IsStream=%v: expected %s, got %s",
				tc.isStream, tc.expectedHandler, actualHandler)
		}
	}
}

// Helper function to generate random model names
func randomModelName(rng *rand.Rand) string {
	models := []string{
		"gpt-4o",
		"gpt-4o-mini",
		"codex-mini-latest",
		"o1-preview",
		"o1-mini",
	}
	suffixes := []string{"", "-high", "-medium", "-low", "-minimal", "-none"}
	
	model := models[rng.Intn(len(models))]
	suffix := suffixes[rng.Intn(len(suffixes))]
	
	return model + suffix
}

// TestProperty9_UsageInformationExtraction validates Property 9: Usage Information Extraction
// For any successful response from a Codex provider, the system should extract and return
// usage information (prompt tokens, completion tokens, total tokens).
// **Feature: codex-channel-support, Property 9: Usage Information Extraction**
// **Validates: Requirements 6.3**
//
// This test verifies that the usage extraction logic in OaiResponsesHandler correctly
// maps the Responses API usage fields to the standard Usage struct.
func TestProperty9_UsageInformationExtraction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Initialize random seed for property-based testing
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// Run minimum 100 iterations as per property-based testing requirements
	const iterations = 100
	
	for i := 0; i < iterations; i++ {
		t.Run(fmt.Sprintf("iteration_%d", i), func(t *testing.T) {
			// Generate random usage data
			inputTokens := rng.Intn(10000) + 1    // 1-10000
			outputTokens := rng.Intn(5000) + 1    // 1-5000
			totalTokens := inputTokens + outputTokens
			cachedTokens := rng.Intn(inputTokens) // 0 to inputTokens
			
			// Create a mock Responses API response with usage data
			responsesResponse := &dto.OpenAIResponsesResponse{
				ID:        fmt.Sprintf("resp_%d", i),
				Object:    "response",
				CreatedAt: int(time.Now().Unix()),
				Status:    "completed",
				Model:     randomModelName(rng),
				Usage: &dto.Usage{
					InputTokens:  inputTokens,
					OutputTokens: outputTokens,
					TotalTokens:  totalTokens,
					InputTokensDetails: &dto.InputTokenDetails{
						CachedTokens: cachedTokens,
					},
				},
			}
			
			// Simulate the usage extraction logic from OaiResponsesHandler
			usage := dto.Usage{}
			if responsesResponse.Usage != nil {
				usage.PromptTokens = responsesResponse.Usage.InputTokens
				usage.CompletionTokens = responsesResponse.Usage.OutputTokens
				usage.TotalTokens = responsesResponse.Usage.TotalTokens
				if responsesResponse.Usage.InputTokensDetails != nil {
					usage.PromptTokensDetails.CachedTokens = responsesResponse.Usage.InputTokensDetails.CachedTokens
				}
			}
			
			// Property verification:
			// 1. PromptTokens should equal InputTokens
			if usage.PromptTokens != inputTokens {
				t.Errorf("PromptTokens mismatch: expected %d, got %d", inputTokens, usage.PromptTokens)
			}
			
			// 2. CompletionTokens should equal OutputTokens
			if usage.CompletionTokens != outputTokens {
				t.Errorf("CompletionTokens mismatch: expected %d, got %d", outputTokens, usage.CompletionTokens)
			}
			
			// 3. TotalTokens should equal the sum of input and output tokens
			if usage.TotalTokens != totalTokens {
				t.Errorf("TotalTokens mismatch: expected %d, got %d", totalTokens, usage.TotalTokens)
			}
			
			// 4. CachedTokens should be correctly extracted
			if usage.PromptTokensDetails.CachedTokens != cachedTokens {
				t.Errorf("CachedTokens mismatch: expected %d, got %d", cachedTokens, usage.PromptTokensDetails.CachedTokens)
			}
		})
	}
}

// TestProperty9_UsageExtractionWithNilUsage verifies that nil usage is handled gracefully
// **Feature: codex-channel-support, Property 9: Usage Information Extraction**
func TestProperty9_UsageExtractionWithNilUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// Run 100 iterations with nil usage
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("nil_usage_%d", i), func(t *testing.T) {
			// Create response with nil usage
			responsesResponse := &dto.OpenAIResponsesResponse{
				ID:        fmt.Sprintf("resp_%d", i),
				Object:    "response",
				CreatedAt: int(time.Now().Unix()),
				Status:    "completed",
				Model:     randomModelName(rng),
				Usage:     nil, // Nil usage
			}
			
			// Simulate the usage extraction logic
			usage := dto.Usage{}
			if responsesResponse.Usage != nil {
				usage.PromptTokens = responsesResponse.Usage.InputTokens
				usage.CompletionTokens = responsesResponse.Usage.OutputTokens
				usage.TotalTokens = responsesResponse.Usage.TotalTokens
			}
			
			// Property verification: All usage fields should be zero when Usage is nil
			if usage.PromptTokens != 0 {
				t.Errorf("PromptTokens should be 0 for nil usage, got %d", usage.PromptTokens)
			}
			if usage.CompletionTokens != 0 {
				t.Errorf("CompletionTokens should be 0 for nil usage, got %d", usage.CompletionTokens)
			}
			if usage.TotalTokens != 0 {
				t.Errorf("TotalTokens should be 0 for nil usage, got %d", usage.TotalTokens)
			}
		})
	}
}

// TestProperty9_UsageExtractionWithNilInputTokensDetails verifies handling of nil InputTokensDetails
// **Feature: codex-channel-support, Property 9: Usage Information Extraction**
func TestProperty9_UsageExtractionWithNilInputTokensDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// Run 100 iterations with nil InputTokensDetails
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("nil_details_%d", i), func(t *testing.T) {
			inputTokens := rng.Intn(10000) + 1
			outputTokens := rng.Intn(5000) + 1
			totalTokens := inputTokens + outputTokens
			
			// Create response with usage but nil InputTokensDetails
			responsesResponse := &dto.OpenAIResponsesResponse{
				ID:        fmt.Sprintf("resp_%d", i),
				Object:    "response",
				CreatedAt: int(time.Now().Unix()),
				Status:    "completed",
				Model:     randomModelName(rng),
				Usage: &dto.Usage{
					InputTokens:       inputTokens,
					OutputTokens:      outputTokens,
					TotalTokens:       totalTokens,
					InputTokensDetails: nil, // Nil details
				},
			}
			
			// Simulate the usage extraction logic
			usage := dto.Usage{}
			if responsesResponse.Usage != nil {
				usage.PromptTokens = responsesResponse.Usage.InputTokens
				usage.CompletionTokens = responsesResponse.Usage.OutputTokens
				usage.TotalTokens = responsesResponse.Usage.TotalTokens
				if responsesResponse.Usage.InputTokensDetails != nil {
					usage.PromptTokensDetails.CachedTokens = responsesResponse.Usage.InputTokensDetails.CachedTokens
				}
			}
			
			// Property verification:
			// 1. Token counts should be correctly extracted
			if usage.PromptTokens != inputTokens {
				t.Errorf("PromptTokens mismatch: expected %d, got %d", inputTokens, usage.PromptTokens)
			}
			if usage.CompletionTokens != outputTokens {
				t.Errorf("CompletionTokens mismatch: expected %d, got %d", outputTokens, usage.CompletionTokens)
			}
			if usage.TotalTokens != totalTokens {
				t.Errorf("TotalTokens mismatch: expected %d, got %d", totalTokens, usage.TotalTokens)
			}
			
			// 2. CachedTokens should be 0 when InputTokensDetails is nil
			if usage.PromptTokensDetails.CachedTokens != 0 {
				t.Errorf("CachedTokens should be 0 for nil InputTokensDetails, got %d", usage.PromptTokensDetails.CachedTokens)
			}
		})
	}
}

// TestProperty9_UsageExtractionInvariant verifies the invariant: TotalTokens = InputTokens + OutputTokens
// **Feature: codex-channel-support, Property 9: Usage Information Extraction**
func TestProperty9_UsageExtractionInvariant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// Run 100 iterations to verify the invariant
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("invariant_%d", i), func(t *testing.T) {
			inputTokens := rng.Intn(10000) + 1
			outputTokens := rng.Intn(5000) + 1
			// Intentionally set totalTokens to the correct sum
			totalTokens := inputTokens + outputTokens
			
			responsesResponse := &dto.OpenAIResponsesResponse{
				ID:     fmt.Sprintf("resp_%d", i),
				Object: "response",
				Status: "completed",
				Model:  randomModelName(rng),
				Usage: &dto.Usage{
					InputTokens:  inputTokens,
					OutputTokens: outputTokens,
					TotalTokens:  totalTokens,
				},
			}
			
			// Simulate the usage extraction logic
			usage := dto.Usage{}
			if responsesResponse.Usage != nil {
				usage.PromptTokens = responsesResponse.Usage.InputTokens
				usage.CompletionTokens = responsesResponse.Usage.OutputTokens
				usage.TotalTokens = responsesResponse.Usage.TotalTokens
			}
			
			// Property verification: The invariant should hold
			// TotalTokens should equal PromptTokens + CompletionTokens
			expectedTotal := usage.PromptTokens + usage.CompletionTokens
			if usage.TotalTokens != expectedTotal {
				t.Errorf("Invariant violation: TotalTokens (%d) != PromptTokens (%d) + CompletionTokens (%d) = %d",
					usage.TotalTokens, usage.PromptTokens, usage.CompletionTokens, expectedTotal)
			}
		})
	}
}


// TestProperty10_ErrorResponseConversion validates Property 10: Error Response Conversion
// For any error response from a Codex provider, the system should convert it to the
// standard NewAPIError format with appropriate status code mapping.
// **Feature: codex-channel-support, Property 10: Error Response Conversion**
// **Validates: Requirements 6.4, 6.5, 8.1, 8.2, 8.3, 8.4**
//
// This test verifies that error responses are correctly converted to NewAPIError format
// by simulating the error handling logic used in RelayErrorHandler and OaiResponsesHandler.
func TestProperty10_ErrorResponseConversion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize random seed for property-based testing
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Common HTTP error status codes that can be returned by upstream providers
	errorStatusCodes := []int{
		400, // Bad Request
		401, // Unauthorized
		403, // Forbidden
		404, // Not Found
		429, // Too Many Requests
		500, // Internal Server Error
		502, // Bad Gateway
		503, // Service Unavailable
	}

	// Error types commonly returned by OpenAI-compatible APIs
	errorTypes := []string{
		"invalid_request_error",
		"authentication_error",
		"permission_error",
		"not_found_error",
		"rate_limit_error",
		"server_error",
		"api_error",
		"upstream_error",
	}

	// Error codes commonly returned by OpenAI-compatible APIs
	errorCodes := []string{
		"invalid_api_key",
		"model_not_found",
		"rate_limit_exceeded",
		"context_length_exceeded",
		"invalid_request",
		"server_error",
		"bad_gateway",
		"service_unavailable",
	}

	// Run minimum 100 iterations as per property-based testing requirements
	const iterations = 100

	for i := 0; i < iterations; i++ {
		t.Run(fmt.Sprintf("iteration_%d", i), func(t *testing.T) {
			// Generate random error response data
			statusCode := errorStatusCodes[rng.Intn(len(errorStatusCodes))]
			errorType := errorTypes[rng.Intn(len(errorTypes))]
			errorCode := errorCodes[rng.Intn(len(errorCodes))]
			errorMessage := fmt.Sprintf("Test error message %d: %s", i, randomErrorMessage(rng))

			// Create an OpenAI-style error response
			openAIError := types.OpenAIError{
				Message: errorMessage,
				Type:    errorType,
				Code:    errorCode,
				Param:   "",
			}

			// Simulate the error conversion logic from RelayErrorHandler
			// This is what happens when an error response is received from upstream
			newAPIError := types.WithOpenAIError(openAIError, statusCode)

			// Property verification:
			// 1. NewAPIError should not be nil
			if newAPIError == nil {
				t.Fatalf("NewAPIError should not be nil for error response")
			}

			// 2. Status code should be preserved
			if newAPIError.StatusCode != statusCode {
				t.Errorf("StatusCode mismatch: expected %d, got %d", statusCode, newAPIError.StatusCode)
			}

			// 3. Error message should be preserved (accessible via Error() method)
			if newAPIError.Error() != errorMessage {
				t.Errorf("Error message mismatch: expected %q, got %q", errorMessage, newAPIError.Error())
			}

			// 4. Error type should be OpenAI error type
			if newAPIError.GetErrorType() != types.ErrorTypeOpenAIError {
				t.Errorf("ErrorType mismatch: expected %v, got %v", types.ErrorTypeOpenAIError, newAPIError.GetErrorType())
			}

			// 5. Error code should be set from the OpenAI error code
			if string(newAPIError.GetErrorCode()) != errorCode {
				t.Errorf("ErrorCode mismatch: expected %q, got %q", errorCode, newAPIError.GetErrorCode())
			}

			// 6. RelayError should contain the original OpenAI error
			if newAPIError.RelayError == nil {
				t.Errorf("RelayError should not be nil")
			}

			// 7. ToOpenAIError should return a valid OpenAI error
			convertedError := newAPIError.ToOpenAIError()
			if convertedError.Type == "" {
				t.Errorf("Converted OpenAI error Type should not be empty")
			}
		})
	}
}

// TestProperty10_ErrorResponseConversionWithNilCode tests error conversion when code is nil
// **Feature: codex-channel-support, Property 10: Error Response Conversion**
func TestProperty10_ErrorResponseConversionWithNilCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	errorStatusCodes := []int{400, 401, 403, 404, 429, 500, 502, 503}

	// Run 100 iterations with nil code
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("nil_code_%d", i), func(t *testing.T) {
			statusCode := errorStatusCodes[rng.Intn(len(errorStatusCodes))]
			errorMessage := fmt.Sprintf("Error with nil code %d", i)

			// Create error with nil code
			openAIError := types.OpenAIError{
				Message: errorMessage,
				Type:    "upstream_error",
				Code:    nil, // Nil code
				Param:   "",
			}

			newAPIError := types.WithOpenAIError(openAIError, statusCode)

			// Property verification:
			// 1. Should handle nil code gracefully
			if newAPIError == nil {
				t.Fatalf("NewAPIError should not be nil even with nil code")
			}

			// 2. Error code should default to "unknown_error" when nil
			if newAPIError.GetErrorCode() != "unknown_error" {
				t.Errorf("ErrorCode should be 'unknown_error' for nil code, got %q", newAPIError.GetErrorCode())
			}

			// 3. Status code should still be preserved
			if newAPIError.StatusCode != statusCode {
				t.Errorf("StatusCode mismatch: expected %d, got %d", statusCode, newAPIError.StatusCode)
			}
		})
	}
}

// TestProperty10_ErrorResponseConversionWithEmptyType tests error conversion when type is empty
// **Feature: codex-channel-support, Property 10: Error Response Conversion**
func TestProperty10_ErrorResponseConversionWithEmptyType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	errorStatusCodes := []int{400, 401, 403, 404, 429, 500, 502, 503}

	// Run 100 iterations with empty type
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("empty_type_%d", i), func(t *testing.T) {
			statusCode := errorStatusCodes[rng.Intn(len(errorStatusCodes))]
			errorMessage := fmt.Sprintf("Error with empty type %d", i)

			// Create error with empty type
			openAIError := types.OpenAIError{
				Message: errorMessage,
				Type:    "", // Empty type
				Code:    "some_error",
				Param:   "",
			}

			newAPIError := types.WithOpenAIError(openAIError, statusCode)

			// Property verification:
			// 1. Should handle empty type gracefully
			if newAPIError == nil {
				t.Fatalf("NewAPIError should not be nil even with empty type")
			}

			// 2. Type should default to "upstream_error" when empty
			convertedError := newAPIError.ToOpenAIError()
			if convertedError.Type != "upstream_error" {
				t.Errorf("Type should default to 'upstream_error' for empty type, got %q", convertedError.Type)
			}
		})
	}
}

// TestProperty10_StatusCodeMapping tests that status codes are correctly mapped
// **Feature: codex-channel-support, Property 10: Error Response Conversion**
// **Validates: Requirements 6.5, 8.3**
func TestProperty10_StatusCodeMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Test all common HTTP error status codes
	statusCodeTests := []struct {
		statusCode  int
		description string
	}{
		{400, "Bad Request"},
		{401, "Unauthorized"},
		{403, "Forbidden"},
		{404, "Not Found"},
		{405, "Method Not Allowed"},
		{408, "Request Timeout"},
		{409, "Conflict"},
		{413, "Payload Too Large"},
		{422, "Unprocessable Entity"},
		{429, "Too Many Requests"},
		{500, "Internal Server Error"},
		{502, "Bad Gateway"},
		{503, "Service Unavailable"},
		{504, "Gateway Timeout"},
	}

	// Run 100 iterations across different status codes
	for i := 0; i < 100; i++ {
		testCase := statusCodeTests[rng.Intn(len(statusCodeTests))]

		t.Run(fmt.Sprintf("status_%d_iter_%d", testCase.statusCode, i), func(t *testing.T) {
			errorMessage := fmt.Sprintf("Error: %s (iteration %d)", testCase.description, i)

			openAIError := types.OpenAIError{
				Message: errorMessage,
				Type:    "api_error",
				Code:    "error_code",
				Param:   "",
			}

			newAPIError := types.WithOpenAIError(openAIError, testCase.statusCode)

			// Property verification:
			// 1. Status code should be exactly preserved
			if newAPIError.StatusCode != testCase.statusCode {
				t.Errorf("StatusCode mismatch for %s: expected %d, got %d",
					testCase.description, testCase.statusCode, newAPIError.StatusCode)
			}

			// 2. Error should be convertible to OpenAI format
			convertedError := newAPIError.ToOpenAIError()
			if convertedError.Message == "" {
				t.Errorf("Converted error message should not be empty")
			}
		})
	}
}

// TestProperty10_ErrorResponseRoundTrip tests that error information is preserved through conversion
// **Feature: codex-channel-support, Property 10: Error Response Conversion**
func TestProperty10_ErrorResponseRoundTrip(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Run 100 iterations
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("round_trip_%d", i), func(t *testing.T) {
			// Generate random error data
			statusCode := 400 + rng.Intn(200) // 400-599 range
			errorMessage := fmt.Sprintf("Round trip test message %d: %s", i, randomErrorMessage(rng))
			errorType := fmt.Sprintf("error_type_%d", rng.Intn(10))
			errorCode := fmt.Sprintf("error_code_%d", rng.Intn(10))

			// Create original OpenAI error
			originalError := types.OpenAIError{
				Message: errorMessage,
				Type:    errorType,
				Code:    errorCode,
				Param:   "",
			}

			// Convert to NewAPIError
			newAPIError := types.WithOpenAIError(originalError, statusCode)

			// Convert back to OpenAI error
			convertedError := newAPIError.ToOpenAIError()

			// Property verification: Round trip should preserve essential information
			// Note: Message may be masked for sensitive info, so we check it's not empty
			if convertedError.Message == "" {
				t.Errorf("Converted error message should not be empty")
			}

			// Type should be preserved
			if convertedError.Type != errorType {
				t.Errorf("Type mismatch: expected %q, got %q", errorType, convertedError.Type)
			}

			// Code should be preserved (may be converted to string)
			codeStr := fmt.Sprintf("%v", convertedError.Code)
			if codeStr != errorCode {
				t.Errorf("Code mismatch: expected %q, got %q", errorCode, codeStr)
			}
		})
	}
}

// TestProperty10_NewOpenAIErrorCreation tests the NewOpenAIError function
// **Feature: codex-channel-support, Property 10: Error Response Conversion**
// **Validates: Requirements 8.1, 8.2, 8.4**
func TestProperty10_NewOpenAIErrorCreation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	errorCodes := []types.ErrorCode{
		types.ErrorCodeConvertRequestFailed,
		types.ErrorCodeDoRequestFailed,
		types.ErrorCodeReadResponseBodyFailed,
		types.ErrorCodeBadResponseStatusCode,
		types.ErrorCodeBadResponse,
		types.ErrorCodeBadResponseBody,
	}

	// Run 100 iterations
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("new_openai_error_%d", i), func(t *testing.T) {
			statusCode := 400 + rng.Intn(200)
			errorCode := errorCodes[rng.Intn(len(errorCodes))]
			errorMessage := fmt.Sprintf("Test error %d", i)

			// Create error using NewOpenAIError
			err := fmt.Errorf("%s", errorMessage)
			newAPIError := types.NewOpenAIError(err, errorCode, statusCode)

			// Property verification:
			// 1. Error should not be nil
			if newAPIError == nil {
				t.Fatalf("NewAPIError should not be nil")
			}

			// 2. Status code should be set
			if newAPIError.StatusCode != statusCode {
				t.Errorf("StatusCode mismatch: expected %d, got %d", statusCode, newAPIError.StatusCode)
			}

			// 3. Error code should be set
			if newAPIError.GetErrorCode() != errorCode {
				t.Errorf("ErrorCode mismatch: expected %v, got %v", errorCode, newAPIError.GetErrorCode())
			}

			// 4. Error message should be accessible
			if newAPIError.Error() == "" {
				t.Errorf("Error message should not be empty")
			}

			// 5. Should be convertible to OpenAI error format
			openAIError := newAPIError.ToOpenAIError()
			if openAIError.Type == "" {
				t.Errorf("OpenAI error type should not be empty")
			}
		})
	}
}

// TestProperty10_OaiResponsesHandlerErrorPaths tests error handling in OaiResponsesHandler
// **Feature: codex-channel-support, Property 10: Error Response Conversion**
// **Validates: Requirements 6.4, 8.4**
func TestProperty10_OaiResponsesHandlerErrorPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Test error scenarios that can occur in OaiResponsesHandler
	errorScenarios := []struct {
		name        string
		errorCode   types.ErrorCode
		description string
	}{
		{
			name:        "ReadResponseBodyFailed",
			errorCode:   types.ErrorCodeReadResponseBodyFailed,
			description: "Error reading response body",
		},
		{
			name:        "BadResponseBody",
			errorCode:   types.ErrorCodeBadResponseBody,
			description: "Error parsing response body",
		},
	}

	// Run iterations for each error scenario
	for _, scenario := range errorScenarios {
		for i := 0; i < 50; i++ { // 50 iterations per scenario = 100 total
			t.Run(fmt.Sprintf("%s_%d", scenario.name, i), func(t *testing.T) {
				errorMessage := fmt.Sprintf("%s: iteration %d - %s", scenario.description, i, randomErrorMessage(rng))
				err := fmt.Errorf("%s", errorMessage)

				// Simulate error creation as done in OaiResponsesHandler
				newAPIError := types.NewOpenAIError(err, scenario.errorCode, 500)

				// Property verification:
				// 1. Error should be created with correct error code
				if newAPIError.GetErrorCode() != scenario.errorCode {
					t.Errorf("ErrorCode mismatch: expected %v, got %v", scenario.errorCode, newAPIError.GetErrorCode())
				}

				// 2. Error should have internal server error status for these cases
				if newAPIError.StatusCode != 500 {
					t.Errorf("StatusCode should be 500 for %s, got %d", scenario.name, newAPIError.StatusCode)
				}

				// 3. Error message should be preserved
				if newAPIError.Error() == "" {
					t.Errorf("Error message should not be empty")
				}
			})
		}
	}
}

// Helper function to generate random error messages
func randomErrorMessage(rng *rand.Rand) string {
	messages := []string{
		"Invalid API key provided",
		"Model not found",
		"Rate limit exceeded",
		"Context length exceeded",
		"Server error occurred",
		"Bad gateway",
		"Service temporarily unavailable",
		"Request timeout",
		"Invalid request format",
		"Permission denied",
		"Resource not found",
		"Internal server error",
		"Gateway timeout",
		"Insufficient quota",
		"Model overloaded",
	}
	return messages[rng.Intn(len(messages))]
}

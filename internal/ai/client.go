package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Layr-Labs/hourglass-avs-template/internal/config"
)

const (
	TogetherAPIURL = "https://api.together.xyz/v1/chat/completions"
)

// Client represents a Together AI client
type Client struct {
	apiKey     string
	httpClient *http.Client
	config     *config.Config
}

// TaskRequest represents an AI task request
type TaskRequest struct {
	TaskType    string  `json:"task_type"`
	Prompt      string  `json:"prompt"`
	Model       string  `json:"model,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
}

// TaskResponse represents an AI task response
type TaskResponse struct {
	Result           string `json:"result"`
	ModelUsed        string `json:"model_used"`
	TokensUsed       int    `json:"tokens_used"`
	ProcessingTimeMs int64  `json:"processing_time_ms"`
}

// TogetherRequest represents the request format for Together AI API
type TogetherRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// TogetherResponse represents the response from Together AI API
type TogetherResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a completion choice
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// NewClient creates a new Together AI client
func NewClient(config *config.Config) *Client {
	return &Client{
		apiKey: config.TogetherAPIKey,
		httpClient: &http.Client{
			Timeout: config.GetTimeout(),
		},
		config: config,
	}
}

// ProcessTask processes an AI task using Together AI
func (c *Client) ProcessTask(ctx context.Context, request *TaskRequest) (*TaskResponse, error) {
	startTime := time.Now()

	// Use default model if not specified
	model := request.Model
	if model == "" {
		model = c.config.DefaultModel
	}

	// Use default values if not specified
	maxTokens := request.MaxTokens
	if maxTokens == 0 {
		maxTokens = c.config.MaxTokens
	}

	temperature := request.Temperature
	if temperature == 0 {
		temperature = c.config.Temperature
	}

	// Create Together AI request
	togetherReq := &TogetherRequest{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: request.Prompt,
			},
		},
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	// Execute request with retries
	var response *TogetherResponse
	var err error

	for attempt := 1; attempt <= c.config.RetryAttempts; attempt++ {
		response, err = c.executeRequest(ctx, togetherReq)
		if err == nil {
			break
		}

		if attempt < c.config.RetryAttempts {
			// Exponential backoff
			backoff := time.Duration(attempt*attempt) * time.Second
			time.Sleep(backoff)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to process task after %d attempts: %w", c.config.RetryAttempts, err)
	}

	// Extract result
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	result := response.Choices[0].Message.Content
	processingTime := time.Since(startTime).Milliseconds()

	return &TaskResponse{
		Result:           result,
		ModelUsed:        response.Model,
		TokensUsed:       response.Usage.TotalTokens,
		ProcessingTimeMs: processingTime,
	}, nil
}

// executeRequest executes a single request to Together AI
func (c *Client) executeRequest(ctx context.Context, request *TogetherRequest) (*TogetherResponse, error) {
	// Marshal request
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", TogetherAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response TogetherResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// ValidateRequest validates a task request
func (c *Client) ValidateRequest(request *TaskRequest) error {
	if request.Prompt == "" {
		return fmt.Errorf("prompt cannot be empty")
	}

	if request.MaxTokens < 0 {
		return fmt.Errorf("max_tokens cannot be negative")
	}

	if request.Temperature < 0 || request.Temperature > 2 {
		return fmt.Errorf("temperature must be between 0 and 2")
	}

	return nil
} 
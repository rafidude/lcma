package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const groqBaseURL = "https://api.groq.com/openai/v1/chat/completions"

// GroqClient handles communication with Groq API
type GroqClient struct {
	apiKey     string
	maxTokens  int
	httpClient *http.Client
}

// GroqMessage represents a single message in the chat
type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GroqRequest represents the request structure for Groq API
type GroqRequest struct {
	Model     string        `json:"model"`
	Messages  []GroqMessage `json:"messages"`
	MaxTokens int           `json:"max_tokens,omitempty"`
}

// GroqResponse represents the response structure from Groq API
type GroqResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      GroqMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type Usage struct {
	QueueTime        float64 `json:"queue_time"`
	PromptTokens     int     `json:"prompt_tokens"`
	PromptTime       float64 `json:"prompt_time"`
	CompletionTokens int     `json:"completion_tokens"`
	CompletionTime   float64 `json:"completion_time"`
	TotalTokens      int     `json:"total_tokens"`
	TotalTime        float64 `json:"total_time"`
}

// NewGroqClient creates a new Groq API client
func NewGroqClient(apiKey string) *GroqClient {
	return &GroqClient{
		apiKey: apiKey,
		// maxTokens: 32768,
		maxTokens: 8000,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateChatCompletion sends a chat completion request to Groq API
func (c *GroqClient) CreateChatCompletion(model string, messages []GroqMessage) (*GroqResponse, error) {
	reqBody := GroqRequest{
		Model:     model,
		Messages:  messages,
		MaxTokens: c.maxTokens,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", groqBaseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var groqResp GroqResponse
	if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &groqResp, nil
}

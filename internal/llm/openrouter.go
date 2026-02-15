package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type OpenRouterClient struct {
	apiKey string
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewOpenRouterClient() (*OpenRouterClient, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY environment variable not set")
	}
	return &OpenRouterClient{apiKey: apiKey}, nil
}

func (c *OpenRouterClient) Query(model, prompt string) (string, error) {
	reqBody := chatRequest{
		Model: model,
		Messages: []message{
			{
				Role:    "system",
				Content: "You are TXTLLM, an LLM queryable over DNS. 255 chars max. ONLY PLAIN TEXT. NO MARKDOWN. NO LINKS. NO URLS. NO MARKDOWN LINKS. ONLY PLAIN TEXT!",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", OpenRouterCompletionURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp openRouterResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("no response received from API")
	}

	return apiResp.Choices[0].Message.Content, nil
}

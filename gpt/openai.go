package gpt

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type chatCompletionRequest struct {
	Model     string    `json:"model"`
	Messages  []message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type message struct {
	Role    string    `json:"role"`
	Content []content `json:"content"`
}

type content struct {
	Type     string   `json:"type"`
	Text     string   `json:"text,omitempty"`
	ImageURL imageURL `json:"image_url,omitempty"`
}

type imageURL struct {
	URL string `json:"url"`
}

type chatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role         string `json:"role"`
			Content      string `json:"content"`
			FinishReason string `json:"finish_reason"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens       int `json:"prompt_tokens"`
		CompletionTokens   int `json:"completion_tokens"`
		TotalTokens        int `json:"total_tokens"`
		PromptTokenDetails struct {
			CachedTokens int `json:"cached_tokens"`
		}
		CompletionTokenDetails struct {
			ReasoningTokens int `json:"reasoning_tokens"`
		}
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

type openAI struct {
	model  string
	apiKey string
}

func NewOpenAI(model, apiKey string) (GPT, error) {
	if model == "" {
		return nil, errors.New("model is not set")
	}
	if apiKey == "" {
		return nil, errors.New("api key is not set")
	}

	return &openAI{
		model:  model,
		apiKey: apiKey,
	}, nil
}

func (g *openAI) Prompt(_ context.Context, prompt string, image []byte) (*Result, error) {
	// build request
	request := chatCompletionRequest{
		Model: g.model,
		Messages: []message{
			{
				Role:    "system",
				Content: []content{{Type: "text", Text: "You are a visual reasoning assistant."}},
			},
			{
				Role: "user",
				Content: []content{
					{Type: "text", Text: prompt},
					{
						Type:     "image_url",
						ImageURL: imageURL{URL: "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)},
					},
				},
			},
		},
		MaxTokens: 300,
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	// send request
	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response := chatCompletionResponse{}
	if err = json.Unmarshal(respBodyBytes, &response); err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 {
		return nil, errors.New("no choices in response")
	}

	content := strings.TrimPrefix(response.Choices[0].Message.Content, "```json")
	content = strings.TrimSuffix(content, "```")

	return &Result{
		Content: content,
		Raw:     response,
	}, nil
}

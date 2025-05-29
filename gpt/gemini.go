package gpt

import (
	"context"
	"errors"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type gemini struct {
	model  string
	apiKey string
}

func NewGemini(model, apiKey string) (GPT, error) {
	if model == "" {
		return nil, errors.New("model is not set")
	}
	if apiKey == "" {
		return nil, errors.New("api key is not set")
	}

	return &gemini{
		model:  model,
		apiKey: apiKey,
	}, nil
}

func (g *gemini) Prompt(ctx context.Context, prompt string, image []byte) (*Result, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(g.apiKey))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	requestParts := []genai.Part{
		genai.Text(prompt),
		genai.ImageData("jpeg", image),
	}

	model := client.GenerativeModel(g.model)
	model.ResponseMIMEType = "application/json"
	resp, err := model.GenerateContent(ctx, requestParts...)
	if err != nil {
		return nil, err
	}

	part, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, errors.New("invalid response type")
	}
	content := strings.TrimPrefix(string(part), "```json")
	content = strings.TrimSuffix(content, "```")

	return &Result{
		Content: content,
		Raw:     resp,
	}, nil
}

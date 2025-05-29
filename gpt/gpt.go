package gpt

import "context"

type GPT interface {
	Prompt(ctx context.Context, prompt string, image []byte) (*Result, error)
}

type Result struct {
	Content string
	Raw     any
}

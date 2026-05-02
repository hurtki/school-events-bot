package ai

import (
	"context"
	"fmt"
)

type NoopGeminiAI struct{}

func NewNoopGeminiAI() *NoopGeminiAI {
	return &NoopGeminiAI{}
}

func (a *NoopGeminiAI) Text(_ context.Context, _ string) (string, error) {
	return "", fmt.Errorf("AI not available")
}

package ai

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/genai"
)

type GeminiAI struct {
	cl    *genai.Client
	model *genai.Model
}

func NewGeminiAI(geminiToken string, modelName string) (*GeminiAI, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	client, err := genai.NewClient(
		ctx,
		&genai.ClientConfig{
			APIKey:  geminiToken,
			Backend: genai.BackendGeminiAPI,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("can't initialize genai client: %w", err)
	}

	m, err := client.Models.Get(ctx, modelName, &genai.GetModelConfig{})
	if err != nil {
		return nil, fmt.Errorf("can't get model '%s': %w", modelName, err)
	}

	return &GeminiAI{
		cl:    client,
		model: m,
	}, nil
}

func (a *GeminiAI) Text(ctx context.Context, prompt string) (string, error) {
	res, err := a.cl.Models.GenerateContent(ctx, a.model.Name, genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("can't generate: %w", err)
	}
	return res.Text(), nil
}

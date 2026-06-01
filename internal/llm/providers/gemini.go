package providers

import (
	"context"
	"fmt"

	"github.com/xSaVageAU/terminal-agent/internal/config"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

func NewGeminiModel(ctx context.Context, modelName string) (model.LLM, error) {
	p := config.CurrentProviders()
	apiKey := p.Gemini.APIKey
	if apiKey == "" {
		return nil, fmt.Errorf("Google API key is not set in config")
	}

	return gemini.NewModel(ctx, modelName, &genai.ClientConfig{
		APIKey: apiKey,
	})
}


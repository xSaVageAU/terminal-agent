package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/xSaVageAU/terminal-agent/internal/config"
	"github.com/xSaVageAU/terminal-agent/internal/llm/providers"
	"google.golang.org/adk/model"
)

type Provider string

const (
	ProviderGemini     Provider = "gemini"
	ProviderOpenRouter Provider = "openrouter"
)

type Config struct {
	Provider Provider
	Model    string
}

type ProviderFactory func(ctx context.Context, modelName string) (model.LLM, error)

var factories = map[Provider]ProviderFactory{
	ProviderGemini:     providers.NewGeminiModel,
	ProviderOpenRouter: providers.NewOpenRouterModel,
}

func LoadFromConfig(cfg *config.Settings) (Config, error) {
	provider, err := resolveProvider(cfg)
	if err != nil {
		return Config{}, err
	}

	p := config.CurrentProviders()
	var modelName string
	switch provider {
	case ProviderOpenRouter:
		modelName = strings.TrimSpace(p.OpenRouter.Model)
		if modelName == "" {
			modelName = "google/gemini-3.1-flash-lite"
		}
	case ProviderGemini:
		modelName = strings.TrimSpace(p.Gemini.Model)
		if modelName == "" {
			modelName = "google/gemini-3.1-flash-lite"
		}
	}

	if modelName == "" {
		return Config{}, fmt.Errorf("model name is empty")
	}

	return Config{Provider: provider, Model: modelName}, nil
}

func resolveProvider(cfg *config.Settings) (Provider, error) {
	if raw := strings.TrimSpace(cfg.Provider); raw != "" {
		switch Provider(strings.ToLower(raw)) {
		case ProviderGemini:
			return ProviderGemini, nil
		case ProviderOpenRouter, "openai":
			return ProviderOpenRouter, nil
		}
	}

	p := config.CurrentProviders()
	if p.OpenRouter.APIKey != "" {
		return ProviderOpenRouter, nil
	}
	if p.Gemini.APIKey != "" {
		return ProviderGemini, nil
	}

	return "", fmt.Errorf("no API key found")
}

func NewModel(ctx context.Context, cfg Config) (model.LLM, error) {
	factory, ok := factories[cfg.Provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider %q", cfg.Provider)
	}
	return factory(ctx, cfg.Model)
}

func (c Config) Describe() string {
	return fmt.Sprintf("%s (%s)", c.Provider, c.Model)
}


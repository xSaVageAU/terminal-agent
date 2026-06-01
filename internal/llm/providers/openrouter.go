package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/xSaVageAU/terminal-agent/internal/config"
	adkopenai "github.com/byebyebruce/adk-go-openai"
	goopenai "github.com/sashabaranov/go-openai"
	"google.golang.org/adk/model"
)

func NewOpenRouterModel(ctx context.Context, modelName string) (model.LLM, error) {
	p := config.CurrentProviders()
	apiKey := p.OpenRouter.APIKey
	if apiKey == "" {
		return nil, fmt.Errorf("OpenRouter API key is not set in config")
	}

	clientCfg := goopenai.DefaultConfig(apiKey)
	clientCfg.BaseURL = "https://openrouter.ai/api/v1"
	if base := strings.TrimSpace(p.OpenRouter.BaseURL); base != "" {
		clientCfg.BaseURL = strings.TrimRight(base, "/")
	}

	m := model.LLM(adkopenai.NewOpenAIModel(modelName, clientCfg))
	return m, nil
}


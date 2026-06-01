package main

import (
	"context"
	"fmt"
	"log"
	"os"

	internalagent "github.com/xSaVageAU/terminal-agent/internal/agent"
	"github.com/xSaVageAU/terminal-agent/internal/config"
	"github.com/xSaVageAU/terminal-agent/internal/llm"
	"github.com/xSaVageAU/terminal-agent/internal/terminal"
	"github.com/xSaVageAU/terminal-agent/tools"

	"google.golang.org/adk/agent/llmagent"
)

func main() {
	cfg, err := config.LoadJSONConfig()
	if err != nil {
		log.Printf("warning: failed to load json config: %v", err)
	}

	// If config is empty (newly generated), write template and exit
	if cfg.Provider == "" {
		fmt.Println("No configuration found. Writing template to ~/.adk-test/config.json and exiting...")
		
		stubSettings := &config.Settings{
			Provider:      "openrouter",
			StreamingMode: "chunk",
			NoColor:       false,
			WorkspaceRoot: "",
		}
		stubProviders := &config.ProviderSettings{
			OpenRouter: config.ProviderConfig{
				Model:  "google/gemini-3.1-flash-lite",
				APIKey: "YOUR_OPENROUTER_API_KEY_HERE",
			},
		}

		if err := config.SaveJSONConfig(stubSettings, stubProviders); err != nil {
			log.Fatalf("failed to write config template: %v", err)
		}
		os.Exit(0)
	}

	ctx := context.Background()

	allTools, err := tools.All()
	if err != nil {
		log.Fatalf("failed to load tools: %v", err)
	}

	llmCfg, err := llm.LoadFromConfig(cfg)
	if err != nil {
		log.Fatalf("failed to load LLM config: %v", err)
	}

	systemPrompt := internalagent.Instruction(llmCfg)

	model, err := llm.NewModel(ctx, llmCfg)
	if err != nil {
		log.Fatalf("failed to create model: %v", err)
	}

	terminal.LogStartup(llmCfg.Describe())

	agentConfig := llmagent.Config{
		Name:        "terminal_explorer",
		Description: "A concise CLI assistant.",
		Model:       model,
		Instruction: systemPrompt,
		Tools:       allTools,
	}

	a, err := llmagent.New(agentConfig)
	if err != nil {
		log.Fatalf("failed to create agent: %v", err)
	}

	if err := terminal.Run(ctx, a); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}


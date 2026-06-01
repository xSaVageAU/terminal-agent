package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PromptModelSelection executes the interactive CLI flow to update provider, model, and keys.
func PromptModelSelection() error {
	cfg := Current()
	p := CurrentProviders()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("==================================================")
	fmt.Println("       Terminal Agent Configuration Tool          ")
	fmt.Println("==================================================")
	fmt.Println()

	// 1. SELECT PROVIDER
	currProvider := strings.TrimSpace(strings.ToLower(cfg.Provider))
	if currProvider == "" {
		currProvider = "openrouter" // sensible default
	}

	fmt.Println("Select LLM Provider:")
	fmt.Println("  1) openrouter")
	fmt.Printf("Enter selection (1) or press Enter to keep [%s]: ", currProvider)

	var chosenProvider string
	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		switch input {
		case "1", "openrouter":
			chosenProvider = "openrouter"
		case "":
			chosenProvider = currProvider
		default:
			lower := strings.ToLower(input)
			if lower == "openrouter" {
				chosenProvider = lower
			} else {
				fmt.Printf("Invalid choice %q, keeping current [%s]\n", input, currProvider)
				chosenProvider = currProvider
			}
		}
	}

	cfg.Provider = chosenProvider
	fmt.Printf("Provider set to: %s\n\n", chosenProvider)

	// 2. CONFIGURE BASED ON PROVIDER
	switch chosenProvider {
	case "openrouter":
		currKey := p.OpenRouter.APIKey
		mask := "empty"
		if currKey != "" {
			mask = "****"
			if len(currKey) > 8 {
				mask = currKey[:4] + "..." + currKey[len(currKey)-4:]
			}
		}
		fmt.Printf("Enter OpenRouter API Key [current: %s]: ", mask)
		if scanner.Scan() {
			keyInput := strings.TrimSpace(scanner.Text())
			if keyInput != "" {
				p.OpenRouter.APIKey = keyInput
				fmt.Println("API Key updated.")
			} else {
				fmt.Println("Keeping current API Key.")
			}
		}

		currModel := strings.TrimSpace(p.OpenRouter.Model)
		if currModel == "" {
			currModel = "google/gemini-3.1-flash-lite"
		}

		fmt.Printf("Enter OpenRouter Model name [current: %s]: ", currModel)
		if scanner.Scan() {
			modelInput := strings.TrimSpace(scanner.Text())
			if modelInput != "" {
				p.OpenRouter.Model = modelInput
				fmt.Printf("Model set to: %s\n", modelInput)
			} else {
				p.OpenRouter.Model = currModel
				fmt.Printf("Keeping model: %s\n", currModel)
			}
		}
	}

	fmt.Println()
	fmt.Println("Saving configuration...")
	if err := SaveJSONConfig(cfg, p); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("Configuration saved successfully!")
	fmt.Println("==================================================")
	return nil
}

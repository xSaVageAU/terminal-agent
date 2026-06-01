package configeditor

import (
	"fmt"
	"strings"

	"charm.land/huh/v2"
	"github.com/xSaVageAU/terminal-agent/internal/config"
)

const defaultModel = "google/gemini-3.1-flash-lite"
const customModelSentinel = "__custom__"

var modelOptions = []string{
	"google/gemini-3.1-flash-lite",
	"google/gemini-2.5-flash",
	"google/gemini-2.5-pro",
	"anthropic/claude-sonnet-4",
	"anthropic/claude-opus-4",
	"openai/gpt-4o",
	"openai/gpt-4o-mini",
	"openai/o3",
	"openai/o3-mini",
	"deepseek/deepseek-chat",
	"deepseek/deepseek-r1",
	"meta-llama/llama-4-maverick",
	"meta-llama/llama-4-scout",
	"Custom…",
}

func modelSelectOpts() []huh.Option[string] {
	opts := make([]huh.Option[string], len(modelOptions))
	for i, m := range modelOptions {
		if m == "Custom…" {
			opts[i] = huh.NewOption("Custom… (type your own)", customModelSentinel)
		} else {
			opts[i] = huh.NewOption(m, m)
		}
	}
	return opts
}

func Run() error {
	cfg, err := config.LoadJSONConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	providers := config.CurrentProviders()

	streamingMode := cfg.StreamingMode
	if streamingMode == "" {
		streamingMode = "chunk"
	}

	provider := cfg.Provider
	if provider == "" {
		provider = "openrouter"
	}

	modelSelect := providers.OpenRouter.Model
	if modelSelect == "" {
		modelSelect = defaultModel
	}
	customModel := ""
	isCustom := true
	for _, m := range modelOptions[:len(modelOptions)-1] {
		if strings.EqualFold(m, modelSelect) {
			isCustom = false
			break
		}
	}
	if isCustom {
		modelSelect = customModelSentinel
		customModel = providers.OpenRouter.Model
	}

	apiKey := providers.OpenRouter.APIKey
	baseURL := providers.OpenRouter.BaseURL
	if baseURL == "" {
		baseURL = "https://openrouter.ai/api/v1"
	}

	theme := huh.ThemeFunc(huh.ThemeCatppuccin)

	for {
		var action string
		menu := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Configuration Menu").
					Options(
						huh.NewOption("Provider (OpenRouter only)", "provider"),
						huh.NewOption("Model", "model"),
						huh.NewOption("API Key & Base URL", "api"),
						huh.NewOption("Streaming Mode", "streaming"),
						huh.NewOption("Workspace Root", "workspace"),
						huh.NewOption("Display Settings", "display"),
						huh.NewOption("Save & Exit", "save"),
					).
					Value(&action),
			),
		)
		menu.WithTheme(theme)

		if err := menu.Run(); err != nil {
			return err
		}
		if action == "" {
			fmt.Println("\nDiscarded.")
			return nil
		}

		switch action {
		case "provider":
			f := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Provider").
						Options(
							huh.NewOption("OpenRouter", "openrouter"),
						).
						Value(&provider),
				),
			)
			f.WithTheme(theme)
			if err := f.Run(); err != nil {
				return err
			}

		case "model":
			f := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Model").
						Description("Pick a known model or choose Custom to type your own.").
						Options(modelSelectOpts()...).
						Value(&modelSelect),
				),
				huh.NewGroup(
					huh.NewInput().
						Title("Custom Model ID").
						Description("Enter the model identifier (e.g. org/model-name).").
						Value(&customModel).
						Placeholder(defaultModel),
				).WithHideFunc(func() bool {
					return modelSelect != customModelSentinel
				}),
			)
			f.WithTheme(theme)
			if err := f.Run(); err != nil {
				return err
			}

		case "api":
			f := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("API Key").
						Description("OpenRouter API key.").
						Value(&apiKey).
						EchoMode(huh.EchoModePassword),
				),
				huh.NewGroup(
					huh.NewInput().
						Title("Base URL").
						Description("OpenRouter API endpoint.").
						Value(&baseURL).
						Placeholder("https://openrouter.ai/api/v1"),
				),
			)
			f.WithTheme(theme)
			if err := f.Run(); err != nil {
				return err
			}

		case "streaming":
			f := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Streaming Mode").
						Description("How LLM responses are streamed.").
						Options(
							huh.NewOption("chunk — word-by-word", "chunk"),
							huh.NewOption("sentence — line-by-line", "sentence"),
							huh.NewOption("none", ""),
						).
						Value(&streamingMode),
				),
			)
			f.WithTheme(theme)
			if err := f.Run(); err != nil {
				return err
			}

		case "workspace":
			f := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Workspace Root").
						Description("Absolute path to restrict file operations (empty = CWD).").
						Value(&cfg.WorkspaceRoot),
				),
			)
			f.WithTheme(theme)
			if err := f.Run(); err != nil {
				return err
			}

		case "display":
			f := huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title("No Color").
						Description("Disable colored terminal output.").
						Value(&cfg.NoColor),
				),
			)
			f.WithTheme(theme)
			if err := f.Run(); err != nil {
				return err
			}

		case "save":
			if modelSelect == customModelSentinel {
				modelSelect = customModel
			}

			cfg.Provider = provider
			cfg.StreamingMode = streamingMode
			newProviders := &config.ProviderSettings{
				OpenRouter: config.ProviderConfig{
					Model:   strings.TrimSpace(modelSelect),
					APIKey:  strings.TrimSpace(apiKey),
					BaseURL: strings.TrimSpace(baseURL),
				},
			}

			if err := config.SaveJSONConfig(cfg, newProviders); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Println("\nConfig saved to ~/.adk-test/")
			return nil
		}
	}
}
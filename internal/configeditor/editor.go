package configeditor

import (
	"errors"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/huh/v2"
	"github.com/xSaVageAU/terminal-agent/internal/config"
)

const defaultModel = "google/gemini-3.1-flash-lite"
const customModelSentinel = "__custom__"

var modelOptions = []string{
	"google/gemini-3.5-flash",
	"google/gemini-3.1-flash-lite",
	"google/gemini-3-flash-preview",
	"google/gemini-3.1-pro-preview",
	"google/gemma-4-31b-it",
	"google/gemma-4-26b-a4b-it",
	"anthropic/claude-opus-4.8",
	"anthropic/claude-opus-4.6",
	"anthropic/claude-sonnet-4.6",
	"anthropic/claude-haiku-4.5",
	"openai/gpt-5.5-pro",
	"openai/gpt-5.5",
	"openai/gpt-5.4",
	"openai/gpt-5.4-mini",
	"openai/gpt-5.4-nano",
	"openai/gpt-oss-120b",
	"openai/gpt-oss-20b",
	"deepseek/deepseek-v4-flash",
	"deepseek/deepseek-v4-pro",
	"tencent/hy3-preview",
	"openrouter/owl-alpha",
	"nvidia/nemotron-3-super-120b-a12b:free",
	"poolside/laguna-m.1:free",
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

	terminalTUI := cfg.TerminalTUI
	if terminalTUI == "" {
		terminalTUI = "classic"
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

	newForm := func(groups ...*huh.Group) *huh.Form {
		f := huh.NewForm(groups...)
		f.WithTheme(theme)
		f.WithShowHelp(true)
		km := huh.NewDefaultKeyMap()
		km.Quit = key.NewBinding(key.WithKeys("ctrl+c", "esc"))
		f.WithKeyMap(km)
		return f
	}

	clearScreen := func() {
		fmt.Print("\033[2J\033[H")
	}

	for {
		var action string
		g := huh.NewGroup(
			huh.NewSelect[string]().
				Title("Configuration Menu").
				Options(
					huh.NewOption("Provider (OpenRouter only)", "provider"),
					huh.NewOption("Model", "model"),
					huh.NewOption("API Key & Base URL", "api"),
					huh.NewOption("Streaming Mode", "streaming"),
					huh.NewOption("Workspace Root", "workspace"),
					huh.NewOption("Display Settings", "display"),
					huh.NewOption("Terminal TUI Theme", "theme"),
					huh.NewOption("Save & Exit", "save"),
				).
				Value(&action),
		)
		menu := newForm(g)

		clearScreen()
		if err := menu.Run(); err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return nil
			}
			return err
		}
		if action == "" {
			return nil
		}

		switch action {
		case "provider":
			g := huh.NewGroup(
				huh.NewSelect[string]().
					Title("Provider").
					Description("Esc to go back").
					Options(
						huh.NewOption("OpenRouter", "openrouter"),
					).
					Value(&provider),
			)
			clearScreen()
			if err := newForm(g).Run(); err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					continue
				}
				return err
			}

		case "model":
			g1 := huh.NewGroup(
				huh.NewSelect[string]().
					Title("Model").
					Description("Pick a known model or choose Custom to type your own. Esc to go back").
					Options(modelSelectOpts()...).
					Value(&modelSelect),
			)
			g2 := huh.NewGroup(
				huh.NewInput().
					Title("Custom Model ID").
					Description("Enter the model identifier (e.g. org/model-name).").
					Value(&customModel).
					Placeholder(defaultModel),
			).WithHideFunc(func() bool {
				return modelSelect != customModelSentinel
			})
			clearScreen()
			if err := newForm(g1, g2).Run(); err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					continue
				}
				return err
			}

		case "api":
			g1 := huh.NewGroup(
				huh.NewInput().
					Title("API Key").
					Description("OpenRouter API key. Esc to go back").
					Value(&apiKey).
					EchoMode(huh.EchoModePassword),
			)
			g2 := huh.NewGroup(
				huh.NewInput().
					Title("Base URL").
					Description("OpenRouter API endpoint.").
					Value(&baseURL).
					Placeholder("https://openrouter.ai/api/v1"),
			)
			clearScreen()
			if err := newForm(g1, g2).Run(); err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					continue
				}
				return err
			}

		case "streaming":
			g := huh.NewGroup(
				huh.NewSelect[string]().
					Title("Streaming Mode").
					Description("How LLM responses are streamed. Esc to go back").
					Options(
						huh.NewOption("chunk — word-by-word", "chunk"),
						huh.NewOption("sentence — line-by-line", "sentence"),
						huh.NewOption("none", ""),
					).
					Value(&streamingMode),
			)
			clearScreen()
			if err := newForm(g).Run(); err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					continue
				}
				return err
			}

		case "workspace":
			g := huh.NewGroup(
				huh.NewInput().
					Title("Workspace Root").
					Description("Absolute path to restrict file operations (empty = CWD). Esc to go back").
					Value(&cfg.WorkspaceRoot),
			)
			clearScreen()
			if err := newForm(g).Run(); err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					continue
				}
				return err
			}

		case "display":
			g := huh.NewGroup(
				huh.NewConfirm().
					Title("No Color").
					Description("Disable colored terminal output. Esc to go back").
					Value(&cfg.NoColor),
			)
			clearScreen()
			if err := newForm(g).Run(); err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					continue
				}
				return err
			}

		case "theme":
			g := huh.NewGroup(
				huh.NewSelect[string]().
					Title("Terminal TUI Theme").
					Description("Choose your interactive terminal interface style. Esc to go back").
					Options(
						huh.NewOption("classic (Standard colorful CLI)", "classic"),
						huh.NewOption("bubbles (Rich Bubble Tea TUI)", "bubbles"),
					).
					Value(&terminalTUI),
			)
			clearScreen()
			if err := newForm(g).Run(); err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					continue
				}
				return err
			}

		case "save":
			if modelSelect == customModelSentinel {
				modelSelect = customModel
			}

			cfg.Provider = provider
			cfg.StreamingMode = streamingMode
			cfg.TerminalTUI = terminalTUI
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

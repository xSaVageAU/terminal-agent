package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ProviderConfig struct {
	Model   string `json:"model"`
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
}

type ProviderSettings struct {
	OpenRouter ProviderConfig `json:"openrouter"`
}

type Settings struct {
	Provider      string `json:"provider"`
	StreamingMode string `json:"streaming_mode"`
	NoColor       bool   `json:"no_color"`
	WorkspaceRoot string `json:"workspace_root"`
	TerminalTUI   string `json:"terminal_tui"`
}

var current *Settings = &Settings{}
var currentProviders *ProviderSettings = &ProviderSettings{}

func Current() *Settings { return current }
func CurrentProviders() *ProviderSettings { return currentProviders }

func LoadJSONConfig() (*Settings, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(home, ".adk-test")
	_ = os.MkdirAll(dir, 0755)

	// Load config.json
	path := filepath.Join(dir, "config.json")
	data, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(data, current); err != nil {
			return nil, fmt.Errorf("invalid config.json: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	// Load providers.json
	pPath := filepath.Join(dir, "providers.json")
	pData, err := os.ReadFile(pPath)
	if err == nil {
		if err := json.Unmarshal(pData, currentProviders); err != nil {
			return nil, fmt.Errorf("invalid providers.json: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	return current, nil
}

func SaveJSONConfig(s *Settings, p *ProviderSettings) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, ".adk-test")
	_ = os.MkdirAll(dir, 0755)

	// Save config.json
	data, _ := json.MarshalIndent(s, "", "  ")
	_ = os.WriteFile(filepath.Join(dir, "config.json"), data, 0644)

	// Save providers.json
	pData, _ := json.MarshalIndent(p, "", "  ")
	return os.WriteFile(filepath.Join(dir, "providers.json"), pData, 0644)
}

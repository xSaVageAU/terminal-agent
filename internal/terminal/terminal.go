package terminal

import (
	"context"
	"log"
	"strings"

	"github.com/xSaVageAU/terminal-agent/internal/config"
	"github.com/xSaVageAU/terminal-agent/internal/terminal/bubbles"
	"github.com/xSaVageAU/terminal-agent/internal/terminal/classic"
	"google.golang.org/adk/agent"
)

// Terminal defines the interface that all terminal/TUI themes must implement.
type Terminal interface {
	Run(ctx context.Context, a agent.Agent) error
}

// New returns the configured Terminal implementation.
func New() Terminal {
	cfg := config.Current()
	if cfg == nil {
		return &classic.Terminal{}
	}

	switch strings.ToLower(strings.TrimSpace(cfg.TerminalTUI)) {
	case "bubbles":
		return &bubbles.Terminal{}
	case "classic", "default", "":
		return &classic.Terminal{}
	default:
		return &classic.Terminal{}
	}
}

// ShouldUse returns true when the app should run the custom terminal instead of ADK's stock console.
func ShouldUse(args []string) bool {
	if len(args) == 0 {
		return true
	}
	return args[0] == "console"
}

// StripConsoleKeyword removes a leading "console" subcommand token if present.
func StripConsoleKeyword(args []string) []string {
	if len(args) > 0 && args[0] == "console" {
		return args[1:]
	}
	return args
}

// LogStartup prints which model is active.
func LogStartup(desc string) {
	log.Printf("Using %s", desc)
}



package terminal

import (
	"fmt"

	"github.com/xSaVageAU/terminal-agent/internal/config"
)

// ANSI styling (Windows 10+ console and Windows Terminal support this).
const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	cyan    = "\033[36m"
	yellow  = "\033[33m"
	green   = "\033[32m"
	red     = "\033[31m"
	magenta = "\033[35m"
	blue    = "\033[34m"
)

func paint(code, s string) string {
	if config.Current().NoColor {
		return s
	}
	return code + s + reset
}

func userPrompt()  { fmt.Print(paint(cyan+bold, "\nYou → ")) }
func agentHeader() { fmt.Print(paint(blue+bold, "\nAgent → ")) }

func printToolCall(name string, args string) {
	fmt.Printf(paint(yellow+bold, "▶ ")+paint(yellow+bold, "%s"), name)
	if args != "" {
		fmt.Printf(paint(dim, "  %s"), args)
	}
	fmt.Println()
}

func printToolResult(name string, body string) {
	fmt.Printf(paint(green+bold, "✓ ")+paint(green+bold, "%s"), name)
	if body != "" {
		fmt.Printf(paint(dim, "  %s"), body)
	}
	fmt.Println()
}

func printConfirmPrompt(toolName, args string) {
	fmt.Printf(paint(magenta+bold, "? %s"), toolName)
	if args != "" {
		fmt.Printf(paint(dim, "  %s"), args)
	}
	fmt.Print(paint(magenta+bold, "\n  Approve? [y/N]: "))
}

func printError(msg string) { fmt.Println(paint(red+bold, "✗ "+msg)) }


package agent

import (
	"fmt"
	"os"
	"runtime"

	"github.com/xSaVageAU/terminal-agent/internal/llm"
)

// Instruction returns the system prompt for the terminal agent.
func Instruction(cfg llm.Config) string {
	osName := runtime.GOOS
	arch := runtime.GOARCH
	numCPU := runtime.NumCPU()

	// Determine shell preference
	shell := "bash"
	if osName == "windows" {
		shell = "powershell"
	}
	if envShell := os.Getenv("SHELL"); envShell != "" {
		shell = envShell
	} else if envComspec := os.Getenv("COMSPEC"); envComspec != "" {
		shell = envComspec
	}

	cwd, _ := os.Getwd()
	home, _ := os.UserHomeDir()

	return fmt.Sprintf(`You are terminal_agent, a friendly and helpful CLI companion.

Runtime: provider=%s, model=%s.
Environment: OS=%s, Arch=%s, CPUs=%d, Shell=%s, CWD=%s, Home=%s.

Directives:
- Be helpful and friendly, while staying efficient.
- Avoid using excessive emojis in chat responses.
- Feel free to provide brief summaries of what you've done to keep the user in the loop.
- Use tools to solve problems, but don't be afraid to explain the "why" when it helps.
- Focus on making the developer's life easier with clear guidance.
- If something is unclear, just ask—I'm here to help!
- Maintain a relaxed, approachable, and professional tone.`,
		cfg.Provider, cfg.Model, osName, arch, numCPU, shell, cwd, home)
}

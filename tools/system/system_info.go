package system

import (
	"runtime"
	"strings"

	"github.com/xSaVageAU/terminal-agent/internal/config"
	"github.com/xSaVageAU/terminal-agent/tools/workspace"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type SystemInfoArgs struct{}

type SystemInfoResult struct {
	OS          string `json:"os"`
	Arch        string `json:"arch"`
	NumCPU      int    `json:"num_cpu"`
	GoVersion   string `json:"go_version"`
	Workspace   string `json:"workspace_root"`
	LLMProvider string `json:"llm_provider_hint"`
}

func systemInfo(_ tool.Context, _ SystemInfoArgs) (SystemInfoResult, error) {
	cfg := config.Current()
	p := config.CurrentProviders()
	provider := strings.TrimSpace(cfg.Provider)
	if provider == "" {
		switch {
		case p.OpenRouter.APIKey != "":
			provider = "openrouter (auto)"
		case p.Gemini.APIKey != "":
			provider = "gemini (auto)"
		default:
			provider = "unknown"
		}
	}
	return SystemInfoResult{
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		NumCPU:      runtime.NumCPU(),
		GoVersion:   runtime.Version(),
		Workspace:   workspace.WorkspaceRoot(),
		LLMProvider: provider,
	}, nil
}

func SystemInfoTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "system_info",
		Description: "Returns OS, CPU count, Go version, workspace path, and configured LLM provider hint.",
	}, systemInfo)
}


package files

import (
	"os"
	"strings"

	"github.com/xSaVageAU/terminal-agent/tools/workspace"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type ReadFileArgs struct {
	Path     string `json:"path" jsonschema:"File path relative to the workspace root."`
	MaxLines *int   `json:"max_lines,omitempty" jsonschema:"Maximum number of lines to return (optional, default 200, max 500)."`
}

type ReadFileResult struct {
	Path      string `json:"path"`
	Content   string `json:"content"`
	Truncated bool   `json:"truncated,omitempty"`
}

func readFile(_ tool.Context, args ReadFileArgs) (ReadFileResult, error) {
	abs, err := workspace.ResolveInWorkspace(args.Path)
	if err != nil {
		return ReadFileResult{}, err
	}

	lock := getFileLock(abs)
	lock.Lock()
	defer lock.Unlock()

	data, err := os.ReadFile(abs)
	if err != nil {
		return ReadFileResult{}, err
	}
	maxVal := 200
	if args.MaxLines != nil {
		maxVal = *args.MaxLines
	}
	max := workspace.ClampInt(maxVal, 200, 500)
	lines := strings.Split(string(data), "\n")
	truncated := len(lines) > max
	if truncated {
		lines = lines[:max]
	}
	return ReadFileResult{
		Path:      abs,
		Content:   strings.Join(lines, "\n"),
		Truncated: truncated,
	}, nil
}

func ReadFileTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "read_file",
		Description: "Reads a text file from the workspace (truncated to max_lines).",
	}, readFile)
}


package files

import (
	"os"
	"path/filepath"

	"github.com/xSaVageAU/terminal-agent/tools/workspace"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type MoveFileArgs struct {
	Source      string `json:"source" jsonschema:"Source file path relative to the workspace root."`
	Destination string `json:"destination" jsonschema:"Destination path relative to the workspace root."`
}

type MoveFileResult struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func moveFile(_ tool.Context, args MoveFileArgs) (MoveFileResult, error) {
	src, err := workspace.ResolveInWorkspace(args.Source)
	if err != nil {
		return MoveFileResult{}, err
	}
	dst, err := workspace.ResolveInWorkspace(args.Destination)
	if err != nil {
		return MoveFileResult{}, err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return MoveFileResult{}, err
	}
	if err := os.Rename(src, dst); err != nil {
		return MoveFileResult{}, err
	}
	return MoveFileResult{Source: src, Destination: dst}, nil
}

func MoveFileTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:                "move_file",
		Description:         "Moves or renames a file/directory. Requires user confirmation.",
		RequireConfirmation: true,
	}, moveFile)
}


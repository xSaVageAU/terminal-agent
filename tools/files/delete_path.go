package files

import (
	"os"

	"github.com/xSaVageAU/terminal-agent/tools/workspace"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type DeletePathArgs struct {
	Path string `json:"path" jsonschema:"File or directory path to delete, relative to workspace root."`
}

type DeletePathResult struct {
	Path string `json:"path"`
}

func deletePath(_ tool.Context, args DeletePathArgs) (DeletePathResult, error) {
	abs, err := workspace.ResolveInWorkspace(args.Path)
	if err != nil {
		return DeletePathResult{}, err
	}
	if err := os.RemoveAll(abs); err != nil {
		return DeletePathResult{}, err
	}
	return DeletePathResult{Path: abs}, nil
}

func DeletePathTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:                "delete_path",
		Description:         "Permanently deletes a file or directory. Requires user confirmation.",
		RequireConfirmation: true,
	}, deletePath)
}


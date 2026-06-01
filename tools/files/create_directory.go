package files

import (
	"os"

	"github.com/xSaVageAU/terminal-agent/tools/workspace"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type CreateDirArgs struct {
	Path string `json:"path" jsonschema:"Directory path relative to the workspace root."`
}

type CreateDirResult struct {
	Path string `json:"path"`
}

func createDirectory(_ tool.Context, args CreateDirArgs) (CreateDirResult, error) {
	abs, err := workspace.ResolveInWorkspace(args.Path)
	if err != nil {
		return CreateDirResult{}, err
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return CreateDirResult{}, err
	}
	return CreateDirResult{Path: abs}, nil
}

func CreateDirectoryTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:                "create_directory",
		Description:         "Creates a directory and any necessary parents. Requires user confirmation.",
		RequireConfirmation: true,
	}, createDirectory)
}


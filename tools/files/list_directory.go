package files

import (
	"os"

	"github.com/xSaVageAU/terminal-agent/tools/workspace"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type ListDirArgs struct {
	Path string `json:"path" jsonschema:"Directory path relative to the workspace root, or empty for the workspace root."`
}

type ListDirEntry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
	Size  int64  `json:"size_bytes,omitempty"`
}

type ListDirResult struct {
	Path    string         `json:"path"`
	Entries []ListDirEntry `json:"entries"`
}

func listDirectory(_ tool.Context, args ListDirArgs) (ListDirResult, error) {
	abs, err := workspace.ResolveInWorkspace(args.Path)
	if err != nil {
		return ListDirResult{}, err
	}
	entries, err := os.ReadDir(abs)
	if err != nil {
		return ListDirResult{}, err
	}
	out := make([]ListDirEntry, 0, len(entries))
	for _, e := range entries {
		item := ListDirEntry{Name: e.Name(), IsDir: e.IsDir()}
		if !e.IsDir() {
			if info, err := e.Info(); err == nil {
				item.Size = info.Size()
			}
		}
		out = append(out, item)
	}
	return ListDirResult{Path: abs, Entries: out}, nil
}

func ListDirectoryTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "list_directory",
		Description: "Lists files and folders in a directory within the agent workspace.",
	}, listDirectory)
}


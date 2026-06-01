package files

import (
	"os"
	"path/filepath"

	"github.com/xSaVageAU/terminal-agent/tools/workspace"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type WriteFileArgs struct {
	Path    string `json:"path" jsonschema:"File path relative to the workspace root."`
	Content string `json:"content" jsonschema:"Full text content to write."`
	Append  bool   `json:"append" jsonschema:"If true, append to the file instead of overwriting."`
}

type WriteFileResult struct {
	Path         string `json:"path"`
	BytesWritten int    `json:"bytes_written"`
	Mode         string `json:"mode"`
}

func writeFile(_ tool.Context, args WriteFileArgs) (WriteFileResult, error) {
	abs, err := workspace.ResolveInWorkspace(args.Path)
	if err != nil {
		return WriteFileResult{}, err
	}
	mode := "overwrite"
	var n int
	if args.Append {
		f, err := os.OpenFile(abs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return WriteFileResult{}, err
		}
		n, err = f.WriteString(args.Content)
		_ = f.Close()
		if err != nil {
			return WriteFileResult{}, err
		}
		mode = "append"
	} else {
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			return WriteFileResult{}, err
		}
		if err := os.WriteFile(abs, []byte(args.Content), 0o644); err != nil {
			return WriteFileResult{}, err
		}
		n = len(args.Content)
	}
	return WriteFileResult{Path: abs, BytesWritten: n, Mode: mode}, nil
}

func WriteFileTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:                "write_file",
		Description:         "Creates or overwrites a text file in the workspace. Requires user confirmation in the terminal.",
		RequireConfirmation: true,
	}, writeFile)
}


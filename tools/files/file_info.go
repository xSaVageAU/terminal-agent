package files

import (
	"os"
	"time"

	"github.com/xSaVageAU/terminal-agent/tools/workspace"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type FileInfoArgs struct {
	Path string `json:"path" jsonschema:"File or directory path relative to the workspace root."`
}

type FileInfoResult struct {
	Path    string `json:"path"`
	IsDir   bool   `json:"is_dir"`
	Size    int64  `json:"size_bytes"`
	ModTime string `json:"mod_time"`
	Mode    string `json:"file_mode"`
}

func fileInfo(_ tool.Context, args FileInfoArgs) (FileInfoResult, error) {
	abs, err := workspace.ResolveInWorkspace(args.Path)
	if err != nil {
		return FileInfoResult{}, err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return FileInfoResult{}, err
	}
	return FileInfoResult{
		Path:    abs,
		IsDir:   info.IsDir(),
		Size:    info.Size(),
		ModTime: info.ModTime().Format(time.RFC3339),
		Mode:    info.Mode().String(),
	}, nil
}

func FileInfoTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "file_info",
		Description: "Returns size, modification time, and type for a file or directory in the workspace.",
	}, fileInfo)
}


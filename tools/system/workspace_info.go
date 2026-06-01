package system

import (
	"runtime"

	"github.com/xSaVageAU/terminal-agent/tools/workspace"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type WorkspaceInfoArgs struct{}

type WorkspaceInfoResult struct {
	Root      string `json:"workspace_root"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	GoVersion string `json:"go_version"`
}

func workspaceInfo(_ tool.Context, _ WorkspaceInfoArgs) (WorkspaceInfoResult, error) {
	return WorkspaceInfoResult{
		Root:      workspace.WorkspaceRoot(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		GoVersion: runtime.Version(),
	}, nil
}

func WorkspaceInfoTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "workspace_info",
		Description: "Returns workspace root path and runtime environment metadata.",
	}, workspaceInfo)
}


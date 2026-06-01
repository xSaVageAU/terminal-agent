// Package tools defines custom ADK function tools for the terminal agent.
package tools

import (
	"github.com/xSaVageAU/terminal-agent/tools/files"
	"github.com/xSaVageAU/terminal-agent/tools/network"
	"github.com/xSaVageAU/terminal-agent/tools/shell"
	"github.com/xSaVageAU/terminal-agent/tools/system"
	"google.golang.org/adk/tool"
)

// All returns every custom tool for registration on the agent.
func All() ([]tool.Tool, error) {
	builders := []func() (tool.Tool, error){
		system.GetTimeTool,
		system.SystemInfoTool,
		system.WorkspaceInfoTool,
		files.ListDirectoryTool,
		files.ReadFileTool,
		files.WriteFileTool,
		files.FileInfoTool,
		files.SearchTextTool,
		files.CreateDirectoryTool,
		files.MoveFileTool,
		network.FetchURLTool,
		shell.RunCommandTool,
		files.DeletePathTool,
		files.PatchFileTool,
	}
	out := make([]tool.Tool, 0, len(builders))
	for _, b := range builders {
		t, err := b()
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}


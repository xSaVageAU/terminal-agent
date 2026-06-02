package files

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xSaVageAU/terminal-agent/tools/workspace"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type PatchFileArgs struct {
	Path string `json:"path" jsonschema:"Path to the file relative to the workspace root."`
	Line int    `json:"line" jsonschema:"1-indexed line number."`
	Mode string `json:"mode" jsonschema:"Mode: 'insert' or 'replace'."`
	Text string `json:"text" jsonschema:"The content to put at that line."`
}

func patchFile(_ tool.Context, args PatchFileArgs) (string, error) {
	abs, err := workspace.ResolveInWorkspace(args.Path)
	if err != nil {
		return "", err
	}

	lock := getFileLock(abs)
	lock.Lock()
	defer lock.Unlock()

	file, err := os.Open(abs)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if args.Line < 1 || args.Line > len(lines)+1 {
		return "", fmt.Errorf("line number %d out of range (max %d)", args.Line, len(lines)+1)
	}

	idx := args.Line - 1
	if args.Mode == "insert" {
		if idx >= len(lines) {
			lines = append(lines, args.Text)
		} else {
			lines = append(lines[:idx], append([]string{args.Text}, lines[idx:]...)...)
		}
	} else if args.Mode == "replace" {
		if idx >= len(lines) {
			return "", fmt.Errorf("cannot replace: line %d does not exist", args.Line)
		}
		lines[idx] = args.Text
	} else {
		return "", fmt.Errorf("invalid mode: %s", args.Mode)
	}

	err = os.WriteFile(abs, []byte(strings.Join(lines, "\n")+"\n"), 0644)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Successfully applied %s at line %d", args.Mode, args.Line), nil
}

func PatchFileTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "patch_file",
		Description: "Insert or replace a line in a file at a specific line number.",
	}, patchFile)
}


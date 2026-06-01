package shell

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/xSaVageAU/terminal-agent/tools/workspace"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type RunCommandArgs struct {
	Command string `json:"command" jsonschema:"Shell command to run. Executed in the workspace directory."`
	Timeout int    `json:"timeout_seconds" jsonschema:"Max runtime in seconds (default 30, max 120)."`
}

type RunCommandResult struct {
	Command  string `json:"command"`
	ExitCode int    `json:"exit_code"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	TimedOut bool   `json:"timed_out,omitempty"`
}

const maxCommandOutput = 64 * 1024

var blockedCommandSubstrings = []string{
	"rm -rf", "rm -fr", "del /f", "del /s", "format ", "mkfs",
	"shutdown", "reboot", ":(){", "dd if=", "> /dev/sd",
}

func runCommand(_ tool.Context, args RunCommandArgs) (RunCommandResult, error) {
	cmdLine := strings.TrimSpace(args.Command)
	if cmdLine == "" {
		return RunCommandResult{}, fmt.Errorf("command must not be empty")
	}
	lower := strings.ToLower(cmdLine)
	for _, blocked := range blockedCommandSubstrings {
		if strings.Contains(lower, blocked) {
			return RunCommandResult{}, fmt.Errorf("command blocked for safety: contains %q", blocked)
		}
	}

	ws, err := workspace.ResolveInWorkspace(".")
	if err != nil {
		return RunCommandResult{}, err
	}

	timeout := workspace.ClampInt(args.Timeout, 30, 120)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", cmdLine)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", cmdLine)
	}
	cmd.Dir = ws

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()
	timedOut := ctx.Err() == context.DeadlineExceeded

	exitCode := 0
	if runErr != nil {
		if exitErr, ok := runErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else if !timedOut {
			return RunCommandResult{}, runErr
		}
	}

	return RunCommandResult{
		Command:  cmdLine,
		ExitCode: exitCode,
		Stdout:   trimOutput(stdout.String()),
		Stderr:   trimOutput(stderr.String()),
		TimedOut: timedOut,
	}, nil
}

func trimOutput(s string) string {
	if len(s) <= maxCommandOutput {
		return s
	}
	return s[:maxCommandOutput] + "\n... (output truncated)"
}

func RunCommandTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:                "run_command",
		Description:         "Runs a shell command in the workspace directory (stdout/stderr capped). Requires user confirmation in the terminal.",
		RequireConfirmation: true,
	}, runCommand)
}


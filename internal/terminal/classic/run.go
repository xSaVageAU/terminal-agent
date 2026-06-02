package classic

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/xSaVageAU/terminal-agent/internal/config"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool/toolconfirmation"
	"google.golang.org/genai"
)

type Terminal struct{}

// Run starts the interactive REPL for the given root agent.
func (t *Terminal) Run(ctx context.Context, root agent.Agent) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	const userID, appName = "console_user", "console_app"

	sessionService := session.InMemoryService()
	resp, err := sessionService.Create(ctx, &session.CreateRequest{
		AppName: appName,
		UserID:  userID,
	})
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}

	r, err := runner.New(runner.Config{
		AppName:        appName,
		Agent:          root,
		SessionService: sessionService,
	})
	if err != nil {
		return fmt.Errorf("create runner: %w", err)
	}

	streamingMode := resolveStreamingMode()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println(paint(dim, "ADK terminal agent — tool calls and results are shown in color."))
	fmt.Println(paint(dim, "Type 'exit' or 'quit' to leave. Type 'yolo' to toggle auto-approval of tool confirmations."))
	userPrompt()

	yoloMode := false

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println()
				return nil
			}
			return err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			userPrompt()
			continue
		}
		if isExit(line) {
			fmt.Println(paint(dim, "Goodbye."))
			return nil
		}
		if strings.EqualFold(line, "yolo") {
			yoloMode = !yoloMode
			if yoloMode {
				fmt.Println(paint(magenta+bold, "⚡ YOLO Mode ON — tool confirmations will be auto-approved."))
			} else {
				fmt.Println(paint(dim, "  YOLO Mode OFF — tool confirmations will be prompted."))
			}
			userPrompt()
			continue
		}

		msg := genai.NewContentFromText(line+"\n", genai.RoleUser)
		if err := runAgentTurn(ctx, r, reader, userID, resp.Session.ID(), msg, streamingMode, true, yoloMode); err != nil {
			printError(err.Error())
		}
		userPrompt()
	}
}

func runAgentTurn(
	ctx context.Context,
	r *runner.Runner,
	reader *bufio.Reader,
	userID, sessionID string,
	input *genai.Content,
	streamingMode agent.StreamingMode,
	showAgentHeader bool,
	yoloMode bool,
) error {
	for {
		state := newRenderState(streamingMode)
		if showAgentHeader {
			agentHeader()
			showAgentHeader = false
		}

		for event, err := range r.Run(ctx, userID, sessionID, input, agent.RunConfig{
			StreamingMode: streamingMode,
		}) {
			state.handleEvent(event, err)
		}

		fmt.Println()
		if len(state.pendingConfirm) == 0 {
			return nil
		}

		parts, err := collectConfirmations(reader, state.pendingConfirm, yoloMode)
		if err != nil {
			return err
		}
		if parts == nil {
			return nil
		}
		input = &genai.Content{Role: genai.RoleUser, Parts: parts}
	}
}

func collectConfirmations(reader *bufio.Reader, calls []*genai.FunctionCall, yoloMode bool) ([]*genai.Part, error) {
	parts := make([]*genai.Part, 0, len(calls))
	for _, fc := range calls {
		orig, err := toolconfirmation.OriginalCallFrom(fc)
		toolName := fc.Name
		if err == nil && orig != nil {
			toolName = orig.Name
		}

		var approved bool
		if yoloMode {
			approved = true
			fmt.Println(paint(magenta+bold, fmt.Sprintf("⚡ %s: auto-approved (YOLO)", toolName)))
		} else {
			fmt.Print(paint(magenta, fmt.Sprintf("  %s: approve? [y/N]: ", toolName)))
			answer, err := reader.ReadString('\n')
			if err != nil {
				return nil, err
			}
			approved = strings.HasPrefix(strings.ToLower(strings.TrimSpace(answer)), "y")
			if !approved {
				fmt.Println(paint(dim, "  (skipped)"))
			}
		}

		parts = append(parts, &genai.Part{
			FunctionResponse: &genai.FunctionResponse{
				Name: toolconfirmation.FunctionCallName,
				ID:   fc.ID,
				Response: map[string]any{
					"confirmed": approved,
				},
			},
		})
	}
	return parts, nil
}

func resolveStreamingMode() agent.StreamingMode {
	cfg := config.Current()
	if cfg == nil {
		return agent.StreamingModeNone
	}
	switch strings.ToLower(strings.TrimSpace(cfg.StreamingMode)) {
	case "sse", "stream", "streaming":
		return agent.StreamingModeSSE
	case "none", "off", "false":
		return agent.StreamingModeNone
	default:
		return agent.StreamingModeNone
	}
}

func isExit(s string) bool {
	switch strings.ToLower(s) {
	case "exit", "quit", "q":
		return true
	}
	return false
}

package bubbles

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool/toolconfirmation"
	"google.golang.org/genai"
)

type ChatStage int

const (
	StageIdle ChatStage = iota
	StageConnecting
	StageProcessing
	StageReplying
	StageConfirm
)

type ChatMessage struct {
	Role    string
	Content string
}

type eventMsg struct {
	event *session.Event
	err   error
}

type doneMsg struct{}

type tickMsg struct{}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

type Terminal struct{}

func (t *Terminal) Run(ctx context.Context, root agent.Agent) error {
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

	app := &App{
		ctx:       ctx,
		agent:     root,
		runner:    r,
		sessionID: resp.Session.ID(),
		userID:    userID,
	}

	ta := textarea.New()
	ta.Placeholder = "Type a message and press Enter..."
	ta.Focus()
	ta.CharLimit = 4096
	ta.SetWidth(80)
	ta.ShowLineNumbers = false
	ta.SetHeight(1)
	ta.KeyMap.InsertNewline.SetEnabled(false)
	app.textarea = ta

	p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	app.program = p
	_, err = p.Run()
	return err
}

type App struct {
	ctx       context.Context
	agent     agent.Agent
	runner    *runner.Runner
	sessionID string
	userID    string

	messages   []ChatMessage
	viewport   viewport.Model
	textarea   textarea.Model
	sending    bool
	err        error
	width      int
	height     int
	ready      bool
	stage      ChatStage
	spinnerIdx int
	program    *tea.Program

	pendingConfirm  []*genai.FunctionCall
	currentTurnText string
	yoloMode        bool
	tickActive      bool
}

func (a *App) Init() tea.Cmd {
	return textarea.Blink
}

func (a *App) startTickIfNeeded() tea.Cmd {
	if a.tickActive {
		return nil
	}
	a.tickActive = true
	return tickCmd()
}

func (a *App) runAgentTurnCmd(input *genai.Content) tea.Cmd {
	return func() tea.Msg {
		go func() {
			for event, err := range a.runner.Run(a.ctx, a.userID, a.sessionID, input, agent.RunConfig{
				StreamingMode: agent.StreamingModeSSE,
			}) {
				a.program.Send(eventMsg{event: event, err: err})
			}
			a.program.Send(doneMsg{})
		}()
		return nil
	}
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		footerHeight := 4

		if !a.ready {
			a.viewport = viewport.New(msg.Width, msg.Height-footerHeight)
			a.ready = true
		} else {
			a.viewport.Width = msg.Width
			a.viewport.Height = msg.Height - footerHeight
		}

		a.textarea.SetWidth(msg.Width)
		a.viewport.SetContent(a.renderChat())
		a.viewport.GotoBottom()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return a, tea.Quit
		case "ctrl+y":
			a.yoloMode = !a.yoloMode
			statusStr := "OFF"
			if a.yoloMode {
				statusStr = "ON"
			}
			a.messages = append(a.messages, ChatMessage{
				Role:    "tool_result",
				Content: fmt.Sprintf("YOLO Mode toggled: %s (auto-approves tool confirmations)", statusStr),
			})
			a.viewport.SetContent(a.renderChat())
			a.viewport.GotoBottom()
			return a, nil
		case "enter":
			if a.stage == StageConfirm {
				// Enter acts as approval denial (Default: No)
				return a, a.handleConfirmationAction(false)
			}
			if a.sending {
				break
			}
			inputStr := strings.TrimSpace(a.textarea.Value())
			if inputStr == "" {
				break
			}
			a.messages = append(a.messages, ChatMessage{Role: "user", Content: inputStr})

			a.textarea.Reset()
			a.err = nil
			a.sending = true
			a.stage = StageConnecting
			a.spinnerIdx = 0
			a.currentTurnText = ""
			a.viewport.SetContent(a.renderChat())
			a.viewport.GotoBottom()

			input := genai.NewContentFromText(inputStr+"\n", genai.RoleUser)
			cmds = append(cmds, a.runAgentTurnCmd(input), a.startTickIfNeeded())
		}

		// Handle confirmation keys in StageConfirm
		if a.stage == StageConfirm {
			switch strings.ToLower(msg.String()) {
			case "y":
				return a, a.handleConfirmationAction(true)
			case "n":
				return a, a.handleConfirmationAction(false)
			}
		}

	case tickMsg:
		if a.sending {
			a.spinnerIdx++
			a.viewport.SetContent(a.renderChat())
			a.viewport.GotoBottom()
			cmds = append(cmds, tickCmd())
		} else {
			a.tickActive = false
		}

	case eventMsg:
		if msg.err != nil {
			a.messages = append(a.messages, ChatMessage{Role: "error", Content: msg.err.Error()})
			a.viewport.SetContent(a.renderChat())
			a.viewport.GotoBottom()
			return a, nil
		}
		if msg.event == nil || msg.event.LLMResponse.Content == nil {
			return a, nil
		}
		for _, p := range msg.event.LLMResponse.Content.Parts {
			if fc := p.FunctionCall; fc != nil {
				if fc.Name == toolconfirmation.FunctionCallName {
					a.pendingConfirm = append(a.pendingConfirm, fc)
				} else {
					argsStr := formatJSON(fc.Args)
					callDesc := fc.Name
					if argsStr != "" {
						callDesc += fmt.Sprintf(" (%s)", argsStr)
					}
					a.messages = append(a.messages, ChatMessage{Role: "tool_call", Content: callDesc})
				}
			} else if fr := p.FunctionResponse; fr != nil {
				if fr.Name != toolconfirmation.FunctionCallName {
					respStr := formatJSON(fr.Response)
					respDesc := fr.Name
					if respStr != "" {
						respDesc += fmt.Sprintf(" -> %s", respStr)
					}
					a.messages = append(a.messages, ChatMessage{Role: "tool_result", Content: respDesc})
				}
			} else if p.Text != "" {
				a.stage = StageReplying
				var newText string
				if a.currentTurnText != "" && strings.HasPrefix(p.Text, a.currentTurnText) {
					newText = p.Text[len(a.currentTurnText):]
				} else {
					newText = p.Text
				}

				if newText != "" {
					a.currentTurnText += newText
					if len(a.messages) > 0 && a.messages[len(a.messages)-1].Role == "assistant" {
						a.messages[len(a.messages)-1].Content += newText
					} else {
						a.messages = append(a.messages, ChatMessage{Role: "assistant", Content: newText})
					}
				}
			}
		}
		a.viewport.SetContent(a.renderChat())
		a.viewport.GotoBottom()

	case doneMsg:
		a.sending = false
		if len(a.pendingConfirm) > 0 {
			if a.yoloMode {
				return a, a.handleConfirmationAction(true)
			}
			a.stage = StageConfirm
			a.textarea.Blur()
		} else {
			a.stage = StageIdle
			a.textarea.Focus()
		}
		a.viewport.SetContent(a.renderChat())
		a.viewport.GotoBottom()
	}

	var passToViewport bool = true
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyPgUp, tea.KeyPgDown, tea.KeyCtrlU, tea.KeyCtrlD, tea.KeyUp, tea.KeyDown:
			passToViewport = true
		default:
			passToViewport = false
		}
	}

	if passToViewport {
		a.viewport, cmd = a.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	if a.stage != StageConfirm {
		a.textarea, cmd = a.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a *App) handleConfirmationAction(approved bool) tea.Cmd {
	parts := make([]*genai.Part, 0, len(a.pendingConfirm))
	for _, fc := range a.pendingConfirm {
		parts = append(parts, &genai.Part{
			FunctionResponse: &genai.FunctionResponse{
				Name: toolconfirmation.FunctionCallName,
				ID:   fc.ID,
				Response: map[string]any{
					"confirmed": approved,
				},
			},
		})

		statusDesc := "skipped"
		if approved {
			statusDesc = "approved"
		}
		a.messages = append(a.messages, ChatMessage{
			Role:    "tool_result",
			Content: fmt.Sprintf("Tool %s (%s)", fc.Name, statusDesc),
		})
	}

	a.pendingConfirm = nil
	a.stage = StageConnecting
	a.sending = true
	a.currentTurnText = ""
	a.textarea.Focus()
	a.viewport.SetContent(a.renderChat())
	a.viewport.GotoBottom()

	input := &genai.Content{Role: genai.RoleUser, Parts: parts}
	return tea.Batch(a.runAgentTurnCmd(input), a.startTickIfNeeded())
}

func (a *App) View() string {
	if !a.ready {
		return "\n  Initializing TUI..."
	}

	divider := dividerStyle.Render(strings.Repeat("─", a.width))
	statusText := "ADK Bubbles TUI"
	if a.yoloMode {
		statusText = "YOLO MODE ON"
	} else if a.stage == StageConfirm {
		statusText = "CONFIRMATION NEEDED"
	}
	status := statusBarStyle.Render(fmt.Sprintf(" %s ", statusText))
	hintText := " Enter: send │ Ctrl+Y: yolo │ Esc: quit "
	if a.stage == StageConfirm {
		hintText = " y: Approve │ n: Deny │ Ctrl+Y: yolo │ Esc: quit "
	}
	hint := hintStyle.Render(hintText)

	inputTextarea := ""
	if a.stage != StageConfirm {
		inputTextarea = inputPromptStyle.Render("› ") + a.textarea.View()
	} else {
		inputTextarea = inputPromptStyle.Render("› ") + lipgloss.NewStyle().Foreground(lipgloss.Color("#565f89")).Render("Waiting for tool confirmation (y/n)...")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		a.viewport.View(),
		divider,
		inputTextarea,
		"",
		status+"  "+hint,
	)
}

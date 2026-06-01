package bubbles

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (a *App) renderChat() string {
	var b strings.Builder
	wrapW := a.viewport.Width

	for _, msg := range a.messages {
		switch msg.Role {
		case "user":
			b.WriteString(userLabelStyle.Render("You"))
			b.WriteString("\n")
			b.WriteString(wrapText(userMsgStyle.Render(msg.Content), wrapW))
			b.WriteString("\n\n")
		case "assistant":
			b.WriteString(assistantLabelStyle.Render("AI"))
			b.WriteString("\n")
			b.WriteString(wrapText(assistantMsgStyle.Render(msg.Content), wrapW))
			b.WriteString("\n\n")
		case "tool_call":
			b.WriteString(toolLabelStyle.Render("▶ " + msg.Content))
			b.WriteString("\n\n")
		case "tool_result":
			b.WriteString(toolResultMsgStyle.Render("✓ " + msg.Content))
			b.WriteString("\n\n")
		case "error":
			b.WriteString(errMsgStyle.Render("✗ " + msg.Content))
			b.WriteString("\n\n")
		}
	}

	if a.sending {
		spinnerFrames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		frame := spinnerFrames[a.spinnerIdx%len(spinnerFrames)]
		spStr := spinnerStyle.Render(frame) + " "

		var statusStr string
		switch a.stage {
		case StageConnecting:
			statusStr = spStr + connectingStyle.Render("Initializing Turn...")
		case StageProcessing:
			statusStr = spStr + processingStyle.Render("AI is thinking...")
		case StageReplying:
			statusStr = spStr + replyingStyle.Render("AI is replying...")
		}

		b.WriteString(thinkingStyle.Render(statusStr))
		b.WriteString("\n")
	}

	if a.stage == StageConfirm {
		b.WriteString(confirmStyle.Render("? Confirm tool executions above? [y/N]: "))
	}

	return b.String()
}

func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}
	style := lipgloss.NewStyle().Width(width)
	return style.Render(text)
}

func formatJSON(v any) string {
	if v == nil {
		return ""
	}

	m, ok := v.(map[string]any)
	if !ok {
		return fmt.Sprintf("%v", v)
	}

	var parts []string
	for k, v := range m {
		if v == nil || v == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s: %v", k, v))
	}
	return joinStrings(parts, "  |  ")
}

func joinStrings(strs []string, sep string) string {
	res := ""
	for i, s := range strs {
		if i > 0 {
			res += sep
		}
		res += s
	}
	return res
}

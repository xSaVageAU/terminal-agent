package bubbles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	userLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7aa2f7")).
			Bold(true)

	assistantLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9ece6a")).
				Bold(true)

	toolLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e0af68")).
			Bold(true)

	userMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#c0caf5")).
			MarginLeft(2)

	assistantMsgStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#c0caf5")).
				MarginLeft(2)

	toolCallMsgStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#e0af68")).
				Italic(true).
				MarginLeft(2)

	toolResultMsgStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#10b981")).
				Italic(true).
				MarginLeft(2)

	thinkingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565f89")).
			Italic(true).
			MarginLeft(2)

	spinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#bb9af7")).
			Bold(true)

	connectingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff9e64")).
			Italic(true)

	processingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0db9d7")).
			Italic(true)

	replyingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9ece6a")).
			Italic(true)

	confirmStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f7768e")).
			Bold(true).
			MarginLeft(2)

	errMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f7768e")).
			MarginLeft(2)

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3b4261"))

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#1a1b26")).
			Background(lipgloss.Color("#7aa2f7")).
			Padding(0, 1)

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565f89")).
			Italic(true)

	inputPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7aa2f7")).
				Bold(true)
)

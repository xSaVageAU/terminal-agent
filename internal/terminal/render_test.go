package terminal

import (
	"bytes"
	"io"
	"os"
	"testing"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/session"
)

func TestHandleText(t *testing.T) {
	tests := []struct {
		name     string
		chunks   []string
		expected string
	}{
		{
			name:     "Accumulated text",
			chunks:   []string{"I", "I can", "I can perform", "I can perform"},
			expected: "I can perform",
		},
		{
			name:     "Chunk text",
			chunks:   []string{"I", " can", " perform", "I can perform"},
			expected: "I can perform",
		},
		{
			name:     "Chunk text with punctuation",
			chunks:   []string{"Hello", ",", " world", "!", "Hello, world!"},
			expected: "Hello, world!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newRenderState(agent.StreamingModeSSE)

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			for _, chunk := range tt.chunks {
				s.handleText(chunk, &session.Event{})
			}

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			actual := buf.String()

			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

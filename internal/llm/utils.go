package llm

import (
	"google.golang.org/adk/model"
)

// withoutStreaming is a placeholder. 
// If specific ADK streaming utilities are unavailable in the environment,
// it acts as an identity function.
func withoutStreaming(m model.LLM) model.LLM {
	return m
}


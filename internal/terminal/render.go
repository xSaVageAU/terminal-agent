package terminal

import (
	"fmt"
	"strings"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool/toolconfirmation"
	"google.golang.org/genai"
)

type renderState struct {
	streamingMode  agent.StreamingMode
	printedText    string
	seenCalls      map[string]struct{}
	seenResponses  map[string]struct{}
	pendingConfirm []*genai.FunctionCall
}

func newRenderState(mode agent.StreamingMode) *renderState {
	return &renderState{
		streamingMode: mode,
		seenCalls:     make(map[string]struct{}),
		seenResponses: make(map[string]struct{}),
	}
}

func (s *renderState) handleEvent(event *session.Event, err error) {
	if err != nil {
		printError(err.Error())
		return
	}
	if event == nil || event.LLMResponse.Content == nil {
		return
	}

	for _, p := range event.LLMResponse.Content.Parts {
		if fc := p.FunctionCall; fc != nil {
			s.handleFunctionCall(fc)
			continue
		}
		if fr := p.FunctionResponse; fr != nil {
			s.handleFunctionResponse(fr)
			continue
		}
		if p.Text != "" {
			s.handleText(p.Text, event)
		}
	}
}

func (s *renderState) handleFunctionCall(fc *genai.FunctionCall) {
	if fc.Name == toolconfirmation.FunctionCallName {
		if _, ok := s.seenCalls[fc.ID]; ok {
			return
		}
		s.seenCalls[fc.ID] = struct{}{}
		s.pendingConfirm = append(s.pendingConfirm, fc)
		orig, err := toolconfirmation.OriginalCallFrom(fc)
		if err != nil {
			printConfirmPrompt("unknown", formatJSON(fc.Args))
			return
		}
		printConfirmPrompt(orig.Name, formatJSON(orig.Args))
		return
	}

	if fc.ID != "" {
		if _, ok := s.seenCalls[fc.ID]; ok {
			return
		}
		s.seenCalls[fc.ID] = struct{}{}
	} else {
		key := fc.Name + formatJSON(fc.Args)
		if _, ok := s.seenCalls[key]; ok {
			return
		}
		s.seenCalls[key] = struct{}{}
	}
	printToolCall(fc.Name, formatJSON(fc.Args))
}

func (s *renderState) handleFunctionResponse(fr *genai.FunctionResponse) {
	if fr.Name == toolconfirmation.FunctionCallName {
		return
	}
	if fr.ID != "" {
		if _, ok := s.seenResponses[fr.ID]; ok {
			return
		}
		s.seenResponses[fr.ID] = struct{}{}
	}
	printToolResult(fr.Name, formatJSON(fr.Response))
}

func (s *renderState) handleText(text string, event *session.Event) {
	if s.printedText != "" && strings.HasPrefix(text, s.printedText) {
		// If text starts with printedText, we only print the new suffix
		newSuffix := text[len(s.printedText):]
		fmt.Print(newSuffix)
		s.printedText += newSuffix
	} else {
		fmt.Print(text)
		s.printedText += text
	}
}

package files

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xSaVageAU/terminal-agent/tools/workspace"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type SearchTextArgs struct {
	Query      string  `json:"query" jsonschema:"Text to search for (case-insensitive substring match)."`
	Path       *string `json:"path,omitempty" jsonschema:"Directory under workspace to search, or empty for entire workspace. (optional)"`
	MaxMatches *int    `json:"max_matches,omitempty" jsonschema:"Maximum matches to return (optional, default 30, max 100)."`
}

type TextMatch struct {
	File string `json:"file"`
	Line int    `json:"line"`
	Text string `json:"text"`
}

type SearchTextResult struct {
	Query     string      `json:"query"`
	Matches   []TextMatch `json:"matches"`
	Truncated bool        `json:"truncated,omitempty"`
}

const maxSearchFileBytes = 512 * 1024

func searchText(_ tool.Context, args SearchTextArgs) (SearchTextResult, error) {
	pathVal := ""
	if args.Path != nil {
		pathVal = *args.Path
	}
	root, err := workspace.ResolveInWorkspace(pathVal)
	if err != nil {
		return SearchTextResult{}, err
	}
	query := strings.ToLower(strings.TrimSpace(args.Query))
	if query == "" {
		return SearchTextResult{}, fmt.Errorf("query must not be empty")
	}
	maxMatchesVal := 30
	if args.MaxMatches != nil {
		maxMatchesVal = *args.MaxMatches
	}
	limit := workspace.ClampInt(maxMatchesVal, 30, 100)

	var matches []TextMatch
	truncated := false

	err = filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if d.IsDir() {
			if path != root && workspace.SkipDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if len(matches) >= limit {
			truncated = true
			return filepath.SkipAll
		}
		info, err := d.Info()
		if err != nil || info.Size() > maxSearchFileBytes {
			return nil
		}
		if isLikelyBinary(path) {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if strings.Contains(strings.ToLower(line), query) {
				rel, _ := filepath.Rel(workspace.WorkspaceRoot(), path)
				matches = append(matches, TextMatch{
					File: rel,
					Line: lineNum,
					Text: strings.TrimSpace(line),
				})
				if len(matches) >= limit {
					truncated = true
					return filepath.SkipAll
				}
			}
		}
		return nil
	})
	if err != nil && err != filepath.SkipAll {
		return SearchTextResult{}, err
	}

	return SearchTextResult{
		Query:     args.Query,
		Matches:   matches,
		Truncated: truncated,
	}, nil
}

func isLikelyBinary(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".exe", ".dll", ".png", ".jpg", ".jpeg", ".gif", ".zip", ".pdf", ".wasm", ".ico":
		return true
	}
	return false
}

func SearchTextTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "search_text",
		Description: "Searches for a text substring in workspace files (skips large/binary files and common vendor folders).",
	}, searchText)
}


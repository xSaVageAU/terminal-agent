package network

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type FetchURLArgs struct {
	URL string `json:"url" jsonschema:"HTTP or HTTPS URL to fetch."`
}

type FetchURLResult struct {
	URL         string `json:"url"`
	Status      int    `json:"status_code"`
	ContentType string `json:"content_type"`
	Body        string `json:"body"`
	Truncated   bool   `json:"truncated,omitempty"`
}

const maxFetchBytes = 32 * 1024

func fetchURL(_ tool.Context, args FetchURLArgs) (FetchURLResult, error) {
	raw := strings.TrimSpace(args.URL)
	if raw == "" {
		return FetchURLResult{}, fmt.Errorf("url is required")
	}
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		return FetchURLResult{}, fmt.Errorf("only http and https URLs are allowed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		return FetchURLResult{}, err
	}
	req.Header.Set("User-Agent", "adk-terminal-agent/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return FetchURLResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxFetchBytes+1))
	if err != nil {
		return FetchURLResult{}, err
	}
	truncated := len(body) > maxFetchBytes
	if truncated {
		body = body[:maxFetchBytes]
	}

	return FetchURLResult{
		URL:         raw,
		Status:      resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		Body:        string(body),
		Truncated:   truncated,
	}, nil
}

func FetchURLTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "fetch_url",
		Description: "Fetches a public HTTP(S) URL and returns status, content type, and body (truncated to 32KB).",
	}, fetchURL)
}

package terminal

import (
	"fmt"
)

func formatJSON(v any) string {
	if v == nil {
		return ""
	}

	m, ok := v.(map[string]any)
	if !ok {
		return fmt.Sprintf("%v", v)
	}

	// Filter out empty/useless fields
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

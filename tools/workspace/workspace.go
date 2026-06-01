package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xSaVageAU/terminal-agent/internal/config"
)

// WorkspaceRoot returns the directory tools may access.
func WorkspaceRoot() string {
	if cfg := config.Current(); cfg != nil && cfg.WorkspaceRoot != "" {
		return cfg.WorkspaceRoot
	}
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

func ResolveInWorkspace(path string) (string, error) {
	root, err := filepath.Abs(WorkspaceRoot())
	if err != nil {
		return "", err
	}
	target := path
	if target == "" || target == "." {
		target = root
	} else if !filepath.IsAbs(target) {
		target = filepath.Join(root, target)
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(root, abs)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("path %q is outside workspace %q", path, root)
	}
	return abs, nil
}

func SkipDir(name string) bool {
	switch name {
	case ".git", "node_modules", "vendor", "bin", "dist", "__pycache__", ".cursor":
		return true
	}
	return false
}

func ClampInt(v, min, max int) int {
	if v <= 0 {
		return min
	}
	if v > max {
		return max
	}
	return v
}


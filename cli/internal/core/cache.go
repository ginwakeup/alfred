package core

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// CreateCacheDir creates and returns a persistent cache directory for your tool.
// It automatically chooses the correct OS-specific location:
// macOS -> ~/Library/Caches/<toolName>
// Linux -> $XDG_CACHE_HOME/<toolName> or ~/.cache/<toolName>
// Windows -> %LocalAppData%\<toolName>
func CreateCacheDir(toolName string) (string, error) {
	var baseDir string

	switch runtime.GOOS {
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot get home dir: %w", err)
		}
		baseDir = filepath.Join(home, "Library", "Caches")
	case "linux":
		baseDir = os.Getenv("XDG_CACHE_HOME")
		if baseDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("cannot get home dir: %w", err)
			}
			baseDir = filepath.Join(home, ".cache")
		}
	case "windows":
		baseDir = os.Getenv("LocalAppData")
		if baseDir == "" {
			return "", fmt.Errorf("LocalAppData not set on Windows")
		}
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	cacheDir := filepath.Join(baseDir, toolName)

	// Ensure the directory exists
	err := os.MkdirAll(cacheDir, 0o755)
	if err != nil {
		return "", fmt.Errorf("cannot create cache dir %s: %w", cacheDir, err)
	}

	return cacheDir, nil
}

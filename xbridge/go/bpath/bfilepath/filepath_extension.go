package bfilepath

import (
	"path/filepath"
	"runtime"
	"strings"
)

// filepathexpansion_IsInDirectory structurally evaluates if the target path is located
// inside the base directory by calculating their relative relationship.
//
// It normalizes path strings to prevent directory traversal exploits (such as "." or "..")
// and enforces case-insensitive evaluation on Windows and macOS platforms.
//
// It returns false if the relative path requires stepping out of the base directory
// boundaries, or if the relative path computation fails.
func filepath_IsInDirectory(baseDir, targetPath string) bool {
	b := filepath.Clean(baseDir)
	t := filepath.Clean(targetPath)

	// Cross-platform case insensitivity for Windows and macOS
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		b = strings.ToLower(b)
		t = strings.ToLower(t)
	}

	rel, err := filepath.Rel(b, t)
	if err != nil {
		return false
	}

	// If the relative path starts with ".." or is "..", it means targetPath
	// escaped the boundaries of baseDir.
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return false
	}

	return true
}

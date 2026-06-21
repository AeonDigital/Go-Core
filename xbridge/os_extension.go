package xbridge

import (
	"path/filepath"
)

// osexpansion_UserLogDir resolves the platform-specific default path
// for user log data using the configured bridge resources.
//
// On Windows, it returns %LocalAppData%\Logs if %LocalAppData% is set,
// otherwise it falls back to $HOME\AppData\Local\Logs.
// On Darwin, it returns $HOME/Library/Logs.
// On other Unix systems, it returns $XDG_STATE_HOME if non-empty,
// otherwise it falls back to $HOME/.local/state.
//
// If the location cannot be determined (for example, $HOME is not defined),
// then it will return an error.
func os_UserLogDir() (string, error) {
	home, err := OSBridge.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch RuntimeBridge.GOOS() {
	case "windows":
		localAppData := OSBridge.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			return filepath.Join(home, "AppData", "Local", "Logs"), nil
		}
		return filepath.Join(localAppData, "Logs"), nil

	case "darwin":
		// Native macOS path for user-space application logs
		return filepath.Join(home, "Library", "Logs"), nil

	default:
		// XDG_STATE_HOME (~/.local/state) - Recommended for persistent history/logs
		xdgState := OSBridge.Getenv("XDG_STATE_HOME")
		if xdgState == "" {
			return filepath.Join(home, ".local", "state"), nil
		}

		// Final fallback conforming to standard paths
		return filepath.Join(xdgState), nil
	}
}

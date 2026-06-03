package xfs

import (
	"os"
	"path/filepath"
)

//
//
// osUserDataDir

func osUserDataDir(home string, appName string) string {
	switch pkgBridgeXFS.GetRuntimeGOOS() {
	case "windows":
		return osUserDataDirWindows(home, appName)
	case "darwin":
		return osUserDataDirDarwin(home, appName)
	default:
		return osUserDataDirLinux(home, appName)
	}
}

func osUserDataDirWindows(home, appName string) string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return filepath.Join(home, "AppData", "Local", appName, "Data")
	}
	return filepath.Join(localAppData, appName, "Data")
}

func osUserDataDirDarwin(home, appName string) string {
	return filepath.Join(home, "Library", "Application Support", appName)
}

func osUserDataDirLinux(home string, appName string) string {
	if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
		return filepath.Join(xdgData, appName)
	}
	return filepath.Join(home, ".local", "share", appName)
}

//
//
// osUserLogDir

func osUserLogDir(home string, appName string) string {
	switch pkgBridgeXFS.GetRuntimeGOOS() {
	case "windows":
		return osUserLogDirWindows(home, appName)
	case "darwin":
		return osUserLogDirDarwin(home, appName)
	default:
		return osUserLogDirLinux(home, appName)
	}
}

func osUserLogDirWindows(home, appName string) string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return filepath.Join(home, "AppData", "Local", appName, "Log")
	}
	return filepath.Join(localAppData, appName, "Log")
}

func osUserLogDirDarwin(home, appName string) string {
	return filepath.Join(home, "Library", "Logs", appName)
}

func osUserLogDirLinux(home, appName string) string {
	if xdgState := os.Getenv("XDG_STATE_HOME"); xdgState != "" {
		return filepath.Join(xdgState, appName)
	}
	if xdgCache := os.Getenv("XDG_CACHE_HOME"); xdgCache != "" {
		return filepath.Join(xdgCache, appName, "log")
	}
	return filepath.Join(home, ".local", "state", appName)
}

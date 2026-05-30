package xfs

import (
	"os"
	"path/filepath"
	"runtime"
)

var fnMockable_CurrentGOOS = runtime.GOOS

var fnMockable_OsUserDataDirWindows = osUserDataDirWindows
var fnMockable_OsUserDataDirDarwin = osUserDataDirDarwin
var fnMockable_OsUserDataDirLinux = osUserDataDirLinux

var fnMockable_OsUserLogDirWindows = osUserLogDirWindows
var fnMockable_OsUserLogDirDarwin = osUserLogDirDarwin
var fnMockable_OsUserLogDirLinux = osUserLogDirLinux

//
// osUserDataDir

func osUserDataDir(home string, appName string) string {
	switch fnMockable_CurrentGOOS {
	case "windows":
		return fnMockable_OsUserDataDirWindows(home, appName)
	case "darwin":
		return fnMockable_OsUserDataDirDarwin(home, appName)
	default:
		return fnMockable_OsUserDataDirLinux(home, appName)
	}
}

func osUserDataDirLinux(home, appName string) string {
	if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
		return filepath.Join(xdgData, appName)
	}
	return filepath.Join(home, ".local", "share", appName)
}

func osUserDataDirDarwin(home, appName string) string {
	return filepath.Join(home, "Library", "Application Support", appName)
}

func osUserDataDirWindows(home, appName string) string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return filepath.Join(home, "AppData", "Local", appName, "Data")
	}
	return filepath.Join(localAppData, appName, "Data")
}

//
// osUserLogDir

func osUserLogDir(home string, appName string) string {
	switch fnMockable_CurrentGOOS {
	case "windows":
		return fnMockable_OsUserLogDirWindows(home, appName)
	case "darwin":
		return fnMockable_OsUserLogDirDarwin(home, appName)
	default:
		return fnMockable_OsUserLogDirLinux(home, appName)
	}
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

func osUserLogDirDarwin(home, appName string) string {
	return filepath.Join(home, "Library", "Logs", appName)
}

func osUserLogDirWindows(home, appName string) string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return filepath.Join(home, "AppData", "Local", appName, "Log")
	}
	return filepath.Join(localAppData, appName, "Log")
}

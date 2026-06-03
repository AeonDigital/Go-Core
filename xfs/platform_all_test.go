package xfs_test

import (
	"path/filepath"
	"testing"

	"github.com/AeonDigital/Go-Core/internal/pkgxmock"
	"github.com/AeonDigital/Go-Core/xfs"
)

func TestOsUserDataDir_Router(t *testing.T) {
	testFunction := "osUserDataDir (Router)"
	const mockAppName = "TestApp"
	const mockHome = "/user/mockhome"

	type inputArgs struct {
		home    string
		appName string
	}

	tests := []pkgxmock.TestCaseXFS{
		{
			Name: "Router: Should route to Windows when OS is windows",
			Input: inputArgs{
				home:    mockHome,
				appName: mockAppName,
			},
			Want: filepath.Join(mockHome, "AppData", "Local", mockAppName, "Data"),
			Env: map[string]string{
				"XDG_DATA_HOME": "",
				"LOCALAPPDATA":  "",
			},
			WantErr:     false,
			ErrContains: "",
			MockFn: func(m *pkgxmock.MockXFS) {
				m.SetReturn.GetRuntimeGOOS("windows")
				m.OnCall.OSUserDataDir(func(home string, appName string) string {
					return filepath.Join(home, "AppData", "Local", appName, "Data")
				})
			},
		},
		{
			Name: "Router: Should route to Darwin when OS is darwin",
			Input: inputArgs{
				home:    mockHome,
				appName: mockAppName,
			},
			Want: filepath.Join(mockHome, "Library", "Application Support", mockAppName),
			Env: map[string]string{
				"XDG_DATA_HOME": "",
				"LOCALAPPDATA":  "",
			},
			WantErr:     false,
			ErrContains: "",
			MockFn: func(m *pkgxmock.MockXFS) {
				m.SetReturn.GetRuntimeGOOS("darwin")
				m.OnCall.OSUserDataDir(func(home string, appName string) string {
					return filepath.Join(home, "Library", "Application Support", appName)
				})
			},
		},
		{
			Name: "Router: Should route to Linux when OS is linux",
			Input: inputArgs{
				home:    mockHome,
				appName: mockAppName,
			},
			Want: filepath.Join(mockHome, ".local", "share", mockAppName),
			Env: map[string]string{
				"XDG_DATA_HOME": "",
				"LOCALAPPDATA":  "",
			},
			WantErr:     false,
			ErrContains: "",
			MockFn: func(m *pkgxmock.MockXFS) {
				m.SetReturn.GetRuntimeGOOS("linux")
				m.OnCall.OSUserDataDir(func(home string, appName string) string {
					return filepath.Join(home, ".local", "share", appName)
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Configure mocks and environment
			defer SetupTestMock(t, tt)()

			// Unpacks the typed parameters from the generic input object safely
			args := tt.Input.(inputArgs)

			// Trigger the business implementation via our test gateway bridge
			got := xfs.ExportOsUserDataDir(args.home, args.appName)

			// Validation
			AssertResult(t, tt, testFunction, got, nil)
		})
	}
}

func TestOsUserDataDir_Linux(t *testing.T) {
	testFunction := "osUserDataDirLinux"
	const mockAppName = "TestApp"
	const mockHome = "/user/mockhome"

	type inputArgs struct {
		home     string
		appName  string
		envKey   string
		envValue string
	}

	tests := []pkgxmock.TestCaseXFS{
		{
			Name: "Should use default path when XDG_DATA_HOME is empty",
			Input: inputArgs{
				home:     mockHome,
				appName:  mockAppName,
				envKey:   "XDG_DATA_HOME",
				envValue: "",
			},
			Want:        filepath.FromSlash("/user/mockhome/.local/share/TestApp"),
			WantErr:     false,
			ErrContains: "",
			MockFn:      nil, // Direct evaluation; no bridge mocks needed
			Env: map[string]string{
				"XDG_DATA_HOME": "",
			},
		},
		{
			Name: "Should respect XDG_DATA_HOME when provided",
			Input: inputArgs{
				home:     mockHome,
				appName:  mockAppName,
				envKey:   "XDG_DATA_HOME",
				envValue: "/custom/xdg/dir",
			},
			Want:        filepath.FromSlash("/custom/xdg/dir/TestApp"),
			WantErr:     false,
			ErrContains: "",
			MockFn:      nil,
			Env: map[string]string{
				"XDG_DATA_HOME": "/custom/xdg/dir",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Configure mocks and environment
			defer SetupTestMock(t, tt)()

			args := tt.Input.(inputArgs)

			// Trigger the business implementation via the test gateway bridge
			got := xfs.ExportOsUserDataDirLinux(args.home, args.appName)

			// Validation
			AssertResult(t, tt, testFunction, got, nil)
		})
	}
}

func TestOsUserDataDir_Darwin(t *testing.T) {
	testFunction := "osUserDataDirDarwin"
	const mockAppName = "TestApp"
	const mockHome = "/user/mockhome"

	type inputArgs struct {
		home    string
		appName string
	}

	tests := []pkgxmock.TestCaseXFS{
		{
			Name: "Should return Apple Application Support layout",
			Input: inputArgs{
				home:    mockHome,
				appName: mockAppName,
			},
			Want:        filepath.FromSlash("/user/mockhome/Library/Application Support/TestApp"),
			WantErr:     false,
			ErrContains: "",
			MockFn:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			defer SetupTestMock(t, tt)()

			args := tt.Input.(inputArgs)

			// Trigger the business implementation via the test gateway bridge
			got := xfs.ExportOsUserDataDirDarwin(args.home, args.appName)

			// Validation
			AssertResult(t, tt, testFunction, got, nil)
		})
	}
}

func TestOsUserDataDir_Windows(t *testing.T) {
	testFunction := "osUserDataDirWindows"
	const mockAppName = "TestApp"
	const mockHome = "/user/mockhome"

	type inputArgs struct {
		home     string
		appName  string
		envKey   string
		envValue string
	}

	tests := []pkgxmock.TestCaseXFS{
		{
			Name: "Should fallback to Home AppData when LOCALAPPDATA is empty",
			Input: inputArgs{
				home:     mockHome,
				appName:  mockAppName,
				envKey:   "LOCALAPPDATA",
				envValue: "",
			},
			Want:        filepath.Join(mockHome, "AppData", "Local", mockAppName, "Data"),
			WantErr:     false,
			ErrContains: "",
			MockFn:      nil,
			Env: map[string]string{
				"LOCALAPPDATA": "",
			},
		},
		{
			Name: "Should use LOCALAPPDATA path when provided",
			Input: inputArgs{
				home:     mockHome,
				appName:  mockAppName,
				envKey:   "LOCALAPPDATA",
				envValue: "/mock/localappdata",
			},
			Want:        filepath.Join("/mock/localappdata", mockAppName, "Data"),
			WantErr:     false,
			ErrContains: "",
			MockFn:      nil,
			Env: map[string]string{
				"LOCALAPPDATA": "/mock/localappdata",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Configure mocks and environment
			defer SetupTestMock(t, tt)()

			args := tt.Input.(inputArgs)

			// Trigger the business implementation via the test gateway bridge
			got := xfs.ExportOsUserDataDirWindows(args.home, args.appName)

			// Validation
			AssertResult(t, tt, testFunction, got, nil)
		})
	}
}

func TestOsUserLogDir_Router(t *testing.T) {
	testFunction := "osUserLogDir (Router)"
	const mockAppName = "LogApp"
	const mockHome = "/user/mockhome"

	type inputArgs struct {
		home    string
		appName string
	}

	tests := []pkgxmock.TestCaseXFS{
		{
			Name: "Router: Should route to Windows when OS is windows",
			Input: inputArgs{
				home:    mockHome,
				appName: mockAppName,
			},
			Want:        filepath.Join(mockHome, "AppData", "Local", mockAppName, "Log"),
			WantErr:     false,
			ErrContains: "",
			MockFn: func(m *pkgxmock.MockXFS) {
				m.SetReturn.GetRuntimeGOOS("windows")
				m.OnCall.OSUserLogDir(func(home string, appName string) string {
					return filepath.Join(home, "AppData", "Local", appName, "Log")
				})
			},
			Env: map[string]string{
				"XDG_STATE_HOME": "",
				"XDG_CACHE_HOME": "",
				"LOCALAPPDATA":   "",
			},
		},
		{
			Name: "Router: Should route to Darwin when OS is darwin",
			Input: inputArgs{
				home:    mockHome,
				appName: mockAppName,
			},
			Want:        filepath.Join(mockHome, "Library", "Logs", mockAppName),
			WantErr:     false,
			ErrContains: "",
			MockFn: func(m *pkgxmock.MockXFS) {
				m.SetReturn.GetRuntimeGOOS("darwin")
				m.OnCall.OSUserLogDir(func(home string, appName string) string {
					return filepath.Join(home, "Library", "Logs", appName)
				})
			},
			Env: map[string]string{
				"XDG_STATE_HOME": "",
				"XDG_CACHE_HOME": "",
				"LOCALAPPDATA":   "",
			},
		},
		{
			Name: "Router: Should route to Linux when OS is linux",
			Input: inputArgs{
				home:    mockHome,
				appName: mockAppName,
			},
			Want:        filepath.Join(mockHome, ".local", "state", mockAppName),
			WantErr:     false,
			ErrContains: "",
			MockFn: func(m *pkgxmock.MockXFS) {
				m.SetReturn.GetRuntimeGOOS("linux")
				m.OnCall.OSUserLogDir(func(home string, appName string) string {
					return filepath.Join(home, ".local", "state", appName)
				})
			},
			Env: map[string]string{
				"XDG_STATE_HOME": "",
				"XDG_CACHE_HOME": "",
				"LOCALAPPDATA":   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Configure mocks and environment
			defer SetupTestMock(t, tt)()

			// Unpacks the typed parameters from the generic input object safely
			args := tt.Input.(inputArgs)

			// Trigger the business implementation via our test gateway bridge
			got := xfs.ExportOsUserLogDir(args.home, args.appName)

			// Validation
			AssertResult(t, tt, testFunction, got, nil)
		})
	}
}

func TestOsUserLogDir_Linux(t *testing.T) {
	testFunction := "osUserLogDirLinux"
	const mockAppName = "LogApp"
	const mockHome = "/user/mockhome"

	type inputArgs struct {
		home    string
		appName string
		env     map[string]string
	}

	tests := []pkgxmock.TestCaseXFS{
		{
			Name: "Linux: Should use default path when environment variables are empty",
			Input: inputArgs{
				home:    mockHome,
				appName: mockAppName,
				env: map[string]string{
					"XDG_STATE_HOME": "",
					"XDG_CACHE_HOME": "",
				},
			},
			Want:        filepath.FromSlash("/user/mockhome/.local/state/LogApp"),
			WantErr:     false,
			ErrContains: "",
			MockFn:      nil, // Direct platform evaluation; no bridge mocks needed
			Env: map[string]string{
				"XDG_STATE_HOME": "",
				"XDG_CACHE_HOME": "",
			},
		},
		{
			Name: "Linux: Should respect XDG_STATE_HOME with highest priority",
			Input: inputArgs{
				home:    mockHome,
				appName: mockAppName,
				env: map[string]string{
					"XDG_STATE_HOME": "/custom/state",
					"XDG_CACHE_HOME": "/custom/cache",
				},
			},
			Want:        filepath.FromSlash("/custom/state/LogApp"),
			WantErr:     false,
			ErrContains: "",
			MockFn:      nil,
			Env: map[string]string{
				"XDG_STATE_HOME": "/custom/state",
				"XDG_CACHE_HOME": "/custom/cache",
			},
		},
		{
			Name: "Linux: Should fallback to XDG_CACHE_HOME if STATE is empty",
			Input: inputArgs{
				home:    mockHome,
				appName: mockAppName,
				env: map[string]string{
					"XDG_STATE_HOME": "",
					"XDG_CACHE_HOME": "/custom/cache",
				},
			},
			Want:        filepath.FromSlash("/custom/cache/LogApp/log"),
			WantErr:     false,
			ErrContains: "",
			MockFn:      nil,
			Env: map[string]string{
				"XDG_STATE_HOME": "",
				"XDG_CACHE_HOME": "/custom/cache",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Configure mocks and environment
			defer SetupTestMock(t, tt)()

			args := tt.Input.(inputArgs)

			// Trigger the business implementation via the test gateway bridge
			got := xfs.ExportOsUserLogDirLinux(args.home, args.appName)

			// Validation
			AssertResult(t, tt, testFunction, got, nil)
		})
	}
}

func TestOsUserLogDir_Darwin(t *testing.T) {
	testFunction := "osUserLogDirDarwin"
	const mockAppName = "LogApp"
	const mockHome = "/user/mockhome"

	type inputArgs struct {
		home    string
		appName string
	}

	tests := []pkgxmock.TestCaseXFS{
		{
			Name: "Darwin: Should return Apple Logs layout",
			Input: inputArgs{
				home:    mockHome,
				appName: mockAppName,
			},
			Want:        filepath.FromSlash("/user/mockhome/Library/Logs/LogApp"),
			WantErr:     false,
			ErrContains: "",
			MockFn:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			defer SetupTestMock(t, tt)()

			args := tt.Input.(inputArgs)

			// Trigger the business implementation via the test gateway bridge
			got := xfs.ExportOsUserLogDirDarwin(args.home, args.appName)

			// Validation
			AssertResult(t, tt, testFunction, got, nil)
		})
	}
}

func TestOsUserLogDir_Windows(t *testing.T) {
	testFunction := "osUserLogDirWindows"
	const mockAppName = "LogApp"
	const mockHome = "/user/mockhome"

	type inputArgs struct {
		home    string
		appName string
		env     map[string]string
	}

	tests := []pkgxmock.TestCaseXFS{
		{
			Name: "Windows: Should fallback to Home AppData when LOCALAPPDATA is empty",
			Input: inputArgs{
				home:    mockHome,
				appName: mockAppName,
				env: map[string]string{
					"LOCALAPPDATA": "",
				},
			},
			Want:        filepath.Join(mockHome, "AppData", "Local", mockAppName, "Log"),
			WantErr:     false,
			ErrContains: "",
			MockFn:      nil,
			Env: map[string]string{
				"LOCALAPPDATA": "",
			},
		},
		{
			Name: "Windows: Should use LOCALAPPDATA path when provided",
			Input: inputArgs{
				home:    mockHome,
				appName: mockAppName,
				env: map[string]string{
					"LOCALAPPDATA": "/mock/localappdata",
				},
			},
			Want:        filepath.Join("/mock/localappdata", mockAppName, "Log"),
			WantErr:     false,
			ErrContains: "",
			MockFn:      nil,
			Env: map[string]string{
				"LOCALAPPDATA": "/mock/localappdata",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Configure mocks and environment
			defer SetupTestMock(t, tt)()

			args := tt.Input.(inputArgs)

			// Trigger the business implementation via the test gateway bridge
			got := xfs.ExportOsUserLogDirWindows(args.home, args.appName)

			// Validation
			AssertResult(t, tt, testFunction, got, nil)
		})
	}
}

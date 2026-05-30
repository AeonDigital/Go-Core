package xfs

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AeonDigital/Go-Core/xerrors"
)

// mockFileInfo implementa estritamente a interface os.FileInfo necessária para o teste
type mockFileInfo struct {
	os.FileInfo      // Incorpora a interface para herdar métodos que não usaremos
	isRegular   bool // Controla o retorno do método Mode().IsRegular()
	isDir       bool
	mode        os.FileMode
	size        int64
}

func (m mockFileInfo) Mode() os.FileMode {
	// Technical Note: Preserves permission bits if explicitly assigned in m.mode
	if m.mode != 0 {
		return m.mode
	}
	if m.isRegular {
		return 0 // A FileMode 0 represents a regular file in Go
	}
	return os.ModeDir // Returns a directory mode if not a regular file
}

func (m mockFileInfo) IsDir() bool { return m.isDir }

func (m mockFileInfo) Size() int64 { return m.size }

// mockFunction substitui temporariamente uma função mockável por uma versão de teste
// e registra automaticamente o seu reset original no ciclo de vida do teste.
func mockFunction[T any](t *testing.T, original *T, mock T) {
	t.Helper()

	// Salva o estado original da variável
	oldValue := *original

	// Injeta o mock
	*original = mock

	// Registra o reset automático para quando o teste atual terminar
	t.Cleanup(func() {
		*original = oldValue
	})
}

// Função auxiliar para resolver caminhos relativos ao PWD do teste com segurança
func mustAbs(t *testing.T, path string) string {
	t.Helper()
	abs, err := filepath.Abs(path)
	if err != nil {
		t.Fatalf("failed to setup absolute path for test: %v", err)
	}
	return abs
}

func TestWrappersCoverage(t *testing.T) {
	// Apenas chama para passar pelas linhas da cobertura
	_, _ = GetUserHomeDir()
	_, _ = GetUserConfigDir()
	_, _ = fnMockable_ReadDir(nil, 0)
}

func TestRetrieveFullPath(t *testing.T) {
	mockHome := filepath.FromSlash("/user/mockhome")

	tests := []struct {
		name         string
		input        string
		want         string
		wantErr      bool
		errContains  string
		customHomeFn func() (string, error)
		customAbsFn  func(string) (string, error)
	}{
		{
			name:        "Empty path string should fail immediately",
			input:       "",
			wantErr:     true,
			errContains: "[MSG: empty string value not allowed][FIELD: path]",
		},
		{
			name:    "Standard path without tilde should turn into absolute",
			input:   "documents/project",
			want:    mustAbs(t, "documents/project"),
			wantErr: false,
		},
		{
			name:  "Tilde alone should resolve exactly to user home directory",
			input: "~",
			// Technical Note: Ensures forward slash uniformity to match production cross-platform sanitization
			want:    filepath.ToSlash(mockHome),
			wantErr: false,
			// Technical Note: Forces the absolute path hook to preserve forward slashes on Windows simulation
			customAbsFn: func(p string) (string, error) { return filepath.ToSlash(p), nil },
		},
		{
			name:        "Tilde with forward slash should append structure to home directory",
			input:       "~/downloads/file.txt",
			want:        filepath.ToSlash(filepath.Join(mockHome, "downloads", "file.txt")),
			wantErr:     false,
			customAbsFn: func(p string) (string, error) { return filepath.ToSlash(p), nil },
		},
		{
			name:        "Tilde with backward slash should append structure to home directory",
			input:       "~\\documents\\report.pdf",
			want:        filepath.ToSlash(filepath.Join(mockHome, "documents", "report.pdf")),
			wantErr:     false,
			customAbsFn: func(p string) (string, error) { return filepath.ToSlash(p), nil },
		},
		{
			name:    "Tilde stuck to letters without slashes should be treated as relative name",
			input:   "~aeon-folder",
			want:    mustAbs(t, "~aeon-folder"),
			wantErr: false,
		},
		{
			name:        "Should fail and wrap error when user home resolution fails",
			input:       "~/downloads",
			wantErr:     true,
			errContains: "[MSG: failed to resolve user home directory][FIELD: home]",
			customHomeFn: func() (string, error) {
				return "", errors.New("permission denied")
			},
		},
		{
			name:        "Should fail when absolute path resolution returns an error",
			input:       "documents/project",
			wantErr:     true,
			errContains: "[MSG: failed to resolve absolute path][FIELD: path]",
			customAbsFn: func(path string) (string, error) {
				return "", errors.New("mocked filesystem breakdown error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Technical Note: Safely captures and restores original mock references after each run
			oldHomeFn := fnMockable_UserHomeDir
			oldAbsFn := fnMockable_FilepathAbs
			defer func() {
				fnMockable_UserHomeDir = oldHomeFn
				fnMockable_FilepathAbs = oldAbsFn
			}()

			if tt.customHomeFn != nil {
				fnMockable_UserHomeDir = tt.customHomeFn
			} else {
				fnMockable_UserHomeDir = func() (string, error) { return mockHome, nil }
			}

			if tt.customAbsFn != nil {
				fnMockable_FilepathAbs = tt.customAbsFn
			} else {
				// Resets to native absolute path behavior for sibling test pipeline executions
				fnMockable_FilepathAbs = filepath.Abs
			}

			// 2. Executa a função
			got, err := RetrieveFullPath(tt.input)

			// 3. Validações
			if (err != nil) != tt.wantErr {
				t.Fatalf("RetrieveFullPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}
			if got != tt.want {
				t.Errorf("RetrieveFullPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetUserDataDir(t *testing.T) {
	const mockAppName = "MyApp"
	mockHome := "/user/mockhome"

	tests := []struct {
		name         string
		customHomeFn func() (string, error)
		want         string
	}{
		{
			name:         "Should use user home directory when available",
			customHomeFn: func() (string, error) { return mockHome, nil },
			want:         osUserDataDir(mockHome, mockAppName),
		},
		{
			name:         "Should fallback to temp directory when home directory fails",
			customHomeFn: func() (string, error) { return "", xerrors.NewErr("failed to retrieve home") },
			want:         osUserDataDir(os.TempDir(), mockAppName),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFunction(t, &fnMockable_GetUserHome, tt.customHomeFn)

			got := GetUserDataDir(mockAppName)

			if got != tt.want {
				t.Errorf("GetUserDataDir() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestOsUserDataDir_Router_AllCases(t *testing.T) {
	tests := []struct {
		name          string
		targetOS      string
		setupMocks    func(t *testing.T, called *bool)
		expectedRoute string
	}{
		{
			name:     "Router: Should dispatch to Windows implementation when OS is windows",
			targetOS: "windows",
			setupMocks: func(t *testing.T, called *bool) {
				mockFunction(t, &fnMockable_OsUserDataDirWindows, func(home, appName string) string {
					*called = true
					return "/mock/windows"
				})
			},
			expectedRoute: "/mock/windows",
		},
		{
			name:     "Router: Should dispatch to Darwin implementation when OS is darwin",
			targetOS: "darwin",
			setupMocks: func(t *testing.T, called *bool) {
				mockFunction(t, &fnMockable_OsUserDataDirDarwin, func(home, appName string) string {
					*called = true
					return "/mock/darwin"
				})
			},
			expectedRoute: "/mock/darwin",
		},
		{
			name:     "Router: Should dispatch to Linux implementation when OS is linux (default)",
			targetOS: "linux",
			setupMocks: func(t *testing.T, called *bool) {
				mockFunction(t, &fnMockable_OsUserDataDirLinux, func(home, appName string) string {
					*called = true
					return "/mock/linux"
				})
			},
			expectedRoute: "/mock/linux",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wasSpecialistCalled := false

			// 1. Altera a variável do SO para enganar o switch do roteador
			mockFunction(t, &fnMockable_CurrentGOOS, tt.targetOS)

			// 2. Configura o espião (spy) para verificar se a rota correta foi acionada
			tt.setupMocks(t, &wasSpecialistCalled)

			// 3. Executa a função roteadora principal
			result := osUserDataDir("/home", "MyApp")

			// 4. Validações estritas
			if !wasSpecialistCalled {
				t.Errorf("Expected router to dispatch call to the %s specialist, but it did not", tt.targetOS)
			}
			if result != tt.expectedRoute {
				t.Errorf("Expected router to return value from dispatched function, got %q, want %q", result, tt.expectedRoute)
			}
		})
	}
}

func TestOsUserDataDir_AllPlatformsSpecialist(t *testing.T) {
	const mockAppName = "TestApp"
	const mockHome = "/user/mockhome"

	// 1. TESTES DA LÓGICA DO LINUX
	t.Run("Linux: Should use default path when XDG_DATA_HOME is empty", func(t *testing.T) {
		t.Setenv("XDG_DATA_HOME", "") // Garante que a variável está vazia

		got := osUserDataDirLinux(mockHome, mockAppName)
		want := filepath.FromSlash("/user/mockhome/.local/share/TestApp")

		if got != want {
			t.Errorf("Linux default layout failed. Got %q, want %q", got, want)
		}
	})

	t.Run("Linux: Should respect XDG_DATA_HOME when provided", func(t *testing.T) {
		t.Setenv("XDG_DATA_HOME", "/custom/xdg/dir")

		got := osUserDataDirLinux(mockHome, mockAppName)
		want := filepath.FromSlash("/custom/xdg/dir/TestApp")

		if got != want {
			t.Errorf("Linux XDG layout failed. Got %q, want %q", got, want)
		}
	})

	// 2. TESTES DA LÓGICA DO MAC (DARWIN)
	t.Run("Darwin: Should return Apple Application Support layout", func(t *testing.T) {
		got := osUserDataDirDarwin(mockHome, mockAppName)
		want := filepath.FromSlash("/user/mockhome/Library/Application Support/TestApp")

		if got != want {
			t.Errorf("Darwin layout failed. Got %q, want %q", got, want)
		}
	})

	// 3. TESTES DA LÓGICA DO WINDOWS
	t.Run("Windows: Should fallback to Home AppData when LOCALAPPDATA is empty", func(t *testing.T) {
		t.Setenv("LOCALAPPDATA", "")

		got := osUserDataDirWindows(mockHome, mockAppName)
		// Como estamos testando no Linux, usamos filepath.Join normal que usará barras normais do Unix,
		// mas validando a semântica da árvore do Windows
		want := filepath.Join(mockHome, "AppData", "Local", mockAppName, "Data")

		if got != want {
			t.Errorf("Windows fallback layout failed. Got %q, want %q", got, want)
		}
	})

	t.Run("Windows: Should use LOCALAPPDATA path when provided", func(t *testing.T) {
		t.Setenv("LOCALAPPDATA", "/mock/localappdata")

		got := osUserDataDirWindows(mockHome, mockAppName)
		want := filepath.Join("/mock/localappdata", mockAppName, "Data")

		if got != want {
			t.Errorf("Windows LOCALAPPDATA layout failed. Got %q, want %q", got, want)
		}
	})
}

func TestGetUserLogDir(t *testing.T) {
	const mockAppName = "LogApp"
	mockHome := "/user/mockhome"

	tests := []struct {
		name         string
		customHomeFn func() (string, error)
		want         string
	}{
		{
			name: "Should use user home directory for logs when available",
			customHomeFn: func() (string, error) {
				return mockHome, nil
			},
			want: osUserLogDir(mockHome, mockAppName),
		},
		{
			name: "Should fallback to temp directory for logs when home directory fails",
			customHomeFn: func() (string, error) {
				return "", xerrors.NewErr("failed to retrieve home")
			},
			want: osUserLogDir(os.TempDir(), mockAppName),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Substitui e agenda o reset automático usando o novo padrão explícito
			mockFunction(t, &fnMockable_GetUserHome, tt.customHomeFn)

			got := GetUserLogDir(mockAppName)

			if got != tt.want {
				t.Errorf("GetUserLogDir() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestOsUserLogDir_Router_AllCases(t *testing.T) {
	const mockAppName = "LogApp"
	const mockHome = "/user/mockhome"

	tests := []struct {
		name          string
		targetOS      string
		setupMocks    func(t *testing.T, called *bool)
		expectedRoute string
	}{
		{
			name:     "Router: Should route to Windows when OS is windows",
			targetOS: "windows",
			setupMocks: func(t *testing.T, called *bool) {
				mockFunction(t, &fnMockable_OsUserLogDirWindows, func(h, a string) string {
					*called = true
					return "/route/windows"
				})
			},
			expectedRoute: "/route/windows",
		},
		{
			name:     "Router: Should route to Darwin when OS is darwin",
			targetOS: "darwin",
			setupMocks: func(t *testing.T, called *bool) {
				mockFunction(t, &fnMockable_OsUserLogDirDarwin, func(h, a string) string {
					*called = true
					return "/route/darwin"
				})
			},
			expectedRoute: "/route/darwin",
		},
		{
			name:     "Router: Should route to Linux when OS is linux",
			targetOS: "linux",
			setupMocks: func(t *testing.T, called *bool) {
				mockFunction(t, &fnMockable_OsUserLogDirLinux, func(h, a string) string {
					*called = true
					return "/route/linux"
				})
			},
			expectedRoute: "/route/linux",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wasSpecialistCalled := false

			// Engana a variável controladora de SO para o switch
			mockFunction(t, &fnMockable_CurrentGOOS, tt.targetOS)

			// Configura o espião correspondente ao cenário
			tt.setupMocks(t, &wasSpecialistCalled)

			result := osUserLogDir(mockHome, mockAppName)

			if !wasSpecialistCalled {
				t.Errorf("Expected router to call %s specialist, but it did not", tt.targetOS)
			}
			if result != tt.expectedRoute {
				t.Errorf("Router delivered wrong path. Got %q, want %q", result, tt.expectedRoute)
			}
		})
	}
}

func TestOsUserLogDir_AllPlatformsSpecialist(t *testing.T) {
	const mockAppName = "LogApp"
	const mockHome = "/user/mockhome"

	// ==========================================
	// 1. LINUX
	// ==========================================
	t.Run("Linux: Should use default path when environment variables are empty", func(t *testing.T) {
		t.Setenv("XDG_STATE_HOME", "")
		t.Setenv("XDG_CACHE_HOME", "")

		got := osUserLogDirLinux(mockHome, mockAppName)
		want := filepath.FromSlash("/user/mockhome/.local/state/LogApp")
		if got != want {
			t.Errorf("Linux default failed. Got %q, want %q", got, want)
		}
	})

	t.Run("Linux: Should respect XDG_STATE_HOME with highest priority", func(t *testing.T) {
		t.Setenv("XDG_STATE_HOME", "/custom/state")
		t.Setenv("XDG_CACHE_HOME", "/custom/cache") // Deve ser ignorado pelo short-circuit do if

		got := osUserLogDirLinux(mockHome, mockAppName)
		want := filepath.FromSlash("/custom/state/LogApp")
		if got != want {
			t.Errorf("Linux XDG_STATE failed. Got %q, want %q", got, want)
		}
	})

	t.Run("Linux: Should fallback to XDG_CACHE_HOME if STATE is empty", func(t *testing.T) {
		t.Setenv("XDG_STATE_HOME", "")
		t.Setenv("XDG_CACHE_HOME", "/custom/cache")

		got := osUserLogDirLinux(mockHome, mockAppName)
		want := filepath.FromSlash("/custom/cache/LogApp/log")
		if got != want {
			t.Errorf("Linux XDG_CACHE failed. Got %q, want %q", got, want)
		}
	})

	// ==========================================
	// 2. MAC (DARWIN)
	// ==========================================
	t.Run("Darwin: Should return Apple Logs layout", func(t *testing.T) {
		got := osUserLogDirDarwin(mockHome, mockAppName)
		want := filepath.FromSlash("/user/mockhome/Library/Logs/LogApp")
		if got != want {
			t.Errorf("Darwin layout failed. Got %q, want %q", got, want)
		}
	})

	// ==========================================
	// 3. WINDOWS
	// ==========================================
	t.Run("Windows: Should fallback to Home AppData when LOCALAPPDATA is empty", func(t *testing.T) {
		t.Setenv("LOCALAPPDATA", "")

		got := osUserLogDirWindows(mockHome, mockAppName)
		want := filepath.Join(mockHome, "AppData", "Local", mockAppName, "Log")
		if got != want {
			t.Errorf("Windows fallback failed. Got %q, want %q", got, want)
		}
	})

	t.Run("Windows: Should use LOCALAPPDATA path when provided", func(t *testing.T) {
		t.Setenv("LOCALAPPDATA", "/mock/localappdata")

		got := osUserLogDirWindows(mockHome, mockAppName)
		want := filepath.Join("/mock/localappdata", mockAppName, "Log")
		if got != want {
			t.Errorf("Windows LOCALAPPDATA failed. Got %q, want %q", got, want)
		}
	})
}

func TestGetParentPath(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		want                 string
		customRetrievePathFn func(string) (string, error)
	}{
		{
			name:  "Empty file path should return empty string immediately",
			input: "",
			want:  "",
			// Nem precisa de mock pois a função tem um guard clause para string vazia
		},
		{
			name:  "Should return empty string if RetrieveFullPath returns an error",
			input: "invalid/path",
			want:  "",
			customRetrievePathFn: func(path string) (string, error) {
				return "", xerrors.NewErr("mocked retrieval failure")
			},
		},
		{
			name:  "Should return the correct parent directory on success",
			input: "~/documents/project/file.txt",
			// filepath.Dir de "/user/mockhome/documents/project/file.txt"
			want: filepath.FromSlash("/user/mockhome/documents/project"),
			customRetrievePathFn: func(path string) (string, error) {
				return filepath.FromSlash("/user/mockhome/documents/project/file.txt"), nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Injeta o mock e agenda o reset automático se o cenário exigir
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}

			got := GetParentPath(tt.input)

			if got != tt.want {
				t.Errorf("GetParentPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExists(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		want                 bool
		customRetrievePathFn func(string) (string, error)
		customOsStatFn       func(string) (os.FileInfo, error)
	}{
		{
			name:  "Should return false if RetrieveFullPath returns an error",
			input: "invalid/path",
			want:  false,
			customRetrievePathFn: func(path string) (string, error) {
				return "", xerrors.NewErr("mocked retrieval failure")
			},
			// os.Stat nem chega a ser chamado aqui
		},
		{
			name:  "Should return false if file or directory does not exist",
			input: "~/missing-file.txt",
			want:  false,
			customRetrievePathFn: func(path string) (string, error) {
				return "/user/mockhome/missing-file.txt", nil
			},
			customOsStatFn: func(name string) (os.FileInfo, error) {
				return nil, os.ErrNotExist // Força o erro de não existente
			},
		},
		{
			name:  "Should return true if file or directory exists successfully",
			input: "~/existing-file.txt",
			want:  true,
			customRetrievePathFn: func(path string) (string, error) {
				return "/user/mockhome/existing-file.txt", nil
			},
			customOsStatFn: func(name string) (os.FileInfo, error) {
				return nil, nil // Nil significa que não houve erro, logo o arquivo existe
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Injeta o mock do RetrieveFullPath se configurado
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}

			// Injeta o mock do os.Stat se configurado
			if tt.customOsStatFn != nil {
				mockFunction(t, &fnMockable_OsStat, tt.customOsStatFn)
			}

			got := Exists(tt.input)

			if got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFile(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		want                 bool
		customRetrievePathFn func(string) (string, error)
		customOsStatFn       func(string) (os.FileInfo, error)
	}{
		{
			name:  "Should return false if RetrieveFullPath returns an error",
			input: "invalid/path",
			want:  false,
			customRetrievePathFn: func(path string) (string, error) {
				return "", xerrors.NewErr("mocked retrieval failure")
			},
		},
		{
			name:  "Should return false if os.Stat returns an error",
			input: "~/missing-file.txt",
			want:  false,
			customRetrievePathFn: func(path string) (string, error) {
				return "/user/mockhome/missing-file.txt", nil
			},
			customOsStatFn: func(name string) (os.FileInfo, error) {
				return nil, os.ErrNotExist
			},
		},
		{
			name:  "Should return false if path exists but it is a directory",
			input: "~/documents",
			want:  false,
			customRetrievePathFn: func(path string) (string, error) {
				return "/user/mockhome/documents", nil
			},
			customOsStatFn: func(name string) (os.FileInfo, error) {
				return mockFileInfo{isRegular: false}, nil
			},
		},
		{
			name:  "Should return true if path exists and it is a regular file",
			input: "~/documents/notes.txt",
			want:  true,
			customRetrievePathFn: func(path string) (string, error) {
				return "/user/mockhome/documents/notes.txt", nil
			},
			customOsStatFn: func(name string) (os.FileInfo, error) {
				return mockFileInfo{isRegular: true}, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Injeta o mock do RetrieveFullPath se configurado
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}

			// Injeta o mock do os.Stat se configurado
			if tt.customOsStatFn != nil {
				mockFunction(t, &fnMockable_OsStat, tt.customOsStatFn)
			}

			got := IsFile(tt.input)

			if got != tt.want {
				t.Errorf("IsFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDir(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		want                 bool
		customRetrievePathFn func(string) (string, error)
		customOsStatFn       func(string) (os.FileInfo, error)
	}{
		{
			name:                 "Should return false if RetrieveFullPath returns an error",
			input:                "invalid/path",
			want:                 false,
			customRetrievePathFn: func(p string) (string, error) { return "", xerrors.NewErr("err") },
		},
		{
			name:                 "Should return false if os.Stat returns an error",
			input:                "~/missing-dir",
			want:                 false,
			customRetrievePathFn: func(p string) (string, error) { return "/mock/dir", nil },
			customOsStatFn:       func(n string) (os.FileInfo, error) { return nil, os.ErrNotExist },
		},
		{
			name:                 "Should return false if path exists but it is a regular file",
			input:                "~/notes.txt",
			want:                 false,
			customRetrievePathFn: func(p string) (string, error) { return "/mock/file", nil },
			// Truque nativo: lê este próprio arquivo de teste para obter um FileInfo de arquivo válido
			customOsStatFn: func(n string) (os.FileInfo, error) { return os.Stat("fs_test.go") },
		},
		{
			name:                 "Should return true if path exists and it is a directory",
			input:                "~/documents",
			want:                 true,
			customRetrievePathFn: func(p string) (string, error) { return "/mock/dir", nil },
			// Truque nativo: lê o diretório atual de trabalho (".") para obter um FileInfo de diretório real
			customOsStatFn: func(n string) (os.FileInfo, error) { return os.Stat(".") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Injeta os mocks utilizando a nossa função genérica e agenda o reset automático
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}
			if tt.customOsStatFn != nil {
				mockFunction(t, &fnMockable_OsStat, tt.customOsStatFn)
			}

			if got := IsDir(tt.input); got != tt.want {
				t.Errorf("IsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEmptyDir(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		want                 bool
		wantErr              bool
		errContains          string
		customRetrievePathFn func(string) (string, error)
		customOsOpenFn       func(string) (*os.File, error)
		customReadDirFn      func(*os.File, int) ([]os.DirEntry, error)
	}{
		{
			name:                 "Should return error if RetrieveFullPath fails",
			input:                "invalid/path",
			want:                 false,
			wantErr:              true,
			errContains:          "retrieval failure",
			customRetrievePathFn: func(p string) (string, error) { return "", xerrors.NewErr("retrieval failure") },
		},
		{
			name:                 "Should return error if os.Open fails",
			input:                "~/missing-dir",
			want:                 false,
			wantErr:              true,
			errContains:          "permission denied",
			customRetrievePathFn: func(p string) (string, error) { return "/mock/dir", nil },
			customOsOpenFn:       func(n string) (*os.File, error) { return nil, os.ErrPermission },
		},
		{
			name:                 "Should return false if directory is NOT empty (err is nil)",
			input:                "~/full-dir",
			want:                 false,
			wantErr:              false,
			customRetrievePathFn: func(p string) (string, error) { return ".", nil },
			customOsOpenFn:       func(n string) (*os.File, error) { return os.Open(".") }, // Abre pasta real apenas para ter um *os.File válido
			customReadDirFn:      func(f *os.File, n int) ([]os.DirEntry, error) { return []os.DirEntry{nil}, nil },
		},
		{
			name:                 "Should return true if directory is empty (err is EOF)",
			input:                "~/empty-dir",
			want:                 true,
			wantErr:              false,
			customRetrievePathFn: func(p string) (string, error) { return ".", nil },
			customOsOpenFn:       func(n string) (*os.File, error) { return os.Open(".") },
			customReadDirFn:      func(f *os.File, n int) ([]os.DirEntry, error) { return nil, io.EOF },
		},
		{
			name:                 "Should return generic error if ReadDir fails with unexpected error",
			input:                "~/broken-dir",
			want:                 false,
			wantErr:              true,
			errContains:          "hardware disk failure",
			customRetrievePathFn: func(p string) (string, error) { return ".", nil },
			customOsOpenFn:       func(n string) (*os.File, error) { return os.Open(".") },
			customReadDirFn:      func(f *os.File, n int) ([]os.DirEntry, error) { return nil, xerrors.NewErr("hardware disk failure") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}
			if tt.customOsOpenFn != nil {
				mockFunction(t, &fnMockable_OsOpen, tt.customOsOpenFn)
			}
			if tt.customReadDirFn != nil {
				mockFunction(t, &fnMockable_ReadDir, tt.customReadDirFn)
			}

			got, err := IsEmptyDir(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("IsEmptyDir() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}
			if got != tt.want {
				t.Errorf("IsEmptyDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFileDirExists(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		want                 bool
		customRetrievePathFn func(string) (string, error)
		customIsDirFn        func(string) bool
	}{
		{
			name:  "Should return false immediately if filePath is empty",
			input: "",
			want:  false,
			// O guard clause resolve direto, não precisa injetar mocks
		},
		{
			name:                 "Should return false if RetrieveFullPath returns an error",
			input:                "invalid/path",
			want:                 false,
			customRetrievePathFn: func(p string) (string, error) { return "", xerrors.NewErr("err") },
		},
		{
			name:                 "Should return false if the parent directory does not exist",
			input:                "~/documents/missing-folder/file.txt",
			want:                 false,
			customRetrievePathFn: func(p string) (string, error) { return "/mock/home/documents/missing-folder/file.txt", nil },
			customIsDirFn:        func(p string) bool { return false },
		},
		{
			name:                 "Should return true if the parent directory exists",
			input:                "~/documents/project/file.txt",
			want:                 true,
			customRetrievePathFn: func(p string) (string, error) { return "/mock/home/documents/project/file.txt", nil },
			customIsDirFn:        func(p string) bool { return true },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}
			if tt.customIsDirFn != nil {
				mockFunction(t, &fnMockable_IsDir, tt.customIsDirFn)
			}

			if got := IsFileDirExists(tt.input); got != tt.want {
				t.Errorf("IsFileDirExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsReadable(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		want                 bool
		customRetrievePathFn func(string) (string, error)
		customOsStatFn       func(string) (os.FileInfo, error)
		customOsOpenFn       func(string) (*os.File, error)
		customReadDirFn      func(*os.File, int) ([]os.DirEntry, error)
		customOsOpenFileFn   func(string, int, os.FileMode) (*os.File, error)
	}{
		{
			name:                 "Should return false if RetrieveFullPath fails",
			input:                "invalid/path",
			want:                 false,
			customRetrievePathFn: func(p string) (string, error) { return "", xerrors.NewErr("err") },
		},
		{
			name:                 "Should return false if os.Stat fails",
			input:                "~/missing-path",
			want:                 false,
			customRetrievePathFn: func(p string) (string, error) { return "/mock", nil },
			customOsStatFn:       func(n string) (os.FileInfo, error) { return nil, os.ErrNotExist },
		},
		{
			name:                 "Directory: Should return false if os.Open fails (permission denied)",
			input:                "~/protected-dir",
			want:                 false,
			customRetrievePathFn: func(p string) (string, error) { return ".", nil },
			customOsStatFn:       func(n string) (os.FileInfo, error) { return os.Stat(".") },
			customOsOpenFn:       func(n string) (*os.File, error) { return nil, os.ErrPermission },
		},
		{
			name:                 "Directory: Should return true if ReadDir succeeds (not empty)",
			input:                "~/populated-dir",
			want:                 true,
			customRetrievePathFn: func(p string) (string, error) { return ".", nil },
			customOsStatFn:       func(n string) (os.FileInfo, error) { return os.Stat(".") },
			customOsOpenFn:       func(n string) (*os.File, error) { return os.Open(".") },
			customReadDirFn:      func(f *os.File, n int) ([]os.DirEntry, error) { return []os.DirEntry{nil}, nil },
		},
		{
			name:                 "Directory: Should return true if ReadDir returns EOF (empty directory)",
			input:                "~/empty-dir",
			want:                 true,
			customRetrievePathFn: func(p string) (string, error) { return ".", nil },
			customOsStatFn:       func(n string) (os.FileInfo, error) { return os.Stat(".") },
			customOsOpenFn:       func(n string) (*os.File, error) { return os.Open(".") },
			customReadDirFn:      func(f *os.File, n int) ([]os.DirEntry, error) { return nil, io.EOF },
		},
		{
			name:                 "Directory: Should return false if ReadDir returns an unexpected error",
			input:                "~/broken-dir",
			want:                 false,
			customRetrievePathFn: func(p string) (string, error) { return ".", nil },
			customOsStatFn:       func(n string) (os.FileInfo, error) { return os.Stat(".") },
			customOsOpenFn:       func(n string) (*os.File, error) { return os.Open(".") },
			customReadDirFn:      func(f *os.File, n int) ([]os.DirEntry, error) { return nil, xerrors.NewErr("disk err") },
		},
		{
			name:                 "File: Should return false if os.OpenFile fails (permission denied)",
			input:                "~/protected-file.txt",
			want:                 false,
			customRetrievePathFn: func(p string) (string, error) { return "/mock/file.txt", nil },
			// Technical Note: Force IsDir() to strictly evaluate as false using a dry stub mock
			customOsStatFn: func(n string) (os.FileInfo, error) {
				return mockFileInfo{isDir: false}, nil
			},
			customOsOpenFileFn: func(n string, f int, m os.FileMode) (*os.File, error) { return nil, os.ErrPermission },
		},
		{
			name:                 "File: Should return true if os.OpenFile succeeds",
			input:                "~/readable-file.txt",
			want:                 true,
			customRetrievePathFn: func(p string) (string, error) { return "/mock/file.txt", nil },
			customOsStatFn: func(n string) (os.FileInfo, error) {
				return mockFileInfo{isDir: false}, nil
			},
			customOsOpenFileFn: func(n string, f int, m os.FileMode) (*os.File, error) {
				// Technical Note: Creates a real, isolated descriptor tracking context to bypass path slips
				tmpFile, err := os.CreateTemp(t.TempDir(), "readable_stub")
				return tmpFile, err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}
			if tt.customOsStatFn != nil {
				mockFunction(t, &fnMockable_OsStat, tt.customOsStatFn)
			}
			if tt.customOsOpenFn != nil {
				mockFunction(t, &fnMockable_OsOpen, tt.customOsOpenFn)
			}
			if tt.customReadDirFn != nil {
				mockFunction(t, &fnMockable_ReadDir, tt.customReadDirFn)
			}
			if tt.customOsOpenFileFn != nil {
				mockFunction(t, &fnMockable_OsOpenFile, tt.customOsOpenFileFn)
			}

			if got := IsReadable(tt.input); got != tt.want {
				t.Errorf("IsReadable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsWritable(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		want                 bool
		customRetrievePathFn func(string) (string, error)
		customOsStatFn       func(string) (os.FileInfo, error)
		customOsCreateTempFn func(string, string) (*os.File, error)
		customOsOpenFileFn   func(string, int, os.FileMode) (*os.File, error)
	}{
		{
			name:                 "Should return false if RetrieveFullPath fails",
			input:                "invalid/path",
			want:                 false,
			customRetrievePathFn: func(p string) (string, error) { return "", xerrors.NewErr("err") },
		},
		{
			name:                 "Should return false if os.Stat fails",
			input:                "~/missing-path",
			want:                 false,
			customRetrievePathFn: func(p string) (string, error) { return "/mock", nil },
			customOsStatFn:       func(n string) (os.FileInfo, error) { return nil, os.ErrNotExist },
		},
		{
			name:                 "Directory: Should return false if os.CreateTemp fails (read-only filesystem)",
			input:                "~/readonly-dir",
			want:                 false,
			customRetrievePathFn: func(p string) (string, error) { return "/mock/dir", nil },
			customOsStatFn: func(n string) (os.FileInfo, error) {
				return mockFileInfo{isDir: true}, nil
			},
			customOsCreateTempFn: func(d, p string) (*os.File, error) { return nil, os.ErrPermission },
		},
		{
			name:                 "Directory: Should return true if os.CreateTemp succeeds",
			input:                "~/writable-dir",
			want:                 true,
			customRetrievePathFn: func(p string) (string, error) { return "/mock/dir", nil },
			customOsStatFn: func(n string) (os.FileInfo, error) {
				return mockFileInfo{isDir: true}, nil
			},
			// Technical Note: Returns a real isolated disk descriptor in memory to safely bypass physical path slips
			customOsCreateTempFn: func(d, p string) (*os.File, error) {
				tmpFile, err := os.CreateTemp(t.TempDir(), "writable_dir_stub")
				return tmpFile, err
			},
		},
		{
			name:                 "File: Should return false if os.OpenFile fails (no write permission)",
			input:                "~/protected-file.txt",
			want:                 false,
			customRetrievePathFn: func(p string) (string, error) { return "/mock/file.txt", nil },
			customOsStatFn: func(n string) (os.FileInfo, error) {
				return mockFileInfo{isDir: false}, nil
			},
			customOsOpenFileFn: func(n string, f int, m os.FileMode) (*os.File, error) { return nil, os.ErrPermission },
		},
		{
			name:                 "File: Should return true if os.OpenFile succeeds",
			input:                "~/writable-file.txt",
			want:                 true,
			customRetrievePathFn: func(p string) (string, error) { return "/mock/file.txt", nil },
			customOsStatFn: func(n string) (os.FileInfo, error) {
				return mockFileInfo{isDir: false}, nil
			},
			// Technical Note: Returns a legitimate ephemeral file instance to simulate a valid active writer stream descriptor
			customOsOpenFileFn: func(n string, f int, m os.FileMode) (*os.File, error) {
				tmpFile, err := os.CreateTemp(t.TempDir(), "writable_file_stub")
				return tmpFile, err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}
			if tt.customOsStatFn != nil {
				mockFunction(t, &fnMockable_OsStat, tt.customOsStatFn)
			}
			if tt.customOsCreateTempFn != nil {
				mockFunction(t, &fnMockable_OsCreateTemp, tt.customOsCreateTempFn)
			}
			if tt.customOsOpenFileFn != nil {
				mockFunction(t, &fnMockable_OsOpenFile, tt.customOsOpenFileFn)
			}

			// Mock padrão para o os.Remove para evitar que tente apagar arquivos do sistema real
			mockFunction(t, &fnMockable_OsRemove, func(name string) error { return nil })

			if got := IsWritable(tt.input); got != tt.want {
				t.Errorf("IsWritable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPermission(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		want                 os.FileMode
		wantErr              bool
		errContains          string
		customRetrievePathFn func(string) (string, error)
		customOsStatFn       func(string) (os.FileInfo, error)
	}{
		{
			name:        "Should return error if RetrieveFullPath fails",
			input:       "invalid/path",
			want:        0,
			wantErr:     true,
			errContains: "[CTX: XFS][MSG: failed to resolve user home directory]",
			// Technical Note: Simulates a real failure leaking natively from the core engine
			customRetrievePathFn: func(p string) (string, error) {
				return "", errors.New("[CTX: XFS][MSG: failed to resolve user home directory]")
			},
		},
		{
			name:                 "Should return error if os.Stat fails",
			input:                "~/missing-path",
			want:                 0,
			wantErr:              true,
			errContains:          "[MSG: failed to read path metadata][FIELD: expandedPath]",
			customRetrievePathFn: func(p string) (string, error) { return "/mock", nil },
			customOsStatFn:       func(n string) (os.FileInfo, error) { return nil, os.ErrNotExist },
		},
		{
			name:                 "Should return correct file permissions on success",
			input:                "~/notes.txt",
			want:                 os.FileMode(0644),
			wantErr:              false,
			customRetrievePathFn: func(p string) (string, error) { return "/mock/notes.txt", nil },
			customOsStatFn: func(n string) (os.FileInfo, error) {
				return mockFileInfo{mode: os.FileMode(0644)}, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}
			if tt.customOsStatFn != nil {
				mockFunction(t, &fnMockable_OsStat, tt.customOsStatFn)
			}

			got, err := GetPermission(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetPermission() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}
			if got != tt.want {
				t.Errorf("GetPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetPermission(t *testing.T) {
	tests := []struct {
		name                 string
		inputPath            string
		inputPerm            os.FileMode
		wantErr              bool
		errContains          string
		customRetrievePathFn func(string) (string, error)
		customOsChmodFn      func(string, os.FileMode) error
	}{
		{
			name:        "Should return error if RetrieveFullPath fails",
			inputPath:   "invalid/path",
			inputPerm:   0644,
			wantErr:     true,
			errContains: "[CTX: XFS][MSG: failed to resolve user home directory]",
			// Technical Note: Simulates a real failure leaking natively from the core engine
			customRetrievePathFn: func(p string) (string, error) {
				return "", errors.New("[CTX: XFS][MSG: failed to resolve user home directory]")
			},
		},
		{
			name:        "Should return error if os.Chmod fails (permission denied)",
			inputPath:   "~/protected-file.txt",
			inputPerm:   0777,
			wantErr:     true,
			errContains: "[MSG: failed to change target filesystem permissions][FIELD: expandedPath]",
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/protected-file.txt", nil
			},
			customOsChmodFn: func(path string, mode os.FileMode) error {
				return os.ErrPermission
			},
		},
		{
			name:      "Should return nil on success",
			inputPath: "~/regular-file.txt",
			inputPerm: 0600,
			wantErr:   false,
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/regular-file.txt", nil
			},
			customOsChmodFn: func(path string, mode os.FileMode) error {
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}
			if tt.customOsChmodFn != nil {
				mockFunction(t, &fnMockable_OsChmod, tt.customOsChmodFn)
			}

			err := SetPermission(tt.inputPath, tt.inputPerm)

			if (err != nil) != tt.wantErr {
				t.Fatalf("SetPermission() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			}
		})
	}
}

func TestCreateFile(t *testing.T) {
	tests := []struct {
		name                 string
		inputPath            string
		inputPerm            []os.FileMode
		wantErr              bool
		errContains          string
		customRetrievePathFn func(string) (string, error)
		customOsOpenFileFn   func(string, int, os.FileMode) (*os.File, error)
	}{
		{
			name:        "Should return error if RetrieveFullPath fails",
			inputPath:   "invalid/path",
			inputPerm:   nil,
			wantErr:     true,
			errContains: "[CTX: XFS][MSG: failed to resolve user home directory]",
			// Technical Note: Simulates a real failure leaking natively from the core engine
			customRetrievePathFn: func(p string) (string, error) {
				return "", errors.New("[CTX: XFS][MSG: failed to resolve user home directory]")
			},
		},
		{
			name:        "Should return error if os.OpenFile fails (permission denied)",
			inputPath:   "~/protected-file.txt",
			inputPerm:   nil,
			wantErr:     true,
			errContains: "[MSG: failed to create targeted file][FIELD: expandedPath]",
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/protected-file.txt", nil
			},
			customOsOpenFileFn: func(n string, f int, m os.FileMode) (*os.File, error) {
				return nil, os.ErrPermission
			},
		},
		{
			name:      "Should succeed using default permissions (0666) when no perm is provided",
			inputPath: "~/default-perm.txt",
			inputPerm: nil,
			wantErr:   false,
			customRetrievePathFn: func(p string) (string, error) {
				tmpFile := filepath.Join(t.TempDir(), "stub_default.txt")
				return tmpFile, nil
			},
			customOsOpenFileFn: func(n string, f int, m os.FileMode) (*os.File, error) {
				if m != 0666 {
					return nil, xerrors.NewErr("expected default permission 0666, got %o", m)
				}
				return os.OpenFile(n, os.O_CREATE|os.O_WRONLY, m)
			},
		},
		{
			name:      "Should succeed using custom permission when provided via variadic parameter",
			inputPath: "~/custom-perm.txt",
			inputPerm: []os.FileMode{0600},
			wantErr:   false,
			customRetrievePathFn: func(p string) (string, error) {
				tmpFile := filepath.Join(t.TempDir(), "stub_custom.txt")
				return tmpFile, nil
			},
			customOsOpenFileFn: func(n string, f int, m os.FileMode) (*os.File, error) {
				if m != 0600 {
					return nil, xerrors.NewErr("expected custom permission 0600, got %o", m)
				}
				return os.OpenFile(n, os.O_CREATE|os.O_WRONLY, m)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}
			if tt.customOsOpenFileFn != nil {
				mockFunction(t, &fnMockable_OsOpenFile, tt.customOsOpenFileFn)
			}

			got, err := CreateFile(tt.inputPath, tt.inputPerm...)

			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}
			if got == nil {
				t.Error("Expected a valid *os.File pointer on success, got nil")
			} else {
				got.Close()
			}
		})
	}
}

func TestCreateDir(t *testing.T) {
	tests := []struct {
		name                 string
		inputPath            string
		inputPerm            []os.FileMode
		wantErr              bool
		errContains          string
		customRetrievePathFn func(string) (string, error)
		customOsMkdirFn      func(string, os.FileMode) error
	}{
		{
			name:        "Should return error if RetrieveFullPath fails",
			inputPath:   "invalid/path",
			inputPerm:   nil,
			wantErr:     true,
			errContains: "[CTX: XFS][MSG: failed to resolve user home directory]",
			// Technical Note: Simulates a real failure leaking natively from the core engine
			customRetrievePathFn: func(p string) (string, error) {
				return "", errors.New("[CTX: XFS][MSG: failed to resolve user home directory]")
			},
		},
		{
			name:        "Should return error if os.Mkdir fails (permission denied)",
			inputPath:   "~/protected-dir",
			inputPerm:   nil,
			wantErr:     true,
			errContains: "[MSG: failed to create targeted directory][FIELD: expandedPath]",
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/protected-dir", nil
			},
			customOsMkdirFn: func(path string, perm os.FileMode) error {
				return os.ErrPermission
			},
		},
		{
			name:      "Should succeed using default permissions (0755) when no perm is provided",
			inputPath: "~/default-dir",
			inputPerm: nil,
			wantErr:   false,
			customRetrievePathFn: func(p string) (string, error) {
				return filepath.Join(t.TempDir(), "sub_default"), nil
			},
			customOsMkdirFn: func(path string, perm os.FileMode) error {
				if perm != 0755 {
					return xerrors.NewErr("expected default permission 0755, got %o", perm)
				}
				return os.Mkdir(path, perm)
			},
		},
		{
			name:      "Should succeed using custom permission when provided via variadic parameter",
			inputPath: "~/custom-dir",
			inputPerm: []os.FileMode{0700},
			wantErr:   false,
			customRetrievePathFn: func(p string) (string, error) {
				return filepath.Join(t.TempDir(), "sub_custom"), nil
			},
			customOsMkdirFn: func(path string, perm os.FileMode) error {
				if perm != 0700 {
					return xerrors.NewErr("expected custom permission 0700, got %o", perm)
				}
				return os.Mkdir(path, perm)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}
			if tt.customOsMkdirFn != nil {
				mockFunction(t, &fnMockable_OsMkdir, tt.customOsMkdirFn)
			}

			// Desempacota o slice para simular o parâmetro variádico (...)
			err := CreateDir(tt.inputPath, tt.inputPerm...)

			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateDir() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			}
		})
	}
}

func TestCreateDirPath(t *testing.T) {
	tests := []struct {
		name                 string
		inputPath            string
		inputPerm            []os.FileMode
		wantErr              bool
		errContains          string
		customRetrievePathFn func(string) (string, error)
		customOsMkdirAllFn   func(string, os.FileMode) error
	}{
		{
			name:        "Should return error immediately if path is empty",
			inputPath:   "",
			inputPerm:   nil,
			wantErr:     true,
			errContains: "[MSG: empty string value not allowed][FIELD: path]",
		},
		{
			name:        "Should return error if RetrieveFullPath fails",
			inputPath:   "invalid/path",
			inputPerm:   nil,
			wantErr:     true,
			errContains: "[CTX: XFS][MSG: failed to resolve user home directory]",
			// Technical Note: Simulates a real failure leaking natively from the core engine
			customRetrievePathFn: func(p string) (string, error) {
				return "", errors.New("[CTX: XFS][MSG: failed to resolve user home directory]")
			},
		},
		{
			name:        "Should return error if os.MkdirAll fails (permission denied)",
			inputPath:   "~/protected-tree",
			inputPerm:   nil,
			wantErr:     true,
			errContains: "[MSG: failed to create targeted directory tree][FIELD: expandedPath]",
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/protected-tree", nil
			},
			customOsMkdirAllFn: func(path string, perm os.FileMode) error {
				return os.ErrPermission
			},
		},
		{
			name:      "Should succeed creating a nested tree using default permissions (0755)",
			inputPath: "~/nested/tree/default",
			inputPerm: nil,
			wantErr:   false,
			customRetrievePathFn: func(p string) (string, error) {
				return filepath.Join(t.TempDir(), "nested", "tree", "default"), nil
			},
			customOsMkdirAllFn: func(path string, perm os.FileMode) error {
				if perm != 0755 {
					return xerrors.NewErr("expected default permission 0755, got %o", perm)
				}
				return os.MkdirAll(path, perm)
			},
		},
		{
			name:      "Should succeed creating a nested tree using custom permissions",
			inputPath: "~/nested/tree/custom",
			inputPerm: []os.FileMode{0700},
			wantErr:   false,
			customRetrievePathFn: func(p string) (string, error) {
				return filepath.Join(t.TempDir(), "nested", "tree", "custom"), nil
			},
			customOsMkdirAllFn: func(path string, perm os.FileMode) error {
				if perm != 0700 {
					return xerrors.NewErr("expected custom permission 0700, got %o", perm)
				}
				return os.MkdirAll(path, perm)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}
			if tt.customOsMkdirAllFn != nil {
				mockFunction(t, &fnMockable_OsMkdirAll, tt.customOsMkdirAllFn)
			}

			// Desempacota o slice para simular o parâmetro variádico (...)
			err := CreateDirPath(tt.inputPath, tt.inputPerm...)

			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateDirPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			}
		})
	}
}

func TestOpenFileWrite(t *testing.T) {
	tests := []struct {
		name                 string
		inputPath            string
		inputTruncate        bool
		inputPerm            []os.FileMode
		wantErr              bool
		errContains          string
		customRetrievePathFn func(string) (string, error)
		customOsOpenFileFn   func(string, int, os.FileMode) (*os.File, error)
	}{
		{
			name:          "Should return error if RetrieveFullPath fails",
			inputPath:     "invalid/path",
			inputTruncate: false,
			inputPerm:     nil,
			wantErr:       true,
			errContains:   "[CTX: XFS][MSG: failed to resolve user home directory]",
			// Technical Note: Simulates a real failure leaking natively from the core engine
			customRetrievePathFn: func(p string) (string, error) {
				return "", errors.New("[CTX: XFS][MSG: failed to resolve user home directory]")
			},
		},
		{
			name:          "Should return error if os.OpenFile fails (permission denied)",
			inputPath:     "~/protected-file.txt",
			inputTruncate: false,
			inputPerm:     nil,
			wantErr:       true,
			errContains:   "[MSG: failed to open targeted file for writing operations][FIELD: expandedPath]",
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/protected-file.txt", nil
			},
			customOsOpenFileFn: func(n string, f int, m os.FileMode) (*os.File, error) {
				return nil, os.ErrPermission
			},
		},
		{
			name:          "Should succeed with APPEND flag and default permission (0644) when truncate is false",
			inputPath:     "~/append-file.txt",
			inputTruncate: false,
			inputPerm:     nil,
			wantErr:       false,
			customRetrievePathFn: func(p string) (string, error) {
				return filepath.Join(t.TempDir(), "append.txt"), nil
			},
			customOsOpenFileFn: func(n string, f int, m os.FileMode) (*os.File, error) {
				if (f & os.O_APPEND) == 0 {
					return nil, xerrors.NewErr("expected flag to contain os.O_APPEND")
				}
				if m != 0644 {
					return nil, xerrors.NewErr("expected default permission 0644, got %o", m)
				}
				return os.OpenFile(n, f, m)
			},
		},
		{
			name:          "Should succeed with TRUNC flag and custom permission when truncate is true",
			inputPath:     "~/trunc-file.txt",
			inputTruncate: true,
			inputPerm:     []os.FileMode{0600},
			wantErr:       false,
			customRetrievePathFn: func(p string) (string, error) {
				return filepath.Join(t.TempDir(), "trunc.txt"), nil
			},
			customOsOpenFileFn: func(n string, f int, m os.FileMode) (*os.File, error) {
				if (f & os.O_TRUNC) == 0 {
					return nil, xerrors.NewErr("expected flag to contain os.O_TRUNC")
				}
				if m != 0600 {
					return nil, xerrors.NewErr("expected custom permission 0600, got %o", m)
				}
				return os.OpenFile(n, f, m)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}
			if tt.customOsOpenFileFn != nil {
				mockFunction(t, &fnMockable_OsOpenFile, tt.customOsOpenFileFn)
			}

			// Desempacota o slice de permissões para suprir o parâmetro variádico (...)
			got, err := OpenFileWrite(tt.inputPath, tt.inputTruncate, tt.inputPerm...)

			if (err != nil) != tt.wantErr {
				t.Fatalf("OpenFileWrite() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}
			if got == nil {
				t.Error("Expected a valid *os.File pointer on success, got nil")
			} else {
				got.Close() // Fecha o descritor de arquivo seguro gerado pelo t.TempDir()
			}
		})
	}
}

func TestOpenFileRead(t *testing.T) {
	tests := []struct {
		name                 string
		inputPath            string
		wantErr              bool
		errContains          string
		customRetrievePathFn func(string) (string, error)
		customOsOpenFileFn   func(string, int, os.FileMode) (*os.File, error)
	}{
		{
			name:        "Should return error if RetrieveFullPath fails",
			inputPath:   "invalid/path",
			wantErr:     true,
			errContains: "[CTX: XFS][MSG: failed to resolve user home directory]",
			// Technical Note: Simulates a real failure leaking natively from the core engine
			customRetrievePathFn: func(p string) (string, error) {
				return "", errors.New("[CTX: XFS][MSG: failed to resolve user home directory]")
			},
		},
		{
			name:        "Should return error if os.OpenFile fails (file not found)",
			inputPath:   "~/missing-file.txt",
			wantErr:     true,
			errContains: "[MSG: failed to open targeted file for reading operations][FIELD: expandedPath]",
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/missing-file.txt", nil
			},
			customOsOpenFileFn: func(n string, f int, m os.FileMode) (*os.File, error) {
				return nil, os.ErrNotExist
			},
		},
		{
			name:      "Should succeed opening a file with RDONLY flag",
			inputPath: "~/readable-file.txt",
			wantErr:   false,
			customRetrievePathFn: func(p string) (string, error) {
				tmpFile := filepath.Join(t.TempDir(), "read_stub.txt")
				if err := os.WriteFile(tmpFile, []byte("data"), 0644); err != nil {
					t.Fatalf("failed to setup test file: %v", err)
				}
				return tmpFile, nil
			},
			customOsOpenFileFn: func(n string, f int, m os.FileMode) (*os.File, error) {
				if f != os.O_RDONLY {
					return nil, xerrors.NewErr("expected flag os.O_RDONLY, got %d", f)
				}
				return os.OpenFile(n, f, m)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}
			if tt.customOsOpenFileFn != nil {
				mockFunction(t, &fnMockable_OsOpenFile, tt.customOsOpenFileFn)
			}

			got, err := OpenFileRead(tt.inputPath)

			if (err != nil) != tt.wantErr {
				t.Fatalf("OpenFileRead() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}
			if got == nil {
				t.Error("Expected a valid *os.File pointer on success, got nil")
			} else {
				got.Close() // Fecha o descritor de arquivo gerado na pasta temporária com segurança
			}
		})
	}
}

func TestDeleteFile(t *testing.T) {
	tests := []struct {
		name                 string
		inputPath            string
		wantErr              bool
		errContains          string
		customRetrievePathFn func(string) (string, error)
		customIsDirFn        func(string) bool
		customOsRemoveFn     func(string) error
	}{
		{
			name:        "Should return error if RetrieveFullPath fails",
			inputPath:   "invalid/path",
			wantErr:     true,
			errContains: "[CTX: XFS][MSG: failed to resolve user home directory]",
			// Technical Note: Simulates a real failure leaking natively from the core engine
			customRetrievePathFn: func(p string) (string, error) {
				return "", errors.New("[CTX: XFS][MSG: failed to resolve user home directory]")
			},
		},
		{
			name:        "Should return error if target path is a directory",
			inputPath:   "~/documents",
			wantErr:     true,
			errContains: "[MSG: cannot use DeleteFile on a directory][FIELD: expandedPath][VALUE: /mock/documents][EXPECTED_TYPE: regular file]",
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/documents", nil
			},
			customIsDirFn: func(p string) bool {
				return true
			},
		},
		{
			name:        "Should return error if os.Remove fails (file protected or missing)",
			inputPath:   "~/locked.txt",
			wantErr:     true,
			errContains: "[MSG: failed to delete file at targeted path][FIELD: expandedPath]",
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/locked.txt", nil
			},
			customIsDirFn: func(p string) bool {
				return false
			},
			customOsRemoveFn: func(name string) error {
				return os.ErrPermission
			},
		},
		{
			name:      "Should succeed and return nil when file is deleted",
			inputPath: "~/temporary.txt",
			wantErr:   false,
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/temporary.txt", nil
			},
			customIsDirFn: func(p string) bool {
				return false
			},
			customOsRemoveFn: func(name string) error {
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}

			// CORREÇÃO: Garante que o IsDir chame a função mockada injetada no cenário,
			// evitando que a lógica real tente buscar metadados no disco físico.
			if tt.customIsDirFn != nil {
				mockFunction(t, &fnMockable_IsDir, tt.customIsDirFn)
			} else {
				mockFunction(t, &fnMockable_IsDir, func(p string) bool { return false })
			}

			if tt.customOsRemoveFn != nil {
				mockFunction(t, &fnMockable_OsRemove, tt.customOsRemoveFn)
			} else {
				mockFunction(t, &fnMockable_OsRemove, func(name string) error { return nil })
			}

			// BLINDAGEM EXTRA: Se por algum motivo o IsDir interno for invocado nativamente,
			// este mock impede o os.Stat de tentar acessar o disco rígido da sua máquina.
			mockFunction(t, &fnMockable_OsStat, func(name string) (os.FileInfo, error) {
				return nil, nil
			})

			err := DeleteFile(tt.inputPath)

			if (err != nil) != tt.wantErr {
				t.Fatalf("DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			}
		})
	}
}

func TestDeleteDir(t *testing.T) {
	tests := []struct {
		name                 string
		inputPath            string
		inputRecursive       bool
		wantErr              bool
		errContains          string
		customRetrievePathFn func(string) (string, error)
		customIsDirFn        func(string) bool
		customOsRemoveFn     func(string) error
		customOsRemoveAllFn  func(string) error
	}{
		{
			name:           "Should return error if RetrieveFullPath fails",
			inputPath:      "invalid/path",
			inputRecursive: false,
			wantErr:        true,
			errContains:    "[CTX: XFS][MSG: failed to resolve user home directory]",
			// Technical Note: Simulates a real failure leaking natively from the core engine
			customRetrievePathFn: func(p string) (string, error) {
				return "", errors.New("[CTX: XFS][MSG: failed to resolve user home directory]")
			},
		},
		{
			name:           "Should return error if target path is not a directory",
			inputPath:      "~/notes.txt",
			inputRecursive: false,
			wantErr:        true,
			errContains:    "[MSG: path is not a directory or does not exist][FIELD: expandedPath][VALUE: /mock/notes.txt][EXPECTED_TYPE: directory]",
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/notes.txt", nil
			},
			customIsDirFn: func(p string) bool {
				return false
			},
		},
		{
			name:           "Non-Recursive: Should return error if os.Remove fails (directory not empty)",
			inputPath:      "~/full-dir",
			inputRecursive: false,
			wantErr:        true,
			errContains:    "[MSG: failed to delete targeted directory hierarchy, recursive='false'][FIELD: expandedPath][TGT: /mock/full-dir]::[ERR: directory not empty]",
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/full-dir", nil
			},
			customIsDirFn: func(p string) bool {
				return true
			},
			customOsRemoveFn: func(name string) error {
				return errors.New("directory not empty")
			},
		},
		{
			name:           "Non-Recursive: Should succeed if directory is empty",
			inputPath:      "~/empty-dir",
			inputRecursive: false,
			wantErr:        false,
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/empty-dir", nil
			},
			customIsDirFn: func(p string) bool {
				return true
			},
			customOsRemoveFn: func(name string) error {
				return nil
			},
		},
		{
			name:           "Recursive: Should return error if os.RemoveAll fails (permission denied)",
			inputPath:      "~/protected-tree",
			inputRecursive: true,
			wantErr:        true,
			errContains:    "[MSG: failed to delete targeted directory hierarchy, recursive='true'][FIELD: expandedPath][TGT: /mock/protected-tree]::[ERR: permission denied]",
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/protected-tree", nil
			},
			customIsDirFn: func(p string) bool {
				return true
			},
			customOsRemoveAllFn: func(path string) error {
				return os.ErrPermission
			},
		},
		{
			name:           "Recursive: Should succeed deleting the entire directory tree",
			inputPath:      "~/nested-tree",
			inputRecursive: true,
			wantErr:        false,
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/nested-tree", nil
			},
			customIsDirFn: func(p string) bool {
				return true
			},
			customOsRemoveAllFn: func(path string) error {
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}
			if tt.customIsDirFn != nil {
				mockFunction(t, &fnMockable_IsDir, tt.customIsDirFn)
			} else {
				mockFunction(t, &fnMockable_IsDir, func(p string) bool { return true })
			}
			if tt.customOsRemoveFn != nil {
				mockFunction(t, &fnMockable_OsRemove, tt.customOsRemoveFn)
			} else {
				mockFunction(t, &fnMockable_OsRemove, func(name string) error { return nil })
			}
			if tt.customOsRemoveAllFn != nil {
				mockFunction(t, &fnMockable_OsRemoveAll, tt.customOsRemoveAllFn)
			} else {
				mockFunction(t, &fnMockable_OsRemoveAll, func(path string) error { return nil })
			}

			err := DeleteDir(tt.inputPath, tt.inputRecursive)

			if (err != nil) != tt.wantErr {
				t.Fatalf("DeleteDir() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			}
		})
	}
}

func TestHasSpaceAvailable(t *testing.T) {
	tests := []struct {
		name                   string
		inputPath              string
		inputBytes             uint64
		want                   bool
		wantErr                bool
		errContains            string
		customRetrievePathFn   func(string) (string, error)
		customExistsFn         func(string) bool
		customGetVolumeSpaceFn func(string, uint64) (bool, error)
	}{
		{
			name:        "Should return error if RetrieveFullPath fails",
			inputPath:   "invalid/path",
			inputBytes:  1024,
			want:        false,
			wantErr:     true,
			errContains: "[CTX: XFS][MSG: failed to resolve user home directory]",
			// Technical Note: Simulates a real failure leaking natively from the core engine
			customRetrievePathFn: func(p string) (string, error) {
				return "", errors.New("[CTX: XFS][MSG: failed to resolve user home directory]")
			},
		},
		{
			name:        "Should return error if it loops up to the root and finds no valid path",
			inputPath:   "~/some/nested/path",
			inputBytes:  1024,
			want:        false,
			wantErr:     true,
			errContains: "[MSG: could not find a valid base volume path][FIELD: expandedPath][TGT: /mock/some/nested/path]",
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/some/nested/path", nil
			},
			customExistsFn: func(p string) bool {
				return false
			},
		},
		{
			name:       "Should break the loop immediately if the expanded path exists",
			inputPath:  "~/existing-dir",
			inputBytes: 500,
			want:       true,
			wantErr:    false,
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/existing-dir", nil
			},
			customExistsFn: func(p string) bool {
				return p == "/mock/existing-dir"
			},
			customGetVolumeSpaceFn: func(path string, bytes uint64) (bool, error) {
				if path != "/mock/existing-dir" {
					return false, errors.New("expected validation track check failure")
				}
				return true, nil
			},
		},
		{
			name:       "Should loop upwards until finding an existing parent directory",
			inputPath:  "/mock/root/dir/missing1/missing2",
			inputBytes: 2000,
			want:       true,
			wantErr:    false,
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/root/dir/missing1/missing2", nil
			},
			customExistsFn: func(p string) bool {
				return p == filepath.FromSlash("/mock/root/dir")
			},
			customGetVolumeSpaceFn: func(path string, bytes uint64) (bool, error) {
				if path != filepath.FromSlash("/mock/root/dir") {
					return false, errors.New("expected validation track check failure")
				}
				return true, nil
			},
		},
		{
			name:        "Should propagate errors returned by getVolumeFreeSpace function",
			inputPath:   "~/existing-dir",
			inputBytes:  100,
			want:        false,
			wantErr:     true,
			errContains: "syscall disk failure",
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/existing-dir", nil
			},
			customExistsFn: func(p string) bool {
				return true
			},
			customGetVolumeSpaceFn: func(path string, bytes uint64) (bool, error) {
				return false, errors.New("syscall disk failure")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}
			if tt.customExistsFn != nil {
				mockFunction(t, &fnMockable_Exists, tt.customExistsFn)
			} else {
				mockFunction(t, &fnMockable_Exists, func(p string) bool { return true })
			}
			if tt.customGetVolumeSpaceFn != nil {
				mockFunction(t, &fnMockable_GetVolumeFreeSpace, tt.customGetVolumeSpaceFn)
			} else {
				mockFunction(t, &fnMockable_GetVolumeFreeSpace, func(p string, b uint64) (bool, error) { return true, nil })
			}

			got, err := HasSpaceAvailable(tt.inputPath, tt.inputBytes)

			if (err != nil) != tt.wantErr {
				t.Fatalf("HasSpaceAvailable() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}
			if got != tt.want {
				t.Errorf("HasSpaceAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFileSize(t *testing.T) {
	tests := []struct {
		name                 string
		inputPath            string
		want                 int64
		wantErr              bool
		errContains          error
		customRetrievePathFn func(string) (string, error)
		customOsStatFn       func(string) (os.FileInfo, error)
	}{
		{
			name:      "Should return error if RetrieveFullPath fails",
			inputPath: "invalid/path",
			want:      0,
			wantErr:   true,
			customRetrievePathFn: func(p string) (string, error) {
				return "", xerrors.NewErr("mocked err")
			},
		},
		{
			name:      "Should return error if os.Stat fails (file not found)",
			inputPath: "~/missing-file.txt",
			want:      0,
			wantErr:   true,
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/missing-file.txt", nil
			},
			customOsStatFn: func(n string) (os.FileInfo, error) {
				return nil, os.ErrNotExist
			},
		},
		{
			name:        "Should return os.ErrInvalid if the path targets a directory",
			inputPath:   "~/documents",
			want:        0,
			wantErr:     true,
			errContains: os.ErrInvalid,
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/documents", nil
			},
			customOsStatFn: func(n string) (os.FileInfo, error) {
				return mockFileInfo{isDir: true}, nil
			},
		},
		{
			name:      "Should return the correct file size on success",
			inputPath: "~/notes.txt",
			want:      int64(4096), // Technical Note: Enforces a predictable in-memory volume size constraint
			wantErr:   false,
			customRetrievePathFn: func(p string) (string, error) {
				return "/mock/notes.txt", nil
			},
			customOsStatFn: func(n string) (os.FileInfo, error) {
				return mockFileInfo{isDir: false, size: 4096}, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}
			if tt.customOsStatFn != nil {
				mockFunction(t, &fnMockable_OsStat, tt.customOsStatFn)
			}

			got, err := GetFileSize(tt.inputPath)

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFileSize() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.errContains != nil {
				if !errors.Is(err, tt.errContains) {
					t.Errorf("Expected error to be %v, got %v", tt.errContains, err)
				}
				return
			}
			if got != tt.want {
				t.Errorf("GetFileSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSameFile(t *testing.T) {
	tests := []struct {
		name                 string
		inputA               string
		inputB               string
		want                 bool
		wantErr              bool
		customRetrievePathFn func(string) (string, error)
		customOsStatFn       func(string) (os.FileInfo, error)
	}{
		{
			name:    "Should return error if RetrieveFullPath fails for pathA",
			inputA:  "invalid/path",
			inputB:  "valid/path",
			want:    false,
			wantErr: true,
			customRetrievePathFn: func(p string) (string, error) {
				if p == "invalid/path" {
					return "", xerrors.NewErr("mocked err A")
				}
				return "/mock/valid", nil
			},
		},
		{
			name:    "Should return error if RetrieveFullPath fails for pathB",
			inputA:  "valid/path",
			inputB:  "invalid/path",
			want:    false,
			wantErr: true,
			customRetrievePathFn: func(p string) (string, error) {
				if p == "invalid/path" {
					return "", xerrors.NewErr("mocked err B")
				}
				return "/mock/valid", nil
			},
		},
		{
			name:    "Should return error if os.Stat fails for pathA",
			inputA:  "~/missing1.txt",
			inputB:  "~/notes.txt",
			want:    false,
			wantErr: true,
			customRetrievePathFn: func(p string) (string, error) {
				return p, nil
			},
			customOsStatFn: func(n string) (os.FileInfo, error) {
				if n == "~/missing1.txt" {
					return nil, os.ErrNotExist
				}
				tmp, err := os.CreateTemp(t.TempDir(), "same_file_stub")
				if err != nil {
					return nil, err
				}
				defer tmp.Close()
				return tmp.Stat()
			},
		},
		{
			name:    "Should return error if os.Stat fails for pathB",
			inputA:  "~/notes.txt",
			inputB:  "~/missing2.txt",
			want:    false,
			wantErr: true,
			customRetrievePathFn: func(p string) (string, error) {
				return p, nil
			},
			customOsStatFn: func(n string) (os.FileInfo, error) {
				if n == "~/missing2.txt" {
					return nil, os.ErrNotExist
				}
				tmp, err := os.CreateTemp(t.TempDir(), "same_file_stub")
				if err != nil {
					return nil, err
				}
				defer tmp.Close()
				return tmp.Stat()
			},
		},
		{
			name:    "Should return true if both paths point to the exact same file",
			inputA:  "~/notes.txt",
			inputB:  "~/notes_alias.txt",
			want:    true,
			wantErr: false,
			customRetrievePathFn: func(p string) (string, error) {
				return p, nil
			},
			// Technical Note: Instantiates a temporary file, captures metadata, and closes the handle immediately to prevent Windows file locking
			customOsStatFn: func() func(string) (os.FileInfo, error) {
				tmpFile, _ := os.CreateTemp(t.TempDir(), "shared_inode_stub")
				info, _ := tmpFile.Stat()
				_ = tmpFile.Close() // Bypasses file tracking locks
				return func(n string) (os.FileInfo, error) {
					return info, nil
				}
			}(),
		},
		{
			name:    "Should return false if paths point to different files",
			inputA:  "~/notes.txt",
			inputB:  "~/documents",
			want:    false,
			wantErr: false,
			customRetrievePathFn: func(p string) (string, error) {
				return p, nil
			},
			// Technical Note: Safely closes both descriptor streams before execution to allow the test framework to run its cleanup
			customOsStatFn: func() func(string) (os.FileInfo, error) {
				tmpA, _ := os.CreateTemp(t.TempDir(), "file_a")
				tmpB, _ := os.CreateTemp(t.TempDir(), "file_b")
				infoA, _ := tmpA.Stat()
				infoB, _ := tmpB.Stat()
				_ = tmpA.Close()
				_ = tmpB.Close()
				return func(n string) (os.FileInfo, error) {
					if n == "~/documents" {
						return infoB, nil
					}
					return infoA, nil
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.customRetrievePathFn != nil {
				mockFunction(t, &fnMockable_RetrieveFullPath, tt.customRetrievePathFn)
			}
			if tt.customOsStatFn != nil {
				mockFunction(t, &fnMockable_OsStat, tt.customOsStatFn)
			}

			got, err := IsSameFile(tt.inputA, tt.inputB)

			if (err != nil) != tt.wantErr {
				t.Fatalf("IsSameFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("IsSameFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

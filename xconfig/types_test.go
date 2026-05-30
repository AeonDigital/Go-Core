package xconfig_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AeonDigital/Go-Core/xconfig"
)

func TestOptions_ValidateExclusivePaths(t *testing.T) {
	parserType := "test-parser"

	// Define all 8 mathematical combinations for the 3 path properties
	tests := []struct {
		name          string
		options       xconfig.Options
		expectError   bool
		expectedError string // Expectativa literal e cirúrgica do erro estruturado
	}{
		{
			name: "Success: no paths provided (zero count)",
			options: xconfig.Options{
				FilePath:   "",
				DirPath:    "",
				ConfigPath: "",
			},
			expectError: false,
		},
		{
			name: "Success: only FilePath provided",
			options: xconfig.Options{
				FilePath:   "/path/to/config.json",
				DirPath:    "",
				ConfigPath: "",
			},
			expectError: false,
		},
		{
			name: "Success: only DirPath provided",
			options: xconfig.Options{
				FilePath:   "",
				DirPath:    "/path/to/dir",
				ConfigPath: "",
			},
			expectError: false,
		},
		{
			name: "Success: only ConfigPath provided",
			options: xconfig.Options{
				FilePath:   "",
				DirPath:    "",
				ConfigPath: "/path/to/generic",
			},
			expectError: false,
		},
		{
			name: "Failure: FilePath and DirPath provided simultaneously",
			options: xconfig.Options{
				FilePath:   "/path/to/config.json",
				DirPath:    "/path/to/dir",
				ConfigPath: "",
			},
			expectError:   true,
			expectedError: "[CTX: XCONFIG][MSG: mutual exclusivity violation (choose only one)][FIELD: ][OPT: FilePath='/path/to/config.json', DirPath='/path/to/dir'][OPTIONS: 'FilePath', 'DirPath', 'ConfigPath']::[ERR: ø]",
		},
		{
			name: "Failure: FilePath and ConfigPath provided simultaneously",
			options: xconfig.Options{
				FilePath:   "/path/to/config.json",
				DirPath:    "",
				ConfigPath: "/path/to/generic",
			},
			expectError:   true,
			expectedError: "[CTX: XCONFIG][MSG: mutual exclusivity violation (choose only one)][FIELD: ][OPT: FilePath='/path/to/config.json', ConfigPath='/path/to/generic'][OPTIONS: 'FilePath', 'DirPath', 'ConfigPath']::[ERR: ø]",
		},
		{
			name: "Failure: DirPath and ConfigPath provided simultaneously",
			options: xconfig.Options{
				FilePath:   "",
				DirPath:    "/path/to/dir",
				ConfigPath: "/path/to/generic",
			},
			expectError:   true,
			expectedError: "[CTX: XCONFIG][MSG: mutual exclusivity violation (choose only one)][FIELD: ][OPT: DirPath='/path/to/dir', ConfigPath='/path/to/generic'][OPTIONS: 'FilePath', 'DirPath', 'ConfigPath']::[ERR: ø]",
		},
		{
			name: "Failure: all three paths provided simultaneously",
			options: xconfig.Options{
				FilePath:   "/path/to/config.json",
				DirPath:    "/path/to/dir",
				ConfigPath: "/path/to/generic",
			},
			expectError:   true,
			expectedError: "[CTX: XCONFIG][MSG: mutual exclusivity violation (choose only one)][FIELD: ][OPT: FilePath='/path/to/config.json', DirPath='/path/to/dir', ConfigPath='/path/to/generic'][OPTIONS: 'FilePath', 'DirPath', 'ConfigPath']::[ERR: ø]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.options.ValidateExclusivePaths(parserType)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected an error but got nil")
				}

				// Validação da igualdade literal do layout de erro gerado pelo utils.go
				if err.Error() != tt.expectedError {
					t.Errorf("error output layout mismatch:\n got:  %q\n want: %q", err.Error(), tt.expectedError)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestOptions_ValidatePrefixAndSeparator(t *testing.T) {
	parserType := "test-parser"

	tests := []struct {
		name          string
		options       xconfig.Options
		expectError   bool
		expectedError string // Expectativa literal e precisa do layout do erro
	}{
		{
			name: "Success: both prefix and separator are provided",
			options: xconfig.Options{
				Prefix:    "APP",
				Separator: "_",
			},
			expectError: false,
		},
		{
			name: "Failure: prefix is missing",
			options: xconfig.Options{
				Prefix:    "",
				Separator: "_",
			},
			expectError:   true,
			expectedError: "[CTX: XCONFIG][MSG: empty string value not allowed][FIELD: Prefix]::[ERR: ø]",
		},
		{
			name: "Failure: separator is missing",
			options: xconfig.Options{
				Prefix:    "APP",
				Separator: "",
			},
			expectError:   true,
			expectedError: "[CTX: XCONFIG][MSG: empty string value not allowed][FIELD: Separator]::[ERR: ø]",
		},
		{
			name: "Failure: both prefix and separator are missing",
			options: xconfig.Options{
				Prefix:    "",
				Separator: "",
			},
			expectError:   true,
			expectedError: "[CTX: XCONFIG][MSG: empty string value not allowed][FIELD: Prefix]::[ERR: ø]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.options.ValidatePrefixAndSeparator(parserType)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected an error but got nil")
				}

				// Validação cirúrgica da string exata gerada pelo motor universal
				if err.Error() != tt.expectedError {
					t.Errorf("error output layout mismatch:\n got:  %q\n want: %q", err.Error(), tt.expectedError)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestOptions_RetrieveConfigFilePaths(t *testing.T) {
	parserType := "test-parser"

	t.Run("Success: scan directory with explicit extensions and strict alphabetical sorting", func(t *testing.T) {
		tmpDir := t.TempDir()

		subDir := filepath.Join(tmpDir, "ignored_subfolder")
		if err := os.Mkdir(subDir, 0755); err != nil {
			t.Fatalf("failed to create dummy subdirectory for code coverage: %v", err)
		}

		filesToCreate := []string{
			"c_config.yml",
			"a_config.yaml",
			"b_config.json",
			"d_config.yaml",
		}

		for _, name := range filesToCreate {
			err := os.WriteFile(filepath.Join(tmpDir, name), []byte(""), 0644)
			if err != nil {
				t.Fatalf("failed to setup test file %s: %v", name, err)
			}
		}

		opts := xconfig.Options{
			DirPath: tmpDir,
		}

		extensions := []string{"yaml", ".yml"}
		result, err := opts.RetrieveConfigFilePaths(extensions, parserType)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		expectedCount := 3
		if len(result) != expectedCount {
			t.Fatalf("expected exactly %d files, got %d", expectedCount, len(result))
		}

		if !strings.HasSuffix(result[0], "a_config.yaml") {
			t.Errorf("expected index 0 to end with a_config.yaml, got %s", result[0])
		}
		if !strings.HasSuffix(result[1], "c_config.yml") {
			t.Errorf("expected index 1 to end with c_config.yml, got %s", result[1])
		}
		if !strings.HasSuffix(result[2], "d_config.yaml") {
			t.Errorf("expected index 2 to end with d_config.yaml, got %s", result[2])
		}
	})

	t.Run("Success: load single valid FilePath directly", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "app.json")

		if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
			t.Fatalf("failed to setup file: %v", err)
		}

		opts := xconfig.Options{
			FilePath: filePath,
		}

		result, err := opts.RetrieveConfigFilePaths([]string{".json"}, parserType)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(result) != 1 || result[0] != filePath {
			t.Errorf("expected path %q, got %v", filePath, result)
		}
	})

	t.Run("Success: auto-detect ConfigPath pointing to a single file", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "config.env")
		if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
			t.Fatalf("failed to setup file: %v", err)
		}

		opts := xconfig.Options{
			ConfigPath: filePath,
		}

		result, err := opts.RetrieveConfigFilePaths([]string{"env"}, parserType)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(result) != 1 || result[0] != filePath {
			t.Errorf("expected path %q, got %v", filePath, result)
		}
	})

	t.Run("Success: auto-detect ConfigPath pointing to a directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "settings.json")
		if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
			t.Fatalf("failed to setup file: %v", err)
		}

		opts := xconfig.Options{
			ConfigPath: tmpDir,
		}

		result, err := opts.RetrieveConfigFilePaths([]string{".json"}, parserType)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(result) != 1 || result[0] != filePath {
			t.Errorf("expected path %q, got %v", filePath, result)
		}
	})

	t.Run("Failure: ConfigPath does not exist on system", func(t *testing.T) {
		opts := xconfig.Options{
			ConfigPath: "/non/existent/path/here",
		}

		_, err := opts.RetrieveConfigFilePaths([]string{".json"}, parserType)
		if err == nil {
			t.Fatalf("expected error due to missing path, got nil")
		}

		// Technical Note: Bypasses rigid string checks to accommodate Windows native syscall messages safely
		errStr := err.Error()
		if !strings.Contains(errStr, "[CTX: XCONFIG][MSG: target resource not found][FIELD: ConfigPath][TGT: /non/existent/path/here]") {
			t.Errorf("error output layout mismatch metadata block, got: %q", errStr)
		}

		hasValidOSError := strings.Contains(errStr, "no such file or directory") ||
			strings.Contains(errStr, "The system cannot find the path specified")

		if !hasValidOSError {
			t.Errorf("expected underlying OS boundary error contract to report missing trackable track, got: %q", errStr)
		}
	})

	t.Run("Failure: DirPath points to an invalid or missing directory", func(t *testing.T) {
		opts := xconfig.Options{
			DirPath: "/invalid/dir/path",
		}

		_, err := opts.RetrieveConfigFilePaths([]string{".json"}, parserType)
		if err == nil {
			t.Fatalf("expected error due to missing folder target, got nil")
		}

		// Technical Note: Adjusts boundary checks to support both POSIX and Windows native directory missing messages
		errStr := err.Error()
		if !strings.Contains(errStr, "[CTX: XCONFIG][MSG: target resource currently unavailable][FIELD: dirPath][TGT: /invalid/dir/path]") {
			t.Errorf("error output layout mismatch metadata block, got: %q", errStr)
		}

		hasValidOSError := strings.Contains(errStr, "no such file or directory") ||
			strings.Contains(errStr, "The system cannot find the path specified")

		if !hasValidOSError {
			t.Errorf("expected underlying OS boundary error contract to report unavailable directory stream, got: %q", errStr)
		}

	})

	t.Run("Failure: directory matches zero file criteria parameters", func(t *testing.T) {
		tmpDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(tmpDir, "notes.txt"), []byte(""), 0644); err != nil {
			t.Fatalf("failed to create layout state: %v", err)
		}

		opts := xconfig.Options{
			DirPath: tmpDir,
		}

		_, err := opts.RetrieveConfigFilePaths([]string{".json"}, parserType)
		if err == nil {
			t.Fatalf("expected error due to zero matches found, got nil")
		}

		expectedError := fmt.Sprintf("[CTX: XCONFIG][MSG: target resource not found][FIELD: configFiles][TGT: DirPath='%s']::[ERR: ø]", tmpDir)

		if err.Error() != expectedError {
			t.Errorf("error output layout mismatch:\n got:  %q\n want: %q", err.Error(), expectedError)
		}
	})
}

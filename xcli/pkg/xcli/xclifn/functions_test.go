package xclifn_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/AeonDigital/Go-Core/sterr/pkg/sterrinterfc"
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xclifn"
)

// TestPrintStdout verifies that messages are correctly written to standard output.
func TestPrintStdout(t *testing.T) {
	// 1. Capture os.Stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe for stdout: %v", err)
	}
	os.Stdout = w

	// 2. Execute the function with a single message (behaves like Println)
	xclifn.PrintStdout("hello world")

	// 3. Execute the function with arguments (behaves like Printf + \n)
	xclifn.PrintStdout("user %s has id %d", "john", 42)

	// Close the writer so we can read from the pipe
	w.Close()
	os.Stdout = oldStdout

	// 4. Read the captured output
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read from stdout pipe: %v", err)
	}
	output := buf.String()

	// 5. Assert the results
	expected1 := "hello world\n"
	expected2 := "user john has id 42\n"

	if !strings.Contains(output, expected1) {
		t.Errorf("Print layout 1 failed.\nExpected to contain: %q\nFull output: %q", expected1, output)
	}
	if !strings.Contains(output, expected2) {
		t.Errorf("Print layout 2 failed.\nExpected to contain: %q\nFull output: %q", expected2, output)
	}
}

// TestPrintStderr verifies that messages are correctly written to standard error output.
func TestPrintStderr(t *testing.T) {
	// 1. Capture os.Stderr
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe for stderr: %v", err)
	}
	os.Stderr = w

	// 2. Execute the function with a single message (behaves like Fprintln)
	xclifn.PrintStderr("an error occurred")

	// 3. Execute the function with arguments (behaves like Fprintf + \n)
	xclifn.PrintStderr("failed to connect to %s on port %d", "localhost", 8080)

	// Close the writer so we can read from the pipe
	w.Close()
	os.Stderr = oldStderr

	// 4. Read the captured output
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read from stderr pipe: %v", err)
	}
	output := buf.String()

	// 5. Assert the results
	expected1 := "an error occurred\n"
	expected2 := "failed to connect to localhost on port 8080\n"

	if !strings.Contains(output, expected1) {
		t.Errorf("PrintStderr layout 1 failed.\nExpected to contain: %q\nFull output: %q", expected1, output)
	}
	if !strings.Contains(output, expected2) {
		t.Errorf("PrintStderr layout 2 failed.\nExpected to contain: %q\nFull output: %q", expected2, output)
	}
}

// TestNewError verifies that errors are correctly created with or without formatting arguments.
func TestNewError(t *testing.T) {
	// 1. Test single message behavior (like errors.New)
	err1 := xclifn.NewError("something went wrong")
	if err1 == nil {
		t.Fatal("expected error to be not nil")
	}
	expected1 := "something went wrong"
	if err1.Error() != expected1 {
		t.Errorf("NewError single message failed.\nExpected: %q\nGot: %q", expected1, err1.Error())
	}

	// 2. Test formatted message behavior (like fmt.Errorf)
	err2 := xclifn.NewError("invalid value %d for field %s", 500, "timeout")
	if err2 == nil {
		t.Fatal("expected formatted error to be not nil")
	}
	expected2 := "invalid value 500 for field timeout"
	if err2.Error() != expected2 {
		t.Errorf("NewError formatted message failed.\nExpected: %q\nGot: %q", expected2, err2.Error())
	}
}

// TestPrintError verifies that error messages are correctly written to standard error output.
func TestPrintError(t *testing.T) {
	// 1. Capture os.Stderr
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe for stderr: %v", err)
	}
	os.Stderr = w

	// 2. Execute with a valid error
	testErr := errors.New("database connection timeout")
	xclifn.PrintError(testErr)

	// 3. Execute with nil error (should not write anything or panic)
	xclifn.PrintError(nil)

	// Close the writer so we can read from the pipe
	w.Close()
	os.Stderr = oldStderr

	// 4. Read the captured output
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read from stderr pipe: %v", err)
	}
	output := buf.String()

	// 5. Assert the results
	expected := "database connection timeout\n"
	if output != expected {
		t.Errorf("PrintError failed.\nExpected: %q\nGot: %q", expected, output)
	}
}

//
//
//

// TestParseRawArgsSuccess validates all valid flag assignment syntax variants
// including standard pairs, assignments with equals, and implicit boolean triggers.
func TestParseRawArgsSuccess(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected map[string]string
	}{
		{
			name:  "standard flag and value pair",
			input: []string{"--output", "json", "-v", "2"},
			expected: map[string]string{
				"output": "json",
				"v":      "2",
			},
		},
		{
			name:  "inline assignment syntax using equals sign",
			input: []string{"--env=production", "-f=file.txt"},
			expected: map[string]string{
				"env": "production",
				"f":   "file.txt",
			},
		},
		{
			name:  "implicit boolean flags combined with assignments",
			input: []string{"--verbose", "-d", "--timeout=30s"},
			expected: map[string]string{
				"verbose": "true",
				"d":       "true",
				"timeout": "30s",
			},
		},
		{
			name:     "empty input returns empty map with no errors",
			input:    []string{},
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := xclifn.ParseRawArgs(tt.input)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}

			if len(res) != len(tt.expected) {
				t.Fatalf("expected map length %d, got %d", len(tt.expected), len(res))
			}

			for k, expectedVal := range tt.expected {
				if actualVal, ok := res[k]; !ok || actualVal != expectedVal {
					t.Errorf("expected flag '%s' to be '%s', got '%s'", k, expectedVal, actualVal)
				}
			}
		})
	}
}

// TestParseRawArgsErrors enforces strict syntax validation boundaries,
// verifying that every single malformed argument maps to its correct sterr instance.
func TestParseRawArgsErrors(t *testing.T) {
	tests := []struct {
		name           string
		input          []string
		expectedErrMsg string
	}{
		{
			name:           "forbidden positional or loose argument",
			input:          []string{"--output", "json", "loose-argument"},
			expectedErrMsg: "invalid argument: 'loose-argument'",
		},
		{
			name:           "empty long flag syntax structure",
			input:          []string{"--"},
			expectedErrMsg: "invalid flag provided: '--'",
		},
		{
			name:           "empty short flag syntax structure",
			input:          []string{"-"},
			expectedErrMsg: "invalid flag provided: '-'",
		},
		{
			name:           "short flag exceeding maximum character length boundary",
			input:          []string{"-longer"},
			expectedErrMsg: "invalid short flag: '-longer'",
		},
		{
			name:           "empty long flag inside an assignment clause",
			input:          []string{"--=value"},
			expectedErrMsg: "invalid flag provided: '--'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := xclifn.ParseRawArgs(tt.input)
			if err == nil {
				t.Fatalf("expected parsing to fail for input %v, but it succeeded", tt.input)
			}

			// Perform interface evaluation to extract sterr capabilities
			cliErr, ok := err.(sterrinterfc.CliError)
			if !ok {
				t.Fatalf("expected error to implement sterrinterfc.CliError, got standard error: %v", err)
			}

			// Validate the exact unifed message produced by SetMessage
			if cliErr.GetUserMessage() != tt.expectedErrMsg {
				t.Errorf("expected user message %q, got %q", tt.expectedErrMsg, cliErr.GetUserMessage())
			}

			if cliErr.GetDevMessage() != tt.expectedErrMsg {
				t.Errorf("expected dev message %q, got %q", tt.expectedErrMsg, cliErr.GetDevMessage())
			}
		})
	}
}

// TestQuantifiablePrivateHelpersDefaults explicitly targets the raw private default branches.
func TestQuantifiablePrivateHelpersDefaults(t *testing.T) {
	type customInt int
	var a customInt = 10
	var b customInt = 20

	// 1. Testa o default de isLess
	if xclifn.ExportIsLess(a, b) {
		t.Errorf("expected false for unhandled type in isLess default branch")
	}

	// 2. Testa o default de isGreater
	if xclifn.ExportIsGreater(a, b) {
		t.Errorf("expected false for unhandled type in isGreater default branch")
	}

	// 3. Testa o default de formatQuantifiable
	g, r := xclifn.ExportFormatQuantifiable(a, b, "2006-01-02")
	if g != "" || r != "" {
		t.Errorf("expected empty strings for unhandled type in formatQuantifiable default branch, got %q and %q", g, r)
	}
}

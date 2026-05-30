package xerrors_test

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/AeonDigital/Go-Core/xerrors"
)

// Custom types declaration to guarantee 100% path coverage
// over custom named primitive types inside reflection blocks.
type customInt int
type customUint uint
type customFloat float64
type customBool bool
type customString string

func TestFormatValueWithQuotes_Coverage(t *testing.T) {
	// Directly tests the internal helper logic through MsgOptions
	// to ensure all primitive kind switches evaluate identically.
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"Nil value", nil, "nil"},
		{"Native Int", 42, "42"},
		{"Custom Int", customInt(-10), "-10"},
		{"Native Uint", uint(100), "100"},
		{"Custom Uint", customUint(50), "50"},
		{"Native Float", 3.14, "3.14"},
		{"Custom Float", customFloat(2.718), "2.718"},
		{"Native Bool", true, "'true'"},
		{"Custom Bool", customBool(false), "'false'"},
		{"Native String", "hello", "'hello'"},
		{"Custom String", customString("world"), "'world'"},
		{"Fallback Complex Layout", []int{1, 2}, "'[1 2]'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We pass a custom blacklist of an empty slice to force MsgOptions
			// to evaluate the items directly through formatValueWithQuotes
			res := xerrors.MsgOptions([]any{}, tt.input)
			if res != tt.expected {
				t.Errorf("formatValueWithQuotes path failed for %s: got %q, want %q", tt.name, res, tt.expected)
			}
		})
	}
}

func TestMsgData_CoverageSuite(t *testing.T) {
	tests := []struct {
		name     string
		args     []any
		expected string
	}{
		{
			name:     "Empty arguments sequence",
			args:     []any{},
			expected: "",
		},
		{
			name:     "Default blacklist filtration (ignores nil, zero, empty string)",
			args:     []any{"k1", "v1", "k2", nil, "k3", 0, "k4", "", "k5", 100},
			expected: "k1='v1', k5=100",
		},
		{
			name:     "Boundary protection against trailing unpaired keys",
			args:     []any{"k1", "v1", "k2"},
			expected: "k1='v1', k2=(MISSING_VALUE)",
		},
		{
			name: "Custom blacklist override via slice leading argument",
			args: []any{
				[]any{"FORBIDDEN", 999}, // Custom blacklist
				"k1", "FORBIDDEN",
				"k2", 999,
				"k3", "VALID",
				"k4", "", // Should be preserved now because it's not in custom blacklist
			},
			expected: "k3='VALID', k4=''",
		},
		{
			name: "Custom blacklist edge case with empty trailing arguments",
			args: []any{
				[]any{nil},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := xerrors.MsgData(tt.args...)
			if got != tt.expected {
				t.Errorf("MsgData() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestMsgOptions_CoverageSuite(t *testing.T) {
	tests := []struct {
		name     string
		args     []any
		expected string
	}{
		{
			name:     "Empty arguments sequence",
			args:     []any{},
			expected: "",
		},
		{
			name:     "Default exclusion rules",
			args:     []any{"opt1", "", 0, nil, "opt2"},
			expected: "'opt1', 'opt2'",
		},
		{
			name: "Custom override rule preserving standard zero-values",
			args: []any{
				[]any{"IGNORE_ME"}, // Custom blacklist
				"opt1",
				"IGNORE_ME",
				"", // Preserved
				0,  // Preserved
			},
			expected: "'opt1', '', 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := xerrors.MsgOptions(tt.args...)
			if got != tt.expected {
				t.Errorf("MsgOptions() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestMsgArraySize_CoverageSuite(t *testing.T) {
	tests := []struct {
		name     string
		args     []any
		expected string
	}{
		{
			name:     "Empty arguments sequence",
			args:     []any{},
			expected: "",
		},
		{
			name:     "Standard structural length pairing (zeros are preserved)",
			args:     []any{"users", 10, "payload", 0},
			expected: "users(10), payload(0)",
		},
		{
			name:     "Boundary protection against trailing unpaired name tokens",
			args:     []any{"users", 5, "orphaned_tag"},
			expected: "users(5), orphaned_tag((MISSING_SIZE))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := xerrors.MsgArraySize(tt.args...)
			if got != tt.expected {
				t.Errorf("MsgArraySize() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestPrint_CoverageSuite(t *testing.T) {
	// Scenario 1: err is nil, function should return immediately doing nothing
	xerrors.Print(nil)

	// Scenario 2: err is present, must capture os.Stderr buffer output stream
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	testErr := errors.New("severe technical system breakdown")
	xerrors.Print(testErr)

	// Close writer stream and restore original system Stderr block
	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	outputStr := buf.String()

	if !strings.Contains(outputStr, "severe technical system breakdown") {
		t.Errorf("Print() output = %q, expected to contain %q", outputStr, testErr.Error())
	}
}

// Declaramos dois tipos inteiros separados para quebrar a igualdade de interface direta.
type customIntA int
type customIntB int

func TestShouldIgnoreValue_ForceDeepNumericReflection(t *testing.T) {
	// Criamos a blacklist usando o tipo customIntB
	customBlacklist := []any{customIntB(100)}

	// Passamos o valor usando o tipo customIntA.
	// 1. A comparação direta (customIntA(100) == customIntB(100)) vai dar FALSE.
	// 2. O código será forçado a entrar no bloco de reflexão.
	// 3. O Kind de ambos será reflect.Int (TRUE).
	// 4. Os valores extraídos vRef.Int() e iRef.Int() serão 100 (TRUE).
	res := xerrors.MsgData(customBlacklist, "test_key", customIntA(100))

	// Se a reflexão interna funcionou e chegou na linha que você quer,
	// o valor será ignorado e a string retornada deve ser vazia ("").
	if res != "" {
		t.Errorf("Expected deep numeric reflection to match different custom types and return empty string, got: %q", res)
	}
}

package xerrors_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/AeonDigital/Go-Core/xerrors"
)

// TestNewError_Table validates that NewError correctly captures all structural fields,
// enforces dynamic component runtime tracing, and handles optional internal errors.
func TestNewError_Table(t *testing.T) {
	tests := []struct {
		name         string
		inputCode    xerrors.ErrorCode
		inputErr     error
		inputMessage string
		inputData    string
		hasInternal  bool
	}{
		{
			name:         "Should capture all metadata and format string with internal error present",
			inputCode:    xerrors.ErrNetworkHTTPConnection,
			inputErr:     errors.New("timeout downstream"),
			inputMessage: "failed to reach balance service",
			inputData:    `{"retry_count": 3}`,
			hasInternal:  true,
		},
		{
			name:         "Should format string cleanly without trailing dividers when internal error is nil",
			inputCode:    xerrors.ErrCreateInstance,
			inputErr:     nil,
			inputMessage: "invalid configuration structural mapping passed",
			inputData:    "nil",
			hasInternal:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Trigger initialization factory
			gotDetailedErr := xerrors.NewError(tt.inputCode, tt.inputErr, tt.inputMessage, tt.inputData)

			if gotDetailedErr == nil {
				t.Fatal("Expected NewError to return a valid non-nil DetailedError interface instance")
			}

			// Assert getters return original properties securely
			if gotDetailedErr.Code() != tt.inputCode {
				t.Errorf("Code() = %q, want %q", gotDetailedErr.Code(), tt.inputCode)
			}
			if gotDetailedErr.Message() != tt.inputMessage {
				t.Errorf("Message() = %q, want %q", gotDetailedErr.Message(), tt.inputMessage)
			}
			if gotDetailedErr.Data() != tt.inputData {
				t.Errorf("Data() = %q, want %q", gotDetailedErr.Data(), tt.inputData)
			}

			// Validates that the runtime component capture tracks this test suite package sequence
			if !strings.Contains(gotDetailedErr.Component(), "xerrors_test") && !strings.Contains(gotDetailedErr.Component(), "TestNewError_Table") {
				t.Errorf("Component() %q expected to include tracking context frame 'xerrors_test' or 'TestNewError_Table'", gotDetailedErr.Component())
			}

			// Validate conditional string serialization outputs based on state presence
			if tt.hasInternal {
				expectedString := "[" + gotDetailedErr.Component() + "] failed to reach balance service: timeout downstream"
				if gotDetailedErr.Error() != expectedString {
					t.Errorf("Error() output layout = %q, want %q", gotDetailedErr.Error(), expectedString)
				}

				debugStr := gotDetailedErr.DebugError()
				if !strings.Contains(debugStr, "--- DETAILED DEBUG ERROR ---") || !strings.Contains(debugStr, "Internal Error: timeout downstream") {
					t.Errorf("DebugError() formatting structure is unexpected:\n%s", debugStr)
				}
			} else {
				expectedString := "[" + gotDetailedErr.Component() + "] invalid configuration structural mapping passed"
				if gotDetailedErr.Error() != expectedString {
					t.Errorf("Error() output layout without error suffix = %q, want %q", gotDetailedErr.Error(), expectedString)
				}
			}
		})
	}
}

// helperErrorFactory simulates an application error wrapper utility function
func helperErrorFactory(code xerrors.ErrorCode, msg string) xerrors.DetailedError {
	// From the developer's standpoint, we are wrapping the instantiation by exactly 1 level.
	// Therefore, we pass WithCallerSkip(1) fluently.
	return xerrors.NewError(code, nil, msg, "").WithCallerSkip(1)
}

func TestError_WithCallerSkip_Success(t *testing.T) {
	// Executes the factory tool to trigger the reflection tracking
	err := helperErrorFactory(xerrors.ErrCreateInstance, "factory error message")

	// Component must bypass helperErrorFactory and point directly to this current test block function frame
	if !strings.Contains(err.Component(), "TestError_WithCallerSkip_Success") {
		t.Errorf("expected component to map 'TestError_WithCallerSkip_Success' due to fluent skip modification, but tracked: %q", err.Component())
	}
}

func TestError_WithCallerSkip_FloorBoundary(t *testing.T) {
	err := xerrors.NewError(xerrors.ErrUnknown, nil, "testing boundary", "")

	// Enforces a negative boundary number to validate internal protection blocks
	err.WithCallerSkip(-5)

	// Component tracking must proceed smoothly and evaluate the active context without crashes
	if !strings.Contains(err.Component(), "TestError_WithCallerSkip_FloorBoundary") {
		t.Errorf("expected component to fall back to current function frame on negative input boundaries, but got: %q", err.Component())
	}
}

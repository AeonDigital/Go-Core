package xerrors_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/AeonDigital/Go-Core/xerrors"
)

func TestNewError_Table(t *testing.T) {
	xerrors.EnableDebugMode()

	tests := []struct {
		name         string
		inputCTX     xerrors.ErrorCode
		inputCode    xerrors.ErrorCode
		inputErr     error
		inputMessage string
		inputInfo    string
		hasInternal  bool
	}{
		{
			name:         "Should capture all metadata and format string with internal error present",
			inputCTX:     xerrors.ErrorCode("USR-AUTH-FLOW"),
			inputCode:    xerrors.XERR_UNEXPECTED_FAIL,
			inputErr:     errors.New("timeout downstream"),
			inputMessage: "failed to reach balance service",
			inputInfo:    `{"retry_count": 3}`,
			hasInternal:  true,
		},
		{
			name:         "Should format string cleanly without trailing dividers when internal error is nil",
			inputCTX:     xerrors.ErrorCode("SYS-DB-INIT"),
			inputCode:    xerrors.XERR_OPERATION_FAILED,
			inputErr:     nil,
			inputMessage: "invalid configuration structural mapping passed",
			inputInfo:    "nil",
			hasInternal:  false,
		},
		{
			name:         "Should fallback to registry message when inputMessage is empty string",
			inputCTX:     xerrors.ErrorCode("FALLBACK-FLOW"),
			inputCode:    xerrors.XERR_UNEXPECTED_FAIL,
			inputErr:     errors.New("something broke"),
			inputMessage: "",
			inputInfo:    "nil",
			hasInternal:  true,
		},
		{
			name:         "Should return empty string if inputMessage is empty and code does not exist in registry",
			inputCTX:     xerrors.ErrorCode("UNKNOWN-FLOW"),
			inputCode:    xerrors.ErrorCode("E999_NON_EXISTENT"),
			inputErr:     nil,
			inputMessage: "",
			inputInfo:    "nil",
			hasInternal:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDetailedErr := xerrors.NewError500(tt.inputCTX, tt.inputCode, tt.inputErr, tt.inputMessage, tt.inputInfo)

			if gotDetailedErr == nil {
				t.Fatal("Expected NewError to return a valid non-nil IError500 interface instance")
			}

			// Assert getters return original properties securely
			if gotDetailedErr.CTX() != tt.inputCTX {
				t.Errorf("CTX() = %q, want %q", gotDetailedErr.CTX(), tt.inputCTX)
			}
			if gotDetailedErr.Code() != tt.inputCode {
				t.Errorf("Code() = %q, want %q", gotDetailedErr.Code(), tt.inputCode)
			}

			// 1. Resolve a mensagem esperada aplicando a regra de fallback se o input for vazio
			expectedMsg := tt.inputMessage
			if tt.inputMessage == "" {
				switch tt.inputCode {
				case xerrors.XERR_UNEXPECTED_FAIL:
					expectedMsg = "unexpected failure"
				case xerrors.XERR_OPERATION_FAILED:
					expectedMsg = "operation failed"
				}
			}

			if gotDetailedErr.Message() != expectedMsg {
				t.Errorf("Message() = %q, want %q", gotDetailedErr.Message(), expectedMsg)
			}
			if gotDetailedErr.Info() != tt.inputInfo {
				t.Errorf("Info() = %q, want %q", gotDetailedErr.Info(), tt.inputInfo)
			}

			// Validates that the runtime component capture tracks this test suite package sequence
			comp := gotDetailedErr.Component()
			if !strings.Contains(comp, "xerrors_test") && !strings.Contains(comp, "TestNewError_Table") {
				t.Errorf("Component() %q expected to include tracking context frame 'xerrors_test' or 'TestNewError_Table'", comp)
			}

			// 2. Alinha os layouts esperados com o formato real produzido pelo método Error() da struct Error500 em modo Debug
			// 2. Alinha os layouts esperados com o formato real produzido pelo método Error() da struct Error500 em modo Debug
			strComponent := "[COMPONENT: " + comp + "]"

			strError := ""
			if tt.inputErr != nil {
				strError = " :: " + tt.inputErr.Error()
			}

			strData := ""
			if tt.inputInfo != "" {
				strData = "[INFO: " + tt.inputInfo + "]"
			}

			expectedString := fmt.Sprintf(
				"[CTX: %s][ERR: %s]%s[MSG: %s]%s%s",
				tt.inputCTX,
				tt.inputCode,
				strComponent,
				expectedMsg,
				strData,
				strError,
			)

			if gotDetailedErr.Error() != expectedString {
				t.Errorf("Error() output layout = %q, want %q", gotDetailedErr.Error(), expectedString)
			}
		})
	}
}

func TestNewError_ProductionMode_String(t *testing.T) {
	// Valida o comportamento do método Error() quando o modo debug está desativado
	xerrors.DisableDebugMode()

	ctx := xerrors.ErrorCode("PROD-FLOW")
	code := xerrors.XERR_UNEXPECTED_FAIL
	msg := "something went wrong"
	err := xerrors.NewError500(ctx, code, errors.New("root error"), msg, "")

	expected := "[CTX: PROD-FLOW][ERR: E4001][MSG: something went wrong]"
	if err.Error() != expected {
		t.Errorf("Error() in production mode = %q, want %q", err.Error(), expected)
	}
}

// helperErrorFactory simulates an application error wrapper utility function
func helperErrorFactory(ctx xerrors.ErrorCode, code xerrors.ErrorCode, msg string) xerrors.IError500 {
	// From the developer's standpoint, we are wrapping the instantiation by exactly 1 level.
	// Therefore, we pass WithCallerSkip(1) fluently.
	return xerrors.NewError500(ctx, code, nil, msg, "").WithCallerSkip(1)
}

func TestError_WithCallerSkip_Success(t *testing.T) {
	err := helperErrorFactory("SKIP-CTX", xerrors.XERR_OPERATION_FAILED, "factory error message")

	// Component must bypass helperErrorFactory and point directly to this current test block function frame
	if !strings.Contains(err.Component(), "TestError_WithCallerSkip_Success") {
		t.Errorf("expected component to map 'TestError_WithCallerSkip_Success' due to fluent skip modification, but tracked: %q", err.Component())
	}
}

func TestError_WithCallerSkip_FloorBoundary(t *testing.T) {
	err := xerrors.NewError500("BOUNDARY-CTX", xerrors.XERR_UNKNOWN, nil, "testing boundary", "")

	// Enforces a negative boundary number to validate internal protection blocks.
	modifiedErr := err.WithCallerSkip(-5)

	// Component tracking must proceed smoothly and evaluate the active context without crashes
	if !strings.Contains(modifiedErr.Component(), "TestError_WithCallerSkip_FloorBoundary") {
		t.Errorf("expected component to fall back to current function frame on negative input boundaries, but got: %q", modifiedErr.Component())
	}

	// To securely verify immutability without stack trace collisions,
	// we use a helper closure execution block to shift the caller frame depth.
	var isolatedErr xerrors.IError500
	func() {
		// Inside this anonymous block, a skip of 0 targets this closure frame,
		// separating its runtime signature from the parent execution flow.
		isolatedErr = err.WithCallerSkip(0)
	}()

	// If the original instance 'err' had been mutated by the earlier negative skip call,
	// 'isolatedErr' would match 'modifiedErr'. Since they point to different stack frames,
	// their runtime component values must be completely different.
	if isolatedErr.Component() == modifiedErr.Component() {
		t.Errorf("immutability check failed: original error was mutated and shares the same frame: %q", modifiedErr.Component())
	}
}

func TestError500_Unwrap(t *testing.T) {
	rootCause := errors.New("database connection lost")

	// Caso 1: Com erro interno (Deve retornar o erro original)
	errWithInternal := xerrors.NewError500("CTX-01", xerrors.XERR_OPERATION_FAILED, rootCause, "failed to query", "")

	// Validação usando a função nativa do Go standard library
	if !errors.Is(errWithInternal, rootCause) {
		t.Errorf("errors.Is() failed: expected to find rootCause inside Error500 structure")
	}

	unwrapped := errors.Unwrap(errWithInternal)
	if unwrapped != rootCause {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, rootCause)
	}

	// Caso 2: Sem erro interno (Deve retornar nil)
	errWithoutInternal := xerrors.NewError500("CTX-02", xerrors.XERR_OPERATION_FAILED, nil, "failed to query", "")

	if errors.Unwrap(errWithoutInternal) != nil {
		t.Errorf("Unwrap() without internal error should return nil, got %v", errors.Unwrap(errWithoutInternal))
	}
}

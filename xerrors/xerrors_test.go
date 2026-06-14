package xerrors_test

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/AeonDigital/Go-Core/xerrors"
)

// Define local error codes for test isolation
const (
	TestCtx  xerrors.ErrorCode = "ERR_TEST"
	TestCode xerrors.ErrorCode = "E9999"
)

func init() {
	// Register a core token specification to exercise the sync.Map fallback layers
	registry := map[xerrors.ErrorCode]xerrors.MetaMessage{
		TestCode: xerrors.NewMetaMessage("base internal test message", "rule_test", []string{"TAG1", "TAG2"}),
	}
	xerrors.RegisterDomainErrors(TestCtx, registry)
}

func TestXError_MessageResolution(t *testing.T) {
	// Scenario A: Explicit custom message provided at creation
	errCustom := xerrors.NewError500(TestCtx, TestCode, nil, "explicit custom message", "")
	if errCustom.Message() != "explicit custom message" {
		t.Errorf("expected explicit message, got: %s", errCustom.Message())
	}

	// Scenario B: No message provided, forcing fallback tracking through sync.Map registry
	// We use an internal framework key matching pattern (ERR_XERR:E9999) to simulate core fallback
	coreRegistry := map[xerrors.ErrorCode]xerrors.MetaMessage{
		TestCode: xerrors.NewMetaMessage("framework fallback message", "", nil),
	}
	xerrors.RegisterDomainErrors(xerrors.XERR_PKGCTX, coreRegistry)

	errFallback := xerrors.NewError500(TestCtx, TestCode, nil, "", "")
	if errFallback.Message() != "framework fallback message" {
		t.Errorf("expected registry fallback message, got: %s", errFallback.Message())
	}
}

func TestXError_ComponentLazyLoadingAndSkip(t *testing.T) {
	xerrors.EnableDebugMode()
	defer xerrors.DisableDebugMode()

	err := xerrors.NewError500(TestCtx, TestCode, nil, "msg", "")

	// 1. Initial execution component path should contain our current function name
	comp1 := err.Component()
	if !strings.Contains(comp1, "TestXError_ComponentLazyLoadingAndSkip") {
		t.Errorf("expected component to map current function, got: %s", comp1)
	}

	// 2. Second read must pull from internal cache directly (exercises mu.RLock branch)
	comp2 := err.Component()
	if comp1 != comp2 {
		t.Errorf("expected cached component to match, initial: %s, second: %s", comp1, comp2)
	}

	// 3. Testing withCallerSkip to re-evaluate frame depth up the stack lines
	errSkipped := err.WithCallerSkip(1)
	compSkipped := errSkipped.Component()
	if compSkipped == comp1 || compSkipped == "" {
		t.Errorf("expected updated frame depth evaluation, got: %s", compSkipped)
	}
}

func TestXError_DeepCopyAndImutability(t *testing.T) {
	errBase := xerrors.NewError500(TestCtx, TestCode, nil, "base", "").WithArgs("val1")

	// Branching off from the same base error concurrently
	errA := errBase.WithArgs("altered1")
	errB := errBase.WithArgs("altered2")

	// Ensure structural isolation (deep copy verified)
	if errA.Error() == errB.Error() {
		t.Error("expected instances to preserve deep state boundaries, but layouts match")
	}
}

func TestXError_FormatLayoutMatrix(t *testing.T) {
	rootErr := errors.New("native standard root cause")

	tests := []struct {
		name          string
		debug         bool
		isOperational bool
		ctx           xerrors.ErrorCode
		code          xerrors.ErrorCode
		err           error
		message       string
		info          string
		args          []any
		expectedParts []string
		notExpected   []string
	}{
		{
			name:          "Error400 (Client Side) - Simple Text Plain",
			debug:         false,
			isOperational: false,
			ctx:           xerrors.XERR_PKGCTX,
			code:          xerrors.XERR_NONE,
			message:       "explicit client text",
			info:          "",
			args:          nil,
			expectedParts: []string{"[CTX: ERR_XERR]", "[MSG: explicit client text]", ":: [ERR: ø]"},
			notExpected:   []string{"[COMPONENT:", "native standard root cause", "[INFO:"},
		},
		{
			name:          "Error400 (Client Side) - Missing Extra Tags with Registered Code",
			debug:         false,
			isOperational: false,
			ctx:           TestCtx,
			code:          TestCode,
			message:       "",
			info:          "",
			args:          nil, // Triggers empty set symbol 'ø' mapping loop
			expectedParts: []string{"[CTX: ERR_TEST]", "[ERR: E9999]", "[TAG1: ø]", "[TAG2: ø]"},
			notExpected:   []string{"[COMPONENT:"},
		},
		{
			name:          "Error500 (Operational) - With Debug and Root Cause",
			debug:         true,
			isOperational: true,
			ctx:           TestCtx,
			code:          TestCode,
			err:           rootErr,
			message:       "server failure",
			info:          "database connection timeout",
			args:          []any{"param1", "param2"},
			expectedParts: []string{"[CTX: ERR_TEST]", "[ERR: E9999]", "[COMPONENT: ", "[MSG: server failure]", "[TAG1: param1]", "[TAG2: param2]", "[INFO: database connection timeout]", ":: native standard root cause"},
			notExpected:   []string{":: [ERR: Wrapped]"},
		},
		{
			name:          "Error500 (Operational) - No Debug with Root Cause",
			debug:         false,
			isOperational: true,
			ctx:           TestCtx,
			code:          TestCode,
			err:           rootErr,
			message:       "server failure",
			info:          "",
			args:          nil,
			expectedParts: []string{"[CTX: ERR_TEST]", "[ERR: E9999]", "[MSG: server failure]", ":: [ERR: Wrapped]"},
			notExpected:   []string{"[COMPONENT:", "native standard root cause"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.debug {
				xerrors.EnableDebugMode()
			} else {
				xerrors.DisableDebugMode()
			}

			var subject error
			if tt.isOperational {
				// Tests Error500 pipeline with full telemetry inputs
				subject = xerrors.NewError500(tt.ctx, tt.code, tt.err, tt.message, tt.info).WithArgs(tt.args...)
			} else {
				// Tests Error400 polymorphic pipeline
				if tt.code == xerrors.XERR_NONE {
					subject = xerrors.NewError400(tt.message).WithArgs(tt.args...)
				} else {
					// Pass context tokens to trigger registry and extra tags loop branch
					subject = xerrors.NewError400(tt.ctx, tt.code).WithArgs(tt.args...)
				}
			}

			output := subject.Error()

			for _, part := range tt.expectedParts {
				if !strings.Contains(output, part) {
					t.Errorf("expected string output to contain [%s], got full string:\n%s", part, output)
				}
			}

			for _, nPart := range tt.notExpected {
				if strings.Contains(output, nPart) {
					t.Errorf("unexpected segment [%s] found in output layout:\n%s", nPart, output)
				}
			}
		})
	}
}

func TestXError_ExtremeConcurrencySafety(t *testing.T) {
	xerrors.EnableDebugMode()
	defer xerrors.DisableDebugMode()

	root := errors.New("shared baseline failure")
	errBase := xerrors.NewError500(TestCtx, TestCode, root, "concurrent-test", "info-payload")

	var wg sync.WaitGroup
	workers := 50
	iterations := 500

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// 1. Concurrent modification chaining (Fluent API verification)
				ForkedErr := errBase.WithArgs(id, j).WithCallerSkip(id % 2)

				// 2. Concurrent intensive reading across multiple cores
				_ = ForkedErr.Error()
				_ = ForkedErr.Component()
				_ = ForkedErr.Message()
			}
		}(i)
	}

	wg.Wait()
	// Test succeeds if no memory corruption or lock race conditions arise.
}

func TestXError_WithCallerSkip_NegativeBoundary(t *testing.T) {
	errBase := xerrors.NewError500(TestCtx, TestCode, nil, "boundary-test", "")

	// Passes a negative skip level to trigger the internal normalization logic (skip = 0)
	errNegative := errBase.WithCallerSkip(-5)

	// Evaluates if the object can still compute its function path without crashing
	comp := errNegative.Component()
	if comp == "" {
		t.Error("expected component to be calculated safely even with negative skip context")
	}
}

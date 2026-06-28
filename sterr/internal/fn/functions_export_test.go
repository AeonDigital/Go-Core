package fn

import (
	"runtime"
	"strings"
	"testing"
)

// TestTraceCallerLocationAllUnknowns forces every single defensive boundary
// inside TraceCallerLocation to execute and achieve 100% full coverage.
func TestTraceCallerLocationAllUnknowns(t *testing.T) {
	// 1. Save original states for clean defer restoration
	oldCaller := callerFunc
	oldFuncForPC := funcForPC
	oldLastIndex := lastIndexFunc
	defer func() {
		callerFunc = oldCaller
		funcForPC = oldFuncForPC
		lastIndexFunc = oldLastIndex
	}()

	// --- TRIGGER FIRST UNKNOWN (if !ok) ---
	callerFunc = func(skip int) (uintptr, string, int, bool) {
		return 0, "", 0, false
	}
	resFirst := TraceCallerLocation(0)
	if resFirst != "unknown::unknown" {
		t.Errorf("expected 'unknown::unknown' when caller ok is false, got '%s'", resFirst)
	}

	// --- TRIGGER SECOND UNKNOWN (if details == nil) ---
	callerFunc = func(skip int) (uintptr, string, int, bool) {
		return 12345, "file.go", 1, true
	}
	funcForPC = func(pc uintptr) *runtime.Func {
		return nil
	}
	resSecond := TraceCallerLocation(0)
	if resSecond != "unknown::unknown" {
		t.Errorf("expected 'unknown::unknown' when details are nil, got '%s'", resSecond)
	}

	// --- TRIGGER THIRD UNKNOWN (if lastDot == -1) ---
	// Restauramos o comportamento real do funcForPC para ele extrair um nome válido do próprio teste
	funcForPC = oldFuncForPC
	callerFunc = oldCaller

	// Forçamos o strings.LastIndex a fingir que não achou nenhum ponto na string
	lastIndexFunc = func(s, substr string) int {
		return -1
	}

	resThird := TraceCallerLocation(0)
	// Como o mock do LastIndex retornou -1, o código vai bater no return "unknown::" + fullName
	if !strings.HasPrefix(resThird, "unknown::") {
		t.Errorf("expected result to start with 'unknown::' when lastDot is -1, got '%s'", resThird)
	}
}

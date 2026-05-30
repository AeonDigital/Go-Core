package xerrors_test

import (
	"testing"

	"github.com/AeonDigital/Go-Core/xerrors"
)

func TestDebugModeFlow(t *testing.T) {
	// Ensure the initial state is disabled (false)
	if xerrors.GetDebugMode() != false {
		t.Errorf("initial GetDebugMode() = true; want false")
	}

	// Test enabling the debug mode
	xerrors.EnableDebugMode()
	if xerrors.GetDebugMode() != true {
		t.Errorf("after EnableDebugMode(), GetDebugMode() = false; want true")
	}

	// Test disabling the debug mode
	xerrors.DisableDebugMode()
	if xerrors.GetDebugMode() != false {
		t.Errorf("after DisableDebugMode(), GetDebugMode() = true; want false")
	}
}

func TestToggleDebugMode(t *testing.T) {
	// Reset to a known state to ensure test isolation
	xerrors.DisableDebugMode()

	// First toggle: false -> true
	xerrors.ToggleDebugMode()
	if xerrors.GetDebugMode() != true {
		t.Errorf("first ToggleDebugMode() = false; want true")
	}

	// Second toggle: true -> false
	xerrors.ToggleDebugMode()
	if xerrors.GetDebugMode() != false {
		t.Errorf("second ToggleDebugMode() = true; want false")
	}
}

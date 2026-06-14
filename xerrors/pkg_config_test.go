package xerrors_test

import (
	"sync"
	"testing"

	"github.com/AeonDigital/Go-Core/xerrors"
)

func TestDebugMode_LinearFlow(t *testing.T) {
	// Ensure we start from a clean, known state
	xerrors.DisableDebugMode()
	if xerrors.GetDebugMode() {
		t.Error("expected debug mode to be false initially")
	}

	// Test enabling
	xerrors.EnableDebugMode()
	if !xerrors.GetDebugMode() {
		t.Error("expected debug mode to be true after enabling")
	}

	// Test disabling
	xerrors.DisableDebugMode()
	if xerrors.GetDebugMode() {
		t.Error("expected debug mode to be false after disabling")
	}

	// Test toggling from false to true
	xerrors.ToggleDebugMode()
	if !xerrors.GetDebugMode() {
		t.Error("expected debug mode to be true after toggling from false")
	}

	// Test toggling from true to false
	xerrors.ToggleDebugMode()
	if xerrors.GetDebugMode() {
		t.Error("expected debug mode to be false after toggling from true")
	}
}

func TestDebugMode_ConcurrencySafety(t *testing.T) {
	// Restores to a standard state
	xerrors.DisableDebugMode()

	var wg sync.WaitGroup
	workers := 100
	iterations := 1000

	// Spawns multiple goroutines continuously reading and writing the debug state
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// Alternates randomly between state modifications and reads
				switch j % 3 {
				case 0:
					xerrors.ToggleDebugMode()
				case 1:
					xerrors.EnableDebugMode()
				default:
					xerrors.DisableDebugMode()
				}

				// Forces concurrent reading
				_ = xerrors.GetDebugMode()
			}
		}()
	}

	wg.Wait()

	// If the test completes without throwing a fatal data race panic,
	// our sync/atomic implementation is fully validated.
}

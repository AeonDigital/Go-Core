package xerrors

import (
	"sync/atomic"
)

var (
	// debugMode uses an atomic int32 (0 for false, 1 for true)
	// to prevent concurrent data races during runtime switches.
	debugMode int32 = 0
)

// GetDebugMode returns the current status of the debug mode.
// When true, errors provide technical details; when false, they are user-friendly.
func GetDebugMode() bool {
	return atomic.LoadInt32(&debugMode) == 1
}

// EnableDebugMode enables technical error details, making logs and outputs
// more comprehensive for debugging purposes.
func EnableDebugMode() {
	atomic.StoreInt32(&debugMode, 1)
}

// DisableDebugMode disables technical error details, switching outputs
// to a user-friendly format suitable for end-users.
func DisableDebugMode() {
	atomic.StoreInt32(&debugMode, 0)
}

// ToggleDebugMode switches the current state of the debug mode
// (enables it if disabled, or disables it if enabled).
func ToggleDebugMode() {
	for {
		current := atomic.LoadInt32(&debugMode)
		var next int32 = 0
		if current == 0 {
			next = 1
		}

		// Performs a Compare-And-Swap loop to guarantee safe state alteration
		if atomic.CompareAndSwapInt32(&debugMode, current, next) {
			break
		}
	}
}

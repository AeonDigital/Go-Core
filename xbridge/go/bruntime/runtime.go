package bruntime

import "runtime"

// IRuntimeBridge defines the unified contract to intercept and mock the Go standard library 'runtime' package.
// This interface decouples system runtime properties and diagnostics from the application logic.
type IRuntimeBridge interface {
	// =========================================================================
	// PLATFORM IDENTIFICATION
	// =========================================================================

	// GOOS is the running program's operating system target:
	// one of darwin, freebsd, linux, and so on.
	// To view possible combinations of GOOS and GOARCH, run "go tool dist list".
	GOOS() string
	// GOARCH is the running program's architecture target:
	// one of 386, amd64, arm, s390x, and so on.
	GOARCH() string
	// NumCPU returns the number of logical CPUs usable by the current process.
	//
	// The set of available CPUs is checked by querying the operating system
	// at process startup. Changes to operating system CPU allocation after
	// process startup are not reflected.
	NumCPU() int
	// Version returns the Go tree's version string.
	// It is either the commit hash and date at the time of the build or,
	// when possible, a release tag like "go1.3".
	Version() string

	// =========================================================================
	// MEMORY, TELEMETRY & MONITORING
	// =========================================================================

	// NumGoroutine returns the number of goroutines that currently exist.
	NumGoroutine() int
	// ReadMemStats populates m with memory allocator statistics.
	//
	// The returned memory allocator statistics are up to date as of the
	// call to ReadMemStats. This is in contrast with a heap profile,
	// which is a snapshot as of the most recently completed garbage
	// collection cycle.
	ReadMemStats(m *runtime.MemStats)
	// GC runs a garbage collection and blocks the caller until the
	// garbage collection is complete. It may also block the entire
	// program.
	GC()
	// Stack formats a stack trace of the calling goroutine into buf
	// and returns the number of bytes written to buf.
	// If all is true, Stack formats stack traces of all other goroutines
	// into buf after the trace for the current goroutine.
	Stack(buf []byte, all bool) int
	// Caller reports file and line number information about function invocations on
	// the calling goroutine's stack. The argument skip is the number of stack frames
	// to ascend, with 0 identifying the caller of Caller. (For historical reasons the
	// meaning of skip differs between Caller and [Callers].) The return values report
	// the program counter, the file name (using forward slashes as path separator, even
	// on Windows), and the line number within the file of the corresponding call.
	// The boolean ok is false if it was not possible to recover the information.
	Caller(skip int) (pc uintptr, file string, line int, ok bool)
}

// sruntimeBridge serves as the private concrete production implementation of IRuntimeBridge,
// routing calls natively to the standard library 'runtime' package.
type sruntimeBridge struct{}

// Ensure at compile time that the private struct implements the public interface.
var _ IRuntimeBridge = sruntimeBridge{}

// NewRuntimeBridge creates and returns a public, mockable IRuntimeBridge instance pointing
// to the production system runtime diagnostics implementation.
func NewRuntime() IRuntimeBridge {
	return sruntimeBridge{}
}

// =========================================================================
// IMPLEMENTATION OF PLATFORM IDENTIFICATION
// =========================================================================

func (sruntimeBridge) GOOS() string    { return runtime.GOOS }
func (sruntimeBridge) GOARCH() string  { return runtime.GOARCH }
func (sruntimeBridge) NumCPU() int     { return runtime.NumCPU() }
func (sruntimeBridge) Version() string { return runtime.Version() }

// =========================================================================
// IMPLEMENTATION OF MEMORY, TELEMETRY & MONITORING
// =========================================================================

func (sruntimeBridge) NumGoroutine() int                { return runtime.NumGoroutine() }
func (sruntimeBridge) ReadMemStats(m *runtime.MemStats) { runtime.ReadMemStats(m) }
func (sruntimeBridge) GC()                              { runtime.GC() }
func (sruntimeBridge) Stack(buf []byte, all bool) int   { return runtime.Stack(buf, all) }
func (sruntimeBridge) Caller(skip int) (uintptr, string, int, bool) {
	return runtime.Caller(skip)
}

// =========================================================================
// GLOBAL PUBLIC INTERFACE INSTANCE
// =========================================================================

// RuntimeBridge is the global public variable used by all core application code
// to query platform architecture, memory metrics, and runtime diagnostics.
// It can be easily hot-swapped at unit test boundaries via generated pkgxmock utilities.
var Runtime IRuntimeBridge = NewRuntime()

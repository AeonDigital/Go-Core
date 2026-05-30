package xerrors

import (
	"fmt"
	"runtime"
)

// Error provides a concrete structural implementation of the DetailedError interface.
// It encapsulates contextual metadata alongside standard system exceptions to maximize
// observability and reduce cognitive friction during troubleshooting pipelines.
type Error struct {
	// code maps the failure to a domain-specific category (e.g., ErrUnknown, ErrDatabase).
	code ErrorCode
	// component lazily holds the fully qualified runtime execution function path string.
	component string
	// err preserves the raw, low-level error (root cause) for wrapping semantics.
	err error
	// callerSkip maintains the runtime execution frame index position for call stack targeting.
	callerSkip int

	// message holds a high-level, human-readable operational summary or a sanitized feedback description.
	message string
	// data provides a space for the engineer to attach raw debugging context chunks or raw payloads.
	data string
}

// NewError initializes a rich DetailedError instance. It pre-configures a default call stack
// frame skip level but postpones the actual heavy runtime inspection until the component property
// is explicitly requested by a logger or console channel.
func NewError(
	code ErrorCode,
	err error,
	message string,
	data string,
) DetailedError {
	return &Error{
		code:       code,
		component:  "",
		err:        err,
		callerSkip: 0,

		message: message,
		data:    data,
	}
}

// Code retrieves the structured domain classification assigned to this error.
func (e *Error) Code() ErrorCode { return e.code }

// Component returns the dynamically reflected runtime namespace where the failure originated.
// It leverages lazy-loading evaluation, executing stack inspection only upon the first request.
func (e *Error) Component() string {
	if e.component == "" {
		e.calculateComponent(3)
	}
	return e.component
}

// calculateComponent captures the execution stack trace at the specified frame depth
// and resolves it into a fully qualified function namespace string (e.g., "package.Function").
// This value is directly assigned to the internal component tracking state.
//
// Parameters:
//   - skip: The exact integer frame depth offset passed to runtime.Callers.
//     A value of 0 targets runtime.Callers itself, while subsequent values climb
//     higher up the active execution stack lines to pinpoint callers outside this library.
func (e *Error) calculateComponent(skip int) {
	progCounter := make([]uintptr, 1)
	runtime.Callers(skip, progCounter)
	frames := runtime.CallersFrames(progCounter)
	frame, _ := frames.Next()

	e.component = frame.Function
}

// WithCallerSkip allows dynamically adjusting the runtime stack frame skip level.
// This resetting behavior automatically clears any previously loaded component metadata state,
// forcing the runtime analyzer to re-evaluate call origins on subsequent property calls.
func (e *Error) WithCallerSkip(skip int) DetailedError {
	if skip < 0 {
		skip = 0
	}
	e.calculateComponent(3 + skip)
	return e
}

// Error satisfies the standard Go library errors handling interface contract specification.
// It builds a summarized formatted diagnostic feedback sequence.
func (e *Error) Error() string {
	// Replaced field access with the method call to safely trigger lazy-loading
	comp := e.Component()
	if e.err != nil {
		return fmt.Sprintf("[%s] %s: %v", comp, e.message, e.err)
	}
	return fmt.Sprintf("[%s] %s", comp, e.message)
}

// Message extracts the human-readable summary detailing the specific breakdown event.
func (e *Error) Message() string { return e.message }

// Data fetches the secondary raw context strings attached to this structural payload container.
func (e *Error) Data() string { return e.data }

// DebugError returns a multi-line diagnostic report containing comprehensive internal tracking metadata.
func (e *Error) DebugError() string {
	strErr := "" +
		"--- DETAILED DEBUG ERROR ---\n" +
		"Code: %s\n" +
		"Component: %s\n" +
		"Internal Error: %v\n" +
		"Message: %s\n" +
		"Data:\n%s\n" +
		"----------------------------"

	// Triggering the Component() method instead of field data prevents empty token prints
	return fmt.Sprintf(
		strErr,
		e.code,
		e.Component(),
		e.err,
		e.message,
		e.data,
	)
}

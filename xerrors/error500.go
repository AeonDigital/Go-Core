package xerrors

import (
	"fmt"
	"runtime"
)

// IError500 standardizes the functional execution capabilities required to manage rich,
// decoupled diagnostic payloads across network borders, microservices, and system boundaries.
type IError500 interface {
	// CTX extracts the structural context identifier or tracking boundary metadata
	// associated with the entry point of this operational flow.
	CTX() ErrorCode

	// Code returns the categorical string constant domain mapping for this failure event.
	Code() ErrorCode

	// Component pinpoints the specific architectural layer or tracking package path.
	Component() string

	// error ensures native alignment with Go standard library errors handling semantics.
	error

	// WithCallerSkip allows dynamically adjusting the runtime stack frame skip level to recompute
	// the component identifier. This is essential when NewError is wrapped inside custom factory
	// utilities or logging helper blocks, ensuring the framework pinpoints the true origin of the failure.
	// It returns the updated IError500 to support fluent API call chaining.
	WithCallerSkip(skip int) IError500

	// Message returns the readable operational summary describing what went wrong.
	Message() string

	// Info returns information for debugging the context.
	Info() string
}

// Error500 provides a concrete structural implementation of the IError500 interface.
// It encapsulates contextual metadata alongside standard system exceptions to maximize
// observability and reduce cognitive friction during troubleshooting pipelines.
type Error500 struct {
	// errCTX defines the domain-specific or structural operational context identifier,
	// typically utilized to trace error execution boundaries across application flows.
	errCTX ErrorCode
	// code maps the failure to a domain-specific category
	errCode ErrorCode
	// component lazily holds the fully qualified runtime execution function path string.
	component string
	// err preserves the raw, low-level error (root cause) for wrapping semantics.
	err error

	// message holds a high-level, human-readable operational summary or a sanitized feedback description.
	message string
	// info provides a space for the engineer to attach raw debugging context information.
	info string
}

// NewError500 initializes a rich IError500 instance. It pre-configures a default call stack
// frame skip level but postpones the actual heavy runtime inspection until the component property
// is explicitly requested by a logger or console channel.
func NewError500(
	errCTX ErrorCode,
	errCode ErrorCode,
	err error,
	message string,
	info string,
) IError500 {
	return &Error500{
		errCTX:    errCTX,
		errCode:   errCode,
		component: "",
		err:       err,

		message: message,
		info:    info,
	}
}

// CTX retrieves the isolated contextual tracking identity assigned to this error payload.
func (e *Error500) CTX() ErrorCode {
	return e.errCTX
}

// Code retrieves the structured domain classification assigned to this error.
func (e *Error500) Code() ErrorCode {
	return e.errCode
}

// Component returns the dynamically reflected runtime namespace where the failure originated.
// It leverages lazy-loading evaluation, executing stack inspection only upon the first request.
func (e *Error500) Component() string {
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
func (e *Error500) calculateComponent(skip int) {
	progCounter := make([]uintptr, 1)
	runtime.Callers(skip, progCounter)
	frames := runtime.CallersFrames(progCounter)
	frame, _ := frames.Next()

	e.component = frame.Function
}

// WithCallerSkip allows dynamically adjusting the runtime stack frame skip level.
// It creates a new instance copy to prevent side effects on the original error.
func (e *Error500) WithCallerSkip(skip int) IError500 {
	if skip < 0 {
		skip = 0
	}

	// Create a shallow copy of the struct to ensure immutability
	clone := *e

	// Reset the component so it forces a recalculation on the next lazy-load read
	clone.component = ""

	// Calculate the component based on the new frame depth context
	clone.calculateComponent(3 + skip)

	return &clone
}

// Message extracts the human-readable summary detailing the specific breakdown event.
func (e *Error500) Message() string {
	if e.message == "" {
		coreKey := ErrorCode(string(XERR_PKGCTX) + ":" + string(e.errCode))
		metaMsg, exists := xerrorMapRegistry[coreKey]
		if exists {
			return metaMsg.message
		}
	}

	return e.message
}

// Info fetches the secondary raw context strings attached to this structural payload container.
func (e *Error500) Info() string {
	return e.info
}

// Error satisfies the standard Go library error handling interface contract specification.
// It dynamically alters its layout density based on the global debugMode state configuration.
func (e *Error500) Error() string {
	strError := ""
	strComponent := ""
	if debugMode {
		strComponent = "[COMPONENT: " + e.Component() + "]"

		if e.err != nil {
			strError = fmt.Sprintf(" :: %v", e.err)
		}
	}

	strData := ""
	if e.info != "" {
		strData = "[INFO: " + e.info + "]"
	}

	return fmt.Sprintf(
		"[CTX: %s][ERR: %s]%s[MSG: %s]%s%s",
		e.errCTX,
		e.errCode,
		strComponent,
		e.Message(),
		strData,
		strError,
	)
}

// Unwrap returns the underlying root cause error, enabling seamless integration
// with standard library features like errors.Is and errors.As.
func (e *Error500) Unwrap() error {
	return e.err
}

package xerrors

// IError500 standardizes server-side, operational diagnostic payloads.
type IError500 interface {
	// CTX extracts the entry point operational flow boundary metadata.
	CTX() ErrorCode

	// Code returns the categorical failure domain classification.
	Code() ErrorCode

	// Component pinpoints the reflective architectural package or execution function path.
	Component() string

	// WithCallerSkip dynamically shifts the stack frame runtime collection depth.
	WithCallerSkip(skip int) IError500

	// WithArgs appends dynamic contextual payloads to map sequentially into metadata extraTags.
	WithArgs(args ...any) IError500

	// Message extracts the human-readable summary detailing the specific failure.
	Message() string

	// Info returns secondary raw debugging contextual payloads.
	Info() string

	// error ensures native integration with Go standard library error semantics.
	error
}

// error500Adapter wraps the internal private engine to satisfy the public IError500 contract.
type error500Adapter struct {
	*xError
}

// NewError500 initializes an operational, stack-traced IError500 instance.
// It leverages lazy-loading for call stack evaluation to remain highly performant.
func NewError500(
	errCTX ErrorCode,
	errCode ErrorCode,
	err error,
	message string,
	info string,
) IError500 {
	return &error500Adapter{
		xError: &xError{
			contextCode:   errCTX,
			errorCode:     errCode,
			underlyingErr: err,
			message:       message,
			info:          info,
			isOperational: true, // Forces system telemetry behavior
		},
	}
}

// CTX retrieves the assigned context tracking identifier.
func (a *error500Adapter) CTX() ErrorCode {
	return a.contextCode
}

// Code retrieves the structured domain classification code.
func (a *error500Adapter) Code() ErrorCode {
	return a.errorCode
}

// Component returns the reflected namespace where the failure originated.
func (a *error500Adapter) Component() string {
	return a.getComponent()
}

// WithCallerSkip adjusts stack frame collection, returning an isolated clone.
func (a *error500Adapter) WithCallerSkip(skip int) IError500 {
	// Calls internal deep-copy mechanism to shield concurrent operations
	clonedEngine := a.withCallerSkip(skip)
	return &error500Adapter{xError: clonedEngine}
}

// WithArgs injects dynamic context tracking variables into the instance safely.
func (a *error500Adapter) WithArgs(args ...any) IError500 {
	if len(args) == 0 {
		return a
	}

	a.mu.RLock()
	// Explicitly copy only data fields to prevent copying the sync.RWMutex value
	clonedEngine := &xError{
		contextCode:   a.contextCode,
		errorCode:     a.errorCode,
		underlyingErr: a.underlyingErr,
		message:       a.message,
		info:          a.info,
		isOperational: a.isOperational,
	}
	a.mu.RUnlock()

	// Deep copy arguments payload to guarantee concurrency safety
	clonedEngine.arguments = make([]any, len(args))
	copy(clonedEngine.arguments, args)

	return &error500Adapter{xError: clonedEngine}
}

// Message extracts the user-friendly operational summary description string.
func (a *error500Adapter) Message() string {
	return a.resolveMessage()
}

// Info fetches unstructured raw secondary metadata context.
func (a *error500Adapter) Info() string {
	return a.info
}

// Error evaluates the final formatting output via strings.Builder.
func (a *error500Adapter) Error() string {
	return a.format()
}

// Unwrap exposes the underlying exception boundary to support errors.Is / errors.As.
func (a *error500Adapter) Unwrap() error {
	return a.underlyingErr
}

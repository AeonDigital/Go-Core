package xerrors

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

// xError serves as the unified, high-performance internal engine for the package.
// It encapsulates all metadata, data tags, and stack-tracing behaviors.
type xError struct {
	// mu protects lazy-loaded fields from data races (specifically the component string)
	mu sync.RWMutex

	// contextCode maps the overall operational flow boundary (e.g., "ERR_XERR")
	contextCode ErrorCode

	// errorCode maps the specific failure classification (e.g., "E1001")
	errorCode ErrorCode

	// component lazily holds the fully qualified runtime execution function path string
	component string

	// underlyingErr preserves the low-level root cause error for wrapping semantics
	underlyingErr error

	// message holds a high-level, human-readable operational summary or dynamic description
	message string

	// info provides an unstructured space for secondary raw debugging context
	info string

	// arguments stores raw dynamic payloads passed for layout building and log tracking
	arguments []any

	// isOperational flags if this instance must compute stack frame tracking (Error500 behavior)
	isOperational bool
}

// resolveMessage safely evaluates and returns the explicit message or fallbacks to the global registry mapping.
func (e *xError) resolveMessage() string {
	if e.message != "" {
		return e.message
	}

	// Fallback lookup using the unified framework pattern (PKGCTX:Code)
	lookupKey := ErrorCode(string(XERR_PKGCTX) + ":" + string(e.errorCode))

	// Safe lock-free read from sync.Map with explicit type assertion
	if value, exists := xerrorMapRegistry.Load(lookupKey); exists {
		if metaMsg, ok := value.(MetaMessage); ok {
			return metaMsg.message
		}
	}

	return ""
}

// getComponent retrieves the lazy-loaded execution path under a read/write mutex boundary.
func (e *xError) getComponent() string {
	e.mu.RLock()
	if e.component != "" {
		defer e.mu.RUnlock()
		return e.component
	}
	e.mu.RUnlock()

	// Escalates to write lock to safely perform the lazy-load calculation
	e.mu.Lock()
	defer e.mu.Unlock()

	// Double-check flag after acquiring lock
	if e.component == "" {
		e.calculateComponent(4) // Adjusted depth frame offset to bypass internal engine calls
	}
	return e.component
}

// calculateComponent inspects the execution stack trace at the specific frame depth.
// Internal execution caller must hold the structural write lock before invoking this method.
func (e *xError) calculateComponent(skip int) {
	progCounter := make([]uintptr, 1)
	runtime.Callers(skip, progCounter)
	frames := runtime.CallersFrames(progCounter)
	frame, _ := frames.Next()

	e.component = frame.Function
}

// withCallerSkip produces a deep-copied clone of the error to isolate concurrent state modifications.
func (e *xError) withCallerSkip(skip int) *xError {
	if skip < 0 {
		skip = 0
	}

	e.mu.RLock()
	// Explicitly copy only data fields to prevent copying the sync.RWMutex value
	clone := &xError{
		contextCode:   e.contextCode,
		errorCode:     e.errorCode,
		underlyingErr: e.underlyingErr,
		message:       e.message,
		info:          e.info,
		isOperational: e.isOperational,
	}
	e.mu.RUnlock()

	// Perform a deep copy of the arguments slice
	if len(e.arguments) > 0 {
		clone.arguments = make([]any, len(e.arguments))
		copy(clone.arguments, e.arguments)
	}

	clone.mu.Lock()
	clone.calculateComponent(4 + skip)
	clone.mu.Unlock()

	return clone
}

// format dynamic constructs the final visual layout string using a highly efficient pre-allocated strings.Builder.
func (e *xError) format() string {
	var builder strings.Builder

	// Step 1: Base Context and Error boundaries
	builder.WriteString("[CTX: ")
	builder.WriteString(string(e.contextCode))
	builder.WriteString("][ERR: ")
	builder.WriteString(string(e.errorCode))
	builder.WriteString("]")

	// Step 2: System Telemetry Integration (Operational/Error500 Capability)
	if e.isOperational && GetDebugMode() {
		builder.WriteString("[COMPONENT: ")
		builder.WriteString(e.getComponent())
		builder.WriteString("]")
	}

	// Step 3: Message Parsing
	builder.WriteString("[MSG: ")
	builder.WriteString(e.resolveMessage())
	builder.WriteString("]")

	// Step 4: Metadata Extra Tags Processing
	lookupKey := ErrorCode(string(e.contextCode) + ":" + string(e.errorCode))

	// Safe lock-free read from sync.Map for specific domain meta tags
	if value, exists := xerrorMapRegistry.Load(lookupKey); exists {
		if meta, ok := value.(MetaMessage); ok && len(meta.extraTags) > 0 {
			for i, tagName := range meta.extraTags {
				builder.WriteString("[")
				builder.WriteString(tagName)
				builder.WriteString(": ")

				// Map arguments sequentially based on tags presence or fallback safely to empty set symbol
				if i < len(e.arguments) && e.arguments[i] != nil {
					builder.WriteString(fmt.Sprintf("%v", e.arguments[i]))
				} else {
					builder.WriteString("ø")
				}
				builder.WriteString("]")
			}
		}
	}

	// Step 5: Unstructured Info Dump
	if e.info != "" {
		builder.WriteString("[INFO: ")
		builder.WriteString(e.info)
		builder.WriteString("]")
	}

	// Step 6: Core System Exception Wrapping (Debug Visibility boundary)
	if e.underlyingErr != nil {
		if e.isOperational && GetDebugMode() {
			builder.WriteString(fmt.Sprintf(" :: %v", e.underlyingErr))
		} else {
			builder.WriteString(" :: [ERR: Wrapped]")
		}
	} else {
		builder.WriteString(" :: [ERR: ø]")
	}

	return builder.String()
}

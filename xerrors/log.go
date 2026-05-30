package xerrors

import (
	"context"
	"log/slog"
	"runtime"
	"strings"
)

// LogOption defines the functional option signature used to configure
// the behavior of structured error logging.
type LogOption func(*logConfig)

// logConfig holds internal configuration options for runtime evaluation.
type logConfig struct {
	callerSkip int
}

// WithLogCallerSkip dynamically adjusts the runtime.Caller frame interception level
// during raw error fallbacks. Pass a positive integer (e.g., 1) if this function
// is wrapped inside another downstream corporate logging utility.
func WithLogCallerSkip(skip int) LogOption {
	return func(c *logConfig) {
		c.callerSkip = skip
	}
}

// LogErr captures a detailed breakdown of the incoming error instance, converts its internal
// components into key-value structural attributes, merges external slog metadata tokens,
// and publishes a structured error stream via the standard library's slog ecosystem.
//
// Unification strategy:
//   - If 'targetErr' satisfies xerrors.DetailedError, it skips raw reflection overhead
//     and extracts the pre-recorded component runtime frame path, code, and debugging data dumps.
//   - If 'targetErr' is a standard Go generic error, it evaluates runtime.Caller backtraces
//     as a fallback parsing layout sequence to locate the failure context.
func Log(
	ctx context.Context,
	targetErr error,
	customAttrs []slog.Attr,
	opts ...LogOption,
) {
	if ctx == nil {
		ctx = context.Background()
	}

	config := logConfig{callerSkip: 0}
	for _, opt := range opts {
		opt(&config)
	}

	var fullName string
	var errorCode string
	var errorData string

	// Step 1: Check if the error is our rich internal DetailedError implementation
	if detailed, ok := targetErr.(DetailedError); ok {
		fullName = detailed.Component()
		errorCode = string(detailed.Code())
		errorData = detailed.Data()
	} else {
		// Step 2: Fallback to raw runtime.Caller tracing if it is just a standard Go error
		pc, _, _, ok := runtime.Caller(1 + config.callerSkip)
		if ok {
			fn := runtime.FuncForPC(pc)
			if fn != nil {
				fullName = fn.Name()
			}
		}
	}

	// Step 3: Parse the functional name into structured semantic boundaries
	pkgName, objectName, methodName := parseRuntimeFuncName(fullName)
	if objectName == "" {
		objectName = "Function"
	}

	// Step 4: Assemble base logging attributes
	attrs := []slog.Attr{
		slog.String("package", pkgName),
		slog.String("object", objectName),
		slog.String("method", methodName),
	}

	// Merge custom developer arguments passed down the pipe
	attrs = append(attrs, customAttrs...)

	if errorCode != "" {
		attrs = append(attrs, slog.String("error_code", errorCode))
	}
	if errorData != "" {
		attrs = append(attrs, slog.String("error_data", errorData))
	}

	if targetErr != nil {
		attrs = append(attrs, slog.String("error", targetErr.Error()))
	}

	// Dispatch structured execution record log event
	slog.LogAttrs(ctx, slog.LevelError, "operation failure detected", attrs...)
}

// parseRuntimeFuncName isolates packaging namespaces from dynamic execution literals.
func parseRuntimeFuncName(fullName string) (string, string, string) {
	var pkgName, objectName, methodName string
	if fullName == "" {
		return pkgName, objectName, methodName
	}

	if lastDot := strings.LastIndex(fullName, "."); lastDot != -1 {
		methodName = fullName[lastDot+1:]
		fullPkgAndObj := fullName[:lastDot]

		if structIdx := strings.Index(fullPkgAndObj, "(*"); structIdx != -1 {
			objectName = fullPkgAndObj[structIdx+2:]
			objectName = strings.ReplaceAll(objectName, ")", "")
			if genIdx := strings.Index(objectName, "["); genIdx != -1 {
				objectName = objectName[:genIdx]
			}
			fullPkgAndObj = fullPkgAndObj[:structIdx]
		}

		fullPkgAndObj = strings.TrimSuffix(fullPkgAndObj, ".")
		if lastSlash := strings.LastIndex(fullPkgAndObj, "/"); lastSlash != -1 {
			pkgName = fullPkgAndObj[lastSlash+1:]
		} else {
			pkgName = fullPkgAndObj
		}
	}

	return pkgName, objectName, methodName
}

package xerrors

import (
	"fmt"
)

// IError400 standardizes validation and client-side failures, allowing
// transport layers to seamlessly extract error codes without string parsing.
type IError400 interface {
	// Code returns the categorical failure domain classification.
	Code() ErrorCode

	// WithArgs appends dynamic contextual payloads to map sequentially into metadata extraTags.
	WithArgs(args ...any) IError400

	// error ensures native alignment with Go standard library error handling semantics.
	error
}

// error400Adapter wraps the internal private engine to satisfy the public IError400 contract.
type error400Adapter struct {
	*xError
}

// NewError400 acts as a highly flexible, polymorphic factory that creates and returns an IError400 instance.
// It preserves low cognitive load for the developer while extracting control tokens from raw telemetry tags.
func NewError400(args ...any) IError400 {
	if len(args) == 0 {
		return &error400Adapter{
			xError: &xError{
				contextCode:   XERR_PKGCTX,
				errorCode:     XERR_NONE,
				isOperational: false,
			},
		}
	}

	var finalCtx ErrorCode = XERR_PKGCTX
	var finalCode ErrorCode = XERR_NONE
	var processedText string
	var cleanedArgs []any

	// Case 1: Evaluate if the first two arguments match an extended domain signature (PKGCTX + Code)
	if len(args) >= 2 {
		pkgCtx, ok1 := args[0].(ErrorCode)
		code, ok2 := args[1].(ErrorCode)

		if ok1 && ok2 {
			composedKey := string(pkgCtx) + ":" + string(code)

			if _, exists := xerrorMapStringToErrorCode.Load(composedKey); exists {
				finalCtx = pkgCtx
				finalCode = code // FIX: Store the isolated token code to match formatting requirements

				// Clean slice: extract control tokens and retain only raw data fields
				if len(args) > 2 {
					cleanedArgs = make([]any, len(args)-2)
					copy(cleanedArgs, args[2:])
				}

				return &error400Adapter{
					xError: &xError{
						contextCode:   finalCtx,
						errorCode:     finalCode,
						arguments:     cleanedArgs,
						isOperational: false,
					},
				}
			}
		}
	}

	// Case 2: Evaluate if only the first argument matches a core framework signature
	if firstCode, ok := args[0].(ErrorCode); ok {
		coreKey := string(XERR_PKGCTX) + ":" + string(firstCode)

		if _, exists := xerrorMapStringToErrorCode.Load(coreKey); exists {
			finalCode = firstCode // FIX: Store the isolated core code token

			// Clean slice: drop the core token and keep subsequent metadata payloads
			if len(args) > 1 {
				cleanedArgs = make([]any, len(args)-1)
				copy(cleanedArgs, args[1:])
			}

			return &error400Adapter{
				xError: &xError{
					contextCode:   XERR_PKGCTX,
					errorCode:     finalCode,
					arguments:     cleanedArgs,
					isOperational: false,
				},
			}
		}
	}

	// Case 3: Fallback to standard string formatting or raw text parsing
	if firstStr, ok := args[0].(string); ok {
		if len(args) > 1 {
			processedText = fmt.Sprintf(firstStr, args[1:]...)
		} else {
			processedText = firstStr
		}
	} else {
		// Safeguard mechanism for unexpected or unmapped raw parameter types
		processedText = fmt.Sprintf("%v", args)
	}

	return &error400Adapter{
		xError: &xError{
			contextCode:   XERR_PKGCTX,
			errorCode:     XERR_NONE,
			message:       processedText,
			isOperational: false,
		},
	}
}

// Code retrieves the structured domain classification code.
func (a *error400Adapter) Code() ErrorCode {
	return a.errorCode
}

// WithArgs injects dynamic validation tracking variables into the instance safely.
func (a *error400Adapter) WithArgs(args ...any) IError400 {
	if len(args) == 0 {
		return a
	}

	a.mu.RLock()
	clonedEngine := &xError{
		contextCode:   a.contextCode,
		errorCode:     a.errorCode,
		underlyingErr: a.underlyingErr,
		message:       a.message,
		info:          a.info,
		isOperational: a.isOperational,
	}
	a.mu.RUnlock()

	clonedEngine.arguments = make([]any, len(args))
	copy(clonedEngine.arguments, args)

	return &error400Adapter{xError: clonedEngine}
}

// Error evaluates the final layout output string via strings.Builder.
func (a *error400Adapter) Error() string {
	return a.format()
}

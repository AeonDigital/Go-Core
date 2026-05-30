package xerrors

import (
	"fmt"
)

// IError400 standardizes validation and client-side failures, allowing
// transport layers to seamlessly extract error codes without string parsing.
type IError400 interface {
	// Code returns the categorical string constant domain mapping for this validation event.
	Code() ErrorCode

	// error ensures native alignment with Go standard library error handling semantics.
	error
}

// error400 provides a high-performance, lightweight implementation of IError400.
// It bypasses runtime stack inspection, making it ideal for high-throughput user input validations.
type error400 struct {
	// code maps the failure to a domain-specific category
	errCode ErrorCode

	// errMSG stores the fully compiled, human-readable validation summary or formatted output.
	errMSG string
}

// NewError400 acts as a highly flexible, polymorphic factory that creates and returns an IError400 instance.
// It abstracts away method overloading limitations in Go by inspecting the type sequence of its variadic arguments.
//
// The evaluation hierarchy operates as follows:
//  1. Composed Domain Token: If the first two arguments are identified as ErrorCode types, they are unified
//     into a namespaced key (e.g., "PKGCTX:E1001") to match domain-extended registries.
//  2. Core Package Token: If only the first argument is an ErrorCode, the factory falls back to the core
//     infrastructure namespace ("ERR_XERR:E1001").
//  3. Plain Text/Standard Format: If no initial ErrorCode tokens are detected, it evaluates the arguments
//     as a standard fmt.Sprintf format string or raw message, defaulting the error code tracking state to XERR_NONE.
func NewError400(args ...any) IError400 {
	if len(args) == 0 {
		return &error400{
			errCode: XERR_NONE,
			errMSG:  "",
		}
	}

	var finalCode ErrorCode = XERR_NONE
	var processedText string

	// Step 1: Evaluate if the first two arguments match an extended domain signature (PKGCTX + Code)
	if len(args) >= 2 {
		pkgCtx, ok1 := args[0].(ErrorCode)
		code, ok2 := args[1].(ErrorCode)

		if ok1 && ok2 {
			composedKey := string(pkgCtx) + ":" + string(code)
			registeredCode, exists := xerrorMapStringToErrorCode[composedKey]
			if exists {
				finalCode = registeredCode

				// Consume the two token parameters and forward the rest to the layout builder
				formatMsg, maskArgs := buildMask(xerrorMapRegistry[registeredCode], args[2:])
				processedText = fmt.Sprintf(formatMsg, maskArgs...)

				return &error400{
					errCode: finalCode,
					errMSG:  processedText,
				}
			}
		}
	}

	// Step 2: Evaluate if only the first argument matches a core framework signature
	firstCode, ok := args[0].(ErrorCode)
	if ok {
		coreKey := string(XERR_PKGCTX) + ":" + string(firstCode)
		registeredCode, exists := xerrorMapStringToErrorCode[coreKey]
		if exists {
			finalCode = registeredCode

			// Consume the single core token and forward the rest to the layout builder
			formatMsg, maskArgs := buildMask(xerrorMapRegistry[registeredCode], args[1:])
			processedText = fmt.Sprintf(formatMsg, maskArgs...)

			return &error400{
				errCode: finalCode,
				errMSG:  processedText,
			}
		}
	}

	// Step 3: Fallback to standard string formatting or raw text parsing
	firstStr, ok := args[0].(string)
	if ok {
		if len(args) > 1 {
			processedText = fmt.Sprintf(firstStr, args[1:]...)
		} else {
			processedText = firstStr
		}
	} else {
		// Safeguard mechanism for unexpected or unmapped raw parameter types
		processedText = fmt.Sprintf("%v", args)
	}

	return &error400{
		errCode: XERR_NONE,
		errMSG:  processedText,
	}
}

// Code retrieves the structured domain classification assigned to this validation error.
func (e *error400) Code() ErrorCode {
	return e.errCode
}

// Error satisfies the standard Go library error handling interface contract specification.
func (e *error400) Error() string {
	return e.errMSG
}

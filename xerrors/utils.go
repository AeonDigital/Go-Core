package xerrors

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// errorRegistry centralizes all specialized error metadata blocks and behavioral definitions
// within the xerrors ecosystem, eliminating verbose switch-case statements and reducing maintenance overhead.
var errorRegistry = map[string]errorMetadata{
	// ============================================================================
	// GROUP 1: PRESENCE, EXISTENCE, AND NULLITY VALIDATIONS
	// ============================================================================
	XERR_FIELD_REQUIRED: {
		defaultMessage: "field required",
		extraTags:      []string{"FIELD"},
	},
	XERR_NIL_NOT_ALLOWED: {
		defaultMessage: "nil pointer value not allowed",
		extraTags:      []string{"FIELD"},
	},
	XERR_EMPTY_NOT_ALLOWED: {
		defaultMessage: "empty string value not allowed",
		extraTags:      []string{"FIELD"},
	},
	XERR_ZERO_NOT_ALLOWED: {
		defaultMessage: "zero numeric value not allowed",
		extraTags:      []string{"FIELD"},
	},
	XERR_ALREADY_EXISTS: {
		defaultMessage: "duplicate restriction violated",
		extraTags:      []string{"FIELD", "VALUE"},
	},
	XERR_NOT_FOUND: {
		defaultMessage: "target resource not found",
		extraTags:      []string{"FIELD", "TGT"},
	},
	XERR_PERMISSION_DENIED: {
		defaultMessage: "resource access permission denied",
		extraTags:      []string{"FIELD", "TGT"},
	},
	XERR_RESOURCE_UNAVAILABLE: {
		defaultMessage: "target resource currently unavailable",
		extraTags:      []string{"FIELD", "TGT"},
	},

	// ============================================================================
	// GROUP 2: NUMERIC, BOUNDARY, AND LIMIT VALIDATIONS
	// ============================================================================
	XERR_INVALID_VALUE: {
		defaultMessage: "invalid value",
		extraTags:      []string{"FIELD", "VALUE", "RULES"},
	},
	XERR_INVALID_VALUE_GT_ZERO: {
		defaultMessage: "invalid numeric",
		defaultRule:    "must be greater than zero (> 0)",
		extraTags:      []string{"FIELD", "VALUE", "RULES"},
	},
	XERR_INVALID_VALUE_GE_ZERO: {
		defaultMessage: "invalid numeric",
		defaultRule:    "must be than or equal to zero (>= 0)",
		extraTags:      []string{"FIELD", "VALUE", "RULES"},
	},
	XERR_INVALID_VALUE_LT_ZERO: {
		defaultMessage: "invalid numeric",
		defaultRule:    "must be less than zero (< 0)",
		extraTags:      []string{"FIELD", "VALUE", "RULES"},
	},
	XERR_INVALID_VALUE_LE_ZERO: {
		defaultMessage: "invalid numeric",
		defaultRule:    "must be less than or equal to zero (<= 0)",
		extraTags:      []string{"FIELD", "VALUE", "RULES"},
	},
	XERR_INVALID_VALUE_OUT_OF_RANGE: {
		defaultMessage: "invalid numeric",
		defaultRule:    "out of range",
		extraTags:      []string{"FIELD", "VALUE", "RULES"},
	},
	XERR_SELECTION_LIMIT_EXCEEDED: {
		defaultMessage: "selection quantity limit exceeded",
		extraTags:      []string{"FIELD", "OPT", "COUNT", "LIMIT"},
	},

	// ============================================================================
	// GROUP 3: STRUCTURE, TYPING, AND CHOICE VALIDATIONS
	// ============================================================================
	XERR_INVALID_FORMAT: {
		defaultMessage: "malformed syntax data structure",
		extraTags:      []string{"FIELD", "EXPECTED", "GIVEN"},
	},
	XERR_INVALID_TYPE: {
		defaultMessage: "type mismatch restriction violated",
		extraTags:      []string{"FIELD", "VALUE", "EXPECTED_TYPE"},
	},
	XERR_INVALID_OPTION: {
		defaultMessage: "invalid option selection",
		extraTags:      []string{"FIELD", "OPT", "OPTIONS"},
	},
	XERR_MUTUAL_EXCLUSIVITY_VIOLATION: {
		defaultMessage: "mutual exclusivity violation (choose only one)",
		extraTags:      []string{"FIELD", "OPT", "OPTIONS"},
	},
	XERR_ASYMMETRIC_SIZES: {
		defaultMessage: "asymmetric collections sizes",
		extraTags:      []string{"FIELDS"},
	},

	// ============================================================================
	// GROUP 4: GENERIC OPERATIONAL FAILURE FALLBACKS
	// ============================================================================
	XERR_INVALID_DATA: {
		defaultMessage: "invalid data",
		extraTags:      []string{"FIELD", "DATA"},
	},
}

// shouldIgnoreValue evaluates if a given option target matches any entries inside the provided blacklist.
// It uses strict interface equality for primitive matches and falls back to deep numeric reflection
// to catch underlying zero-value variations (e.g., customInt vs native int) across numeric boundaries.
func shouldIgnoreValue(val any, blacklist []any) bool {
	for _, ign := range blacklist {
		if val == ign {
			return true
		}
		if val != nil && ign != nil {
			vRef := reflect.ValueOf(val)
			iRef := reflect.ValueOf(ign)
			if vRef.Kind() == iRef.Kind() && vRef.IsValid() && iRef.IsValid() {
				if vRef.Kind() >= reflect.Int && vRef.Kind() <= reflect.Int64 {
					if vRef.Int() == iRef.Int() {
						return true
					}
				}
			}
		}
	}
	return false
}

// formatValueWithQuotes evaluates an item's underlying type primitive via reflection
// and stringifies it, dynamically enclosing strings and complex structures within single quotes.
func formatValueWithQuotes(item any) string {
	if item == nil {
		return "nil"
	}

	v := reflect.ValueOf(item)
	kind := v.Kind()

	switch {
	case kind >= reflect.Int && kind <= reflect.Uint64:
		return fmt.Sprintf("%d", item)
	case kind == reflect.Float32 || kind == reflect.Float64:
		return fmt.Sprintf("%v", item)
	case kind == reflect.Bool:
		return fmt.Sprintf("'%t'", item)
	default:
		return fmt.Sprintf("'%v'", item)
	}
}

// MENSAGENS DE ERRO COMUNS CONFORME O ERR EVOCADO
// Print writes the error's string representation directly to the standard error stream (os.Stderr).
// It acts as a lightweight, zero-allocation wrapper around fmt.Fprintln, providing a swift
// debugging mechanism that bypasses context allocations or structured logging handler configurations.
// If the incoming error target is nil, the execution block returns early without writing to the stream.
func Print(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err.Error())
}

// MsgData processes a variadic sequence of alternating key-value arguments and stringifies them
// into a highly predictable, comma-separated operational layout (e.g., "key1=val1, key2=val2").
//
// Exclusion Mechanics:
//   - By default, it completely omits any pair where the VALUE resolves to nil, zero (0), or an empty string ("").
//   - If the very first argument is provided as a []any slice, it overrides the default value blacklist,
//     enforcing only the elements specified within that slice as the new filtration parameters.
func MsgData(args ...any) string {
	if len(args) == 0 {
		return ""
	}

	ignoredValues := []any{nil, 0, ""}
	startingIdx := 0

	// Check if the developer provided a custom blacklist slice override in the first index
	if customIgnored, ok := args[0].([]any); ok {
		ignoredValues = customIgnored
		startingIdx = 1
	}

	remainingArgs := args[startingIdx:]
	if len(remainingArgs) == 0 {
		return ""
	}

	var parts []string
	for i := 0; i < len(remainingArgs); i += 2 {
		key := fmt.Sprintf("%v", remainingArgs[i])

		// Boundary protection if the developer provided a key but forgot the value
		if i+1 >= len(remainingArgs) {
			parts = append(parts, fmt.Sprintf("%s=(MISSING_VALUE)", key))
			break
		}

		val := remainingArgs[i+1]
		// Evaluates the paired value against our unified blacklist rule
		if shouldIgnoreValue(val, ignoredValues) {
			continue
		}

		parts = append(parts, fmt.Sprintf("%s=%s", key, formatValueWithQuotes(val)))
	}

	return strings.Join(parts, ", ")
}

// MsgOptions aggregates a variadic sequence of arguments into a clean, comma-separated
// text representation where strings and complex types are dynamically enclosed in single quotes.
//
// Exclusion Mechanics:
//   - By default, it ignores explicit nil values, zeros (0), and empty strings ("") to prevent visual clutter.
//   - If the very first argument is provided as a []any slice, it overrides the default exclusion rule.
func MsgOptions(args ...any) string {
	if len(args) == 0 {
		return ""
	}

	ignoredValues := []any{nil, 0, ""}
	startingIdx := 0

	if customIgnored, ok := args[0].([]any); ok {
		ignoredValues = customIgnored
		startingIdx = 1
	}

	var compiledParts []string

	for i := startingIdx; i < len(args); i++ {
		item := args[i]
		if shouldIgnoreValue(item, ignoredValues) {
			continue
		}

		if item == nil {
			compiledParts = append(compiledParts, "nil")
			continue
		}

		compiledParts = append(compiledParts, formatValueWithQuotes(item))
	}

	return strings.Join(compiledParts, ", ")
}

// MsgArraySize processes a variadic sequence of alternating key-length arguments
// and builds a unified diagnostic string mapping collections to their dimensions.
//
// Arguments must follow a strict sequential layout: string (name) paired with an integer (size).
// Numeric zeros (0) are explicitly preserved as they represent critical operational states.
func MsgArraySize(args ...any) string {
	if len(args) == 0 {
		return ""
	}

	var parts []string
	for i := 0; i < len(args); i += 2 {
		name := fmt.Sprintf("%v", args[i])

		// Boundary protection against trailing unpaired name tokens
		if i+1 >= len(args) {
			parts = append(parts, fmt.Sprintf("%s((MISSING_SIZE))", name))
			break
		}

		size := args[i+1]
		parts = append(parts, fmt.Sprintf("%s(%v)", name, size))
	}

	return strings.Join(parts, ", ")
}

// NewErr acts as a lightweight, polymorphic factory that creates and returns a standard Go error instance.
// It features an advanced internal normalization engine tailored for the corporate layout triggered by the XERR token.
// If the incoming message matches XERR, it dynamically reconstructs the structural string output based on
// the number of arguments provided, automatically backfilling missing slots with the empty set marker "ø".
// NewErr acts as a lightweight, polymorphic factory that creates and returns a standard Go error instance.
// It intercepts specialized corporate XERR_ tokens to dynamically delegate the formatting sequence
// to dedicated underlying layout builder functions.
func NewErr(message string, args ...any) error {
	// If it is just a plain, regular message with zero parameters, return immediately via errors.New
	if len(args) == 0 {
		return errors.New(message)
	}

	formatMsg := message
	nArgs := []any{}

	if meta, exists := errorRegistry[message]; exists {
		formatMsg, args = buildMask(meta, args)
	}

	for _, val := range args {
		nArgs = append(nArgs, val)
	}

	return fmt.Errorf(formatMsg, nArgs...)
}

// buildMask dynamically constructs a standardized visual layout string and normalizes
// the corresponding arguments slice for a given corporate XERR_ token specification.
//
// It processes the error metadata through a multi-step pipeline:
//  1. Appends the mandatory contextual tracking segment ([CTX: %v]).
//  2. Evaluates the human-readable message (MSG), backfilling with defaultMessage if empty or "ø" if omitted.
//  3. Intercepts Group 2 validation boundaries to inject defaultRule behaviors when missing.
//  4. Iterates through semantic extraTags (e.g., FIELD, TGT, VALUE), mapping available arguments or safe-guarding missing ones with the mathematical empty set ("ø").
//  5. Detects a trailing root cause error contract behind a double-colon boundary, safely applying the native Go wrapping verb (%w) to preserve errors.Unwrap execution.
func buildMask(meta errorMetadata, args []any) (string, []any) {
	// Base allocations for high-performance string layout compilation
	var builder strings.Builder
	builder.WriteString("[CTX: %v]")

	// Step 1: Normalize the human-friendly message (MSG) if provided as an empty string
	if len(args) > 1 {
		if msgStr, ok := args[1].(string); ok && msgStr == "" {
			args[1] = meta.defaultMessage
		}
		builder.WriteString("[MSG: %v]")
	} else {
		builder.WriteString("[MSG: ø]")
	}

	// Step 2: Automatically inject the default mathematical rule if it was omitted by the developer
	// This specifically protects Group 2 boundaries (e.g., XERR_INVALID_VALUE_GT_ZERO)
	expectedArgsCount := 2 + len(meta.extraTags)
	if meta.defaultRule != "" {
		if len(args) == expectedArgsCount-1 {
			// Developer provided values for all extra tags except the last one (the rule)
			args = append(args, meta.defaultRule)
		} else if len(args) >= expectedArgsCount {
			// Developer provided an argument, but if it is an empty string, apply fallback
			ruleIdx := expectedArgsCount - 1
			if ruleStr, ok := args[ruleIdx].(string); ok && ruleStr == "" {
				args[ruleIdx] = meta.defaultRule
			}
		}
	}

	// Step 3: Iterate through the dynamic extra tags array to construct block structures
	for i, tagName := range meta.extraTags {
		argIdx := 2 + i // Account for offset: index 0 is CTX, index 1 is MSG

		if len(args) > argIdx {
			// Argument is present for this tag frame segment
			builder.WriteString(fmt.Sprintf("[%s: %%v]", tagName))
		} else {
			// Target argument is missing, apply the mathematical empty set safeguard character
			builder.WriteString(fmt.Sprintf("[%s: ø]", tagName))
		}
	}

	// Step 4: Isolate the low-level system error root cause behind the double-colon boundary
	// Check if there is an extra trailing argument representing the native error contract
	if len(args) > expectedArgsCount {
		lastArgIdx := len(args) - 1
		if rootErr, ok := args[lastArgIdx].(error); ok {
			// Fulfills the standard wrapping contract using the native percentage-w verb safely
			args[lastArgIdx] = rootErr
			builder.WriteString("::[ERR: %w]")
		} else {
			builder.WriteString("::[ERR: %v]")
		}
	} else {
		builder.WriteString("::[ERR: ø]")
	}

	return builder.String(), args
}

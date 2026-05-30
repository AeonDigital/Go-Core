package xerrors

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

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

// buildMask dynamically constructs a standardized visual layout string and normalizes
// the corresponding arguments slice for a given corporate XERR_ token specification.
//
// It processes the error metadata through a multi-step pipeline:
//  1. Appends the mandatory contextual tracking segment ([CTX: %v]).
//  2. Evaluates the human-readable message (MSG), backfilling with defaultMessage if empty or "ø" if omitted.
//  3. Intercepts Group 2 validation boundaries to inject fieldRule behaviors when missing.
//  4. Iterates through semantic extraTags (e.g., FIELD, TGT, VALUE), mapping available arguments or safe-guarding missing ones with the mathematical empty set ("ø").
//  5. Detects a trailing root cause error contract behind a double-colon boundary, safely applying the native Go wrapping verb (%w) to preserve the execution of errors.Unwrap.
func buildMask(meta MetaMessage, args []any) (string, []any) {
	var builder strings.Builder
	builder.WriteString("[CTX: %v]")

	// Step 1: Isolate the trailing root cause error contract to prevent positional slippage
	var rootErr error
	hasRootErr := false

	if len(args) > 0 {
		lastIdx := len(args) - 1
		if e, ok := args[lastIdx].(error); ok {
			rootErr = e
			hasRootErr = true
			// Temporarily slice out the error to normalize inner data tags safely
			args = args[:lastIdx]
		}
	}

	// Structural Safeguard: Dynamically normalize the slice capacity to match
	// the combined payload threshold of the current core token metadata.
	// This locks in position predictability and eliminates index out of range panic states.
	minExpectedLen := 2 + len(meta.extraTags)
	for len(args) < minExpectedLen {
		args = append(args, "ø")
	}

	// Step 2: Normalize the human-friendly message (MSG)
	if msgStr, ok := args[1].(string); ok && msgStr == "" {
		args[1] = meta.message
	} else if args[1] == "ø" && meta.message != "" {
		args[1] = meta.message
	}
	builder.WriteString("[MSG: %v]")

	// Step 3: Enforce/Inject Group 2 boundary constraints (fieldRule)
	if meta.fieldRule != "" {
		ruleIdx := minExpectedLen - 1
		if ruleStr, ok := args[ruleIdx].(string); ok && ruleStr == "" {
			args[ruleIdx] = meta.fieldRule
		} else if args[ruleIdx] == "ø" {
			args[ruleIdx] = meta.fieldRule
		}
	}

	// Step 4: Map extra tags safely
	for _, tagName := range meta.extraTags {
		builder.WriteString("[")
		builder.WriteString(tagName)
		builder.WriteString(": %v]")
	}

	// Step 5: Re-attach the root cause error behind the double-colon boundary
	if hasRootErr {
		args = append(args, rootErr)
		builder.WriteString(" :: [ERR: %w]")
	} else {
		builder.WriteString(" :: [ERR: ø]")
	}

	return builder.String(), args
}

// shouldIgnoreValue evaluates if a given option target matches any entries inside the provided blacklist.
// It uses strict interface equality for primitive matches and falls back to deep numeric reflection
// to catch underlying zero-value variations (e.g., a customInt vs native int) across numeric boundaries.
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

// formatValueWithQuotes evaluates an item's underlying primitive kind via reflection
// and stringifies its raw extracted value, dynamically enclosing boolean, string,
// and complex structural shapes within single quotes while keeping numeric formats clean.
func formatValueWithQuotes(item any) string {
	if item == nil {
		return "nil"
	}

	v := reflect.ValueOf(item)
	kind := v.Kind()

	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		// Extracts the true primitive int64 even from custom named types (e.g., type CustomInt int)
		return fmt.Sprintf("%d", v.Int())

	case kind >= reflect.Uint && kind <= reflect.Uint64:
		// Extracts the true primitive uint64 safely
		return fmt.Sprintf("%d", v.Uint())

	case kind == reflect.Float32 || kind == reflect.Float64:
		// Extracts the true primitive float64 safely
		return fmt.Sprintf("%v", v.Float())

	case kind == reflect.Bool:
		// Extracts the true primitive bool safely and encloses it in quotes
		return fmt.Sprintf("'%t'", v.Bool())

	case kind == reflect.String:
		// Extracts the true primitive string safely and encloses it in quotes
		return fmt.Sprintf("'%s'", v.String())

	default:
		// Fallback for complex structures, slices, maps, and pointers
		return fmt.Sprintf("'%v'", item)
	}
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

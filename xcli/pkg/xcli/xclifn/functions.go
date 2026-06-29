package xclifn

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AeonDigital/Go-Core/sterr/pkg/sterr"
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xcliintfc"
)

/*
  ARCHITECTURE & SCOPE LIMITATION:
  functions.go groups stateless, decoupled utility behaviors and pure computational
  routines required to back up the central proposal of the package.

  Design Constraints:
  - Every function placed here should ideally be deterministic (same input produces same output).
  - No global state mutation or complex side-effects are allowed within these routines.
  - If routines grow complex or introduce stateful context, split them into dedicated files.
*/

// Insert standalone functions or mathematical algorithms below.

//
//
//

// PrintStdout writes a message to the standard output.
//
// Arguments:
//   - message: The text string to be printed or used as a format base.
//   - args: Optional variadic parameters to populate the formatting placeholders.
//
// Error & Panic Natures:
//   - Complex Errors: Silently fails to print if the underlying os.Stdout stream
//     is closed or encounters a hardware write block.
func PrintStdout(message string, args ...any) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stdout, message)
		return
	}

	format := message
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stdout, format, args...)
}

// PrintStderr writes a message to the standard error output.
//
// Arguments:
//   - message: The text string to be printed or used as a format base.
//   - args: Optional variadic parameters to populate the formatting placeholders.
//
// Error & Panic Natures:
//   - Complex Errors: Silently fails to print if the underlying os.Stderr stream
//     is closed or encounters a hardware write block.
func PrintStderr(message string, args ...any) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, message)
		return
	}

	format := message
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
}

// NewError creates and returns an error with a formatted message.
//
// Arguments:
//   - message: The text string to be used directly or as a format base.
//   - args: Optional variadic parameters to populate the formatting placeholders.
//
// Returns:
//   - error: A standard Go error wrapping the generated message string.
func NewError(message string, args ...any) error {
	if len(args) == 0 {
		return errors.New(message)
	}

	return fmt.Errorf(message, args...)
}

// PrintError writes an error's message directly to the standard error output.
func PrintError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

//
//
//

// ParseRawArgs extracts flags and their associated raw values from terminal arguments.
//
// Arguments:
//   - rawArgs: A slice containing strictly the flag segments of the command line,
//     with all command and subcommand nodes already resolved and stripped out.
//
// Returns:
//   - map[string]string: A map linking the clean flag name (without dashes) to its raw string value.
//   - error: Returns an error if any argument violates the strict flag syntax rules.
//
// Error & Panic Natures:
//   - Complex Errors: Returns an sterr.CliError via SetMessage if a positional argument
//     is found, if a flag name is empty, or if a short flag exceeds length boundaries.
func ParseRawArgs(rawArgs []string) (map[string]string, error) {
	rawFlags := make(map[string]string)
	i := 0

	for i < len(rawArgs) {
		arg := rawArgs[i]

		// 1. Check for assignment syntax (--flag=value)
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			flagPart := parts[0]
			valuePart := parts[1]

			cleanFlag, err := extractFlagName(flagPart)
			if err != nil {
				return nil, err
			}

			rawFlags[cleanFlag] = valuePart
			i++
			continue
		}

		// 2. Check if it is a standard flag pointer (--flag or -f)
		if strings.HasPrefix(arg, "-") {
			cleanFlag, err := extractFlagName(arg)
			if err != nil {
				return nil, err
			}

			// Peek next argument to see if it is the value or another flag
			if i+1 < len(rawArgs) && !strings.HasPrefix(rawArgs[i+1], "-") {
				rawFlags[cleanFlag] = rawArgs[i+1]
				i += 2 // Consumed both flag and value
			} else {
				// Boolean flag trigger (no explicit value provided)
				rawFlags[cleanFlag] = "true"
				i++
			}
			continue
		}

		// 3. Any loose argument that does not start with "-" and has no "=" is a positional argument
		return nil, sterr.New().
			SetMessage("invalid argument: '%s'", arg)
	}

	return rawFlags, nil
}

// extractFlagName removes dashes and validates the length of a flag string.
func extractFlagName(flag string) (string, error) {
	// If it was detected as a long flag
	if strings.HasPrefix(flag, "--") {
		clean := flag[2:]
		if clean == "" {
			return "", sterr.New().SetMessage("invalid flag provided: '--'")
		}
		return clean, nil
	}

	// If it falls here, it is a short flag (guaranteed by the main loop structure)
	clean := flag[1:]
	if clean == "" {
		return "", sterr.New().SetMessage("invalid flag provided: '-'")
	}
	if len(clean) > 3 {
		return "", sterr.New().SetMessage("invalid short flag: '%s'", flag)
	}
	return clean, nil
}

//
//
//

// xcliintfc.Quantifiable

// isLess safely performs a compile-time and runtime dispatch to compare ordering constraints.
func isLess[T xcliintfc.Quantifiable](a, b T) bool {
	switch v := any(a).(type) {
	case time.Time:
		return v.Before(any(b).(time.Time))
	case int:
		return v < any(b).(int)
	case float64:
		return v < any(b).(float64)
	case time.Duration:
		return v < any(b).(time.Duration)
	default:
		return false
	}
}

// isGreater safely performs a compile-time and runtime dispatch to compare ordering constraints.
func isGreater[T xcliintfc.Quantifiable](a, b T) bool {
	switch v := any(a).(type) {
	case time.Time:
		return v.After(any(b).(time.Time))
	case int:
		return v > any(b).(int)
	case float64:
		return v > any(b).(float64)
	case time.Duration:
		return v > any(b).(time.Duration)
	default:
		return false
	}
}

// formatQuantifiable helper isolates terminal serialization layouts for quantifiable primitives.
func formatQuantifiable[T xcliintfc.Quantifiable](val, limit T, timeLayout string) (string, string) {
	switch v := any(val).(type) {
	case int:
		l := any(limit).(int)
		return fmt.Sprintf("%d", v), fmt.Sprintf("%d", l)
	case float64:
		l := any(limit).(float64)
		return fmt.Sprintf("%g", v), fmt.Sprintf("%g", l)
	case time.Duration:
		l := any(limit).(time.Duration)
		return v.String(), l.String()
	case time.Time:
		l := any(limit).(time.Time)
		return v.Format(timeLayout), l.Format(timeLayout)
	default:
		return "", ""
	}
}

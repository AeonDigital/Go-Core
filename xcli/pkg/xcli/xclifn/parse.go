package xclifn

import (
	"encoding/json"
	"net/mail"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/AeonDigital/Go-Core/sterr/pkg/sterr"
)

// ParseString evaluates and cleans a raw string flag value from the terminal.
//
// Arguments:
//   - val: The raw string representation from the terminal.
//
// Returns:
//   - string: The trimmed and standardized target string.
//   - error: Returns nil as text translation natively always succeeds.
func ParseString(val string) (string, error) {
	return strings.TrimSpace(val), nil
}

// ParseBool converts a raw string flag value into a native Go boolean.
//
// Arguments:
//   - val: The raw string representation from the terminal (e.g., "true", "false").
//
// Returns:
//   - bool: The parsed boolean value state.
//   - error: Returns an sterr.CliError if the conversion fails.
//
// Error & Panic Natures:
//   - Complex Errors: Returns an error if the string is not a valid boolean expression.
func ParseBool(val string) (bool, error) {
	// Standardize to lowercase for resilient parsing
	clean := strings.ToLower(strings.TrimSpace(val))
	if clean == "true" || clean == "1" {
		return true, nil
	}
	if clean == "false" || clean == "0" {
		return false, nil
	}

	return false, sterr.New().
		SetMessage("invalid boolean value: '%s'", val)
}

// ParseInt converts a raw string flag value into a native Go integer.
//
// Arguments:
//   - val: The raw string representation from the terminal.
//
// Returns:
//   - int: The parsed integer primitive.
//   - error: Returns an sterr.CliError if the conversion fails.
//
// Error & Panic Natures:
//   - Complex Errors: Returns an error if the string contains non-numeric characters.
func ParseInt(val string) (int, error) {
	clean := strings.TrimSpace(val)
	res, err := strconv.Atoi(clean)
	if err != nil {
		return 0, sterr.New().
			SetMessage("invalid integer value: '%s'", val)
	}
	return res, nil
}

// ParseFloat converts a raw string flag value into a native Go float64.
//
// Arguments:
//   - val: The raw string representation from the terminal.
//
// Returns:
//   - float64: The parsed floating-point decimal.
//   - error: Returns an sterr.CliError if the conversion fails.
//
// Error & Panic Natures:
//   - Complex Errors: Returns an error if the string is not a valid decimal representation.
func ParseFloat(val string) (float64, error) {
	clean := strings.TrimSpace(val)
	res, err := strconv.ParseFloat(clean, 64)
	if err != nil {
		return 0.0, sterr.New().
			SetMessage("invalid float value: '%s'", val)
	}
	return res, nil
}

// ParseDuration converts a raw string flag value into a native Go time.Duration.
//
// Arguments:
//   - val: The raw string representation from the terminal (e.g., "30s", "2h45m").
//
// Returns:
//   - time.Duration: The parsed duration object.
//   - error: Returns an sterr.CliError if the conversion fails.
//
// Error & Panic Natures:
//   - Complex Errors: Returns an error if the string format does not match Go duration syntax rules.
func ParseDuration(val string) (time.Duration, error) {
	clean := strings.TrimSpace(val)
	res, err := time.ParseDuration(clean)
	if err != nil {
		return 0, sterr.New().
			SetMessage("invalid duration value: '%s'. Expected unit format (e.g., '30s', '1h30m')", val)
	}
	return res, nil
}

// ParseDate converts a raw string flag value into a native Go time.Time object using YYYY-MM-DD layout.
//
// Arguments:
//   - val: The raw string representation from the terminal.
//
// Returns:
//   - time.Time: The parsed time object anchored at midnight.
//   - error: Returns an sterr.CliError if the conversion fails.
//
// Error & Panic Natures:
//   - Complex Errors: Returns an error if the string layout does not strictly match 'YYYY-MM-DD'.
func ParseDate(val string) (time.Time, error) {
	clean := strings.TrimSpace(val)
	res, err := time.Parse("2006-01-02", clean)
	if err != nil {
		return time.Time{}, sterr.New().
			SetMessage("invalid date value: '%s'. Expected format: 'YYYY-MM-DD'", val)
	}
	return res, nil
}

// ParseTime converts a raw string flag value into a native Go time.Time object.
//
// It flexibly accepts both full clock layout 'HH:MM:SS' and short layout 'HH:MM'.
// The resulting object will be anchored at the neutral Year Zero (0000-01-01).
//
// Arguments:
//   - val: The raw string representation from the terminal.
//
// Returns:
//   - time.Time: The parsed time object containing only clock metrics.
//   - error: Returns an sterr.CliError if the conversion fails.
//
// Error & Panic Natures:
//   - Complex Errors: Returns an error if the string layout does not match 'HH:MM:SS' or 'HH:MM'.
func ParseTime(val string) (time.Time, error) {
	clean := strings.TrimSpace(val)

	// Try the full clock format first (HH:MM:SS)
	if res, err := time.Parse("15:04:05", clean); err == nil {
		return res, nil
	}

	// Fallback to the short clock format (HH:MM)
	if res, err := time.Parse("15:04", clean); err == nil {
		return res, nil
	}

	return time.Time{}, sterr.New().
		SetMessage("invalid time value: '%s'. Expected format: 'HH:MM' or 'HH:MM:SS'", val)
}

// ParseDateTime converts a raw string flag value into a native Go time.Time object using YYYY-MM-DD HH:MM:SS layout.
//
// Arguments:
//   - val: The raw string representation from the terminal.
//
// Returns:
//   - time.Time: The fully parsed calendar and clock time object.
//   - error: Returns an sterr.CliError if the conversion fails.
//
// Error & Panic Natures:
//   - Complex Errors: Returns an error if the string layout does not strictly match 'YYYY-MM-DD HH:MM:SS'.
func ParseDateTime(val string) (time.Time, error) {
	clean := strings.TrimSpace(val)
	res, err := time.Parse("2006-01-02 15:04:05", clean)
	if err != nil {
		return time.Time{}, sterr.New().
			SetMessage("invalid datetime value: '%s'. Expected format: 'YYYY-MM-DD HH:MM:SS'", val)
	}
	return res, nil
}

// ParseEmail validates a raw string flag value ensuring it adheres to strict email syntax formats.
//
// Arguments:
//   - val: The raw string representation from the terminal.
//
// Returns:
//   - string: The validated and sanitized email address string.
//   - error: Returns an sterr.CliError if the validation fails.
//
// Error & Panic Natures:
//   - Complex Errors: Returns an error if the payload fails RFC 5322 address formatting rules.
func ParseEmail(val string) (string, error) {
	clean := strings.TrimSpace(val)
	addr, err := mail.ParseAddress(clean)
	if err != nil || addr.Address != clean {
		return "", sterr.New().
			SetMessage("invalid email address: '%s'", val)
	}
	return addr.Address, nil
}

// ParseJSON validates if a raw string flag value is a syntactically correct JSON payload.
//
// Arguments:
//   - val: The raw string representation from the terminal.
//
// Returns:
//   - string: The validated raw JSON text string.
//   - error: Returns an sterr.CliError if the format evaluation fails.
//
// Error & Panic Natures:
//   - Complex Errors: Returns an error if the string payload contains broken or unparsable JSON syntax.
func ParseJSON(val string) (string, error) {
	clean := strings.TrimSpace(val)
	if !json.Valid([]byte(clean)) {
		return "", sterr.New().
			SetMessage("invalid JSON format provided")
	}
	return clean, nil
}

// ParseFilename evaluates if a string is a syntactically valid standalone filename.
//
// Arguments:
//   - val: The raw string representation from the terminal.
//
// Returns:
//   - string: The cleaned semantic filename string layout.
//   - error: Returns an sterr.CliError if directory slashes are detected.
//
// Error & Panic Natures:
//   - Complex Errors: Returns an error if the payload is empty or contains path separator bounds.
func ParseFilename(val string) (string, error) {
	clean := strings.TrimSpace(val)
	if clean == "" {
		return "", sterr.New().
			SetMessage("invalid filename: empty")
	}

	cleanedPath := filepath.Clean(clean)

	if strings.Contains(cleanedPath, "/") || strings.Contains(cleanedPath, "\\") {
		return "", sterr.New().
			SetMessage("invalid filename: '%s'. It must not contain directory slashes", val)
	}

	return cleanedPath, nil
}

// ParseDirname evaluates if a string is a syntactically valid directory name layout.
//
// Arguments:
//   - val: The raw string representation from the terminal.
//
// Returns:
//   - string: The cleaned semantic directory name layout text.
//   - error: Returns an sterr.CliError if the payload is empty.
func ParseDirname(val string) (string, error) {
	clean := strings.TrimSpace(val)
	if clean == "" {
		return "", sterr.New().
			SetMessage("invalid directory: empty")
	}

	return filepath.Clean(clean), nil
}

// ParsePath evaluates if a string is a grammatically valid filesystem layout component.
//
// Arguments:
//   - val: The raw string representation from the terminal.
//
// Returns:
//   - string: The cleaned and standardized semantic path layout text.
//   - error: Returns an sterr.CliError if the payload is empty.
func ParsePath(val string) (string, error) {
	clean := strings.TrimSpace(val)
	if clean == "" {
		return "", sterr.New().
			SetMessage("invalid path: empty")
	}

	return filepath.Clean(clean), nil
}

// ParseURLStandard validates if a raw string payload adheres to basic URL architecture criteria.
//
// Arguments:
//   - val: The raw string representation from the terminal.
//
// Returns:
//   - string: The validated raw URL text string layout.
//   - error: Returns an sterr.CliError if formatting rules fail evaluation.
func ParseURLStandard(val string) (string, error) {
	clean := strings.TrimSpace(val)
	_, err := url.Parse(clean)
	if err != nil {
		return "", sterr.New().
			SetMessage("invalid URL: '%s'", val)
	}

	return clean, nil
}

// ParseFullURL validates if a raw string payload is a fully qualified URL locator domain.
//
// Arguments:
//   - val: The raw string representation from the terminal.
//
// Returns:
//   - string: The validated full URL text string layout.
//   - error: Returns an sterr.CliError if the protocol scheme or host layer is absent.
func ParseFullURL(val string) (string, error) {
	clean := strings.TrimSpace(val)
	parsed, err := url.Parse(clean)
	if err != nil {
		return "", sterr.New().
			SetMessage("invalid URL: '%s'", val)
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return "", sterr.New().
			SetMessage("invalid full URL: '%s'. Scheme and Host are mandatory", val)
	}

	return clean, nil
}

// ParseRelativeURL validates if a raw string payload is a strict relative target locator layout.
//
// Arguments:
//   - val: The raw string representation from the terminal.
//
// Returns:
//   - string: The validated relative URL text string layout.
//   - error: Returns an sterr.CliError if a protocol scheme or host layer is present.
func ParseRelativeURL(val string) (string, error) {
	clean := strings.TrimSpace(val)
	parsed, err := url.Parse(clean)
	if err != nil {
		return "", sterr.New().
			SetMessage("invalid URL: '%s'", val)
	}

	if parsed.Scheme != "" || parsed.Host != "" {
		return "", sterr.New().
			SetMessage("invalid relative URL: '%s'. Protocol Scheme and Host must be absent", val)
	}

	return clean, nil
}

// ParseFilepath normalizes and cleans a raw string payload for physical file evaluation tracks.
//
// Arguments:
//   - val: The raw string representation from the terminal.
//
// Returns:
//   - string: The cleaned physical file path layout text.
//   - error: Returns an sterr.CliError if the payload is empty.
func ParseFilepath(val string) (string, error) {
	clean := strings.TrimSpace(val)
	if clean == "" {
		return "", sterr.New().
			SetMessage("invalid filepath: empty")
	}

	return filepath.Clean(clean), nil
}

// ParseDirpath normalizes and cleans a raw string payload for physical directory evaluation tracks.
//
// Arguments:
//   - val: The raw string representation from the terminal.
//
// Returns:
//   - string: The cleaned physical directory path layout text.
//   - error: Returns an sterr.CliError if the payload is empty.
func ParseDirpath(val string) (string, error) {
	clean := strings.TrimSpace(val)
	if clean == "" {
		return "", sterr.New().
			SetMessage("invalid dirpath: empty")
	}

	return filepath.Clean(clean), nil
}

package xclifn

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/AeonDigital/Go-Core/sterr/pkg/sterr"
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xcliintfc"
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xclistruc"
)

var (
	// osStatBridge redirects telemetry access checks to allow full unit test injection coverage.
	osStatBridge = os.Stat

	// osOpenBridge redirects file reading privilege checks to allow structural error tracking.
	osOpenBridge = os.Open

	// osCreateBridge redirects directory writing privilege checks to isolate test mock environments.
	osCreateBridge = os.Create

	// osOpenFileBridge redirects file appending privilege checks to force low-level permission blocks.
	osOpenFileBridge = os.OpenFile
)

// newValidationError constructs a beautifully formatted, structured terminal error layout.
//
// Arguments:
//   - flag: The command line flag identifier label causing the failure (e.g., "--output").
//   - msg: A concise string describing the specific validation rule break nature.
//   - given: The raw or incorrectly formatted value provided by the user in the terminal.
//   - expected: A descriptive text detailing the compliance constraints layout criteria.
//   - rule: The exact technical boundary setting enforced by the flag schema spec.
//
// Returns:
//   - error: A structured sterr instance containing a highly readable visual error block.
//
// Error & Panic Natures:
//   - Complex Errors: Orchestrates a multi-line visual block for terminal outputs.
//     The first line establishes the contextual error anchor. Subsequent diagnostic metrics
//     (Rule, Given, Expected) are dynamically appended only if their string payload is
//     not empty, being padded with exactly six spaces to isolate the error anatomy
//     and maximize direct raw readability for the operator.
func newValidationError(
	flag string,
	msg string,
	given string,
	expected string,
	rule string,
) error {
	firstLine := fmt.Sprintf("[ERR] %s : %s", flag, msg)
	indent := "      "
	var subLines []string

	if rule != "" {
		subLines = append(subLines, fmt.Sprintf("%sRule: %s", indent, rule))
	}
	if given != "" {
		subLines = append(subLines, fmt.Sprintf("%sGiven: %s", indent, given))
	}
	if expected != "" {
		subLines = append(subLines, fmt.Sprintf("%sExpected: %s", indent, expected))
	}

	finalMsg := firstLine
	if len(subLines) > 0 {
		finalMsg = finalMsg + "\n" + strings.Join(subLines, "\n")
	}

	return sterr.New().SetMessage(finalMsg)
}

// ValidateStringLimits enforces length boundaries and regular expression pattern matching on a text value.
//
// Arguments:
//   - flag: The command line flag identifier label (e.g., "--name", "-n").
//   - val: The already parsed and cleaned target string wrapped in a generic interface.
//   - spec: The complete core flag layout containing the validation metadata guidelines.
//
// Returns:
//   - error: Returns an sterr.CliError if any constraint condition is violated.
//
// Error & Panic Natures:
//   - Complex Errors: Extracts and evaluates MinLength and MaxLength tracking unicode
//     character counts via rune evaluation bounds. Evaluates custom regex patterns
//     triggering a formatted terminal layout error if the expression evaluation fails.
func ValidateStringLimits(
	flag string,
	val string,
	spec xclistruc.Flag,
) error {
	length := utf8.RuneCountInString(val)

	if spec.MinLength != nil && length < *spec.MinLength {
		return newValidationError(
			flag, "string too short", val, "", fmt.Sprintf("min length: %d", *spec.MinLength),
		)
	}

	if spec.MaxLength != nil && length > *spec.MaxLength {
		return newValidationError(
			flag, "string too long", val, "", fmt.Sprintf("max length: %d", *spec.MaxLength),
		)
	}

	if spec.Regex != "" {
		matched, err := regexp.MatchString(spec.Regex, val)
		if err != nil {
			return sterr.New().
				SetMessage("invalid regex pattern: %s", spec.Regex)
		}
		if !matched {
			return newValidationError(
				flag, "invalid pattern value", val, "", fmt.Sprintf("regex: %s", spec.Regex),
			)
		}
	}

	return nil
}

// ValidateNumberLimits verifies if a quantifiable value resides between inclusive mathematical bounds.
//
// Arguments:
//   - flag: The command line flag identifier label (e.g., "--count", "-c").
//   - val: The native converted quantifiable primitive instance to evaluate (int, float64, time.Duration, time.Time).
//   - spec: The complete core flag layout containing the validation metadata guidelines.
//
// Returns:
//   - error: Returns an sterr.CliError if any mathematical boundary evaluation fails.
func ValidateNumberLimits[T xcliintfc.Quantifiable](
	flag string,
	val T,
	spec xclistruc.Flag,
) error {
	layout := "2006-01-02 15:04:05"

	// 1. Evaluate inclusive lower bound constraints
	if spec.Min != nil {
		if minVal, ok := spec.Min.(T); ok {
			if isLess(val, minVal) {
				givenStr, ruleStr := formatQuantifiable(val, minVal, layout)
				return newValidationError(
					flag, "out of range", givenStr, "", fmt.Sprintf("min: %s", ruleStr),
				)
			}
		}
	}

	// 2. Evaluate inclusive upper bound constraints
	if spec.Max != nil {
		if maxVal, ok := spec.Max.(T); ok {
			if isGreater(val, maxVal) {
				givenStr, ruleStr := formatQuantifiable(val, maxVal, layout)
				return newValidationError(
					flag, "out of range", givenStr, "", fmt.Sprintf("max: %s", ruleStr),
				)
			}
		}
	}

	return nil
}

// ValidateArrayLimits enforces capacity item layout boundaries on a strongly-typed slice.
//
// Arguments:
//   - flag: The command line flag identifier label (e.g., "--ids", "-i").
//   - totalItems: The calculated length metric of the active slice matrix instance.
//   - spec: The complete core flag layout containing the validation metadata guidelines.
//
// Returns:
//   - error: Returns an sterr.CliError if the slice layout violates capacity boundaries.
func ValidateArrayLimits(
	flag string,
	totalItems int,
	spec xclistruc.Flag,
) error {
	if spec.MinItems != nil && totalItems < *spec.MinItems {
		return newValidationError(
			flag, "out of range", fmt.Sprintf("%d items", totalItems), "", fmt.Sprintf("min items: %d", *spec.MinItems),
		)
	}

	if spec.MaxItems != nil && totalItems > *spec.MaxItems {
		return newValidationError(
			flag, "out of range", fmt.Sprintf("%d items", totalItems), "", fmt.Sprintf("max items: %d", *spec.MaxItems),
		)
	}

	return nil
}

// ValidateDiskResources performs OS evaluation tracks ensuring system capabilities match requirements.
//
// Arguments:
//   - flag: The command line flag identifier label (e.g., "--src", "-s").
//   - val: The cleaned physical path target string wrapped in a generic interface.
//   - spec: The complete core flag layout containing the validation metadata guidelines.
//
// Returns:
//   - error: Returns an sterr.CliError if any filesystem constraint or capability check fails.
//
// Error & Panic Natures:
//   - Complex Errors: Performs physical disk lookups evaluating spec.MustExist boundaries.
//     Inferes dynamically if the target must be a directory or a file matching the underlying
//     spec.Type metadata configuration. Triggers a formatted error block indicating
//     read or write access permission violations by testing low-level OS capabilities.
func ValidateDiskResources(
	flag string,
	val string,
	spec xclistruc.Flag,
) error {
	info, statErr := osStatBridge(val)

	// 1. Evaluate mandatory disk existence rule bounds
	if spec.MustExist && statErr != nil {
		if os.IsNotExist(statErr) {
			return newValidationError(
				flag, "path resource does not exist", val, "", "resource must exist",
			)
		}
		return sterr.New().
			SetMessage("[ERR] %s : disk telemetry access block: %v", flag, statErr)
	}

	typeStr := string(spec.Type)

	// 2. Validate target type mapping consistency
	if statErr == nil {
		if typeStr == "dirpath" && !info.IsDir() {
			return newValidationError(
				flag, "resource location is a file, directory required", val, "", "directory type constraint",
			)
		}
		if typeStr == "filepath" && info.IsDir() {
			return newValidationError(
				flag, "resource location is a directory, file required", val, "", "file type constraint",
			)
		}
	}

	// 3. Evaluate Access Permissions matrix layers
	if spec.Access != "" && statErr == nil {
		mode := string(spec.Access)
		isDir := info.IsDir()

		// Technical read capabilities confirmation lookups hook
		if mode == "read" || mode == "readwrite" {
			file, err := osOpenBridge(val)
			if err != nil {
				return newValidationError(
					flag, "missing read permission access privileges", val, "", "read access required",
				)
			}
			file.Close()
		}

		// Technical write capabilities confirmation lookups hook
		if mode == "write" || mode == "readwrite" {
			var err error
			if isDir {
				// Test directory write capacity by creating a temporary validation trace file node
				testFile := filepath.Join(val, ".xcli_perm_test")
				var f *os.File
				f, err = osCreateBridge(testFile)
				if err == nil {
					f.Close()
					os.Remove(testFile)
				}
			} else {
				// Test file write capacity by attempting to open it in write-append mode safely
				var f *os.File
				f, err = osOpenFileBridge(val, os.O_WRONLY, 0666)
				if err == nil {
					f.Close()
				}
			}

			if err != nil {
				return newValidationError(
					flag, "missing write permission access privileges", val, "", "write access required",
				)
			}
		}
	}

	return nil
}

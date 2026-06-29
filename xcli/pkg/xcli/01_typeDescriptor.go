package xcli

import (
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xclistruc"
)

// ============================================================================
// 1. GENERIC DESCRIPTOR FOR PRIMITIVE/STRUCTURED TYPES
// ============================================================================

// TypedDescriptor implements the xcliintfc.ValueParser interface, orchestrating
// type conversion and multi-layered validation tracks for a single primitive type.
type TypedDescriptor[T any] struct {
	// Parser conversions primitive engine pointer.
	Parser func(val string) (T, error)

	// Optional specific limit validator hook.
	LimitValidator func(flag string, val T, spec xclistruc.Flag) error

	// Optional structural disk capability validator hook.
	DiskValidator func(flag string, val T, spec xclistruc.Flag) error
}

// ParseAndValidate orchestrates the complete strongly-typed pipeline for single terminal values.
//
// It handles primitive mapping, standard constraint limits enforcement, and specialized
// physical file system privilege validation tracks using compile-time type dispatch.
//
// Arguments:
//   - flag: The command line flag identifier label causing the failure (e.g., "--output").
//   - raw: The raw unparsed text payload extracted from the command arguments collection.
//   - spec: The complete core flag layout containing the validation metadata guidelines.
//
// Returns:
//   - any: The fully converted, type-safe native Go primitive instance.
//   - error: Returns an sterr.CliError capturing visual telemetry diagnostics if any rule fails.
//
// Error & Panic Natures:
//   - Complex Errors: Evaluates sequential pipeline layers. Fails early if the core type parser
//     throws an execution blocker. Sequentially triggers LimitValidator and DiskValidator bounds
//     hooks only if they are actively configured inside the global registration matrix.
func (td TypedDescriptor[T]) ParseAndValidate(
	flag string,
	raw string,
	spec xclistruc.Flag,
) (any, error) {
	// 1. Transform raw terminal input string into native concrete type T
	parsedVal, err := td.Parser(raw)
	if err != nil {
		return nil, err
	}

	// 2. Enforce limits constraints bounds check if configured
	if td.LimitValidator != nil {
		if err := td.LimitValidator(flag, parsedVal, spec); err != nil {
			return nil, err
		}
	}

	// 3. Enforce low-level physical disk capability criteria checks if configured
	if td.DiskValidator != nil {
		if err := td.DiskValidator(flag, parsedVal, spec); err != nil {
			return nil, err
		}
	}

	return parsedVal, nil
}

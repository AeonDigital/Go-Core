package xcli

import (
	"strings"

	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xclifn"
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xclistruc"
)

// ============================================================================
// 2. DESCRITOR GENÉRICO PARA ARRAYS / SLICES ([]type)
// ============================================================================

// ArrayDescriptor implements the xcliintfc.ValueParser interface, orchestrating
// tokenizer splitting, iterative primitive parsing, and capacity validation for slices.
type ArrayDescriptor[T any] struct {
	// ElementParser references the single item parser function pointer.
	ElementParser func(val string) (T, error)

	// Optional item-level limit validator hook applied to every slice element.
	ItemLimitValidator func(flag string, val T, spec xclistruc.Flag) error

	// Optional item-level physical disk validator hook applied to every slice element.
	ItemDiskValidator func(flag string, val T, spec xclistruc.Flag) error
}

// validateCapacity enforces slice size limits against spec guidelines without reflection.
//
// It delegates collection item metrics checking directly to the specialized core validation
// package layers to isolate error visual formatting structures.
//
// Arguments:
//   - flag: The command line flag identifier label causing the failure (e.g., "--ids").
//   - total: The total number of items parsed and identified inside the active slice collection.
//   - spec: The complete core flag layout containing the validation metadata guidelines.
//
// Returns:
//   - error: Returns an sterr.CliError if the count violates configured boundaries.
//
// Error & Panic Natures:
//   - Complex Errors: Routes length counters down to the package validation infrastructure,
//     triggering a multi-line visual terminal error layout if minItems or maxItems are violated.
func (ad ArrayDescriptor[T]) ValidateCapacity(flag string, total int, spec xclistruc.Flag) error {
	return xclifn.ValidateArrayLimits(flag, total, spec)
}

// ParseAndValidate orchestrates split tokenization loops over array inputs.
//
// It splits terminal string blocks, validates total capacities metrics, and iteratively
// invokes primitive element translation and multi-layered restriction rules layers.
//
// Arguments:
//   - flag: The command line flag identifier label causing the failure (e.g., "--ids").
//   - raw: The raw unparsed text payload extracted from the command arguments collection.
//   - spec: The complete core flag layout containing the validation metadata guidelines.
//
// Returns:
//   - any: The fully converted, type-safe native Go slice matching the target element layout ([]T).
//   - error: Returns an sterr.CliError capturing visual telemetry diagnostics if any rule fails.
//
// Error & Panic Natures:
//   - Complex Errors: Evaluates high-level collection boundary lengths via custom helpers.
//     Iterates over individual string tokens throwing errors if a single primitive element parsing
//     or item-level constraint assertion fails evaluation targets.
func (ad ArrayDescriptor[T]) ParseAndValidate(flag string, raw string, spec xclistruc.Flag) (any, error) {
	cleanRaw := strings.TrimSpace(raw)

	// Handle completely blank array scenarios safely
	if cleanRaw == "" {
		emptySlice := make([]T, 0)
		if err := ad.ValidateCapacity(flag, 0, spec); err != nil {
			return nil, err
		}
		return emptySlice, nil
	}

	// Tokenize input by splitting the raw text using standard commas
	tokens := strings.Split(cleanRaw, ",")
	totalItems := len(tokens)

	// 1. Enforce high-level array capacity criteria boundaries
	if err := ad.ValidateCapacity(flag, totalItems, spec); err != nil {
		return nil, err
	}

	// Instantiate the final type-safe native concrete slice
	sliceResult := make([]T, totalItems)

	// 2. Iteratively process, convert, and validate every extracted token element
	for i, token := range tokens {
		itemVal, err := ad.ElementParser(token)
		if err != nil {
			return nil, err
		}

		if ad.ItemLimitValidator != nil {
			if err := ad.ItemLimitValidator(flag, itemVal, spec); err != nil {
				return nil, err
			}
		}

		if ad.ItemDiskValidator != nil {
			if err := ad.ItemDiskValidator(flag, itemVal, spec); err != nil {
				return nil, err
			}
		}

		sliceResult[i] = itemVal
	}

	return sliceResult, nil
}

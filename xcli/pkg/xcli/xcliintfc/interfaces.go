package xcliintfc

import (
	"time"

	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xclistruc"
)

/*
	ARCHITECTURE & SCOPE LIMITATION:
	interfaces.go defines the structural contracts and behaviors required
	to decouple high-level logic from concrete low-level implementations.

	Design Constraints:
	- Keep interfaces small and highly focused, adhering strictly to the Single
		Responsibility Principle (ideally 1 to 3 methods per interface).
	- Interfaces should be defined where they are consumed, not where they are implemented.
	- If an interface demands a massive set of methods or forms a distinct module contract,
		extract it into its own file inside this folder.
*/

// Insert abstract interface definitions below.

//
//
//

// ValueParser encapsulates the unified execution contract required to process,
// transform, and validate raw command-line terminal arguments into strongly-typed
// domain instances.
//
// Technical Role:
//   - Acts as an inversion-of-control boundary separating the generic string-based
//     CLI router from the underlying native Go representation mapping layers.
type ValueParser interface {
	// ParseAndValidate orchestrates the complete parsing, conversion, and constraint
	// evaluation lifecycle for a single command-line option value.
	//
	// Arguments:
	//   - flag: The primary diagnostic label representing the active terminal token (e.g., "--timeout").
	//   - raw: The raw unparsed text payload extracted from the command arguments collection.
	//   - spec: The complete metadata layout structure defining constraints, defaults, and validations.
	//
	// Returns:
	//   - any: The fully converted, type-safe native Go primitive or slice instance.
	//   - error: Returns an sterr.CliError capturing visual telemetry diagnostics if any rule fails.
	ParseAndValidate(flag string, raw string, spec xclistruc.Flag) (any, error)
}

// Quantifiable defines the strict set of types capable of undergoing inclusive
// mathematical boundary evaluation tracks within the xcli routing motor.
type Quantifiable interface {
	~int | ~float64 | time.Duration | time.Time
}

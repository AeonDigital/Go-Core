package xerrors

import (
	"github.com/AeonDigital/Go-Core/xerrors"
)

/*
  ARCHITECTURE & SCOPE LIMITATION:
  xerrors.go acts as the package error registry. It declares domain-specific
  error codes mapped tightly to corporate layout schemas, registering them into
  the global Go-Core xerrors ecosystem at initialization runtime.

  Design Constraints:
  - Every constant belonging to the error family must use the explicit xerrors.ErrorCode type.
  - All errors should follow the "EXXXX" naming standard, where "E" stands for Error
    and the first numerical digit points to the specific semantic family/group.
	- Do not implement complex logic here; this is strictly a registry for
    error matching and tracing.
  - If errors belong to a highly isolated, distinct subsystem (e.g., external API failures),
    group them into a separate file within this package.
*/

// Insert sentinel errors and custom error types below.

//
//
//

const (
	// XERR_NONE serves as a local empty fallback marker to simplify assertion testing.
	XERR_NONE xerrors.ErrorCode = ""

	//
	//
	// X_ANCHOR_PKGCTX_START

	// XERR_PKGCTX defines the error code namespace for this package.
	XERR_PKGCTX xerrors.ErrorCode = "ERR_XCLI"

	// X_ANCHOR_PKGCTX_END
	//
	//

	//
	//
	// X_ANCHOR_CONSTANTS_START

	// ============================================================================
	// === FAMILY: 1 | TITLE: GENERAL AND SYSTEM FALLBACKS
	// ============================================================================

	// XERR_UNKNOWN serves as the fallback categorization for untracked exceptions within this domain.
	// Format expects: CTX, MSG, [error]
	XERR_UNKNOWN xerrors.ErrorCode = "E1001"

	// X_ANCHOR_CONSTANTS_END
	//
	//
)

// xerrorLocalRegistry maps error codes to their structural metadata boundaries.
var xerrorLocalRegistry = map[xerrors.ErrorCode]xerrors.MetaMessage{
	//
	//
	// X_ANCHOR_REGISTRY_START

	// ============================================================================
	// === FAMILY: 1
	// ============================================================================

	XERR_UNKNOWN: xerrors.NewMetaMessage(
		"unexpected internal xcli error encountered",
		"",
		[]string{},
	),

	// X_ANCHOR_REGISTRY_END
	//
	//
}

func init() {
	// Automatically register local errors into the centralized tracking engine upon instantiation
	xerrors.RegisterDomainErrors(XERR_PKGCTX, xerrorLocalRegistry)
}

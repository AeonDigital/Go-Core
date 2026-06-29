package xtests

import (
	"github.com/AeonDigital/Go-Core/xerrors"
)

/*
  ARCHITECTURE & SCOPE LIMITATION:
  01_xerrors.go acts as the package error registry. It declares domain-specific
  error codes mapped tightly to corporate layout schemas, registering them into
  the global Go-Core xerrors ecosystem at initialization runtime.

  Design Constraints:
  - Every constant belonging to the error family must use the explicit xerrors.ErrorCode type.
  - All errors should follow the "EXXXX" naming standard, where "E" stands for Error
    and the first numerical digit points to the specific semantic family/group.
*/

const (
	// XERR_NONE serves as a local empty fallback marker to simplify assertion testing.
	XERR_NONE xerrors.ErrorCode = ""

	//
	//
	// X_ANCHOR_PKGCTX_START

	// XERR_PKGCTX defines the unique tracking string boundary separating this package from cross-domain collisions.
	XERR_PKGCTX xerrors.ErrorCode = "ERR_XTESTS"

	// XERR_PKGCTX_AUTH defines a specialized sub-context scope derived from XERR_PKGCTX.
	XERR_PKGCTX_AUTH xerrors.ErrorCode = "ERR_XTESTS_AUTH"

	// XERR_PKGCTX_VALID defines a specialized sub-context scope derived from XERR_PKGCTX.
	XERR_PKGCTX_VALID xerrors.ErrorCode = "ERR_XTESTS_VALID"

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
	// INVALID_COMMAND_NAME belongs to Family 1.
	// use in validation of command name not match with any available expected option.
	// Format expects: CTX, MSG, GIVEN, [error]
	INVALID_COMMAND_NAME xerrors.ErrorCode = "E1002"
	// ANOTHER_ERROR_CODE belongs to Family 1.
	// test position
	// Format expects: CTX, MSG, [error]
	ANOTHER_ERROR_CODE xerrors.ErrorCode = "E1003"
	// ONE_MORE_TIME_ON_FAMILY_ONE belongs to Family 1.
	// i beleave this time will work
	// Format expects: CTX, MSG, [error]
	ONE_MORE_TIME_ON_FAMILY_ONE xerrors.ErrorCode = "E1004"
	// YYYHHH belongs to Family 1.
	// yyyhhh
	// Format expects: CTX, MSG, YYYHHH, [error]
	YYYHHH xerrors.ErrorCode = "E1005"

	// ============================================================================
	// === FAMILY: 2 | TITLE: LLM ERRORS FAMILY
	// ============================================================================

	// INVALID_LLMRESPONSE belongs to Family 2.
	// use when llm connected give back an invalid response
	// Format expects: CTX, MSG, HASH, [error]
	INVALID_LLMRESPONSE xerrors.ErrorCode = "E2001"
	// ANOTHER_BRICK_IN_THE_WALL belongs to Family 2.
	// fff
	// Format expects: CTX, MSG, EEE, [error]
	ANOTHER_BRICK_IN_THE_WALL xerrors.ErrorCode = "E2002"
	// HGFDFSDS belongs to Family 2.
	// efdgfd
	// Format expects: CTX, MSG, FERWFVFFHG, [error]
	HGFDFSDS xerrors.ErrorCode = "E2003"

	// ============================================================================
	// === FAMILY: 3 | TITLE: TESTES
	// ============================================================================

	// TESTES belongs to Family 3.
	// testes
	// Format expects: CTX, MSG, TESTES, [error]
	TESTES xerrors.ErrorCode = "E3001"
	// UTFUTF belongs to Family 3.
	// utfutf
	// Format expects: CTX, MSG, UTFUTF, [error]
	UTFUTF xerrors.ErrorCode = "E3002"
	// QQQWWW belongs to Family 3.
	// qqqwww
	// Format expects: CTX, MSG, QQQWWW, [error]
	QQQWWW xerrors.ErrorCode = "E3003"
	// OOOOOO belongs to Family 3.
	// oooooo
	// Format expects: CTX, MSG, OOOOOO, [error]
	OOOOOO xerrors.ErrorCode = "E3004"

	// ============================================================================
	// === FAMILY: 4 | TITLE: EEEWWW
	// ============================================================================

	// EEEWWW belongs to Family 4.
	// eeewww
	// Format expects: CTX, MSG, EEEWWW, [error]
	EEEWWW xerrors.ErrorCode = "E4001"
	// WWWEEE belongs to Family 4.
	// wwweee
	// Format expects: CTX, MSG, WWWEEE, [error]
	WWWEEE xerrors.ErrorCode = "E4002"

	// ============================================================================
	// === FAMILY: 5 | TITLE: RRREEE
	// ============================================================================

	// RRREEE belongs to Family 5.
	// rrreee
	// Format expects: CTX, MSG, RRREEE, [error]
	RRREEE xerrors.ErrorCode = "E5001"
	// EEERRR belongs to Family 5.
	// eeerrr
	// Format expects: CTX, MSG, EEERRR, [error]
	EEERRR xerrors.ErrorCode = "E5002"

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
		"unexpected internal xtests error encountered",
		"",
		[]string{},
	),
	INVALID_COMMAND_NAME: xerrors.NewMetaMessage(
		"received invalid command name",
		"",
		[]string{"GIVEN"},
	),
	ANOTHER_ERROR_CODE: xerrors.NewMetaMessage(
		"Another error code",
		"",
		[]string{},
	),
	ONE_MORE_TIME_ON_FAMILY_ONE: xerrors.NewMetaMessage(
		"one more error inserted for family one",
		"",
		[]string{},
	),
	YYYHHH: xerrors.NewMetaMessage(
		"yyyhhh",
		"",
		[]string{"YYYHHH"},
	),

	// ============================================================================
	// === FAMILY: 2
	// ============================================================================

	INVALID_LLMRESPONSE: xerrors.NewMetaMessage(
		"invalid llm response",
		"",
		[]string{"HASH"},
	),
	ANOTHER_BRICK_IN_THE_WALL: xerrors.NewMetaMessage(
		"sss",
		"",
		[]string{"EEE"},
	),
	HGFDFSDS: xerrors.NewMetaMessage(
		"hgfbfdfde",
		"",
		[]string{"FERWFVFFHG"},
	),

	// ============================================================================
	// === FAMILY: 3
	// ============================================================================

	TESTES: xerrors.NewMetaMessage(
		"testes",
		"",
		[]string{"TESTES"},
	),
	UTFUTF: xerrors.NewMetaMessage(
		"utfutf",
		"",
		[]string{"UTFUTF"},
	),
	QQQWWW: xerrors.NewMetaMessage(
		"qqqwww",
		"",
		[]string{"QQQWWW"},
	),
	OOOOOO: xerrors.NewMetaMessage(
		"oooooo",
		"",
		[]string{"OOOOOO"},
	),

	// ============================================================================
	// === FAMILY: 4
	// ============================================================================

	EEEWWW: xerrors.NewMetaMessage(
		"eeewww",
		"",
		[]string{"EEEWWW"},
	),
	WWWEEE: xerrors.NewMetaMessage(
		"wwweee",
		"",
		[]string{"WWWEEE"},
	),

	// ============================================================================
	// === FAMILY: 5
	// ============================================================================

	RRREEE: xerrors.NewMetaMessage(
		"rrreee",
		"",
		[]string{"RRREEE"},
	),
	EEERRR: xerrors.NewMetaMessage(
		"eeerrr",
		"",
		[]string{"EEERRR"},
	),

	// X_ANCHOR_REGISTRY_END
	//
	//
}

func init() {
	// Automatically register local errors into the centralized tracking engine upon instantiation
	xerrors.RegisterDomainErrors(XERR_PKGCTX, xerrorLocalRegistry)
}

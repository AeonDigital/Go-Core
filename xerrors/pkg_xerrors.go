package xerrors

import "sync"

// ErrorCode defines a domain-specific string representation for error classification.
type ErrorCode string

// xerrorMapRegistry centralizes specialized error metadata definitions.
// Converted to sync.Map to ensure lock-free concurrent reads during application uptime.
var xerrorMapRegistry sync.Map // Internally stores map[ErrorCode]MetaMessage

// xerrorMapStringToErrorCode maps a unified string context back to its ErrorCode constant.
var xerrorMapStringToErrorCode sync.Map // Internally stores map[string]ErrorCode

const (
	XERR_NONE   ErrorCode = ""
	XERR_PKGCTX ErrorCode = "ERR_XERR"

	// XERR_UNKNOWN serves as the general fallback categorization for untracked exceptions.
	XERR_UNKNOWN ErrorCode = "E0001"

	/*
	   THE XERR FAMILY CONCEPT (POLYMORPHIC STRUCTURAL ERROR TOKENS)

	   The XERR_ constants serve as non-formatting, declarative tokens designed exclusively
	   to trigger the internal polymorphic parsing engine within the NewError400 factory.

	   Philosophy, Structural Uniformity & Visual Predictability:
	   Traditional Go formatting layout strings (e.g., via fmt.Errorf) require strict compilation-time
	   matching of verb quantities (%s, %w), which introduces human friction, parameter slippage,
	   and static analysis warnings (go vet). The XERR family bypasses these constraints by serving
	   as unified semantic identifiers mapped to strict corporate layout signatures.

	   To maximize performance and optimize high-throughput user validation tracking, this
	   architecture encapsulates failures under the lightweight IError400 interface, bypassing
	   expensive runtime stack inspection while enforcing clean transport-layer code extraction.

	   Architecture, Namespace Isolation & Positional Robustness:
	   Error categorization is dynamically split across bracket-enclosed visual tracking blocks
	   ([CTX], [MSG], and sequential [extraTags] such as FIELD, VALUE, EXPECTED, or RULES). Packages
	   can seamlessly extend this registry via RegisterDomainErrors using a package context token
	   (PKGCTX) to prevent cross-domain key collisions (e.g., "PKGCTX:CODE").

	   When intercepted inside NewError400, the buildMask formatting engine automatically extracts
	   any trailing native Go root cause error contract behind a double-colon boundary (:: [ERR: %w])
	   to protect standard library errors.Unwrap mechanics. The engine then safely normalizes the
	   underlying slice capacity to eliminate index out of range panic states. Missing human-readable
	   messages are backfilled with corporate map defaults, and omitted extra parameters are safely
	   substituted with the mathematical empty set marker ("ø") to guarantee unbreakable layout
	   predictability across system console output blocks, grepping tools, and log aggregators.
	*/

	// ============================================================================
	// GROUP 1: PRESENCE, EXISTENCE, AND NULLITY VALIDATIONS
	// Shared Layout Structure Token: [CTX: %v][MSG: %v][FIELD: %v][extraTags...] :: [ERR: %w]
	// ============================================================================

	// XERR_FIELD_REQUIRED belongs to Group 1 (Presence & Nullity).
	// Targets missing parameters, empty payloads, or fields whose total absence matches an 'undefined' state.
	// Format expects: CTX, MSG, FIELD, [error]
	XERR_FIELD_REQUIRED ErrorCode = "E1001"

	// XERR_NIL_NOT_ALLOWED belongs to Group 1 (Presence & Nullity).
	// Targets cases where a field or reference pointer is explicitly supplied but contains a forbidden nil value.
	// Format expects: CTX, MSG, FIELD, [error]
	XERR_NIL_NOT_ALLOWED ErrorCode = "E1002"

	// XERR_EMPTY_NOT_ALLOWED belongs to Group 1 (Presence & Nullity).
	// Targets fields that are allocated and non-nil, but their textual contents resolve to an empty string ("").
	// Format expects: CTX, MSG, FIELD, [error]
	XERR_EMPTY_NOT_ALLOWED ErrorCode = "E1003"

	// XERR_ZERO_NOT_ALLOWED belongs to Group 1 (Presence & Nullity).
	// Targets scenarios where numeric primitives, lengths, or uninitialized value types resolve to a forbidden zero state (0).
	// Format expects: CTX, MSG, FIELD, [error]
	XERR_ZERO_NOT_ALLOWED ErrorCode = "E1004"

	// XERR_ALREADY_EXISTS belongs to Group 1 (Presence & Nullity).
	// Targets uniqueness constraint violations where a valid field value cannot be accepted because it duplicates an active record.
	// Format expects: CTX, MSG, FIELD, VALUE, [error]
	XERR_ALREADY_EXISTS ErrorCode = "E1005"

	// XERR_NOT_FOUND belongs to Group 1 (Presence & Nullity).
	// Targets cases where a perfectly valid field lookup identifier fails to map to an active resource or file path.
	// Format expects: CTX, MSG, FIELD, TGT, [error]
	XERR_NOT_FOUND ErrorCode = "E1006"

	// XERR_PERMISSION_DENIED belongs to Group 1 (Presence & Nullity).
	// Targets security contract breaches where the application lacks the OS credentials, RBAC tokens,
	// or read/write privileges required to interact with the target resource.
	// Format expects: CTX, MSG, FIELD, TGT, [error]
	XERR_PERMISSION_DENIED ErrorCode = "E1007"

	// XERR_RESOURCE_UNAVAILABLE belongs to Group 1 (Presence & Nullity).
	// Targets IO blockages, hardware failures, timeout sequences, or networking disruptions
	// that prevent communication with an otherwise structurally valid target endpoint or file stream.
	// Format expects: CTX, MSG, FIELD, TGT, [error]
	XERR_RESOURCE_UNAVAILABLE ErrorCode = "E1008"

	// XERR_RESOURCE_CORRUPTED belongs to Group 1 (Presence & Nullity).
	// Targets integrity structural failures where the resource (file, payload, or state)
	// is physically present but its inner bytes break syntax rules, cryptographic checksums, or validation schemas.
	// Format expects: CTX, MSG, FIELD, TGT, [error]
	XERR_RESOURCE_CORRUPTED ErrorCode = "E1009"

	// XERR_RESOURCE_NOT_FOUND belongs to Group 1 (Presence & Nullity).
	// Targets system I/O, file systems, or directories that do not exist at the specified target path.
	// Format expects: CTX, MSG, FIELD, TGT, [error]
	XERR_RESOURCE_NOT_FOUND ErrorCode = "E1010"

	// ============================================================================
	// GROUP 2: NUMERIC, BOUNDARY, AND LIMIT VALIDATIONS
	// Shared Layout Structure Token: [CTX: %v][MSG: %v][FIELD: %v][VALUE: %v][RULES: %v] :: [ERR: %w]
	// ============================================================================

	// XERR_INVALID_VALUE belongs to Group 2 (Numeric & Boundaries).
	// General fallback for values that satisfy basic structural parsing but fail specialized domain business rules.
	// Format expects: CTX, MSG, FIELD, VALUE, RULES, [error]
	XERR_INVALID_VALUE ErrorCode = "E2001"

	// XERR_INVALID_VALUE_GT_ZERO belongs to Group 2 (Numeric & Boundaries).
	// Enforces that an evaluated mathematical property must be strictly greater than zero (> 0).
	// Format expects: CTX, MSG, FIELD, VALUE, [RULES], [error]
	XERR_INVALID_VALUE_GT_ZERO ErrorCode = "E2002"

	// XERR_INVALID_VALUE_GE_ZERO belongs to Group 2 (Numeric & Boundaries).
	// Enforces that an evaluated mathematical property must be greater than or equal to zero (>= 0).
	// Format expects: CTX, MSG, FIELD, VALUE, [RULES], [error]
	XERR_INVALID_VALUE_GE_ZERO ErrorCode = "E2003"

	// XERR_INVALID_VALUE_LT_ZERO belongs to Group 2 (Numeric & Boundaries).
	// Enforces that an evaluated mathematical property must be strictly less than zero (< 0).
	// Format expects: CTX, MSG, FIELD, VALUE, [RULES], [error]
	XERR_INVALID_VALUE_LT_ZERO ErrorCode = "E2004"

	// XERR_INVALID_VALUE_LE_ZERO belongs to Group 2 (Numeric & Boundaries).
	// Enforces that an evaluated mathematical property must be less than or equal to zero (<= 0).
	// Format expects: CTX, MSG, FIELD, VALUE, [RULES], [error]
	XERR_INVALID_VALUE_LE_ZERO ErrorCode = "E2005"

	// XERR_INVALID_VALUE_OUT_OF_RANGE belongs to Group 2 (Numeric & Boundaries).
	// Enforces that numbers, calendar dates, or generic offsets must stay enclosed within explicit low-high thresholds.
	// Format expects: CTX, MSG, FIELD, VALUE, [RULES], [error]
	XERR_INVALID_VALUE_OUT_OF_RANGE ErrorCode = "E2006"

	// XERR_SELECTION_LIMIT_EXCEEDED belongs to Group 2 (Numeric & Boundaries).
	// Targets scenarios where the number of selected items breaks cardinality boundaries or exceeds the maximum quantity constraints.
	// Format expects: CTX, MSG, FIELD, OPT, COUNT, LIMIT, [error]
	XERR_SELECTION_LIMIT_EXCEEDED ErrorCode = "E2007"

	// ============================================================================
	// GROUP 3: STRUCTURE, TYPING, AND CHOICE VALIDATIONS
	// Shared Layout Structure Token: [CTX: %v][MSG: %v][FIELD: %v][GIVEN/VALUE: %v][EXPECTED/OPTIONS: %v] :: [ERR: %w]
	// ============================================================================

	// XERR_INVALID_FORMAT belongs to Group 3 (Structure & Choices).
	// Targets syntax anomalies where string shapes break regex validations, structural encoding, or lexical requirements.
	// Format expects: CTX, MSG, FIELD, GIVEN, EXPECTED, [error]
	XERR_INVALID_FORMAT ErrorCode = "E3001"

	// XERR_INVALID_FORMAT_MARSHAL belongs to Group 3 (Structure & Choices).
	// Targets serialization anomalies where structural objects or typed domain entities fail to transform into target data shapes.
	// Format expects: CTX, MSG, FIELD, GIVEN, EXPECTED, [error]
	XERR_INVALID_FORMAT_MARSHAL ErrorCode = "E3002"

	// XERR_INVALID_FORMAT_UNMARSHAL belongs to Group 3 (Structure & Choices).
	// Targets deserialization anomalies where raw payload inputs break type structural specifications during structural unmarshaling.
	// Format expects: CTX, MSG, FIELD, GIVEN, EXPECTED, [error]
	XERR_INVALID_FORMAT_UNMARSHAL ErrorCode = "E3003"

	// XERR_INVALID_FORMAT_PARSE belongs to Group 3 (Structure & Choices).
	// Targets conversion anomalies where string primitives or unstructured slices fail lexical conversion into valid data primitives.
	// Format expects: CTX, MSG, FIELD, GIVEN, EXPECTED, [error]
	XERR_INVALID_FORMAT_PARSE ErrorCode = "E3004"

	// XERR_INVALID_TYPE belongs to Group 3 (Structure & Choices).
	// Targets type mismatch exceptions triggered during interface assertions, reflection mapping, or payload unmarshaling.
	// Format expects: CTX, MSG, FIELD, VALUE, EXPECTED_TYPE, [error]
	XERR_INVALID_TYPE ErrorCode = "E3005"

	// XERR_INVALID_OPTION belongs to Group 3 (Structure & Choices).
	// Targets invalid parameters outside a restrictive list of valid options or mutual exclusivity boundary contract violations.
	// Format expects: CTX, MSG, FIELD, OPT, OPTIONS, [error]
	XERR_INVALID_OPTION ErrorCode = "E3006"

	// XERR_MUTUAL_EXCLUSIVITY_VIOLATION belongs to Group 3 (Structure & Choices).
	// Targets structural contract breaches where choosing a specific field or option strictly invalidates the co-existence of others.
	// Format expects: CTX, MSG, FIELD, OPT, OPTIONS, [error]
	XERR_MUTUAL_EXCLUSIVITY_VIOLATION ErrorCode = "E3007"

	// XERR_ASYMMETRIC_SIZES belongs to Group 3 (Structure & Choices).
	// Targets structural contract breaches where interdependent collections fail to match linear sequence length.
	// Format expects: CTX, MSG, FIELDS, [error]
	XERR_ASYMMETRIC_SIZES ErrorCode = "E3008"

	// ============================================================================
	// GROUP 4: GENERIC OPERATIONAL FAILURE FALLBACKS
	// Shared Layout Structure Token: [CTX: %v][MSG: %v][FIELDS/DATA: %v] :: [ERR: %w]
	// ============================================================================

	// XERR_UNEXPECTED_FAIL belongs to Group 4 (Generic Operational Fallbacks).
	// Targets severe, non-deterministic system breakdowns, unmapped runtime panic states, or logic violations
	// that breach system stability invariants. Designed to act as a global strategic catch-all mechanism.
	// Format expects: CTX, MSG, DATA, [error]
	XERR_UNEXPECTED_FAIL ErrorCode = "E4001"

	// XERR_OPERATION_FAILED belongs to Group 4 (Generic Operational Fallbacks).
	// Targets state machine disruptions, unexecutable business commands, or processing
	// sequences that fail to complete their logic due to internal routine faults.
	// Format expects: CTX, MSG, FIELD, DATA, [error]
	XERR_OPERATION_FAILED ErrorCode = "E4002"
)

// xerrorDomainMapRegistry centralizes the core validation error metadata block and default
// corporate layout mapping definitions specific to the framework's own runtime context.
var xerrorDomainMapRegistry = map[ErrorCode]MetaMessage{

	// ============================================================================
	// GROUP 1: PRESENCE, EXISTENCE, AND NULLITY VALIDATIONS
	// ============================================================================

	XERR_FIELD_REQUIRED: {
		message:   "field required",
		extraTags: []string{"FIELD"},
	},
	XERR_NIL_NOT_ALLOWED: {
		message:   "nil pointer value not allowed",
		extraTags: []string{"FIELD"},
	},
	XERR_EMPTY_NOT_ALLOWED: {
		message:   "empty string value not allowed",
		extraTags: []string{"FIELD"},
	},
	XERR_ZERO_NOT_ALLOWED: {
		message:   "zero numeric value not allowed",
		extraTags: []string{"FIELD"},
	},
	XERR_ALREADY_EXISTS: {
		message:   "duplicate restriction violated",
		extraTags: []string{"FIELD", "VALUE"},
	},
	XERR_NOT_FOUND: {
		message:   "target resource not found",
		extraTags: []string{"FIELD", "TGT"},
	},
	XERR_PERMISSION_DENIED: {
		message:   "resource access permission denied",
		extraTags: []string{"FIELD", "TGT"},
	},
	XERR_RESOURCE_UNAVAILABLE: {
		message:   "target resource currently unavailable",
		extraTags: []string{"FIELD", "TGT"},
	},
	XERR_RESOURCE_CORRUPTED: {
		message:   "target resource is corrupted",
		extraTags: []string{"FIELD", "TGT"},
	},
	XERR_RESOURCE_NOT_FOUND: {
		message:   "target system resource or path not found",
		extraTags: []string{"FIELD", "TGT"},
	},

	// ============================================================================
	// GROUP 2: NUMERIC, BOUNDARY, AND LIMIT VALIDATIONS
	// ============================================================================

	XERR_INVALID_VALUE: {
		message:   "invalid value",
		extraTags: []string{"FIELD", "VALUE", "RULES"},
	},
	XERR_INVALID_VALUE_GT_ZERO: {
		message:   "invalid numeric",
		fieldRule: "must be greater than zero (> 0)",
		extraTags: []string{"FIELD", "VALUE", "RULES"},
	},
	XERR_INVALID_VALUE_GE_ZERO: {
		message:   "invalid numeric",
		fieldRule: "must be greater than or equal to zero (>= 0)",
		extraTags: []string{"FIELD", "VALUE", "RULES"},
	},
	XERR_INVALID_VALUE_LT_ZERO: {
		message:   "invalid numeric",
		fieldRule: "must be less than zero (< 0)",
		extraTags: []string{"FIELD", "VALUE", "RULES"},
	},
	XERR_INVALID_VALUE_LE_ZERO: {
		message:   "invalid numeric",
		fieldRule: "must be less than or equal to zero (<= 0)",
		extraTags: []string{"FIELD", "VALUE", "RULES"},
	},
	XERR_INVALID_VALUE_OUT_OF_RANGE: {
		message:   "invalid numeric",
		fieldRule: "out of range",
		extraTags: []string{"FIELD", "VALUE", "RULES"},
	},
	XERR_SELECTION_LIMIT_EXCEEDED: {
		message:   "selection quantity limit exceeded",
		extraTags: []string{"FIELD", "OPT", "COUNT", "LIMIT"},
	},

	// ============================================================================
	// GROUP 3: STRUCTURE, TYPING, AND CHOICE VALIDATIONS
	// ============================================================================

	XERR_INVALID_FORMAT: {
		message:   "malformed syntax data structure",
		extraTags: []string{"FIELD", "GIVEN", "EXPECTED"},
	},
	XERR_INVALID_FORMAT_MARSHAL: {
		message:   "serialization failed; marshal",
		extraTags: []string{"FIELD", "GIVEN", "EXPECTED"},
	},
	XERR_INVALID_FORMAT_UNMARSHAL: {
		message:   "deserialization failed; unmarshal",
		extraTags: []string{"FIELD", "GIVEN", "EXPECTED"},
	},
	XERR_INVALID_FORMAT_PARSE: {
		message:   "parse failed",
		extraTags: []string{"FIELD", "GIVEN", "EXPECTED"},
	},
	XERR_INVALID_TYPE: {
		message:   "type mismatch restriction violated",
		extraTags: []string{"FIELD", "VALUE", "EXPECTED_TYPE"},
	},
	XERR_INVALID_OPTION: {
		message:   "invalid option selection",
		extraTags: []string{"FIELD", "OPT", "OPTIONS"},
	},
	XERR_MUTUAL_EXCLUSIVITY_VIOLATION: {
		message:   "mutual exclusivity violation (choose only one)",
		extraTags: []string{"FIELD", "OPT", "OPTIONS"},
	},
	XERR_ASYMMETRIC_SIZES: {
		message:   "asymmetric collection sizes",
		extraTags: []string{"FIELDS"},
	},

	// ============================================================================
	// GROUP 4: GENERIC OPERATIONAL FAILURE FALLBACKS
	// ============================================================================

	XERR_OPERATION_FAILED: {
		message:   "operation failed",
		extraTags: []string{"FIELD", "DATA"},
	},
	XERR_UNEXPECTED_FAIL: {
		message:   "unexpected failure",
		extraTags: []string{"DATA"},
	},
}

func init() {
	RegisterDomainErrors(XERR_PKGCTX, xerrorDomainMapRegistry)
}

// RegisterDomainErrors injects custom domain error configurations into the centralized core registry.
// It uses sync.Map capabilities to safely register codes even if called concurrently during uptime.
func RegisterDomainErrors(pkgCtx ErrorCode, customRegistry map[ErrorCode]MetaMessage) {
	for code, meta := range customRegistry {

		// Construct the unified corporate namespace pattern
		fullCodeStr := string(pkgCtx) + ":" + string(code)
		fullErrorCode := ErrorCode(fullCodeStr)

		// Integrity check using LoadOrStore to protect core framework boundaries from overwrites
		// If the key does not exist, it stores the value and returns (nil, false)
		_, loaded := xerrorMapRegistry.LoadOrStore(fullErrorCode, meta)
		if !loaded {
			xerrorMapStringToErrorCode.Store(fullCodeStr, fullErrorCode)
		}
	}
}

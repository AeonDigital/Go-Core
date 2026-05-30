package xerrors

// ErrorCode defines a domain-specific string representation for error classification.
type ErrorCode string

const (
	// ErrUnknown serves as the general fallback categorization for untracked exceptions.
	ErrUnknown ErrorCode = "UNKNOWN_ERROR"

	/*
	   THE XERR FAMILY CONCEPT (POLYMORPHIC STRUCTURAL ERROR TOKENS)

	   The XERR_ constants serve as non-formatting, declarative tokens designed exclusively
	   to trigger the internal polymorphic parsing engine within the NewErr factory.

	   Philosophy, Structural Uniformity & Visual Predictability:
	   Traditional Go formatting layout strings (e.g., via fmt.Errorf) require strict compilation-time
	   matching of verb quantities (%s, %w), which introduces human friction, parameter slippage,
	   and static analysis warnings (go vet). The XERR family bypasses these constraints by serving
	   as unified semantic identifiers.

	   To minimize cognitive friction, this architecture consolidates any variable validation target—be it
	   a web form field, a JSON property, an internal function argument, or a map key—under the uniform concept
	   of a "FIELD". The layout splits error categorization across four highly optimized, bracket-enclosed
	   visual tracking blocks ([FIELD], [VALUE/RULES], [OPT/OPTIONS], or [DATA]).

	   When intercepted inside NewErr, the engine dynamically dismantles, cleans, and reconstructs
	   the target format sequence according to the exact slice size of incoming arguments. Missing
	   contextual fields are safely backfilled using the mathematical empty set marker ("ø") to guarantee
	   unbreakable visual predictability across system console output blocks, grepping tools, and log aggregators.
	*/

	// ============================================================================
	// GROUP 1: PRESENCE, EXISTENCE, AND NULLITY VALIDATIONS
	// Shared Layout Structure Token: [CTX: %v][MSG: %v][FIELD: %v]::[ERR: %v]
	// ============================================================================

	// XERR_FIELD_REQUIRED belongs to Group 1 (Presence & Nullity).
	// Targets missing parameters, empty payloads, or fields whose total absence matches an 'undefined' state.
	// Format expects: CTX, MSG, FIELD, [error]
	XERR_FIELD_REQUIRED = "[XERR_FIELD_REQUIRED]"

	// XERR_NIL_NOT_ALLOWED belongs to Group 1 (Presence & Nullity).
	// Targets cases where a field or reference pointer is explicitly supplied but contains a forbidden nil value.
	// Format expects: CTX, MSG, FIELD, [error]
	XERR_NIL_NOT_ALLOWED = "[XERR_NIL_NOT_ALLOWED]"

	// XERR_EMPTY_NOT_ALLOWED belongs to Group 1 (Presence & Nullity).
	// Targets fields that are allocated and non-nil, but their textual contents resolve to an empty string ("").
	// Format expects: CTX, MSG, FIELD, [error]
	XERR_EMPTY_NOT_ALLOWED = "[XERR_EMPTY_NOT_ALLOWED]"

	// XERR_ZERO_NOT_ALLOWED belongs to Group 1 (Presence & Nullity).
	// Targets scenarios where numeric primitives, lengths, or uninitialized value types resolve to a forbidden zero state (0).
	// Format expects: CTX, MSG, FIELD, [error]
	XERR_ZERO_NOT_ALLOWED = "[XERR_ZERO_NOT_ALLOWED]"

	// XERR_ALREADY_EXISTS belongs to Group 1 (Presence & Nullity).
	// Targets uniqueness constraint violations where a valid field value cannot be accepted because it duplicates an active record.
	// Format expects: CTX, MSG, FIELD, [error]
	XERR_ALREADY_EXISTS = "[XERR_ALREADY_EXISTS]"

	// XERR_NOT_FOUND belongs to Group 1 (Presence & Nullity).
	// Targets cases where a perfectly valid field lookup identifier fails to map to an active resource or file path.
	// Format expects: CTX, MSG, FIELD, TGT, [error]
	XERR_NOT_FOUND = "[XERR_NOT_FOUND]"

	// XERR_PERMISSION_DENIED belongs to Group 1 (Presence & Nullity).
	// Targets security contract breaches where the application lacks the OS credentials, RBAC tokens,
	// or read/write privileges required to interact with the target resource.
	// Format expects: CTX, MSG, FIELD, TGT, [error]
	XERR_PERMISSION_DENIED = "[XERR_PERMISSION_DENIED]"

	// XERR_RESOURCE_UNAVAILABLE belongs to Group 1 (Presence & Nullity).
	// Targets IO blockages, hardware failures, timeout sequences, or networking disruptions
	// that prevent communication with an otherwise structurally valid target endpoint or file stream.
	// Format expects: CTX, MSG, FIELD, TGT, [error]
	XERR_RESOURCE_UNAVAILABLE = "[XERR_RESOURCE_UNAVAILABLE]"

	// ============================================================================
	// GROUP 2: NUMERIC, BOUNDARY, AND LIMIT VALIDATIONS
	// Shared Layout Structure Token: [CTX: %v][MSG: %v][FIELD: %v][VALUE: %v][RULES: %v]::[ERR: %v]
	// ============================================================================

	// XERR_INVALID_VALUE belongs to Group 2 (Numeric & Boundaries).
	// General fallback for values that satisfy basic structural parsing but fail specialized domain business rules.
	// Format expects: CTX, MSG, FIELD, VALUE, RULES, [error]
	XERR_INVALID_VALUE = "[XERR_INVALID_VALUE]"

	// XERR_INVALID_VALUE_GT_ZERO belongs to Group 2 (Numeric & Boundaries).
	// Enforces that a evaluated mathematical property must be strictly greater than zero (> 0).
	// Format expects: CTX, MSG, FIELD, VALUE, RULES, [error]
	XERR_INVALID_VALUE_GT_ZERO = "[XERR_INVALID_VALUE_GT_ZERO]"

	// XERR_INVALID_VALUE_GE_ZERO belongs to Group 2 (Numeric & Boundaries).
	// Enforces that a evaluated mathematical property must be greater than or equal to zero (>= 0).
	// Format expects: CTX, MSG, FIELD, VALUE, RULES, [error]
	XERR_INVALID_VALUE_GE_ZERO = "[XERR_INVALID_VALUE_GE_ZERO]"

	// XERR_INVALID_VALUE_LT_ZERO belongs to Group 2 (Numeric & Boundaries).
	// Enforces that a evaluated mathematical property must be strictly less than zero (< 0).
	// Format expects: CTX, MSG, FIELD, VALUE, RULES, [error]
	XERR_INVALID_VALUE_LT_ZERO = "[XERR_INVALID_VALUE_LT_ZERO]"

	// XERR_INVALID_VALUE_LE_ZERO belongs to Group 2 (Numeric & Boundaries).
	// Enforces that a evaluated mathematical property must be less than or equal to zero (<= 0).
	// Format expects: CTX, MSG, FIELD, VALUE, RULES, [error]
	XERR_INVALID_VALUE_LE_ZERO = "[XERR_INVALID_VALUE_LE_ZERO]"

	// XERR_INVALID_VALUE_OUT_OF_RANGE belongs to Group 2 (Numeric & Boundaries).
	// Enforces that numbers, calendar dates, or generic offsets must stay enclosed within explicit low-high thresholds.
	// Format expects: CTX, MSG, FIELD, VALUE, RULES, [error]
	XERR_INVALID_VALUE_OUT_OF_RANGE = "[XERR_INVALID_VALUE_OUT_OF_RANGE]"

	// XERR_SELECTION_LIMIT_EXCEEDED belongs to Group 2 (Numeric & Boundaries).
	// Targets scenarios where the number of selected items breaks cardinality boundaries or exceeds the maximum quantity constraints.
	// Format expects: CTX, MSG, FIELD, OPT, COUNT, LIMIT, [error]
	XERR_SELECTION_LIMIT_EXCEEDED = "[XERR_SELECTION_LIMIT_EXCEEDED]"

	// ============================================================================
	// GROUP 3: STRUCTURE, TYPING, AND CHOICE VALIDATIONS
	// Shared Layout Structure Token: [CTX: %v][MSG: %v][FIELD: %v][VALUE: %v][EXPECTED/OPTIONS: %v]::[ERR: %v]
	// ============================================================================

	// XERR_INVALID_FORMAT belongs to Group 3 (Structure & Choices).
	// Targets syntax anomalies where string shapes break regex validations, structural encoding, or lexical requirements.
	// Format expects: CTX, MSG, FIELD, EXPECTED, GIVEN, [error]
	XERR_INVALID_FORMAT = "[XERR_INVALID_FORMAT]"

	XERR_MSG_INVALID_FORMAT_MARSHAL   = "marshal failed"
	XERR_MSG_INVALID_FORMAT_UNMARSHAL = "unmarshal failed"
	XERR_MSG_INVALID_FORMAT_PARSE     = "parse failed"

	// XERR_INVALID_TYPE belongs to Group 3 (Structure & Choices).
	// Targets type mismatch exceptions triggered during interface assertions, reflection mapping, or payload unmarshaling.
	// Format expects: CTX, MSG, FIELD, VALUE, EXPECTED_TYPE, [error]
	XERR_INVALID_TYPE = "[XERR_INVALID_TYPE]"

	// XERR_INVALID_OPTION belongs to Group 3 (Structure & Choices).
	// Targets invalid parameters outside a restrictive list of valid options or mutual exclusivity boundary contract violations.
	// Format expects: CTX, MSG, FIELD, OPT, OPTIONS, [error]
	XERR_INVALID_OPTION = "[XERR_INVALID_OPTION]"

	// XERR_MUTUAL_EXCLUSIVITY_VIOLATION belongs to Group 3 (Structure & Choices).
	// Targets structural contract breaches where choosing a specific field or option strictly invalidates the co-existence of others.
	// Format expects: CTX, MSG, FIELD, OPT, OPTIONS, [error]
	XERR_MUTUAL_EXCLUSIVITY_VIOLATION = "[XERR_MUTUAL_EXCLUSIVITY_VIOLATION]"

	// XERR_ASYMMETRIC_SIZES belongs to Group 3 (Structure & Choices).
	// Targets structural contract breaches where interdependent collections fail to match linear sequence lengths.
	// Format expects: CTX, MSG, FIELDS, [error]
	XERR_ASYMMETRIC_SIZES = "[XERR_ASYMMETRIC_SIZES]"

	// ============================================================================
	// GROUP 4: GENERIC OPERATIONAL FAILURE FALLBACKS
	// Shared Layout Structure Token: [CTX: %v][MSG: %v][FIELD: %v][DATA: %v]::[ERR: %v]
	// ============================================================================

	// XERR_INVALID_DATA belongs to Group 4 (Generic Operational Fallbacks).
	// Standard layout fallback intended for wide operational dumps or key-value diagnostic pairs formatted via DataMsg.
	// Format expects: CTX, MSG, FIELD, DATA, [error]
	XERR_INVALID_DATA = "[XERR_INVALID_DATA]"
)

// errorMetadata encapsulates static corporate fallback strings and semantic tagging rules
// required to build standardized, predictable error formatting structures.
type errorMetadata struct {
	defaultMessage string
	defaultRule    string
	extraTags      []string
}

// DetailedError standardizes the functional execution capabilities required to manage rich,
// decoupled diagnostic payloads across network borders, microservices, and system boundaries.
type DetailedError interface {
	// Code returns the categorical string constant domain mapping for this failure event.
	Code() ErrorCode

	// Component pinpoints the specific architectural layer or tracking package path.
	Component() string

	// error ensures native alignment with Go standard library errors handling semantics.
	error

	// WithCallerSkip allows dynamically adjusting the runtime stack frame skip level to recompute
	// the component identifier. This is essential when NewError is wrapped inside custom factory
	// utilities or logging helper blocks, ensuring the framework pinpoints the true origin of the failure.
	// It returns the updated DetailedError to support fluent API call chaining.
	WithCallerSkip(skip int) DetailedError

	// Message returns the readable operational summary describing what went wrong.
	Message() string

	// Data returns raw debugging contexts such as data dumps or corrupted payloads.
	Data() string

	// DebugError computes a highly readable formatted diagnostic block string ready for loggers or consoles.
	DebugError() string
}

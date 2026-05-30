package xerrors

import (
	"errors"
	"testing"
)

func TestFormatValueWithQuotes_DirectNilAccess(t *testing.T) {
	// Como este teste está no 'package xerrors', ele ignora as travas do MsgOptions
	// e consegue chamar a função interna diretamente com um nil puro.
	res := formatValueWithQuotes(nil)

	if res != "nil" {
		t.Errorf("expected formatValueWithQuotes(nil) to return 'nil', got %q", res)
	}
}

func TestBuildMask_DirectDirectCoverage(t *testing.T) {
	rootCause := errors.New("low level io timeout")

	tests := []struct {
		name           string
		meta           MetaMessage
		inputArgs      []any
		expectedLayout string
		expectedArgLen int
		checkLastError bool
	}{
		{
			name: "Scenario 1: Complete arguments with root error present (Group 2 Rules)",
			meta: MetaMessage{
				message:   "numeric constraint broken",
				fieldRule: "must be > 0",
				extraTags: []string{"FIELD", "VALUE", "RULES"},
			},
			inputArgs: []any{
				"FLOW-CTX",            // Index 0: CTX
				"custom msg override", // Index 1: MSG
				"age_field",           // Tag 1: FIELD
				-5,                    // Tag 2: VALUE
				"",                    // Tag 3: RULES (Empty string forces fieldRule fallback)
				rootCause,             // Trailing root cause error
			},
			expectedLayout: "[CTX: %v][MSG: %v][FIELD: %v][VALUE: %v][RULES: %v] :: [ERR: %w]",
			expectedArgLen: 6, // CTX, MSG, FIELD, VALUE, RULES, error
			checkLastError: true,
		},
		{
			name: "Scenario 2: Omitted message and completely omitted fieldRule",
			meta: MetaMessage{
				message:   "fallback registry message",
				fieldRule: "must be >= 0",
				extraTags: []string{"FIELD", "VALUE", "RULES"},
			},
			inputArgs: []any{
				"FLOW-CTX", // Index 0: CTX only
			},
			// CORREÇÃO: A nova lógica gera os colchetes com %v e injeta a mensagem do mapa + ø nos argumentos
			expectedLayout: "[CTX: %v][MSG: %v][FIELD: %v][VALUE: %v][RULES: %v] :: [ERR: ø]",
			expectedArgLen: 5, // CTX, MSG (meta.message), FIELD (ø), VALUE (ø), RULES (ø)
			checkLastError: false,
		},
		{
			name: "Scenario 3: FieldRule completely omitted by developer but root error is passed",
			meta: MetaMessage{
				message:   "boundary failure",
				fieldRule: "must be less than zero (< 0)",
				extraTags: []string{"FIELD", "VALUE", "RULES"},
			},
			inputArgs: []any{
				"FLOW-CTX",      // Index 0: CTX
				"failed check",  // Index 1: MSG
				"balance_field", // Tag 1: FIELD
				150,             // Tag 2: VALUE
				// Note: Index 5 (RULES) is completely omitted here, leaving the error right after VALUE
				rootCause, // Trailing root cause error
			},
			// The engine must intercept the missing rule, inject it, and push the error to the correct final index
			expectedLayout: "[CTX: %v][MSG: %v][FIELD: %v][VALUE: %v][RULES: %v] :: [ERR: %w]",
			expectedArgLen: 6, // CTX, MSG, FIELD, VALUE, meta.fieldRule, error
			checkLastError: true,
		},
		{
			name: "Scenario 4: Valid trailing parameter that is NOT an error contract",
			meta: MetaMessage{
				message:   "simple informational failure",
				extraTags: []string{"FIELD"},
			},
			inputArgs: []any{
				"FLOW-CTX",
				"plain error message",
				"username",
				"this string is not a go error interface", // Trailing arg but not an error type
			},
			expectedLayout: "[CTX: %v][MSG: %v][FIELD: %v] :: [ERR: ø]",
			expectedArgLen: 4, // Keeps all arguments but maps the visual layout to ERR: ø
			checkLastError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLayout, gotArgs := buildMask(tt.meta, tt.inputArgs)

			// Assert that the generated format string template matches the corporate expectation perfectly
			if gotLayout != tt.expectedLayout {
				t.Errorf("buildMask() layout = %q, want %q", gotLayout, tt.expectedLayout)
			}

			// Assert that arguments slice size is correctly normalized (backfilled or sliced properly)
			if len(gotArgs) != tt.expectedArgLen {
				t.Errorf("buildMask() normalized arguments slice length = %d, want %d", len(gotArgs), tt.expectedArgLen)
			}

			// If the test case included a root cause error, verify it was re-attached at the final slot
			if tt.checkLastError {
				lastIdx := len(gotArgs) - 1
				if _, ok := gotArgs[lastIdx].(error); !ok {
					t.Errorf("Expected trailing normalized argument at index %d to be an error interface type, got %T", lastIdx, gotArgs[lastIdx])
				}
			}
		})
	}
}

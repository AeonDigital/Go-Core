package xerrors_test

import (
	"strings"
	"testing"

	"github.com/AeonDigital/Go-Core/xerrors"
)

// init registers a fake domain error to simulate an external package extension
// and guarantee 100% path coverage inside NewError400.
func init() {
	xerrors.RegisterDomainErrors(
		"ERR_MOCK",
		map[xerrors.ErrorCode]xerrors.MetaMessage{
			"E9999": xerrors.NewMetaMessage(
				"mock validation failed",
				"",
				[]string{"FIELD", "VALUE"},
			),
		},
	)
}

func TestNewError400_CoverageSuite(t *testing.T) {
	tests := []struct {
		name         string
		inputArgs    []any
		expectedCode xerrors.ErrorCode
		containsText []string
	}{
		{
			name:         "Path 1: Should return an empty message and XERR_NONE when args are omitted",
			inputArgs:    []any{},
			expectedCode: xerrors.XERR_NONE,
			containsText: []string{""},
		},
		{
			name: "Path 2: Should evaluate namespaced domain token when first two args are ErrorCodes",
			inputArgs: []any{
				xerrors.ErrorCode("ERR_MOCK"),
				xerrors.ErrorCode("E9999"),
				"TX-ID",             // CTX
				"",                  // MSG (triggers default fallback)
				"email_field",       // FIELD tag
				"invalid_payload_v", // VALUE tag
			},
			expectedCode: xerrors.ErrorCode("ERR_MOCK:E9999"),
			containsText: []string{"[CTX: TX-ID]", "[MSG: mock validation failed]", "[FIELD: email_field]", "[VALUE: invalid_payload_v]"},
		},
		{
			name: "Path 3: Should evaluate core framework token when only the first arg is an ErrorCode",
			inputArgs: []any{
				xerrors.XERR_FIELD_REQUIRED,
				"AUTH-FLOW", // CTX
				"required",  // MSG
				"password",  // FIELD tag
			},
			expectedCode: xerrors.ErrorCode("ERR_XERR:E1001"),
			containsText: []string{"[CTX: AUTH-FLOW]", "[MSG: required]", "[FIELD: password]"},
		},
		{
			name: "Path 4A: Should format string like standard Sprintf when first arg is plain text",
			inputArgs: []any{
				"the user metadata for id %d is invalid: %s",
				42,
				"malformed_email",
			},
			expectedCode: xerrors.XERR_NONE,
			containsText: []string{"the user metadata for id 42 is invalid: malformed_email"},
		},
		{
			name: "Path 4B: Should handle fallback gracefully when first arg is a single plain string",
			inputArgs: []any{
				"plain unformatted static text error",
			},
			expectedCode: xerrors.XERR_NONE,
			containsText: []string{"plain unformatted static text error"},
		},
		{
			name: "Path 4C: Should handle fallback gracefully when initial parameter is an unmapped type",
			inputArgs: []any{
				12345, // Not a string, nor an ErrorCode
				"extra data",
			},
			expectedCode: xerrors.XERR_NONE,
			containsText: []string{"12345", "extra data"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := xerrors.NewError400(tt.inputArgs...)

			if gotErr == nil {
				t.Fatal("Expected NewError400 to return a valid non-nil IError400 interface instance")
			}

			// Validate if the ErrorCode mapping contract was respected
			if gotErr.Code() != tt.expectedCode {
				t.Errorf("Code() = %q, want %q", gotErr.Code(), tt.expectedCode)
			}

			// Validate if the generated error message string contains the expected structural tokens
			errString := gotErr.Error()
			for _, textSegment := range tt.containsText {
				if !strings.Contains(errString, textSegment) {
					t.Errorf("Error() string %q was expected to contain segment %q but it did not", errString, textSegment)
				}
			}
		})
	}
}

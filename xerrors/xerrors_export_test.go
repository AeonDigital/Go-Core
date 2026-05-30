package xerrors

const (
	ErrCreateInstance        ErrorCode = "CREATE_INSTANCE_ERROR"
	ErrNetworkHTTPConnection ErrorCode = "NETWORK_HTTP_CONNECTION_ERROR"

	// XERR_TEST_COMPLEX serves exclusively as a high-complexity testing token to validate
	// all inner branches, loops, and backfill safeguards inside the universal mask engine.
	XERR_TEST_COMPLEX = "[XERR_TEST_COMPLEX]"
)

// RegisterTestTokenInRegistry injects our custom high-complexity test token into the
// active errorRegistry map, ensuring 100% test isolation from real production variables.
func RegisterTestTokenInRegistry() {
	errorRegistry[XERR_TEST_COMPLEX] = errorMetadata{
		defaultMessage: "test default message layout",
		defaultRule:    "test default rule requirement",
		extraTags:      []string{"TAG1", "TAG2", "TAG3"},
	}
}

// ParseRuntimeFuncNameForTest expõe a função privada de parsing para a suíte de testes de caixa-preta.
func ParseRuntimeFuncNameForTest(fullName string) (string, string, string) {
	return parseRuntimeFuncName(fullName)
}

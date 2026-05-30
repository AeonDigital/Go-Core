package xerrors

var (
	debugMode bool = false
)

// GetDebugMode returns the current status of the debug mode.
// When true, errors provide technical details; when false, they are user-friendly.
func GetDebugMode() bool {
	return debugMode
}

// EnableDebugMode enables technical error details, making logs and outputs
// more comprehensive for debugging purposes.
func EnableDebugMode() {
	debugMode = true
}

// DisableDebugMode disables technical error details, switching outputs
// to a user-friendly format suitable for end-users.
func DisableDebugMode() {
	debugMode = false
}

// ToggleDebugMode switches the current state of the debug mode
// (enables it if disabled, or disables it if enabled).
func ToggleDebugMode() {
	debugMode = !debugMode
}

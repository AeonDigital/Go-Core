package sterrinterfc

// CliError defines the structural behavioral contract for light, terminal-friendly errors.
type CliError interface {
	// error implements the native Go error interface.
	error

	// SetDevMessage overrides the current technical diagnostic text.
	SetDevMessage(format string, args ...any) CliError

	// SetUserMessage overrides the current end-user friendly instruction.
	SetUserMessage(format string, args ...any) CliError

	// AppendDevMessage concatenates a text payload into the current developer message.
	AppendDevMessage(format string, args ...any) CliError

	// AppendUserMessage concatenates a text payload into the current end-user message.
	AppendUserMessage(format string, args ...any) CliError

	// AppendLNDevMessage concatenates a text payload appending an explicit newline character at the end.
	AppendLNDevMessage(format string, args ...any) CliError

	// AppendLNUserMessage concatenates a text payload appending an explicit newline character at the end.
	AppendLNUserMessage(format string, args ...any) CliError

	// ClearDevMessage purges the developer message content, resetting it to empty string.
	ClearDevMessage() CliError

	// ClearUserMessage purges the end-user message content, resetting it to empty string.
	ClearUserMessage() CliError

	// WithDepth forces the engine to recalculate the runtime trace stack location using an offset.
	WithDepth(additionalDepth int) CliError

	// GetFunction returns the qualified target name tracking the package and functional scope.
	GetFunction() string

	// GetDevMessage returns the technical diagnostics text string.
	GetDevMessage() string

	// GetUserMessage returns the actionable text instruction designed for human end-users.
	GetUserMessage() string

	// HasDevMessage verifies if a technical message payload has been populated.
	HasDevMessage() bool

	// HasUserMessage verifies if a human-friendly instruction payload has been populated.
	HasUserMessage() bool

	// HasErrors checks if any descriptive message fields contain active content tracking failures.
	HasErrors() bool
}

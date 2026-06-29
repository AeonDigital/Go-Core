package xcli

import (
	"github.com/AeonDigital/Go-Core/sterr/pkg/sterr"
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xclifn"
)

// Router manages the CLI lifecycle by registering the root command,
// navigating the command tree, and triggering argument parsing and validation.
type Router struct {
	// Root is the entry point of the CLI application.
	Root *Command
}

// NewRouter initializes a new CLI router with a defined root command.
func NewRouter(root *Command) *Router {
	return &Router{Root: root}
}

// Run reads the raw terminal arguments, resolves the command path,
// validates the flags, checks for help triggers, and executes the final command hook.
//
// Arguments:
//   - rawArgs: The slice of strings representing terminal inputs (usually os.Args[1:]).
//
// Returns:
//   - error: Returns an sterr.CliError if a command is not found or flag parsing/validation fails.
func (r *Router) Run(rawArgs []string) error {
	if r.Root == nil {
		return sterr.New().
			SetMessage("root command is not registered")
	}

	currentCmd := r.Root
	argIndex := 0

	// Phase 1: Navigate the command tree hierarchy
	for argIndex < len(rawArgs) {
		arg := rawArgs[argIndex]

		if len(arg) > 0 && arg[0] == '-' {
			break
		}

		if arg == "help" || arg == "--help" || arg == "-h" {
			return currentCmd.TriggerHelp()
		}

		if nextCmd, exists := currentCmd.Subcommands[arg]; exists {
			currentCmd = nextCmd
			argIndex++
			continue
		}

		return sterr.New().
			SetMessage("unknown command: '%s' for scope '%s'", arg, currentCmd.Name)
	}

	for _, remainingArg := range rawArgs[argIndex:] {
		if remainingArg == "--help" || remainingArg == "-h" {
			return currentCmd.TriggerHelp()
		}
	}

	// Phase 2: Extract raw tokens from the clean flags boundary block
	flagTokens := rawArgs[argIndex:]
	rawFlagsMap, err := xclifn.ParseRawArgs(flagTokens)
	if err != nil {
		return err
	}

	// Phase 3: Hydrate and Validate Flag Values directly through the Command instance
	ctxValues, err := currentCmd.ValidateAndHydrateFlags(rawFlagsMap)
	if err != nil {
		return err
	}

	// Phase 4: Execute the command business logic hook
	if currentCmd.Run == nil {
		return currentCmd.TriggerHelp()
	}

	return currentCmd.Run(ctxValues)
}

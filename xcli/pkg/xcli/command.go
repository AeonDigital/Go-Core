package xcli

import (
	"fmt"

	"github.com/AeonDigital/Go-Core/sterr/pkg/sterr"
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xclifn"
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xclistruc"
)

// Command represents a node in the CLI command tree.
// Each command forms an isolated scope and executes its own business logic.
type Command struct {
	// Name is the string that triggers this command in the terminal.
	Name string

	// ShortDescription is a brief one-line summary used in general help listings.
	ShortDescription string

	// LongDescription is a detailed explanation shown when help is requested
	// specifically for this command. If left empty, ShortDescription will be used.
	LongDescription string

	// Flags is the list of exclusive options accepted strictly by this command.
	Flags []xclistruc.Flag

	// Subcommands holds the next layer of commands, indexed by their execution Name.
	// Due to context isolation, children do not inherit flags from their parents.
	Subcommands map[string]*Command

	// Run is the execution hook containing the command's business logic.
	// It receives the Context containing all parsed, typed, and validated flags.
	Run func(ctx *xclistruc.FlagValues) error
}

// ValidateAndHydrateFlags loops through registered constraints performing types translation and bounds enforcement.
//
// It checks flag mandatory presence, injects fallback defaults allocation blocks, and leverages
// the package-level central registry directory to resolve runtime type translations dynamically.
//
// Arguments:
//   - rawFlags: The extracted command line layout map pairing flag strings to text values.
//
// Returns:
//   - *xclistruc.FlagValues: A populated type-safe context directory containing translated Go instances.
//   - error: Returns an sterr.CliError capturing visual telemetry diagnostics if any validation path fails.
//
// Error & Panic Natures:
//   - Complex Errors: Fails immediately if a required token is missing from raw inputs.
//     Queries the global type registry throwing errors if a type metadata token is unmapped.
//     Performs zombie variable tracking blocks, failing if unregistered flag tokens are found.
func (c *Command) ValidateAndHydrateFlags(
	rawFlags map[string]string,
) (
	*xclistruc.FlagValues,
	error,
) {
	ctxValues := xclistruc.NewFlagValues()

	// Track which keys from terminal inputs were actually processed to detect unmapped items later
	processedRawKeys := make(map[string]bool)

	for _, spec := range c.Flags {
		// Enforce visual notation representation matching user input options
		flagLabel := "--" + spec.LongName

		// Locate the raw value token inside the parsed terminal mapping block
		var rawValue string
		var found bool

		if val, ok := rawFlags[spec.LongName]; ok {
			rawValue = val
			found = true
			processedRawKeys[spec.LongName] = true
		} else if spec.ShortName != "" {
			if val, ok := rawFlags[spec.ShortName]; ok {
				rawValue = val
				found = true
				processedRawKeys[spec.ShortName] = true
			}
		}

		// 1. Mandatory Presence enforcement check
		if spec.Required && !found {
			return nil, sterr.New().
				SetMessage("[ERR] %s : required", flagLabel)
		}

		// 2. Default values allocation fallback logic
		if !found {
			if spec.DefaultValue != nil {
				ctxValues.SetInternalValue(spec.LongName, spec.DefaultValue)
			}
			continue
		}

		// 3. Dynamic Type-Safe Parsing and Unified Multi-Layered Validation Engine Dispatch via Global Registry
		parserEngine, exists := GlobalTypeRegistry[spec.Type]
		if !exists {
			return nil, sterr.New().
				SetMessage("unsupported flag type: '%s'", spec.Type)
		}

		// Execute the complete type transformation and domain boundary enforcement pipeline
		typedValue, parseErr := parserEngine.ParseAndValidate(flagLabel, rawValue, spec)
		if parseErr != nil {
			return nil, parseErr
		}

		ctxValues.SetInternalValue(spec.LongName, typedValue)
	}

	// 4. Security Check: Block unregistered positional arguments or zombie variables sent by mistake
	for k := range rawFlags {
		if !processedRawKeys[k] {
			return nil, sterr.New().
				SetMessage("invalid argument: unrecognized flag identifier token provided: '%s'", k)
		}
	}

	return ctxValues, nil
}

// TriggerHelp intercepts the flow and renders the automatic command documentation.
//
// It evaluates operational descriptions metadata fields and directly serializes aligned visual
// tables mapping usage definitions, child subcommands, and flags guidelines to standard stdout tracks.
//
// Returns:
//   - error: Returns a structured error tracking instance if writing stdout triggers telemetry blocks.
func (c *Command) TriggerHelp() error {
	// Determine the best description to display based on availability
	description := c.LongDescription
	if description == "" {
		description = c.ShortDescription
	}

	// Render Command Usage and Description
	xclifn.PrintStdout("Usage:")
	xclifn.PrintStdout("  %s [subcommand] [--flags]\n", c.Name)

	if description != "" {
		xclifn.PrintStdout("Description:")
		xclifn.PrintStdout("  %s\n", description)
	}

	// Render Available Subcommands if the command has children
	if len(c.Subcommands) > 0 {
		xclifn.PrintStdout("Available Subcommands:")
		for _, sub := range c.Subcommands {
			xclifn.PrintStdout("  %-15s %s", sub.Name, sub.ShortDescription)
		}
		xclifn.PrintStdout("")
	}

	// Render Flags Configuration if the command accepts any options
	if len(c.Flags) > 0 {
		xclifn.PrintStdout("Flags:")
		for _, flag := range c.Flags {
			// Format the flag triggers (e.g., "--output, -out")
			syntaxBuilder := "--" + flag.LongName
			if flag.ShortName != "" {
				syntaxBuilder += ", -" + flag.ShortName
			}

			// Format the type and required/metadata properties
			metaBuilder := string(flag.Type)
			if flag.Required {
				metaBuilder += " (required)"
			} else if flag.DefaultValue != nil {
				metaBuilder += fmt.Sprintf(" (default: %v)", flag.DefaultValue)
			}

			// Determine the flag description fallback
			flagDesc := flag.LongDescription
			if flagDesc == "" {
				flagDesc = flag.ShortDescription
			}

			// Print structured, aligned row for high human scannability
			xclifn.PrintStdout("  %-25s %-25s %s", syntaxBuilder, "["+metaBuilder+"]", flagDesc)
		}
		xclifn.PrintStdout("")
	}

	return nil
}

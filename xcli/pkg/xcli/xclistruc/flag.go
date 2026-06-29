package xclistruc

import (
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xcliconstt"
)

// Flag defines a command-line option, its data type, and validation constraints.
type Flag struct {
	// LongName is the full identifier invoked via double dashes (e.g., "output" for --output).
	LongName string

	// ShortName is the shorthand identifier invoked via a single dash (e.g., "out" for -out).
	// It supports up to 3 characters. Flag bundling/agglutination is not supported.
	ShortName string

	// ShortDescription is a brief explanation of the flag purpose used in the help menu.
	ShortDescription string

	// LongDescription is a detailed explanation of the flag usage (optional, defaults to ShortDescription).
	LongDescription string

	// Required determines if the command execution must fail if this flag is absent.
	Required bool

	// DefaultValue is the fallback value applied if the flag is omitted. Ignored if Required is true.
	DefaultValue any

	// Type specifies the data format for conversion and validation (e.g., "string", "int", "email", "[]int").
	Type xcliconstt.FlagType

	//
	//
	// --- Numeric and Temporal Validations (int, float, duration, date, time, datetime) ---

	// Min defines the inclusive lower bound allowed for quantifiable types. Ignored if nil.
	Min any

	// Max defines the inclusive upper bound allowed for quantifiable types. Ignored if nil.
	Max any

	//
	//
	// --- String Validations ---

	// MinLength defines the minimum number of characters required for string types. Ignored if nil.
	MinLength *int

	// MaxLength defines the maximum number of characters allowed for string types. Ignored if nil.
	MaxLength *int

	// Regex is a regular expression pattern that the raw string value must match. Ignored if empty.
	Regex string

	//
	//
	// --- Array Validations ([]type) ---

	// MinItems defines the minimum number of elements required when the flag is an array. Ignored if nil.
	MinItems *int

	// MaxItems defines the maximum number of elements allowed when the flag is an array. Ignored if nil.
	MaxItems *int

	//
	//
	// --- System and Network ---

	// MustExist enforces that the filepath or dirpath must already exist on disk.
	// If false, the validator will only check syntax and write permissions.
	MustExist bool

	// Access defines the permission level (read, write, readwrite) required for the resource.
	// If left empty, the router will skip permission checks and only validate syntax/existence.
	Access xcliconstt.AccessMode
}

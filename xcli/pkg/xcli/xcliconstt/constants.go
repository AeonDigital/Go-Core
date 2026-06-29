package xcliconstt

/*
	ARCHITECTURE & SCOPE LIMITATION:
	constants.go centralizes immutable, read-only values and global literals
	exclusively required to configure or support the package logic.

	Design Constraints:
	- Only truly constant and stateless values (string, int, time.Duration primitives)
		are allowed here.
	- Never declare mutable global variables or pointers inside this package.
	- If a distinct domain domain context (e.g., custom error codes or CLI flag keys)
		grows large, split those constants into a separate, context-named file here.
*/

// Insert configuration constants, structs and config functions below.

//
//
//

// FlagType defines the strict set of data formats supported by the router
// for automatic parsing, type conversion, and validation.
type FlagType string

const (
	//
	// --- Primitive Types ---

	TypeString FlagType = "string"
	TypeInt    FlagType = "int"
	TypeFloat  FlagType = "float"
	TypeBool   FlagType = "bool"

	//
	// --- Structured Data ---

	TypeJSON     FlagType = "json"
	TypeDuration FlagType = "duration"
	TypeDate     FlagType = "date"
	TypeTime     FlagType = "time"
	TypeDateTime FlagType = "datetime"
	TypeEmail    FlagType = "email"

	//
	// --- System and Network Validations ---

	TypePath        FlagType = "path"
	TypeFilename    FlagType = "filename"
	TypeFilepath    FlagType = "filepath"
	TypeDirname     FlagType = "dirname"
	TypeDirpath     FlagType = "dirpath"
	TypeURL         FlagType = "url"
	TypeFullURL     FlagType = "fullurl"
	TypeRelativeURL FlagType = "relativeurl"

	//
	// --- Array / Slice Formats ---

	// --- Slice of Primitive Types

	TypeStringArray FlagType = "[]string"
	TypeIntArray    FlagType = "[]int"
	TypeFloatArray  FlagType = "[]float"
	TypeBoolArray   FlagType = "[]bool"

	// --- Slice of Structured Data

	TypeDurationArray FlagType = "[]duration"
	TypeDateArray     FlagType = "[]date"
	TypeTimeArray     FlagType = "[]time"
	TypeDateTimeArray FlagType = "[]datetime"
	TypeEmailArray    FlagType = "[]email"

	// --- Slice of System and Network Validations

	TypePathArray        FlagType = "[]path"
	TypeFilenameArray    FlagType = "[]filename"
	TypeFilepathArray    FlagType = "[]filepath"
	TypeDirnameArray     FlagType = "[]dirname"
	TypeDirpathArray     FlagType = "[]dirpath"
	TypeURLArray         FlagType = "[]url"
	TypeFullURLArray     FlagType = "[]fullurl"
	TypeRelativeURLArray FlagType = "[]relativeurl"
)

// AccessMode defines the required filesystem permissions needed to validate a path flag.
type AccessMode string

const (
	// ModeRead ensures the application can open and read the file or directory.
	ModeRead AccessMode = "read"

	// ModeWrite ensures the application can modify or create the file or directory.
	ModeWrite AccessMode = "write"

	// ModeReadWrite ensures both read and write capabilities are available.
	ModeReadWrite AccessMode = "readwrite"
)

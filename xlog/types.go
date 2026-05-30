package xlog

// DefaultTimeFormat defines the layout string used for timestamp rendering.
const DefaultTimeFormat = "YYYY-MM-DD HH:mm:ss"

// DefaultTimeZone defines the default fallback timezone configuration.
const DefaultTimeZone = "UTC"

// LogLevel represents the custom severity thresholds supported by this logger.
type LogLevel string

const (
	// LevelAll enables all available logging severities.
	LevelAll LogLevel = "all"
	// LevelNone completely silences the logger output.
	LevelNone LogLevel = "none"
	// LevelInfo filters logs to only show information, warning, and error messages.
	LevelInfo LogLevel = "info"
	// LevelWarn filters logs to only show warning and error messages.
	LevelWarn LogLevel = "warn"
	// LevelError restricts logs to error messages only.
	LevelError LogLevel = "error"
)

// ANSI sequence codes utilized for command-line interface terminal colorization.
const (
	ansiReset  = "\033[0m"
	ansiCyan   = "\033[36m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiRed    = "\033[31m"
	ansiGray   = "\033[90m"
)

// CLIColorPalette manages the specific ANSI string colors associated
// with each logging level and metadata element inside CLI environments.
type CLIColorPalette struct {
	Debug    string `json:"debug"`
	Info     string `json:"info"`
	Warn     string `json:"warn"`
	Error    string `json:"error"`
	DateTime string `json:"datetime"`
}

// DefaultCLIPalette provides the standard out-of-the-box color definitions
// for readable command line outputs.
var DefaultCLIPalette = &CLIColorPalette{
	Debug:    ansiCyan,
	Info:     ansiGreen,
	Warn:     ansiYellow,
	Error:    ansiRed,
	DateTime: ansiGray,
}

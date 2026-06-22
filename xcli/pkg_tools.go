package xcli

import (
	"fmt"
	"os"
	"strings"
)

// PrintLog writes a message to the standard output.
// If no extra arguments (args) are provided, it behaves like fmt.Println, printing the message directly.
// If extra arguments are provided, it treats the message as a format string and behaves like fmt.Printf, automatically appending a newline character.
func Print(message string, args ...any) {
	if len(args) == 0 {
		// Ensures a clean, predictable line-break output layout without complex formatting overhead
		fmt.Fprintln(os.Stdout, message)
		return
	}

	// Normalizes trailing boundary spacing dynamically to prevent double-newline visual distortions
	format := message
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stdout, format, args...)
}

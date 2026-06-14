package xerrors

import (
	"fmt"
	"os"
)

// Print writes the error's string representation directly to the standard error stream (os.Stderr).
// It acts as a lightweight, zero-allocation wrapper around fmt.Fprintln, providing a swift
// debugging mechanism that bypasses context allocations or structured logging handler configurations.
// If the incoming error target is nil, the execution block returns early without writing to the stream.
func Print(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err.Error())
}

package xclifn

import "os"

// ExportNewValidationError exposes the private newValidationError function strictly for external testing packages.
func ExportNewValidationError(
	flag string,
	msg string,
	given string,
	expected string,
	rule string,
) error {
	return newValidationError(flag, msg, given, expected, rule)
}

// SetDiskBridgesMocks allows external test suites to inject custom simulated OS behaviors.
func SetDiskBridgesMocks(
	stat func(string) (os.FileInfo, error),
	open func(string) (*os.File, error),
	create func(string) (*os.File, error),
	openFile func(string, int, os.FileMode) (*os.File, error),
) {
	if stat != nil {
		osStatBridge = stat
	}
	if open != nil {
		osOpenBridge = open
	}
	if create != nil {
		osCreateBridge = create
	}
	if openFile != nil {
		osOpenFileBridge = openFile
	}
}

// ResetDiskBridges restores all internal filesystem bridges to their original native OS behaviors.
func ResetDiskBridges() {
	osStatBridge = os.Stat
	osOpenBridge = os.Open
	osCreateBridge = os.Create
	osOpenFileBridge = os.OpenFile
}

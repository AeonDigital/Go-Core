package xfs

// SetBridgeXFSForTest overrides the internal package bridge during testing.
func SetBridgeXFSForTest(b IPkgBridgeXFS) func() {
	old := pkgBridgeXFS
	pkgBridgeXFS = b
	return func() {
		pkgBridgeXFS = old
	}
}

// ExportOsUserDataDirLinux exposes the private platform-specific logic
// strictly to the black-box xfs_test package execution suite.
func ExportOsUserDataDirLinux(home string, appName string) string {
	return osUserDataDirLinux(home, appName)
}

// ExportOsUserDataDirDarwin exposes the private platform-specific logic
// strictly to the black-box xfs_test package execution suite.
func ExportOsUserDataDirDarwin(home string, appName string) string {
	return osUserDataDirDarwin(home, appName)
}

// ExportOsUserDataDirWindows exposes the private platform-specific logic
// strictly to the black-box xfs_test package execution suite.
func ExportOsUserDataDirWindows(home string, appName string) string {
	return osUserDataDirWindows(home, appName)
}

// ExportOsUserLogDir exposes the private routing mechanism to xfs_test.
func ExportOsUserLogDir(home string, appName string) string {
	return osUserLogDir(home, appName)
}

// ExportOsUserDataDir exposes the private routing mechanism to xfs_test.
func ExportOsUserDataDir(home string, appName string) string {
	return osUserDataDir(home, appName)
}

// ExportOsUserLogDirLinux exposes the private specialist function to xfs_test.
func ExportOsUserLogDirLinux(home string, appName string) string {
	return osUserLogDirLinux(home, appName)
}

// ExportOsUserLogDirDarwin exposes the private specialist function to xfs_test.
func ExportOsUserLogDirDarwin(home string, appName string) string {
	return osUserLogDirDarwin(home, appName)
}

// ExportOsUserLogDirWindows exposes the private specialist function to xfs_test.
func ExportOsUserLogDirWindows(home string, appName string) string {
	return osUserLogDirWindows(home, appName)
}

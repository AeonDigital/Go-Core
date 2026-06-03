package xfs

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/AeonDigital/Go-Core/xerrors"
)

const xerror_CTX = "XFS"

// RetrieveFullPath returns the fully resolved absolute path, or an error if
// system paths cannot be determined.
//
// Replaces the tilde prefix ("~") with the user's home directory and
// strictly resolves the final result into a clean, absolute filesystem path.
// If the path is relative, it anchors it to the current working directory, removing
// any redundancies like "." or "..".
func RetrieveFullPath(path string) (string, error) {
	if path == "" {
		return "", xerrors.NewErr(
			xerrors.XERR_EMPTY_NOT_ALLOWED,
			xerror_CTX,
			"",
			"path",
		)
	}

	// Corrigido para garantir compatibilidade cross-platform de barras
	path = strings.ReplaceAll(path, "\\", "/")

	var resolvedPath string

	if strings.HasPrefix(path, "~") {
		home, err := pkgBridgeXFS.GetUserHomeDir()
		if err != nil {
			return "", xerrors.NewErr(
				xerrors.XERR_RESOURCE_UNAVAILABLE,
				xerror_CTX,
				"failed to resolve user home directory",
				"home",
				"~",
				err,
			)
		}

		if path == "~" {
			resolvedPath = home
		} else if path[1] == '/' {
			resolvedPath = filepath.Join(home, path[2:])
		} else {
			resolvedPath = path
		}
	} else {
		resolvedPath = path
	}

	// Usando a variável mockável aqui
	absolutePath, err := pkgBridgeXFS.FilepathAbs(resolvedPath)
	if err != nil {
		return "", xerrors.NewErr(
			xerrors.XERR_INVALID_VALUE,
			xerror_CTX,
			"failed to resolve absolute path",
			"path",
			path,
			"must be a valid filesystem trackable sequence",
			err,
		)
	}

	return absolutePath, nil
}

// ============================================================================
// GROUP: SPECIAL DIRECTORIES
// ============================================================================

// GetUserHomeDir retrieves the absolute filesystem path to the current user's home directory.
// It acts as a wrapper around Go's native standard library call, ensuring cross-platform compatibility.
// Returns the home directory path string, or an error if the system environment variables or user profile cannot be resolved.
func GetUserHomeDir() (string, error) {
	return pkgBridgeXFS.GetUserHomeDir()
}

// GetUserConfigDir retrieves the absolute filesystem path to the current user's configuration directory.
// It wraps Go's native standard library call to return the appropriate standard path based on the operating system guidelines.
// Returns the configuration directory path string, or an error if the path cannot be determined.
func GetUserConfigDir() (string, error) {
	return pkgBridgeXFS.GetUserConfigDir()
}

// GetUserDataDir determines the standard platform-specific directory for storing application data.
// It accepts the application name (appName) to append to the base system path.
// It automatically falls back to the system's temporary directory if the user's home directory cannot be resolved.
// Returns the fully constructed and cleaned absolute path tailored to the target operating system guidelines.
func GetUserDataDir(appName string) (string, error) {
	home, err := pkgBridgeXFS.GetUserHomeDir()
	if err != nil {
		home = os.TempDir()
	}

	return pkgBridgeXFS.OSUserDataDir(home, appName), nil
}

// GetUserLogDir determines the standard platform-specific directory for storing application log files.
// It accepts the application name (appName) to append to the base system path.
// It automatically falls back to the system's temporary directory if the user's home directory cannot be resolved.
// Returns the fully constructed and cleaned absolute path tailored to the target operating system guidelines.
func GetUserLogDir(appName string) (string, error) {
	home, err := pkgBridgeXFS.GetUserHomeDir()
	if err != nil {
		home = os.TempDir()
	}

	return pkgBridgeXFS.OSUserLogDir(home, appName), nil
}

// GetParentPath extracts and returns the parent directory path from a given file path.
// It automatically expands the tilde prefix ("~") using the RetrieveFullPath function before isolating the layout.
// Returns the parent directory string, or an empty string if the input path is empty or if path expansion fails.
func GetParentPath(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return "", err
	}

	return filepath.Dir(expandedPath), nil
}

// ============================================================================
// GROUP: EXISTENCE AND TYPES
// ============================================================================

// Exists checks if a file or directory exists at the given path.
// It automatically expands the tilde prefix ("~") using the RetrieveFullPath function before executing the check.
// Returns true if the target path is found on the system, or false if it does not exist or if path expansion fails.
func Exists(path string) bool {
	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return false
	}

	_, err = pkgBridgeXFS.OsStat(expandedPath)
	return err == nil
}

// IsFile checks if the specified path exists and points to a regular file.
// It automatically expands the tilde prefix ("~") using the RetrieveFullPath function before executing the check.
// Returns true if the path is a regular file, or false if it does not exist, is a directory/symlink, or if system calls fail.
func IsFile(path string) bool {
	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return false
	}

	info, err := pkgBridgeXFS.OsStat(expandedPath)
	if err != nil {
		return false
	}

	return info.Mode().IsRegular()
}

// IsDir checks if the specified path exists and points to a directory.
// It automatically expands the tilde prefix ("~") using the RetrieveFullPath function before executing the check.
// Returns true if the path is a directory, or false if it does not exist, is a regular file/symlink, or if system calls fail.
func IsDir(path string) bool {
	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return false
	}

	info, err := pkgBridgeXFS.OsStat(expandedPath)
	if err != nil {
		return false
	}

	return info.IsDir()
}

// IsEmptyDir checks if the specified directory exists and contains no files or subdirectories.
// It automatically expands the tilde prefix ("~") before evaluating the target directory path.
// Returns true if the directory is completely empty, or false along with an error if the path does not exist, is not a directory, or if read permissions are denied.
func IsEmptyDir(path string) (bool, error) {
	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return false, err
	}

	f, err := pkgBridgeXFS.OsOpen(expandedPath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = pkgBridgeXFS.ReadDir(f, 1)
	if err == nil {
		return false, nil
	}

	if errors.Is(err, io.EOF) {
		return true, nil
	}

	return false, err
}

// IsFileDirExists extracts the parent directory path from a given file path and checks if it exists as a valid directory.
// It automatically expands the tilde prefix ("~") before isolating the parent directory layout.
// Returns true if the parent directory structure exists on the system, or false if the path is empty, expansion fails, or the directory does not exist.
func IsFileDirExists(filePath string) bool {
	if filePath == "" {
		return false
	}

	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(filePath)
	if err != nil {
		return false
	}

	parentDir := filepath.Dir(expandedPath)
	return pkgBridgeXFS.IsDir(parentDir)
}

// // ============================================================================
// // GROUP: PERMISSIONS
// // ============================================================================

// IsReadable checks if the current application process has sufficient permissions to read the file or directory at the specified path.
// It automatically expands the tilde prefix ("~") and robustly handles both regular files and directory permission boundaries across different operating systems.
// Returns true if the path can be successfully read, or false if permissions are denied or if the path does not exist.
func IsReadable(path string) bool {
	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return false
	}

	info, err := pkgBridgeXFS.OsStat(expandedPath)
	if err != nil {
		return false
	}

	if info.IsDir() {
		f, err := pkgBridgeXFS.OsOpen(expandedPath)
		if err != nil {
			return false
		}
		defer f.Close()

		_, err = pkgBridgeXFS.ReadDir(f, 1)

		return err == nil || errors.Is(err, io.EOF)
	}

	file, err := pkgBridgeXFS.OsOpenFile(expandedPath, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

// IsWritable checks if the current application process has sufficient permissions to modify or write to the file or directory at the specified path.
// It automatically expands the tilde prefix ("~") and tests directories by attempting to create a temporary file inside them, or files by opening them in write-only append mode.
// Returns true if the path can be written to, or false if permissions are denied, if the path does not exist, or if the test operations fail.
func IsWritable(path string) bool {
	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return false
	}

	info, err := pkgBridgeXFS.OsStat(expandedPath)
	if err != nil {
		return false
	}

	if info.IsDir() {
		tempFile, err := pkgBridgeXFS.OsCreateTemp(expandedPath, ".fsutil_test_")
		if err != nil {
			return false
		}

		defer pkgBridgeXFS.OsRemove(tempFile.Name())
		defer tempFile.Close()

		return true
	}

	file, err := pkgBridgeXFS.OsOpenFile(expandedPath, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

// GetPermission retrieves the system permission bits (os.FileMode) for the specified file or directory.
// It automatically expands the tilde prefix ("~") using the RetrieveFullPath function before executing the check.
// Returns the file mode permissions, or an error if path expansion fails or if the path does not exist.
func GetPermission(path string) (os.FileMode, error) {
	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return 0, err
	}

	info, err := pkgBridgeXFS.OsStat(expandedPath)
	if err != nil {
		return 0, xerrors.NewErr(
			xerrors.XERR_NOT_FOUND,
			xerror_CTX,
			"failed to read path metadata",
			"expandedPath",
			expandedPath,
			err,
		)
	}

	return info.Mode().Perm(), nil
}

// SetPermission updates the filesystem permission bits (os.FileMode) for the specified file or directory.
// It automatically expands the tilde prefix ("~") using the RetrieveFullPath function before applying the changes.
// Returns nil on success, or an error if path expansion fails, the path does not exist, or the operation is denied.
func SetPermission(path string, perm os.FileMode) error {
	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return err
	}

	err = pkgBridgeXFS.OsChmod(expandedPath, perm)
	if err != nil {
		return xerrors.NewErr(
			xerrors.XERR_PERMISSION_DENIED,
			xerror_CTX,
			"failed to change target filesystem permissions",
			"expandedPath",
			expandedPath,
			err,
		)
	}

	return nil
}

// // ============================================================================
// // GROUP: CREATE, EDIT and DELETE
// // ============================================================================

// CreateFile creates a new file at the specified path, truncating it if it already exists.
// It automatically expands the tilde prefix ("~") before executing the operation.
// An optional os.FileMode can be passed; if omitted, it defaults to standard system permissions (0666).
// Returns a pointer to the opened file object (*os.File) ready for writing, or an error if path expansion or file creation fails.
func CreateFile(path string, perm ...os.FileMode) (*os.File, error) {
	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return nil, err
	}

	filePerm := os.FileMode(0666)
	if len(perm) > 0 {
		filePerm = perm[0]
	}

	file, err := os.OpenFile(expandedPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, filePerm)
	if err != nil {
		return nil, xerrors.NewErr(
			xerrors.XERR_RESOURCE_UNAVAILABLE,
			xerror_CTX,
			"failed to create targeted file",
			"expandedPath",
			expandedPath,
			err,
		)
	}

	return file, nil
}

// CreateDir creates a single directory at the specified path.
// It automatically expands the tilde prefix ("~") before creating the folder and fails if the parent directory structure does not exist.
// An optional os.FileMode can be passed; if omitted, it defaults to standard permissions (0755).
// Returns nil on success, or an error if path expansion fails or the directory cannot be created.
func CreateDir(path string, perm ...os.FileMode) error {
	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return err
	}

	dirPerm := os.FileMode(0755)
	if len(perm) > 0 {
		dirPerm = perm[0]
	}

	err = pkgBridgeXFS.OsMkdir(expandedPath, dirPerm)
	if err != nil {
		return xerrors.NewErr(
			xerrors.XERR_RESOURCE_UNAVAILABLE,
			xerror_CTX,
			"failed to create targeted directory",
			"expandedPath",
			expandedPath,
			err,
		)
	}

	return nil
}

// CreateDirPath creates a directory path along with any necessary parent directories.
// It automatically expands the tilde prefix ("~") before creating the folder structure.
// An optional os.FileMode can be passed; if omitted, it defaults to standard permissions (0755).
// If the target path already exists, it does nothing and returns nil.
func CreateDirPath(path string, perm ...os.FileMode) error {
	if path == "" {
		return xerrors.NewErr(
			xerrors.XERR_EMPTY_NOT_ALLOWED,
			xerror_CTX,
			"",
			"path",
		)
	}

	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return err
	}

	dirPerm := os.FileMode(0755)
	if len(perm) > 0 {
		dirPerm = perm[0]
	}

	err = pkgBridgeXFS.OsMkdirAll(expandedPath, dirPerm)
	if err != nil {
		return xerrors.NewErr(
			xerrors.XERR_RESOURCE_UNAVAILABLE,
			xerror_CTX,
			"failed to create targeted directory tree",
			"expandedPath",
			expandedPath,
			err,
		)
	}

	return nil
}

// OpenFileWrite opens a file for writing, creating it with the specified or default permissions (0644) if it does not exist.
// It automatically expands the tilde prefix ("~") before opening the path.
// If truncate is true, the file's existing content is cleared upon opening; otherwise, new data is appended to the end.
// An optional os.FileMode can be passed to override the default creation permissions.
// Returns a pointer to the opened file object (*os.File) ready for writing, or an error if path expansion or file opening fails.
func OpenFileWrite(path string, truncate bool, perm ...os.FileMode) (*os.File, error) {
	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return nil, err
	}

	// 0644: Standard read/write permissions for the owner, read-only for others
	filePerm := os.FileMode(0644)
	if len(perm) > 0 {
		filePerm = perm[0]
	}

	// Flag Explanations:
	// os.O_WRONLY: Open exclusively for writing
	// os.O_CREATE: Create the file if it does not exist
	// os.O_APPEND: Write data at the end of the file (log style)
	// os.O_TRUNC: Truncate entire file content on open file
	useFlag := os.O_APPEND
	if truncate {
		useFlag = os.O_TRUNC
	}

	file, err := os.OpenFile(expandedPath, os.O_WRONLY|os.O_CREATE|useFlag, filePerm)
	if err != nil {
		return nil, xerrors.NewErr(
			xerrors.XERR_RESOURCE_UNAVAILABLE,
			xerror_CTX,
			"failed to open targeted file for writing operations",
			"expandedPath",
			expandedPath,
			err,
		)
	}

	return file, nil
}

// OpenFileRead opens an existing file exclusively in read-only mode.
// It automatically expands the tilde prefix ("~") using the RetrieveFullPath function before attempting to open the path.
// Returns a pointer to the opened file object (*os.File) ready for reading, or an error if path expansion fails or if the file does not exist.
func OpenFileRead(path string) (*os.File, error) {
	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return nil, err
	}

	// os.O_RDONLY: Open the file read-only.
	file, err := os.OpenFile(expandedPath, os.O_RDONLY, 0)
	if err != nil {
		return nil, xerrors.NewErr(
			xerrors.XERR_NOT_FOUND,
			xerror_CTX,
			"failed to open targeted file for reading operations",
			"expandedPath",
			expandedPath,
			err,
		)
	}

	return file, nil
}

// DeleteFile removes a regular file at the specified path.
// It automatically expands the tilde prefix ("~") before performing the check and deletion.
// It explicitly fails if the targeted path is a directory to prevent accidental directory structural loss.
// Returns nil on success, or an error if path expansion fails, the path points to a directory, or the system removal operation is denied.
func DeleteFile(path string) error {
	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return err
	}

	if pkgBridgeXFS.IsDir(expandedPath) {
		return xerrors.NewErr(
			xerrors.XERR_INVALID_TYPE,
			xerror_CTX,
			"cannot use DeleteFile on a directory",
			"expandedPath",
			expandedPath,
			"regular file",
		)
	}

	err = pkgBridgeXFS.OsRemove(expandedPath)
	if err != nil {
		return xerrors.NewErr(
			xerrors.XERR_RESOURCE_UNAVAILABLE,
			xerror_CTX,
			"failed to delete file at targeted path",
			"expandedPath",
			expandedPath,
			err,
		)
	}

	return nil
}

// DeleteDir removes a directory at the specified path, with optional recursive deletion of all its contents.
// It automatically expands the tilde prefix ("~") before executing any filesystem validation or removal steps.
// If recursive is true, it uses os.RemoveAll to forcefully delete the folder and everything inside it; otherwise, it uses os.Remove and fails if the directory is not completely empty.
// Returns nil on success, or an error if path expansion fails, the path is not a directory, or the system removal operation is denied.
func DeleteDir(path string, recursive bool) error {
	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return err
	}

	if !pkgBridgeXFS.IsDir(expandedPath) {
		return xerrors.NewErr(
			xerrors.XERR_INVALID_TYPE,
			xerror_CTX,
			"path is not a directory or does not exist",
			"expandedPath",
			expandedPath,
			"directory",
		)
	}

	if recursive {
		err = pkgBridgeXFS.OsRemoveAll(expandedPath)
	} else {
		err = pkgBridgeXFS.OsRemove(expandedPath)
	}

	if err != nil {
		msgWithData := "failed to delete targeted directory hierarchy, " + xerrors.MsgData("recursive", recursive)
		return xerrors.NewErr(
			xerrors.XERR_RESOURCE_UNAVAILABLE,
			xerror_CTX,
			msgWithData,
			"expandedPath",
			expandedPath,
			err,
		)
	}

	return nil
}

// // ============================================================================
// // GROUP: METADATA AND INTEGRITY
// // ============================================================================

// HasSpaceAvailable checks if the storage volume containing the specified path has at least the required number of bytes free.
// It accepts a filesystem path (path) and the minimum required space in bytes (requiredBytes).
// If the target path does not exist, it traverses upward to find the closest existing parent directory to accurately target the underlying volume.
// Returns true if the volume has sufficient space available for the current user, or false along with an error if the path expansion, UTF-16 conversion, or Win32 GetDiskFreeSpaceEx API call fails.
func HasSpaceAvailable(path string, requiredBytes uint64) (bool, error) {
	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return false, err
	}

	currentPath := expandedPath
	for {
		if pkgBridgeXFS.Exists(currentPath) {
			break
		}

		parent := filepath.Dir(currentPath)
		if parent == currentPath {
			return false, xerrors.NewErr(
				xerrors.XERR_NOT_FOUND,
				xerror_CTX,
				"could not find a valid base volume path",
				"expandedPath",
				expandedPath,
			)
		}
		currentPath = parent
	}

	return pkgBridgeXFS.GetVolumeFreeSpace(currentPath, requiredBytes)
}

// GetFileSize returns the size of the file in bytes at the specified path.
// It automatically expands the tilde prefix ("~") before checking the path.
// It explicitly returns an error if the path points to a directory or if the file does not exist.
// Returns the file size in bytes as an int64, or an error if path expansion, system calls, or validations fail.
func GetFileSize(path string) (int64, error) {
	expandedPath, err := pkgBridgeXFS.RetrieveFullPath(path)
	if err != nil {
		return 0, err
	}

	info, err := pkgBridgeXFS.OsStat(expandedPath)
	if err != nil {
		return 0, err
	}

	if info.IsDir() {
		return 0, os.ErrInvalid // Directories do not have a standard "file size"
	}

	return info.Size(), nil
}

// IsSameFile checks if two different paths point to the exact same physical file or directory on the storage device.
// It automatically expands the tilde prefix ("~") on both inputs and resolves any symbolic links encountered.
// It relies on Go's native os.SameFile mechanism to compare low-level system identities like inodes or file IDs.
// Returns true if both paths target the identical filesystem resource, or false along with an error if metadata retrieval fails.
func IsSameFile(pathA, pathB string) (bool, error) {
	expandedA, err := pkgBridgeXFS.RetrieveFullPath(pathA)
	if err != nil {
		return false, err
	}

	expandedB, err := pkgBridgeXFS.RetrieveFullPath(pathB)
	if err != nil {
		return false, err
	}

	// os.Stat automatically follows symbolic links
	infoA, err := pkgBridgeXFS.OsStat(expandedA)
	if err != nil {
		return false, err
	}

	infoB, err := pkgBridgeXFS.OsStat(expandedB)
	if err != nil {
		return false, err
	}

	// os.SameFile compares the underlying system-specific file identities (inodes/file IDs)
	return os.SameFile(infoA, infoB), nil
}

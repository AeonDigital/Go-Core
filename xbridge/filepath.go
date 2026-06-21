package xbridge

import (
	"path/filepath"
	"strings"
)

// IFilepathBridge defines the unified contract to intercept and mock the Go standard library 'path/filepath' package.
// This interface decouples filesystem path manipulations and string parsing from the application logic.
type IFilepathBridge interface {
	// =========================================================================
	// PATH MANIPULATION & CLEANING
	// =========================================================================

	// Abs returns an absolute representation of path.
	// If the path is not absolute it will be joined with the current
	// working directory to turn it into an absolute path. The absolute
	// path name for a given file is not guaranteed to be unique.
	// Abs calls [Clean] on the result.
	Abs(path string) (string, error)
	// Clean returns the shortest path name equivalent to path
	// by purely lexical processing. It applies the following rules
	// iteratively until no further processing can be done:
	//
	//  1. Replace multiple [Separator] elements with a single one.
	//  2. Eliminate each . path name element (the current directory).
	//  3. Eliminate each inner .. path name element (the parent directory)
	//     along with the non-.. element that precedes it.
	//  4. Eliminate .. elements that begin a rooted path:
	//     that is, replace "/.." by "/" at the beginning of a path,
	//     assuming Separator is '/'.
	//
	// The returned path ends in a slash only if it represents a root directory,
	// such as "/" on Unix or `C:\` on Windows.
	//
	// Finally, any occurrences of slash are replaced by Separator.
	//
	// If the result of this process is an empty string, Clean
	// returns the string ".".
	//
	// On Windows, Clean does not modify the volume name other than to replace
	// occurrences of "/" with `\`.
	// For example, Clean("//host/share/../x") returns `\\host\share\x`.
	//
	// See also Rob Pike, “Lexical File Names in Plan 9 or
	// Getting Dot-Dot Right,”
	// https://9p.io/sys/doc/lexnames.html
	Clean(path string) string
	// Dir returns all but the last element of path, typically the path's directory.
	// After dropping the final element, Dir calls [Clean] on the path and trailing
	// slashes are removed.
	// If the path is empty, Dir returns ".".
	// If the path consists entirely of separators, Dir returns a single separator.
	// The returned path does not end in a separator unless it is the root directory.
	Dir(path string) string
	// Base returns the last element of path.
	// Trailing path separators are removed before extracting the last element.
	// If the path is empty, Base returns ".".
	// If the path consists entirely of separators, Base returns a single separator.
	Base(path string) string
	// Ext returns the file name extension used by path.
	// The extension is the suffix beginning at the final dot
	// in the final element of path; it is empty if there is
	// no dot.
	Ext(path string) string
	// FromSlash returns the result of replacing each slash ('/') character
	// in path with a separator character. Multiple slashes are replaced
	// by multiple separators.
	//
	// See also the Localize function, which converts a slash-separated path
	// as used by the io/fs package to an operating system path.
	FromSlash(path string) string
	// ToSlash returns the result of replacing each separator character
	// in path with a slash ('/') character. Multiple separators are
	// replaced by multiple slashes.
	ToSlash(path string) string

	// =========================================================================
	// PATH COMPOSITION & DECOMPOSITION
	// =========================================================================

	// Join joins any number of path elements into a single path,
	// separating them with an OS specific [Separator]. Empty elements
	// are ignored. The result is Cleaned. However, if the argument
	// list is empty or all its elements are empty, Join returns
	// an empty string.
	// On Windows, the result will only be a UNC path if the first
	// non-empty element is a UNC path.
	Join(elem ...string) string
	// Split splits path immediately following the final [Separator],
	// separating it into a directory and file name component.
	// If there is no Separator in path, Split returns an empty dir
	// and file set to path.
	// The returned values have the property that path = dir+file.
	Split(path string) (dir, file string)
	// SplitList splits a list of paths joined by the OS-specific [ListSeparator],
	// usually found in PATH or GOPATH environment variables.
	// Unlike strings.Split, SplitList returns an empty slice when passed an empty
	// string.
	SplitList(path string) []string

	// =========================================================================
	// PATH MATCHING & EVALUATION
	// =========================================================================

	// EvalSymlinks returns the path name after the evaluation of any symbolic
	// links.
	// If path is relative the result will be relative to the current directory,
	// unless one of the components is an absolute symbolic link.
	// EvalSymlinks calls [Clean] on the result.
	EvalSymlinks(path string) (string, error)
	// Glob returns the names of all files matching pattern or nil
	// if there is no matching file. The syntax of patterns is the same
	// as in [Match]. The pattern may describe hierarchical names such as
	// /usr/*/bin/ed (assuming the [Separator] is '/').
	//
	// Glob ignores file system errors such as I/O errors reading directories.
	// The only possible returned error is [ErrBadPattern], when pattern
	// is malformed.
	Glob(pattern string) (matches []string, err error)
	// Match reports whether name matches the shell file name pattern.
	// The pattern syntax is:
	//
	//	pattern:
	//		{ term }
	//	term:
	//		'*'         matches any sequence of non-Separator characters
	//		'?'         matches any single non-Separator character
	//		'[' [ '^' ] { character-range } ']'
	//		            character class (must be non-empty)
	//		c           matches character c (c != '*', '?', '\\', '[')
	//		'\\' c      matches character c (except on Windows)
	//
	//	character-range:
	//		c           matches character c (c != '\\', '-', ']')
	//		'\\' c      matches character c (except on Windows)
	//		lo '-' hi   matches character c for lo <= c <= hi
	//
	// Path segments in the pattern must be separated by [Separator].
	//
	// Match requires pattern to match all of name, not just a substring.
	// The only possible returned error is [ErrBadPattern], when pattern
	// is malformed.
	//
	// On Windows, escaping is disabled. Instead, '\\' is treated as
	// path separator.
	Match(pattern, name string) (matched bool, err error)
	// Rel returns a relative path that is lexically equivalent to targPath when
	// joined to basePath with an intervening separator. That is,
	// [Join](basePath, Rel(basePath, targPath)) is equivalent to targPath itself.
	//
	// The returned path will always be relative to basePath, even if basePath and
	// targPath share no elements. Rel calls [Clean] on the result.
	//
	// An error is returned if targPath can't be made relative to basePath
	// or if knowing the current working directory would be necessary to compute it.
	Rel(basepath, targpath string) (string, error)

	// =========================================================================
	// PATH ANALYSIS
	// =========================================================================

	// IsAbs reports whether the path is absolute.
	IsAbs(path string) bool
	// HasPrefix evaluates if the target path string starts with the given prefix literal.
	// This method is preserved strictly for historical compatibility with existing codebases.
	//
	// Deprecated: HasPrefix performs a literal string check; it does not respect
	// path boundaries and does not ignore case on case-insensitive filesystems.
	// Use IsInDirectory for secure path confinement validations instead.
	HasPrefix(p, prefix string) bool
	// IsInDirectory structurally evaluates if the target path is located inside
	// the base directory, safely respecting path boundaries and cross-platform casing.
	IsInDirectory(baseDir, targetPath string) bool
	// VolumeName returns leading volume name.
	// Given "C:\foo\bar" it returns "C:" on Windows.
	// Given "\\host\share\foo" it returns "\\host\share".
	// On other platforms it returns "".
	VolumeName(path string) string
}

// sfilepathBridge serves as the private concrete production implementation of IFilepathBridge,
// routing calls natively to the standard library 'path/filepath' package.
type sfilepathBridge struct{}

// Ensure at compile time that the private struct implements the public interface.
var _ IFilepathBridge = sfilepathBridge{}

// =========================================================================
// IMPLEMENTATION OF PATH MANIPULATION & CLEANING
// =========================================================================

func (sfilepathBridge) Abs(path string) (string, error) { return filepath.Abs(path) }
func (sfilepathBridge) Clean(path string) string        { return filepath.Clean(path) }
func (sfilepathBridge) Dir(path string) string          { return filepath.Dir(path) }
func (sfilepathBridge) Base(path string) string         { return filepath.Base(path) }
func (sfilepathBridge) Ext(path string) string          { return filepath.Ext(path) }
func (sfilepathBridge) FromSlash(path string) string    { return filepath.FromSlash(path) }
func (sfilepathBridge) ToSlash(path string) string      { return filepath.ToSlash(path) }

// =========================================================================
// IMPLEMENTATION OF PATH COMPOSITION & DECOMPOSITION
// =========================================================================

func (sfilepathBridge) Join(elem ...string) string { return filepath.Join(elem...) }
func (sfilepathBridge) Split(path string) (string, string) {
	return filepath.Split(path)
}
func (sfilepathBridge) SplitList(path string) []string { return filepath.SplitList(path) }

// =========================================================================
// IMPLEMENTATION OF PATH MATCHING & EVALUATION
// =========================================================================

func (sfilepathBridge) EvalSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}
func (sfilepathBridge) Glob(pattern string) ([]string, error) { return filepath.Glob(pattern) }
func (sfilepathBridge) Match(pattern, name string) (bool, error) {
	return filepath.Match(pattern, name)
}
func (sfilepathBridge) Rel(basepath, targpath string) (string, error) {
	return filepath.Rel(basepath, targpath)
}

// =========================================================================
// IMPLEMENTATION OF PATH ANALYSIS
// =========================================================================

func (sfilepathBridge) IsAbs(path string) bool          { return filepath.IsAbs(path) }
func (sfilepathBridge) HasPrefix(p, prefix string) bool { return strings.HasPrefix(p, prefix) }
func (sfilepathBridge) IsInDirectory(baseDir, targetPath string) bool {
	return filepath_IsInDirectory(baseDir, targetPath)
}
func (sfilepathBridge) VolumeName(path string) string { return filepath.VolumeName(path) }

// =========================================================================
// GLOBAL PUBLIC INTERFACE INSTANCE
// =========================================================================

// FilepathBridge is the global public variable used by all core application code
// to perform cross-platform filesystem path evaluations and manipulations.
// It can be easily hot-swapped at unit test boundaries via generated pkgxmock utilities.
var FilepathBridge IFilepathBridge = sfilepathBridge{}

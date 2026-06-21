package bfs

import (
	"io/fs"
	"os"
)

// IFileBridge defines the unified contract to intercept and mock an open file instance.
// This allows you to completely mock file behaviors (Read, Write, Close) in memory.
type IFileBridge interface {
	// =========================================================================
	// DATA STREAMING (READ & WRITE)
	// =========================================================================

	// Read reads up to len(b) bytes from the File and stores them in b.
	// It returns the number of bytes read and any error encountered. At end of file, Read returns 0, io.EOF.
	Read(b []byte) (n int, err error)
	// ReadAt reads len(b) bytes from the File starting at byte offset off.
	// It returns the number of bytes read and the error, if any.
	// ReadAt always returns a non-nil error when n < len(b). At end of file, that error is io.EOF.
	ReadAt(b []byte, off int64) (n int, err error)
	// Write writes len(b) bytes from b to the File.
	// It returns the number of bytes written and an error, if any. Write returns a non-nil error when n != len(b).
	Write(b []byte) (n int, err error)
	// WriteAt writes len(b) bytes to the File starting at byte offset off.
	// It returns the number of bytes written and an error, if any. WriteAt returns a non-nil error when n != len(b).
	//
	// If file was opened with the O_APPEND flag, WriteAt returns an error.
	WriteAt(b []byte, off int64) (n int, err error)
	// WriteString is like Write, but writes the contents of string s rather than a slice of bytes.
	WriteString(s string) (n int, err error)

	// =========================================================================
	// POSITIONING & LIFECYCLE
	// =========================================================================

	// Seek sets the offset for the next Read or Write on file to offset, interpreted
	// according to whence: 0 means relative to the origin of the file, 1 means
	// relative to the current offset, and 2 means relative to the end.
	// It returns the new offset and an error, if any.
	// The behavior of Seek on a file opened with [O_APPEND] is not specified.
	Seek(offset int64, whence int) (ret int64, err error)
	// Close closes the File, rendering it unusable for I/O.
	// On files that support File.SetDeadline, any pending I/O operations will be canceled and return
	// immediately with an ErrClosed error.
	// Close will return an error if it has already been called.
	Close() error
	// Sync commits the current contents of the file to stable storage.
	// Typically, this means flushing the file system's in-memory copy of recently written data to disk.
	Sync() error

	// =========================================================================
	// METADATA & DESCRIPTORS
	// =========================================================================

	// Stat returns the FileInfo structure describing file.
	// If there is an error, it will be of type *PathError.
	Stat() (IFileInfoBridge, error)
	// Fd returns the system file descriptor or handle referencing the open file.
	// If f is closed, the descriptor becomes invalid.
	// If f is garbage collected, a finalizer may close the descriptor,
	// making it invalid; see [runtime.SetFinalizer] for more information on when
	// a finalizer might be run.
	//
	// Do not close the returned descriptor; that could cause a later
	// close of f to close an unrelated descriptor.
	//
	// Fd's behavior differs on some platforms:
	//
	//   - On Unix and Windows, [File.SetDeadline] methods will stop working.
	//   - On Windows, the file descriptor will be disassociated from the
	//     Go runtime I/O completion port if there are no concurrent I/O
	//     operations on the file.
	//
	// For most uses prefer the f.SyscallConn method.
	Fd() uintptr
	// Name returns the name of the file as presented to Open.
	//
	// It is safe to call Name after [Close].
	Name() string

	// =========================================================================
	// DIRECTORY ITERATION
	// =========================================================================

	// ReadDir reads the contents of the directory associated with the file f
	// and returns a slice of [DirEntry] values in directory order.
	// Subsequent calls on the same file will yield later DirEntry records in the directory.
	//
	// If n > 0, ReadDir returns at most n DirEntry records.
	// In this case, if ReadDir returns an empty slice, it will return an error explaining why.
	// At the end of a directory, the error is [io.EOF].
	//
	// If n <= 0, ReadDir returns all the DirEntry records remaining in the directory.
	// When it succeeds, it returns a nil error (not io.EOF).
	ReadDir(n int) ([]fs.DirEntry, error)
}

// sFileBridge wraps a concrete *os.File to satisfy the IFileBridge interface.
type sFileBridge struct {
	file *os.File
}

// Ensure at compile time that the private struct implements the public interface.
var _ IFileBridge = (*sFileBridge)(nil)

// NewFileBridge wraps a concrete *os.File and returns it as a public, mockable IFileBridge interface.
// This is the production entry point used when converting native os.File descriptors.
func NewFile(f *os.File) IFileBridge {
	if f == nil {
		return nil
	}
	return &sFileBridge{file: f}
}

func (f *sFileBridge) Read(b []byte) (int, error)               { return f.file.Read(b) }
func (f *sFileBridge) ReadAt(b []byte, off int64) (int, error)  { return f.file.ReadAt(b, off) }
func (f *sFileBridge) Write(b []byte) (int, error)              { return f.file.Write(b) }
func (f *sFileBridge) WriteAt(b []byte, off int64) (int, error) { return f.file.WriteAt(b, off) }
func (f *sFileBridge) WriteString(s string) (int, error)        { return f.file.WriteString(s) }

func (f *sFileBridge) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}
func (f *sFileBridge) Close() error { return f.file.Close() }
func (f *sFileBridge) Sync() error  { return f.file.Sync() }

func (f *sFileBridge) Stat() (IFileInfoBridge, error) {
	fi, err := f.file.Stat()
	if err != nil {
		return nil, err
	}
	return &sFileInfoBridge{info: fi}, nil
}
func (f *sFileBridge) Fd() uintptr  { return f.file.Fd() }
func (f *sFileBridge) Name() string { return f.file.Name() }

func (f *sFileBridge) ReadDir(n int) ([]fs.DirEntry, error) { return f.file.ReadDir(n) }

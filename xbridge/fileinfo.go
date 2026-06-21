package xbridge

import (
	"io/fs"
	"os"
	"time"
)

// IFileInfoBridge defines the unified contract to intercept and mock file metadata descriptors.
// This interface mirrors the standard library 'os.FileInfo' to enable seamless isolation in unit tests.
type IFileInfoBridge interface {
	// base name of the file
	Name() string
	// length in bytes for regular files; system-dependent for others
	Size() int64
	// file mode bits
	Mode() os.FileMode
	// modification time
	ModTime() time.Time
	// abbreviation for Mode().IsDir()
	IsDir() bool
	// underlying data source (can return nil)
	Sys() any
}

// sFileInfoBridge wraps a concrete fs.FileInfo to satisfy the IFileInfoBridge interface.
type sFileInfoBridge struct {
	info fs.FileInfo
}

// Ensure at compile time that the private struct implements the public interface.
var _ IFileInfoBridge = (*sFileInfoBridge)(nil)

func (s *sFileInfoBridge) Name() string       { return s.info.Name() }
func (s *sFileInfoBridge) Size() int64        { return s.info.Size() }
func (s *sFileInfoBridge) Mode() os.FileMode  { return s.info.Mode() }
func (s *sFileInfoBridge) ModTime() time.Time { return s.info.ModTime() }
func (s *sFileInfoBridge) IsDir() bool        { return s.info.IsDir() }
func (s *sFileInfoBridge) Sys() any           { return s.info.Sys() }

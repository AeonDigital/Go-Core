package bio

import "io"

// IIOBridge defines the unified contract to intercept and mock the Go standard library 'io' package.
// This interface decouples data streaming, reading, and writing operations from the application logic.
type IIOBridge interface {
	// =========================================================================
	// DATA COPYING
	// =========================================================================

	// Copy copies from src to dst until either EOF is reached
	// on src or an error occurs. It returns the number of bytes
	// copied and the first error encountered while copying, if any.
	//
	// A successful Copy returns err == nil, not err == EOF.
	// Because Copy is defined to read from src until EOF, it does
	// not treat an EOF from Read as an error to be reported.
	//
	// If src implements [WriterTo],
	// the copy is implemented by calling src.WriteTo(dst).
	// Otherwise, if dst implements [ReaderFrom],
	// the copy is implemented by calling dst.ReadFrom(src).
	Copy(dst io.Writer, src io.Reader) (written int64, err error)
	// CopyBuffer is identical to Copy except that it stages through the
	// provided buffer (if one is required) rather than allocating a
	// temporary one. If buf is nil, one is allocated; otherwise if it has
	// zero length, CopyBuffer panics.
	//
	// If either src implements [WriterTo] or dst implements [ReaderFrom],
	// buf will not be used to perform the copy.
	CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error)
	// CopyN copies n bytes (or until an error) from src to dst.
	// It returns the number of bytes copied and the earliest
	// error encountered while copying.
	// On return, written == n if and only if err == nil.
	//
	// If dst implements [ReaderFrom], the copy is implemented using it.
	CopyN(dst io.Writer, src io.Reader, n int64) (written int64, err error)

	// =========================================================================
	// UTILITY READERS & WRITERS
	// =========================================================================

	// ReadAll reads from r until an error or EOF and returns the data it read.
	// A successful call returns err == nil, not err == EOF. Because ReadAll is
	// defined to read from src until EOF, it does not treat an EOF from Read
	// as an error to be reported.
	ReadAll(r io.Reader) ([]byte, error)
	// ReadAtLeast reads from r into buf until it has read at least min bytes.
	// It returns the number of bytes copied and an error if fewer bytes were read.
	// The error is EOF only if no bytes were read.
	// If an EOF happens after reading fewer than min bytes,
	// ReadAtLeast returns [ErrUnexpectedEOF].
	// If min is greater than the length of buf, ReadAtLeast returns [ErrShortBuffer].
	// On return, n >= min if and only if err == nil.
	// If r returns an error having read at least min bytes, the error is dropped.
	ReadAtLeast(r io.Reader, buf []byte, min int) (n int, err error)
	// ReadFull reads exactly len(buf) bytes from r into buf.
	// It returns the number of bytes copied and an error if fewer bytes were read.
	// The error is EOF only if no bytes were read.
	// If an EOF happens after reading some but not all the bytes,
	// ReadFull returns [ErrUnexpectedEOF].
	// On return, n == len(buf) if and only if err == nil.
	// If r returns an error having read at least len(buf) bytes, the error is dropped.
	ReadFull(r io.Reader, buf []byte) (n int, err error)
	// WriteString writes the contents of the string s to w, which accepts a slice of bytes.
	// If w implements [StringWriter], [StringWriter.WriteString] is invoked directly.
	// Otherwise, [Writer.Write] is called exactly once.
	WriteString(w io.Writer, s string) (n int, err error)

	// =========================================================================
	// FACTORIES AND WRAPPERS
	// =========================================================================

	// LimitReader returns a Reader that reads from r
	// but stops with EOF after n bytes.
	// The underlying implementation is a *LimitedReader.
	LimitReader(r io.Reader, n int64) io.Reader
	// MultiReader returns a Reader that's the logical concatenation of
	// the provided input readers. They're read sequentially. Once all
	// inputs have returned EOF, Read will return EOF.  If any of the readers
	// return a non-nil, non-EOF error, Read will return that error.
	MultiReader(readers ...io.Reader) io.Reader
	// MultiWriter creates a writer that duplicates its writes to all the
	// provided writers, similar to the Unix tee(1) command.
	//
	// Each write is written to each listed writer, one at a time.
	// If a listed writer returns an error, that overall write operation
	// stops and returns the error; it does not continue down the list.
	MultiWriter(writers ...io.Writer) io.Writer
	// TeeReader returns a [Reader] that writes to w what it reads from r.
	// All reads from r performed through it are matched with
	// corresponding writes to w. There is no internal buffering -
	// the write must complete before the read completes.
	// Any error encountered while writing is reported as a read error.
	TeeReader(r io.Reader, w io.Writer) io.Reader
	// NopCloser returns a [ReadCloser] with a no-op Close method wrapping
	// the provided [Reader] r.
	// If r implements [WriterTo], the returned [ReadCloser] will implement [WriterTo]
	// by forwarding calls to r.
	NopCloser(r io.Reader) io.ReadCloser
}

// sioBridge serves as the private concrete production implementation of IIOBridge,
// routing calls natively to the standard library 'io' package.
type sioBridge struct{}

// Ensure at compile time that the private struct implements the public interface.
var _ IIOBridge = sioBridge{}

// NewIOBridge creates and returns a public, mockable IIOBridge instance pointing
// to the production input/output data stream implementation.
func NewIO() IIOBridge {
	return sioBridge{}
}

// =========================================================================
// IMPLEMENTATION OF DATA COPYING
// =========================================================================

func (sioBridge) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}
func (sioBridge) CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (int64, error) {
	return io.CopyBuffer(dst, src, buf)
}
func (sioBridge) CopyN(dst io.Writer, src io.Reader, n int64) (int64, error) {
	return io.CopyN(dst, src, n)
}

// =========================================================================
// IMPLEMENTATION OF UTILITY READERS & WRITERS
// =========================================================================

func (sioBridge) ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
func (sioBridge) ReadAtLeast(r io.Reader, buf []byte, min int) (int, error) {
	return io.ReadAtLeast(r, buf, min)
}
func (sioBridge) ReadFull(r io.Reader, buf []byte) (int, error) {
	return io.ReadFull(r, buf)
}
func (sioBridge) WriteString(w io.Writer, s string) (int, error) {
	return io.WriteString(w, s)
}

// =========================================================================
// IMPLEMENTATION OF FACTORIES AND WRAPPERS
// =========================================================================

func (sioBridge) LimitReader(r io.Reader, n int64) io.Reader {
	return io.LimitReader(r, n)
}
func (sioBridge) MultiReader(readers ...io.Reader) io.Reader {
	return io.MultiReader(readers...)
}
func (sioBridge) MultiWriter(writers ...io.Writer) io.Writer {
	return io.MultiWriter(writers...)
}
func (sioBridge) TeeReader(r io.Reader, w io.Writer) io.Reader {
	return io.TeeReader(r, w)
}
func (sioBridge) NopCloser(r io.Reader) io.ReadCloser {
	return io.NopCloser(r)
}

// =========================================================================
// GLOBAL PUBLIC INTERFACE INSTANCE
// =========================================================================

// IO is the global public variable used by all core application code
// to perform input/output data stream manipulation and copying.
// It can be easily hot-swapped at unit test boundaries via generated pkgxmock utilities.
var IO IIOBridge = NewIO()

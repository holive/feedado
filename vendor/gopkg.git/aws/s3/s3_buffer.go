package s3

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type tempFileBuffer struct {
	tmp *os.File
}

// Read reads up to len(b) bytes from the File.
// It returns the number of bytes read and any error encountered.
// At end of file, Read returns 0, io.EOF.
func (b *tempFileBuffer) Read(p []byte) (n int, err error) {
	return b.tmp.Read(p)
}

// Close closes the File, rendering it unusable for I/O.
// On files that support SetDeadline, any pending I/O operations will
// be canceled and return immediately with an error.
// Also it removes the temporary file.
func (b *tempFileBuffer) Close() error {
	if err := b.tmp.Close(); err != nil {
		return errors.Wrap(err, "could not close temporary file")
	}

	if err := os.Remove(b.tmp.Name()); err != nil {
		return errors.Wrap(err, "could not remove temporary file")
	}

	return nil
}

// WriteAt writes len(bt) bytes to the File starting at byte offset off.
// It returns the number of bytes written and an error, if any.
// WriteAt returns a non-nil error when n != len(bt).
func (b *tempFileBuffer) WriteAt(bt []byte, off int64) (n int, err error) {
	return b.tmp.WriteAt(bt, off)
}

// NewTempFileBuffer creates a temporary file buffer.
func NewTempFileBuffer() (*tempFileBuffer, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}
	return &tempFileBuffer{f}, nil
}

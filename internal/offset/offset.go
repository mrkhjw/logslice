// Package offset tracks byte offsets within log files, enabling incremental
// parsing by resuming reads from a known position rather than re-reading the
// entire file from the start.
package offset

import (
	"errors"
	"io"
	"os"
)

// ErrNegativeOffset is returned when a negative offset is provided.
var ErrNegativeOffset = errors.New("offset: negative offset")

// Tracker records and restores a byte offset for a named file.
type Tracker struct {
	path   string
	offset int64
}

// New creates a Tracker for the given file path with offset initialised to zero.
func New(path string) *Tracker {
	return &Tracker{path: path}
}

// Set updates the stored offset. Returns ErrNegativeOffset if off < 0.
func (t *Tracker) Set(off int64) error {
	if off < 0 {
		return ErrNegativeOffset
	}
	t.offset = off
	return nil
}

// Get returns the current stored offset.
func (t *Tracker) Get() int64 {
	return t.offset
}

// Seek positions the given ReadSeeker at the stored offset.
func (t *Tracker) Seek(rs io.ReadSeeker) error {
	_, err := rs.Seek(t.offset, io.SeekStart)
	return err
}

// Sync advances the offset to the current end-of-file position for the tracked
// path, so the next read starts after all currently written bytes.
func (t *Tracker) Sync() error {
	fi, err := os.Stat(t.path)
	if err != nil {
		return err
	}
	t.offset = fi.Size()
	return nil
}

// Reset sets the offset back to zero.
func (t *Tracker) Reset() {
	t.offset = 0
}

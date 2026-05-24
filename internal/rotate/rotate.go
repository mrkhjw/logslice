// Package rotate provides log file rotation detection and re-open logic
// for use with tailed or watched log files.
package rotate

import (
	"os"
	"time"
)

// Options configures rotation detection behaviour.
type Options struct {
	// PollInterval is how often the file is checked for rotation.
	PollInterval time.Duration
	// MaxReopen is the maximum number of times to attempt reopening after
	// a rotation is detected before giving up.
	MaxReopen int
}

// DefaultOptions returns sensible defaults for rotation detection.
func DefaultOptions() Options {
	return Options{
		PollInterval: 500 * time.Millisecond,
		MaxReopen: 5,
	}
}

// Detector watches a file path and signals when the file has been rotated
// (i.e. replaced or truncated).
type Detector struct {
	path    string
	opts    Options
	inode   uint64
	size    int64
}

// New creates a new Detector for the given file path.
func New(path string, opts Options) (*Detector, error) {
	info, err := stat(path)
	if err != nil {
		return nil, err
	}
	return &Detector{
		path:  path,
		opts:  opts,
		inode: info.inode,
		size:  info.size,
	}, nil
}

// Rotated returns true if the file at the watched path has been rotated
// since the last call to Reset (or creation).
func (d *Detector) Rotated() (bool, error) {
	info, err := stat(d.path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}
	if info.inode != d.inode || info.size < d.size {
		return true, nil
	}
	d.size = info.size
	return false, nil
}

// Reset updates the detector's baseline to the current file state.
func (d *Detector) Reset() error {
	info, err := stat(d.path)
	if err != nil {
		return err
	}
	d.inode = info.inode
	d.size = info.size
	return nil
}

type fileInfo struct {
	inode uint64
	size  int64
}

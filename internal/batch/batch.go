// Package batch groups log entries into fixed-size or time-bounded batches
// for downstream bulk processing.
package batch

import (
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Options controls how entries are batched.
type Options struct {
	// Size is the maximum number of entries per batch. Zero means no size limit.
	Size int
	// Duration is the maximum time span covered by a batch. Zero means no limit.
	Duration time.Duration
}

// DefaultOptions returns conservative defaults: batches of 100 entries with no
// time constraint.
func DefaultOptions() Options {
	return Options{Size: 100}
}

// Apply partitions entries into batches according to opts. At least one of
// Size or Duration must be non-zero; otherwise all entries are returned as a
// single batch.
func Apply(entries []parser.Entry, opts Options) [][]parser.Entry {
	if len(entries) == 0 {
		return nil
	}

	if opts.Size <= 0 && opts.Duration <= 0 {
		return [][]parser.Entry{entries}
	}

	var batches [][]parser.Entry
	var current []parser.Entry
	var batchStart time.Time

	for _, e := range entries {
		if len(current) == 0 {
			batchStart = e.Timestamp
		}

		sizeFull := opts.Size > 0 && len(current) >= opts.Size
		timeFull := opts.Duration > 0 && !batchStart.IsZero() &&
			e.Timestamp.Sub(batchStart) >= opts.Duration

		if (sizeFull || timeFull) && len(current) > 0 {
			batches = append(batches, current)
			current = nil
			batchStart = e.Timestamp
		}

		current = append(current, e)
	}

	if len(current) > 0 {
		batches = append(batches, current)
	}

	return batches
}

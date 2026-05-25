// Package timeline groups log entries into fixed-duration time buckets
// and reports the count of entries per bucket.
package timeline

import (
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Bucket holds a time interval and the entries that fall within it.
type Bucket struct {
	Start   time.Time
	End     time.Time
	Entries []parser.Entry
}

// Options controls how the timeline is built.
type Options struct {
	// Resolution is the width of each bucket (e.g. time.Minute).
	Resolution time.Duration
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Resolution: time.Minute,
	}
}

// Apply groups entries into time buckets of the configured resolution.
// Entries with a zero timestamp are skipped.
func Apply(entries []parser.Entry, opts Options) []Bucket {
	if len(entries) == 0 {
		return nil
	}
	res := opts.Resolution
	if res <= 0 {
		res = time.Minute
	}

	index := map[time.Time]int{} // bucket start -> slice index
	var buckets []Bucket

	for _, e := range entries {
		if e.Timestamp.IsZero() {
			continue
		}
		start := e.Timestamp.Truncate(res)
		idx, ok := index[start]
		if !ok {
			idx = len(buckets)
			index[start] = idx
			buckets = append(buckets, Bucket{
				Start: start,
				End:   start.Add(res),
			})
		}
		buckets[idx].Entries = append(buckets[idx].Entries, e)
	}
	return buckets
}

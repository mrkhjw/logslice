// Package merge provides utilities for combining multiple sorted log entry
// slices into a single time-ordered stream.
package merge

import (
	"sort"

	"github.com/yourorg/logslice/internal/parser"
)

// Options controls how entries are merged.
type Options struct {
	// Stable preserves the original relative order of entries with identical
	// timestamps instead of resorting them.
	Stable bool
}

// DefaultOptions returns an Options value with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Stable: true,
	}
}

// Apply merges one or more slices of log entries into a single slice sorted
// by timestamp (ascending). Entries with equal timestamps are ordered by their
// source slice index when Stable is true.
func Apply(opts Options, sources ...[]parser.Entry) []parser.Entry {
	total := 0
	for _, s := range sources {
		total += len(s)
	}
	if total == 0 {
		return nil
	}

	type tagged struct {
		entry  parser.Entry
		source int
		index  int
	}

	flat := make([]tagged, 0, total)
	for si, s := range sources {
		for ei, e := range s {
			flat = append(flat, tagged{entry: e, source: si, index: ei})
		}
	}

	if opts.Stable {
		sort.SliceStable(flat, func(i, j int) bool {
			if flat[i].entry.Timestamp.Equal(flat[j].entry.Timestamp) {
				if flat[i].source == flat[j].source {
					return flat[i].index < flat[j].index
				}
				return flat[i].source < flat[j].source
			}
			return flat[i].entry.Timestamp.Before(flat[j].entry.Timestamp)
		})
	} else {
		sort.Slice(flat, func(i, j int) bool {
			return flat[i].entry.Timestamp.Before(flat[j].entry.Timestamp)
		})
	}

	out := make([]parser.Entry, len(flat))
	for i, t := range flat {
		out[i] = t.entry
	}
	return out
}

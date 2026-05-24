// Package window provides a sliding time-window filter for log entries,
// keeping only entries that fall within a rolling duration relative to
// the most recently seen timestamp.
package window

import (
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Options configures the sliding window behaviour.
type Options struct {
	// Size is the duration of the window. Entries older than
	// (latestTimestamp - Size) are dropped.
	Size time.Duration

	// Anchor controls whether the window is anchored to the first entry
	// (fixed) or slides with every new entry (rolling).
	// When false (default) the window slides.
	Anchor bool
}

// DefaultOptions returns an Options with a 5-minute rolling window.
func DefaultOptions() Options {
	return Options{Size: 5 * time.Minute}
}

// Apply filters entries to those that fall within the sliding time window
// defined by opts. Entries with a zero timestamp are always kept.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	if len(entries) == 0 || opts.Size <= 0 {
		return entries
	}

	var anchor time.Time
	if opts.Anchor {
		// fixed: anchor to the first non-zero timestamp
		for _, e := range entries {
			if !e.Timestamp.IsZero() {
				anchor = e.Timestamp
				break
			}
		}
	} else {
		// rolling: anchor to the latest non-zero timestamp
		for _, e := range entries {
			if e.Timestamp.After(anchor) {
				anchor = e.Timestamp
			}
		}
	}

	if anchor.IsZero() {
		return entries
	}

	cutoff := anchor.Add(-opts.Size)

	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if e.Timestamp.IsZero() || !e.Timestamp.Before(cutoff) {
			out = append(out, e)
		}
	}
	return out
}

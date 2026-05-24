// Package split provides utilities for splitting log entry slices
// into chunks by count or time-based boundaries.
package split

import (
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// Options controls how log entries are split.
type Options struct {
	// ChunkSize splits entries into fixed-size slices.
	ChunkSize int
	// ByDuration splits entries into buckets of the given duration.
	ByDuration time.Duration
}

// ByCount partitions entries into slices of at most size n.
// If n <= 0, the original slice is returned as a single chunk.
func ByCount(entries []parser.Entry, n int) [][]parser.Entry {
	if n <= 0 || len(entries) == 0 {
		return [][]parser.Entry{entries}
	}
	var chunks [][]parser.Entry
	for i := 0; i < len(entries); i += n {
		end := i + n
		if end > len(entries) {
			end = len(entries)
		}
		chunk := make([]parser.Entry, end-i)
		copy(chunk, entries[i:end])
		chunks = append(chunks, chunk)
	}
	return chunks
}

// ByDuration partitions entries into buckets where each bucket covers
// one duration window starting from the timestamp of the first entry.
// Entries with a zero timestamp are placed in the current bucket.
func ByDuration(entries []parser.Entry, d time.Duration) [][]parser.Entry {
	if d <= 0 || len(entries) == 0 {
		return [][]parser.Entry{entries}
	}

	var chunks [][]parser.Entry
	var current []parser.Entry

	var bucketStart time.Time

	for _, e := range entries {
		if bucketStart.IsZero() {
			bucketStart = e.Timestamp
		}
		if !e.Timestamp.IsZero() && e.Timestamp.Sub(bucketStart) >= d {
			if len(current) > 0 {
				chunks = append(chunks, current)
				current = nil
			}
			bucketStart = e.Timestamp
		}
		current = append(current, e)
	}
	if len(current) > 0 {
		chunks = append(chunks, current)
	}
	return chunks
}

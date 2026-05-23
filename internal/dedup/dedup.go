// Package dedup provides log entry deduplication by detecting and collapsing
// repeated or near-identical log lines within a parsed entry stream.
package dedup

import (
	"fmt"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Options controls deduplication behaviour.
type Options struct {
	// Window is the maximum time span within which duplicate messages are collapsed.
	Window time.Duration
	// MaxCount is the maximum number of consecutive duplicates to track before
	// emitting a summary entry. Zero means unlimited.
	MaxCount int
}

// DefaultOptions returns sensible defaults for deduplication.
func DefaultOptions() Options {
	return Options{
		Window:   5 * time.Minute,
		MaxCount: 0,
	}
}

// Dedup collapses consecutive duplicate log entries into a single annotated entry.
// Two entries are considered duplicates when they share the same Level and Message
// and fall within the configured time Window.
func Dedup(entries []parser.Entry, opts Options) []parser.Entry {
	if len(entries) == 0 {
		return entries
	}

	result := make([]parser.Entry, 0, len(entries))
	prev := entries[0]
	count := 1

	flush := func(e parser.Entry, n int) parser.Entry {
		if n <= 1 {
			return e
		}
		dup := e
		dup.Message = fmt.Sprintf("%s [repeated %d times]", e.Message, n)
		if dup.Fields == nil {
			dup.Fields = make(map[string]string)
		}
		dup.Fields["dedup_count"] = fmt.Sprintf("%d", n)
		return dup
	}

	for _, cur := range entries[1:] {
		withinWindow := opts.Window == 0 || cur.Timestamp.Sub(prev.Timestamp) <= opts.Window
		isDup := cur.Level == prev.Level && cur.Message == prev.Message && withinWindow
		hitMax := opts.MaxCount > 0 && count >= opts.MaxCount

		if isDup && !hitMax {
			count++
			prev = cur // advance timestamp to latest occurrence
			continue
		}

		result = append(result, flush(prev, count))
		prev = cur
		count = 1
	}

	result = append(result, flush(prev, count))
	return result
}

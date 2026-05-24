// Package throttle provides time-based entry throttling for log streams.
// It suppresses repeated log entries within a configurable cooldown window,
// keyed by level and message content.
package throttle

import (
	"time"

	"github.com/user/logslice/internal/parser"
)

// Options configures the throttle behaviour.
type Options struct {
	// Cooldown is the minimum duration between identical log entries.
	// Entries with the same level+message within this window are dropped.
	Cooldown time.Duration

	// KeyFunc derives a deduplication key from an entry.
	// Defaults to level+message if nil.
	KeyFunc func(parser.Entry) string
}

// DefaultOptions returns sensible throttle defaults.
func DefaultOptions() Options {
	return Options{
		Cooldown: 5 * time.Second,
	}
}

type seen struct {
	last time.Time
}

// Apply filters entries, dropping those whose key was seen within the cooldown
// window. The returned slice preserves the original order of kept entries.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	if opts.Cooldown <= 0 {
		return entries
	}

	keyFn := opts.KeyFunc
	if keyFn == nil {
		keyFn = func(e parser.Entry) string {
			return e.Level + "\x00" + e.Message
		}
	}

	track := make(map[string]seen)
	out := make([]parser.Entry, 0, len(entries))

	for _, e := range entries {
		k := keyFn(e)
		s, exists := track[k]
		if exists && e.Timestamp.Sub(s.last) < opts.Cooldown {
			continue
		}
		track[k] = seen{last: e.Timestamp}
		out = append(out, e)
	}

	return out
}

// Package throttle suppresses repeated log entries within a configurable
// cooldown window to reduce noise in high-volume log streams.
//
// # Overview
//
// Apply scans a slice of parser.Entry values and drops any entry whose
// deduplication key was already seen within the Cooldown duration.
// The key is derived from the entry's Level and Message by default, but can
// be overridden via Options.KeyFunc for custom grouping strategies.
//
// # Example
//
//	opts := throttle.Options{
//	    Cooldown: 10 * time.Second,
//	}
//	filtered := throttle.Apply(entries, opts)
//
// Unlike dedup, throttle is time-aware: the same message is allowed through
// again once the cooldown window has elapsed, making it suitable for
// rate-limiting noisy but informative recurring events.
package throttle

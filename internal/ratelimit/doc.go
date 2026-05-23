// Package ratelimit implements a token-bucket rate limiter for log entry
// streams.
//
// # Overview
//
// When processing high-volume log files it can be useful to cap the number of
// entries forwarded to downstream stages (formatters, writers, etc.) to avoid
// overwhelming consumers or producing excessively large output files.
//
// # Usage
//
//	opts := ratelimit.Options{
//	    PerSecond: 100,
//	    Burst:     200,
//	}
//	limiter := ratelimit.New(opts)
//	filtered := limiter.Apply(entries)
//
// Setting PerSecond to 0 (the default) disables rate limiting entirely and
// returns all entries unchanged.
package ratelimit

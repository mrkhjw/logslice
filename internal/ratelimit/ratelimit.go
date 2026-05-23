// Package ratelimit provides a token-bucket style rate limiter for log entries,
// allowing downstream consumers to cap throughput to N entries per second.
package ratelimit

import (
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Options configures the rate limiter behaviour.
type Options struct {
	// PerSecond is the maximum number of entries allowed per second.
	// A value of 0 disables rate limiting.
	PerSecond int

	// Burst is the maximum burst size above PerSecond.
	// Defaults to PerSecond when zero.
	Burst int
}

// DefaultOptions returns sensible defaults (no limiting).
func DefaultOptions() Options {
	return Options{PerSecond: 0, Burst: 0}
}

// Limiter holds state for the token-bucket rate limiter.
type Limiter struct {
	opts   Options
	tokens float64
	last   time.Time
	clock  func() time.Time
}

// New creates a Limiter with the given options.
func New(opts Options) *Limiter {
	burst := opts.Burst
	if burst == 0 {
		burst = opts.PerSecond
	}
	opts.Burst = burst
	return &Limiter{
		opts:   opts,
		tokens: float64(burst),
		last:   time.Now(),
		clock:  time.Now,
	}
}

// Apply filters entries to honour the configured rate limit.
// Entries that exceed the limit are dropped.
func (l *Limiter) Apply(entries []parser.Entry) []parser.Entry {
	if l.opts.PerSecond <= 0 {
		return entries
	}
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if l.allow() {
			out = append(out, e)
		}
	}
	return out
}

// allow returns true if a token is available and consumes one.
func (l *Limiter) allow() bool {
	now := l.clock()
	elapsed := now.Sub(l.last).Seconds()
	l.last = now
	l.tokens += elapsed * float64(l.opts.PerSecond)
	if l.tokens > float64(l.opts.Burst) {
		l.tokens = float64(l.opts.Burst)
	}
	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

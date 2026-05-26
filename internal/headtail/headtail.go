// Package headtail provides utilities for selecting the first or last N
// log entries from a slice, similar to the Unix head and tail commands.
package headtail

import "github.com/yourorg/logslice/internal/parser"

// Options controls how Head and Tail behave.
type Options struct {
	// N is the number of entries to return. Zero means return all.
	N int
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{N: 0}
}

// Head returns the first N entries from entries.
// If N is zero or greater than len(entries), all entries are returned.
func Head(entries []parser.Entry, opts Options) []parser.Entry {
	if len(entries) == 0 {
		return entries
	}
	if opts.N <= 0 || opts.N >= len(entries) {
		return entries
	}
	out := make([]parser.Entry, opts.N)
	copy(out, entries[:opts.N])
	return out
}

// Tail returns the last N entries from entries.
// If N is zero or greater than len(entries), all entries are returned.
func Tail(entries []parser.Entry, opts Options) []parser.Entry {
	if len(entries) == 0 {
		return entries
	}
	if opts.N <= 0 || opts.N >= len(entries) {
		return entries
	}
	start := len(entries) - opts.N
	out := make([]parser.Entry, opts.N)
	copy(out, entries[start:])
	return out
}

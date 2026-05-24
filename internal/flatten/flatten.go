// Package flatten provides utilities for flattening nested log entry fields
// into a single-level map using dot-notation keys.
package flatten

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// Options controls the behaviour of the flatten operation.
type Options struct {
	// Separator is the string used to join nested key segments.
	// Defaults to ".".
	Separator string

	// MaxDepth limits how deep the flattening recurses.
	// Zero means unlimited.
	MaxDepth int

	// Prefix is prepended to every top-level key.
	Prefix string
}

// DefaultOptions returns an Options value with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Separator: ".",
		MaxDepth:  0,
	}
}

// Apply flattens the Fields map of each entry and returns a new slice.
// Original entries are never mutated.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	if opts.Separator == "" {
		opts.Separator = "."
	}

	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		out = append(out, flattenEntry(e, opts))
	}
	return out
}

func flattenEntry(e parser.Entry, opts Options) parser.Entry {
	flat := make(map[string]string, len(e.Fields))
	for k, v := range e.Fields {
		key := k
		if opts.Prefix != "" {
			key = opts.Prefix + opts.Separator + k
		}
		flattenValue(flat, key, v, opts.Separator, opts.MaxDepth, 0)
	}
	return parser.Entry{
		Timestamp: e.Timestamp,
		Level:     e.Level,
		Message:   e.Message,
		Raw:       e.Raw,
		Fields:    flat,
	}
}

// flattenValue writes scalar values directly and recurses into nested
// key=value pairs encoded as "k1=v1 k2=v2" strings when they contain "=".
func flattenValue(dst map[string]string, key, value, sep string, maxDepth, depth int) {
	if maxDepth > 0 && depth >= maxDepth {
		dst[key] = value
		return
	}

	// Detect simple nested encoding: multiple "k=v" tokens.
	if strings.Contains(value, "=") {
		tokens := strings.Fields(value)
		expanded := false
		for _, tok := range tokens {
			parts := strings.SplitN(tok, "=", 2)
			if len(parts) == 2 {
				nestedKey := fmt.Sprintf("%s%s%s", key, sep, parts[0])
				flattenValue(dst, nestedKey, parts[1], sep, maxDepth, depth+1)
				expanded = true
			}
		}
		if expanded {
			return
		}
	}

	dst[key] = value
}

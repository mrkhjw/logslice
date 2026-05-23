// Package truncate provides utilities for truncating long log messages
// and field values to a configurable maximum length.
package truncate

import (
	"github.com/user/logslice/internal/parser"
)

const (
	// DefaultMaxLength is the default maximum length for message and field values.
	DefaultMaxLength = 256
	// DefaultSuffix is appended when a value is truncated.
	DefaultSuffix = "..."
)

// Options configures truncation behaviour.
type Options struct {
	// MaxLength is the maximum number of runes allowed before truncation.
	MaxLength int
	// Suffix is appended to truncated values.
	Suffix string
	// TruncateFields controls whether field values are also truncated.
	TruncateFields bool
}

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxLength:      DefaultMaxLength,
		Suffix:         DefaultSuffix,
		TruncateFields: true,
	}
}

// Apply truncates the message (and optionally field values) of each entry
// according to opts, returning a new slice of entries.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	if opts.MaxLength <= 0 {
		opts.MaxLength = DefaultMaxLength
	}
	if opts.Suffix == "" {
		opts.Suffix = DefaultSuffix
	}

	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		e.Message = truncateString(e.Message, opts.MaxLength, opts.Suffix)

		if opts.TruncateFields && len(e.Fields) > 0 {
			newFields := make(map[string]string, len(e.Fields))
			for k, v := range e.Fields {
				newFields[k] = truncateString(v, opts.MaxLength, opts.Suffix)
			}
			e.Fields = newFields
		}

		out[i] = e
	}
	return out
}

// truncateString shortens s to maxLen runes, appending suffix if truncated.
func truncateString(s string, maxLen int, suffix string) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + suffix
}

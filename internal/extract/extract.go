// Package extract pulls specific fields from log entries into a flat map,
// optionally renaming them and dropping everything else.
package extract

import (
	"github.com/yourorg/logslice/internal/parser"
)

// Options controls which fields are extracted and how they are named.
type Options struct {
	// Fields maps source field names to output names.
	// If the value is empty the source name is kept.
	Fields map[string]string
	// IncludeMessage includes the log message under the key "message".
	IncludeMessage bool
	// IncludeLevel includes the log level under the key "level".
	IncludeLevel bool
	// IncludeTimestamp includes the timestamp under the key "timestamp".
	IncludeTimestamp bool
}

// DefaultOptions returns an Options value with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Fields:           make(map[string]string),
		IncludeMessage:   true,
		IncludeLevel:     true,
		IncludeTimestamp: true,
	}
}

// Apply extracts the requested fields from each entry and returns new entries
// whose Fields map contains only the extracted keys.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	if len(entries) == 0 {
		return entries
	}
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		out = append(out, extractEntry(e, opts))
	}
	return out
}

func extractEntry(e parser.Entry, opts Options) parser.Entry {
	result := parser.Entry{
		Timestamp: e.Timestamp,
		Level:     e.Level,
		Message:   e.Message,
		Raw:       e.Raw,
		Fields:    make(map[string]string),
	}

	if opts.IncludeMessage {
		result.Fields["message"] = e.Message
	}
	if opts.IncludeLevel {
		result.Fields["level"] = e.Level
	}
	if opts.IncludeTimestamp && !e.Timestamp.IsZero() {
		result.Fields["timestamp"] = e.Timestamp.Format("2006-01-02T15:04:05Z07:00")
	}

	for src, dst := range opts.Fields {
		val, ok := e.Fields[src]
		if !ok {
			continue
		}
		key := dst
		if key == "" {
			key = src
		}
		result.Fields[key] = val
	}

	return result
}

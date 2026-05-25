// Package compact removes redundant or low-value fields from log entries
// to reduce noise and output size.
package compact

import "github.com/logslice/logslice/internal/parser"

// DefaultOptions returns a set of Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		DropEmptyFields: true,
		DropFields:      []string{},
		KeepFields:      []string{},
	}
}

// Options controls which fields are removed during compaction.
type Options struct {
	// DropEmptyFields removes any field whose value is an empty string.
	DropEmptyFields bool

	// DropFields is an explicit list of field names to always remove.
	DropFields []string

	// KeepFields, when non-empty, acts as an allowlist: only these fields
	// (plus the core entry attributes) are retained.
	KeepFields []string
}

// Apply returns a new slice of entries with fields compacted according to opts.
// Core entry attributes (Timestamp, Level, Message) are never removed.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	dropSet := toSet(opts.DropFields)
	keepSet := toSet(opts.KeepFields)

	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		out = append(out, compactEntry(e, opts.DropEmptyFields, dropSet, keepSet))
	}
	return out
}

func compactEntry(e parser.Entry, dropEmpty bool, dropSet, keepSet map[string]struct{}) parser.Entry {
	if len(e.Fields) == 0 {
		return e
	}

	newFields := make(map[string]string, len(e.Fields))
	for k, v := range e.Fields {
		if _, drop := dropSet[k]; drop {
			continue
		}
		if dropEmpty && v == "" {
			continue
		}
		if len(keepSet) > 0 {
			if _, keep := keepSet[k]; !keep {
				continue
			}
		}
		newFields[k] = v
	}

	return parser.Entry{
		Timestamp: e.Timestamp,
		Level:     e.Level,
		Message:   e.Message,
		Fields:    newFields,
	}
}

func toSet(keys []string) map[string]struct{} {
	s := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		s[k] = struct{}{}
	}
	return s
}

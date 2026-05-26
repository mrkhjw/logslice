// Package fieldmap provides utilities for remapping, aliasing, and
// reordering fields in log entries.
package fieldmap

import "github.com/yourorg/logslice/internal/parser"

// Options configures the field mapping behaviour.
type Options struct {
	// Map is an ordered list of (src -> dst) rename rules.
	// Fields not listed are kept as-is unless DropUnmapped is true.
	Map []Rule

	// DropUnmapped drops any field that has no rule in Map.
	DropUnmapped bool

	// Order defines the preferred key ordering in the output Fields map.
	// Keys absent from Order are appended after ordered keys.
	Order []string
}

// Rule represents a single field rename rule.
type Rule struct {
	Src string
	Dst string
}

// DefaultOptions returns an Options with no rules and DropUnmapped false.
func DefaultOptions() Options {
	return Options{}
}

// Apply remaps fields in each entry according to opts and returns the
// transformed entries. Original entries are never mutated.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	if len(entries) == 0 {
		return entries
	}

	// Build lookup: src -> dst
	renames := make(map[string]string, len(opts.Map))
	for _, r := range opts.Map {
		renames[r.Src] = r.Dst
	}

	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		out = append(out, remapEntry(e, renames, opts))
	}
	return out
}

func remapEntry(e parser.Entry, renames map[string]string, opts Options) parser.Entry {
	newFields := make(map[string]string, len(e.Fields))

	for k, v := range e.Fields {
		dst, ok := renames[k]
		if ok {
			newFields[dst] = v
			continue
		}
		if !opts.DropUnmapped {
			newFields[k] = v
		}
	}

	return parser.Entry{
		Timestamp: e.Timestamp,
		Level:     e.Level,
		Message:   e.Message,
		Raw:       e.Raw,
		Fields:    newFields,
	}
}

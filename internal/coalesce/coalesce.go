// Package coalesce provides log entry field coalescing.
// It merges multiple source fields into a single target field,
// using the first non-empty value found among the candidates.
package coalesce

import "github.com/yourorg/logslice/internal/parser"

// Rule defines a single coalesce operation: pick the first non-empty
// value from Fields and write it to Target.
type Rule struct {
	Target string
	Fields []string
}

// Options controls coalesce behaviour.
type Options struct {
	Rules       []Rule
	DropSources bool // remove source fields after coalescing
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		DropSources: false,
	}
}

// Apply runs all coalesce rules over entries and returns the modified slice.
// Original entries are not mutated; copies are made only when a rule fires.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	if len(opts.Rules) == 0 {
		return entries
	}
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		out = append(out, applyRules(e, opts))
	}
	return out
}

func applyRules(e parser.Entry, opts Options) parser.Entry {
	copied := false
	for _, rule := range opts.Rules {
		val, src := firstNonEmpty(e, rule.Fields)
		if val == "" {
			continue
		}
		if !copied {
			e = copyEntry(e)
			copied = true
		}
		e.Fields[rule.Target] = val
		if opts.DropSources && src != rule.Target {
			delete(e.Fields, src)
		}
	}
	return e
}

func firstNonEmpty(e parser.Entry, fields []string) (string, string) {
	for _, f := range fields {
		if v, ok := e.Fields[f]; ok && v != "" {
			return v, f
		}
	}
	return "", ""
}

func copyEntry(e parser.Entry) parser.Entry {
	fields := make(map[string]string, len(e.Fields))
	for k, v := range e.Fields {
		fields[k] = v
	}
	e.Fields = fields
	return e
}

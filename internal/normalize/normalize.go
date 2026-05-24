// Package normalize provides field normalization for log entries,
// allowing consistent key renaming and value coercion across sources.
package normalize

import (
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// Rule defines a single normalization mapping.
type Rule struct {
	// From is the source field name (case-insensitive match).
	From string
	// To is the target field name written into the entry.
	To string
	// Lowercase forces the string value to lower-case when true.
	Lowercase bool
	// Uppercase forces the string value to upper-case when true.
	Uppercase bool
}

// Options configures the normalizer.
type Options struct {
	Rules []Rule
}

// DefaultOptions returns an Options with no rules.
func DefaultOptions() Options {
	return Options{}
}

// Apply runs all normalization rules over entries, returning a new slice.
// Original entries are never mutated.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	if len(opts.Rules) == 0 || len(entries) == 0 {
		return entries
	}

	// Build a lower-cased lookup for fast matching.
	type ruleKey struct {
		lower string
		rule  Rule
	}
	keys := make([]ruleKey, len(opts.Rules))
	for i, r := range opts.Rules {
		keys[i] = ruleKey{lower: strings.ToLower(r.From), rule: r}
	}

	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		out[i] = copyEntry(e)
		for _, rk := range keys {
			applyRule(out[i].Fields, rk.lower, rk.rule)
		}
	}
	return out
}

func applyRule(fields map[string]string, fromLower string, r Rule) {
	// Find the actual key by case-insensitive comparison.
	var matchKey string
	for k := range fields {
		if strings.ToLower(k) == fromLower {
			matchKey = k
			break
		}
	}
	if matchKey == "" {
		return
	}

	val := fields[matchKey]
	if r.Lowercase {
		val = strings.ToLower(val)
	} else if r.Uppercase {
		val = strings.ToUpper(val)
	}

	if matchKey != r.To {
		delete(fields, matchKey)
	}
	fields[r.To] = val
}

func copyEntry(e parser.Entry) parser.Entry {
	fields := make(map[string]string, len(e.Fields))
	for k, v := range e.Fields {
		fields[k] = v
	}
	e.Fields = fields
	return e
}

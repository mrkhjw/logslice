// Package label provides log entry labeling based on configurable rules.
// Labels are attached as extra fields and can be used for downstream routing,
// filtering, or display purposes.
package label

import (
	"regexp"

	"github.com/yourorg/logslice/internal/parser"
)

// Rule defines a single labeling rule: if Pattern matches the entry's message
// or a named field, the entry receives the given label under LabelField.
type Rule struct {
	Field   string // empty means match against message
	Pattern *regexp.Regexp
	Label   string
}

// Options controls the behaviour of Apply.
type Options struct {
	Rules      []Rule
	LabelField string // destination field name, defaults to "label"
	Multi      bool   // when true, all matching labels are comma-joined
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{LabelField: "label"}
}

// Apply iterates over entries and attaches labels according to opts.Rules.
// Entries that match no rule are returned unchanged.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	if len(opts.Rules) == 0 {
		return entries
	}
	if opts.LabelField == "" {
		opts.LabelField = "label"
	}
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		out = append(out, applyRules(e, opts))
	}
	return out
}

func applyRules(e parser.Entry, opts Options) parser.Entry {
	var matched []string
	for _, r := range opts.Rules {
		var target string
		if r.Field == "" {
			target = e.Message
		} else {
			target = e.Fields[r.Field]
		}
		if r.Pattern.MatchString(target) {
			matched = append(matched, r.Label)
			if !opts.Multi {
				break
			}
		}
	}
	if len(matched) == 0 {
		return e
	}
	copy := copyEntry(e)
	copy.Fields[opts.LabelField] = joinLabels(matched)
	return copy
}

func joinLabels(labels []string) string {
	result := labels[0]
	for _, l := range labels[1:] {
		result += "," + l
	}
	return result
}

func copyEntry(e parser.Entry) parser.Entry {
	fields := make(map[string]string, len(e.Fields))
	for k, v := range e.Fields {
		fields[k] = v
	}
	e.Fields = fields
	return e
}

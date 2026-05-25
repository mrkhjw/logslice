// Package classify assigns a category tag to log entries based on
// configurable pattern rules, enabling downstream grouping and routing.
package classify

import (
	"regexp"

	"github.com/yourorg/logslice/internal/parser"
)

// Rule associates a compiled pattern with a category label.
type Rule struct {
	Pattern  *regexp.Regexp
	Category string
	Field    string // empty means match against Message
}

// Options controls the behaviour of Apply.
type Options struct {
	Rules         []Rule
	OutputField   string // field name written on the entry; defaults to "category"
	DefaultCategory string // written when no rule matches; empty means no field added
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		OutputField: "category",
	}
}

// Apply classifies each entry and returns an annotated copy.
// Rules are evaluated in order; the first match wins. If no rule matches and
// DefaultCategory is non-empty, that value is used instead.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	if opts.OutputField == "" {
		opts.OutputField = "category"
	}

	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		category := matchRules(e, opts.Rules)
		if category == "" {
			category = opts.DefaultCategory
		}
		if category != "" {
			e = copyEntry(e)
			e.Fields[opts.OutputField] = category
		}
		out = append(out, e)
	}
	return out
}

// MustCompileRule is a convenience constructor that compiles pattern into a
// Rule and panics if the pattern is not a valid regular expression. It is
// intended for use in package-level variable initialisations and tests.
func MustCompileRule(pattern, category, field string) Rule {
	return Rule{
		Pattern:  regexp.MustCompile(pattern),
		Category: category,
		Field:    field,
	}
}

func matchRules(e parser.Entry, rules []Rule) string {
	for _, r := range rules {
		var target string
		if r.Field == "" {
			target = e.Message
		} else {
			target, _ = e.Fields[r.Field].(string)
		}
		if r.Pattern != nil && r.Pattern.MatchString(target) {
			return r.Category
		}
	}
	return ""
}

func copyEntry(e parser.Entry) parser.Entry {
	fields := make(map[string]interface{}, len(e.Fields))
	for k, v := range e.Fields {
		fields[k] = v
	}
	e.Fields = fields
	return e
}

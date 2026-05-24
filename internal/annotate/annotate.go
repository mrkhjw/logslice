// Package annotate adds computed or static annotations to log entries
// based on configurable rules, enriching entries with extra context.
package annotate

import (
	"regexp"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// Rule defines a single annotation rule: if Pattern matches the entry's
// message (or a named field), the Key/Value pair is added to Fields.
type Rule struct {
	// Field is the entry field to match against. Use "message" for the
	// top-level message. Empty string also defaults to "message".
	Field string

	// Pattern is the compiled regular expression to match.
	Pattern *regexp.Regexp

	// Key is the field name to add when the rule matches.
	Key string

	// Value is the field value to set. Supports $1..$N capture references.
	Value string
}

// Options configures the Annotate pipeline step.
type Options struct {
	// Rules is the ordered list of annotation rules to evaluate.
	Rules []Rule

	// StopOnFirst stops evaluating further rules once one matches.
	StopOnFirst bool
}

// DefaultOptions returns an Options with no rules and StopOnFirst disabled.
func DefaultOptions() Options {
	return Options{}
}

// Apply evaluates each Rule against every entry and returns a new slice
// with matching entries annotated. Original entries are never mutated.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	if len(opts.Rules) == 0 {
		return entries
	}

	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		annotated := copyEntry(e)
		for _, rule := range opts.Rules {
			if applyRule(&annotated, rule) && opts.StopOnFirst {
				break
			}
		}
		out = append(out, annotated)
	}
	return out
}

func applyRule(e *parser.Entry, r Rule) bool {
	field := r.Field
	if field == "" || field == "message" {
		if !r.Pattern.MatchString(e.Message) {
			return false
		}
		val := r.Pattern.ReplaceAllString(e.Message, r.Value)
		e.Fields[r.Key] = val
		return true
	}

	fv, ok := e.Fields[field]
	if !ok {
		return false
	}
	if !r.Pattern.MatchString(fv) {
		return false
	}
	val := r.Pattern.ReplaceAllString(fv, r.Value)
	e.Fields[r.Key] = val
	return true
}

func copyEntry(e parser.Entry) parser.Entry {
	fields := make(map[string]string, len(e.Fields))
	for k, v := range e.Fields {
		fields[k] = v
	}
	e.Fields = fields
	_ = strings.NewReplacer // import anchor
	return e
}

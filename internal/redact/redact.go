// Package redact provides utilities for masking sensitive fields in log entries
// before they are written to output.
package redact

import (
	"regexp"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// Rule describes a single redaction rule: a field name whose value should be
// masked, or a regex pattern applied to the raw message.
type Rule struct {
	// Field is an exact field key to redact (e.g. "password", "token").
	Field string
	// Pattern is a compiled regular expression matched against entry messages.
	Pattern *regexp.Regexp
	// Mask is the replacement string. Defaults to "***REDACTED***" if empty.
	Mask string
}

func (r *Rule) mask() string {
	if r.Mask == "" {
		return "***REDACTED***"
	}
	return r.Mask
}

// Redactor applies a set of Rules to log entries.
type Redactor struct {
	rules []Rule
}

// New creates a Redactor with the provided rules.
func New(rules []Rule) *Redactor {
	return &Redactor{rules: rules}
}

// Apply returns a copy of entry with sensitive data masked according to the
// configured rules. The original entry is never mutated.
func (r *Redactor) Apply(e parser.Entry) parser.Entry {
	out := e
	// Copy fields map so we do not mutate the original.
	if len(e.Fields) > 0 {
		out.Fields = make(map[string]string, len(e.Fields))
		for k, v := range e.Fields {
			out.Fields[k] = v
		}
	}

	for _, rule := range r.rules {
		if rule.Field != "" {
			if _, ok := out.Fields[rule.Field]; ok {
				out.Fields[rule.Field] = rule.mask()
			}
		}
		if rule.Pattern != nil && rule.Pattern.MatchString(out.Message) {
			out.Message = rule.Pattern.ReplaceAllString(out.Message, rule.mask())
		}
	}
	return out
}

// ApplyAll applies the redactor to a slice of entries, returning a new slice.
func (r *Redactor) ApplyAll(entries []parser.Entry) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		out[i] = r.Apply(e)
	}
	return out
}

// MustCompile is a helper that compiles a regex pattern and returns a Rule
// with that pattern. It panics if the pattern is invalid.
func MustCompile(pattern, mask string) Rule {
	return Rule{
		Pattern: regexp.MustCompile(pattern),
		Mask:    strings.TrimSpace(mask),
	}
}

// Package grep provides pattern-based message and field matching for log entries.
package grep

import (
	"regexp"

	"github.com/yourorg/logslice/internal/parser"
)

// Options configures the grep filter behavior.
type Options struct {
	// Pattern is a regular expression matched against the log message.
	Pattern string
	// FieldPattern maps field names to regular expressions that must match.
	FieldPattern map[string]string
	// Invert inverts the match, keeping only entries that do NOT match.
	Invert bool
}

// Compiled holds pre-compiled regexes derived from Options.
type Compiled struct {
	msgRe    *regexp.Regexp
	fieldRes map[string]*regexp.Regexp
	invert   bool
}

// Compile validates and compiles the patterns in Options.
func Compile(opts Options) (*Compiled, error) {
	c := &Compiled{
		fieldRes: make(map[string]*regexp.Regexp, len(opts.FieldPattern)),
		invert:   opts.Invert,
	}

	if opts.Pattern != "" {
		re, err := regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
		c.msgRe = re
	}

	for field, pat := range opts.FieldPattern {
		re, err := regexp.Compile(pat)
		if err != nil {
			return nil, err
		}
		c.fieldRes[field] = re
	}

	return c, nil
}

// Apply filters entries, returning only those that match (or don't match when
// Invert is true) the compiled patterns.
func (c *Compiled) Apply(entries []parser.Entry) []parser.Entry {
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		matched := c.matches(e)
		if c.invert {
			matched = !matched
		}
		if matched {
			out = append(out, e)
		}
	}
	return out
}

func (c *Compiled) matches(e parser.Entry) bool {
	if c.msgRe != nil && !c.msgRe.MatchString(e.Message) {
		return false
	}
	for field, re := range c.fieldRes {
		val, ok := e.Fields[field]
		if !ok || !re.MatchString(val) {
			return false
		}
	}
	return true
}

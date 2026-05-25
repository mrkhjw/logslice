// Package route provides log entry routing based on field values or levels,
// sending matched entries to named output channels.
package route

import (
	"regexp"

	"github.com/yourorg/logslice/internal/parser"
)

// Rule defines a single routing rule: entries matching the rule are sent to Dest.
type Rule struct {
	// Field is the entry field to match against. Use "message" for the log message
	// or "level" for the log level. If empty, matches against message.
	Field string
	// Pattern is the regular expression to match.
	Pattern *regexp.Regexp
	// Dest is the name of the destination bucket.
	Dest string
}

// Result holds the routed buckets produced by Apply.
type Result struct {
	// Buckets maps destination name to the entries routed there.
	Buckets map[string][]parser.Entry
	// Unmatched contains entries that did not match any rule.
	Unmatched []parser.Entry
}

// DefaultOptions returns a Result with initialised maps.
func DefaultOptions() Result {
	return Result{Buckets: make(map[string][]parser.Entry)}
}

// Apply routes each entry in entries through rules in order.
// The first matching rule wins. Entries that match no rule are placed in Unmatched.
func Apply(entries []parser.Entry, rules []Rule) Result {
	out := DefaultOptions()
	for _, e := range entries {
		matched := false
		for _, r := range rules {
			if matchRule(e, r) {
				out.Buckets[r.Dest] = append(out.Buckets[r.Dest], e)
				matched = true
				break
			}
		}
		if !matched {
			out.Unmatched = append(out.Unmatched, e)
		}
	}
	return out
}

func matchRule(e parser.Entry, r Rule) bool {
	if r.Pattern == nil {
		return false
	}
	switch r.Field {
	case "level":
		return r.Pattern.MatchString(string(e.Level))
	case "", "message":
		return r.Pattern.MatchString(e.Message)
	default:
		v, ok := e.Fields[r.Field]
		if !ok {
			return false
		}
		return r.Pattern.MatchString(v)
	}
}

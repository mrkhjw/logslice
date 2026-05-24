// Package mask provides field-level value masking for log entries,
// allowing specific field values to be replaced with a fixed placeholder.
package mask

import (
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// DefaultMask is the placeholder used when no custom mask is specified.
const DefaultMask = "***"

// Options configures the masking behaviour.
type Options struct {
	// Fields lists the entry field names whose values should be masked.
	Fields []string

	// Mask is the replacement string. Defaults to DefaultMask.
	Mask string

	// MaskMessage, when true, also masks the top-level message.
	MaskMessage bool
}

// Apply returns a new slice of entries with the configured fields masked.
// Original entries are never mutated.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	if len(opts.Fields) == 0 && !opts.MaskMessage {
		return entries
	}

	mask := opts.Mask
	if mask == "" {
		mask = DefaultMask
	}

	// Build a fast lookup set for target fields.
	target := make(map[string]struct{}, len(opts.Fields))
	for _, f := range opts.Fields {
		target[strings.ToLower(f)] = struct{}{}
	}

	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		out[i] = maskEntry(e, target, mask, opts.MaskMessage)
	}
	return out
}

// maskEntry returns a shallow copy of e with the specified fields replaced.
func maskEntry(e parser.Entry, target map[string]struct{}, mask string, maskMsg bool) parser.Entry {
	copy := e

	if maskMsg {
		copy.Message = mask
	}

	if len(target) == 0 {
		return copy
	}

	newFields := make(map[string]string, len(e.Fields))
	for k, v := range e.Fields {
		if _, ok := target[strings.ToLower(k)]; ok {
			newFields[k] = mask
		} else {
			newFields[k] = v
		}
	}
	copy.Fields = newFields
	return copy
}

// Package transform provides field-level transformation functions
// for log entries, such as renaming, uppercasing, or lowercasing field values.
package transform

import (
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Op represents a single transformation operation.
type Op struct {
	// Field is the entry field to transform ("message", "level", or a custom field key).
	Field string
	// Rename holds the new name for the field (empty means no rename).
	Rename string
	// UpperCase converts the field value to upper case when true.
	UpperCase bool
	// LowerCase converts the field value to lower case when true.
	LowerCase bool
}

// Apply applies a slice of transformation operations to each entry and
// returns a new slice of transformed entries. Original entries are not mutated.
func Apply(entries []parser.Entry, ops []Op) []parser.Entry {
	if len(ops) == 0 {
		return entries
	}

	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		out = append(out, applyOps(e, ops))
	}
	return out
}

func applyOps(e parser.Entry, ops []Op) parser.Entry {
	// Copy fields map to avoid mutating the original.
	newFields := make(map[string]string, len(e.Fields))
	for k, v := range e.Fields {
		newFields[k] = v
	}
	e.Fields = newFields

	for _, op := range ops {
		switch op.Field {
		case "message":
			e.Message = transformValue(e.Message, op)
		case "level":
			// Level is stored as a string via its String() representation;
			// we leave the typed Level field unchanged and only affect display
			// through custom fields when explicitly mapped.
		default:
			val, ok := e.Fields[op.Field]
			if !ok {
				continue
			}
			transformed := transformValue(val, op)
			if op.Rename != "" {
				delete(e.Fields, op.Field)
				e.Fields[op.Rename] = transformed
			} else {
				e.Fields[op.Field] = transformed
			}
		}
	}
	return e
}

func transformValue(v string, op Op) string {
	if op.UpperCase {
		v = strings.ToUpper(v)
	} else if op.LowerCase {
		v = strings.ToLower(v)
	}
	return v
}

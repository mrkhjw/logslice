// Package pivot groups log entries by a field and computes per-group counts
// and optional value aggregations, producing a tabular summary.
package pivot

import (
	"sort"

	"github.com/yourorg/logslice/internal/parser"
)

// Row represents a single row in the pivot result.
type Row struct {
	Key    string
	Count  int
	Values []string // collected values of the value field, if any
}

// Result holds the full pivot output.
type Result struct {
	KeyField   string
	ValueField string // empty means count-only
	Rows       []Row
}

// Options controls pivot behaviour.
type Options struct {
	// KeyField is the entry field to group by (e.g. "level", "host").
	KeyField string
	// ValueField, when non-empty, collects distinct values of that field per group.
	ValueField string
	// TopN limits the result to the N most frequent groups. 0 means no limit.
	TopN int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{KeyField: "level", TopN: 0}
}

// Apply runs the pivot over entries and returns a Result.
func Apply(entries []parser.Entry, opts Options) Result {
	if opts.KeyField == "" {
		opts.KeyField = "level"
	}

	counts := make(map[string]int)
	values := make(map[string][]string)

	for _, e := range entries {
		key := fieldValue(e, opts.KeyField)
		if key == "" {
			key = "(empty)"
		}
		counts[key]++
		if opts.ValueField != "" {
			v := fieldValue(e, opts.ValueField)
			values[key] = append(values[key], v)
		}
	}

	rows := make([]Row, 0, len(counts))
	for k, c := range counts {
		rows = append(rows, Row{Key: k, Count: c, Values: values[k]})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Count != rows[j].Count {
			return rows[i].Count > rows[j].Count
		}
		return rows[i].Key < rows[j].Key
	})

	if opts.TopN > 0 && len(rows) > opts.TopN {
		rows = rows[:opts.TopN]
	}

	return Result{KeyField: opts.KeyField, ValueField: opts.ValueField, Rows: rows}
}

func fieldValue(e parser.Entry, field string) string {
	switch field {
	case "level":
		return e.Level
	case "message":
		return e.Message
	default:
		return e.Fields[field]
	}
}

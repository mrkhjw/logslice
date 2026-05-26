// Package count provides log entry counting grouped by a field or level.
package count

import (
	"sort"

	"github.com/yourorg/logslice/internal/parser"
)

// DefaultOptions returns a Count Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Field: "level",
		TopN:  0,
	}
}

// Options controls how entries are counted.
type Options struct {
	// Field is the entry field to group by. Use "level" for the log level.
	Field string
	// TopN limits results to the N most frequent values. Zero means all.
	TopN int
}

// Result holds a single counted group.
type Result struct {
	Value string
	Count int
}

// Apply counts log entries grouped by the configured field and returns a
// slice of Result values sorted by count descending.
func Apply(entries []parser.Entry, opts Options) []Result {
	if opts.Field == "" {
		opts.Field = "level"
	}

	counts := make(map[string]int, 16)
	for _, e := range entries {
		v := fieldValue(e, opts.Field)
		counts[v]++
	}

	results := make([]Result, 0, len(counts))
	for val, n := range counts {
		results = append(results, Result{Value: val, Count: n})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Count != results[j].Count {
			return results[i].Count > results[j].Count
		}
		return results[i].Value < results[j].Value
	})

	if opts.TopN > 0 && opts.TopN < len(results) {
		results = results[:opts.TopN]
	}
	return results
}

func fieldValue(e parser.Entry, field string) string {
	if field == "level" {
		return e.Level
	}
	if v, ok := e.Fields[field]; ok {
		return v
	}
	return ""
}

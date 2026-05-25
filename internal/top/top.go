// Package top provides frequency analysis over log entry fields,
// returning the top N most common values for a given field.
package top

import (
	"sort"

	"github.com/yourorg/logslice/internal/parser"
)

// Result holds a field value and its occurrence count.
type Result struct {
	Value string
	Count int
}

// Options configures the Top analysis.
type Options struct {
	// Field is the entry field to analyse. Use "level" or "message" for
	// built-in fields, or any key present in entry.Fields.
	Field string

	// N is the maximum number of results to return. Zero means return all.
	N int
}

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Field: "level",
		N:     10,
	}
}

// Apply counts occurrences of each unique value of opts.Field across entries
// and returns up to opts.N results ordered by descending count.
func Apply(entries []parser.Entry, opts Options) []Result {
	if len(entries) == 0 {
		return nil
	}

	counts := make(map[string]int, 32)
	for _, e := range entries {
		v := fieldValue(e, opts.Field)
		counts[v]++
	}

	results := make([]Result, 0, len(counts))
	for v, c := range counts {
		results = append(results, Result{Value: v, Count: c})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Count != results[j].Count {
			return results[i].Count > results[j].Count
		}
		return results[i].Value < results[j].Value
	})

	if opts.N > 0 && len(results) > opts.N {
		results = results[:opts.N]
	}
	return results
}

func fieldValue(e parser.Entry, field string) string {
	switch field {
	case "level":
		return e.Level
	case "message":
		return e.Message
	default:
		if v, ok := e.Fields[field]; ok {
			return v
		}
		return ""
	}
}

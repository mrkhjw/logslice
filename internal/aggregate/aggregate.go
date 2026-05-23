// Package aggregate provides log entry grouping and counting by field values.
package aggregate

import (
	"fmt"
	"sort"

	"github.com/yourorg/logslice/internal/parser"
)

// Result holds the aggregation output for a single group key.
type Result struct {
	Key   string
	Value string
	Count int
}

// Options controls how aggregation is performed.
type Options struct {
	// Field is the entry field name to group by (e.g. "level", "host").
	Field string
	// TopN limits results to the N most frequent groups. 0 means no limit.
	TopN int
}

// Apply groups entries by the specified field and returns sorted results.
func Apply(entries []parser.Entry, opts Options) []Result {
	counts := make(map[string]int)

	for _, e := range entries {
		val := fieldValue(e, opts.Field)
		counts[val]++
	}

	results := make([]Result, 0, len(counts))
	for val, count := range counts {
		results = append(results, Result{
			Key:   opts.Field,
			Value: val,
			Count: count,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Count != results[j].Count {
			return results[i].Count > results[j].Count
		}
		return results[i].Value < results[j].Value
	})

	if opts.TopN > 0 && len(results) > opts.TopN {
		results = results[:opts.TopN]
	}

	return results
}

// fieldValue extracts a comparable string value from an entry for the given field name.
func fieldValue(e parser.Entry, field string) string {
	switch field {
	case "level":
		return e.Level
	case "message":
		return e.Message
	default:
		if v, ok := e.Fields[field]; ok {
			return fmt.Sprintf("%v", v)
		}
		return ""
	}
}

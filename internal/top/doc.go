// Package top provides frequency analysis over a slice of log entries.
//
// It counts how often each unique value of a chosen field appears and returns
// the top N results ordered by descending frequency. Ties are broken
// alphabetically so output is deterministic.
//
// Usage:
//
//	opts := top.Options{Field: "host", N: 5}
//	results := top.Apply(entries, opts)
//	for _, r := range results {
//		fmt.Printf("%s: %d\n", r.Value, r.Count)
//	}
//
// Built-in field names "level" and "message" map directly to the
// corresponding parser.Entry fields. Any other name is looked up in
// entry.Fields; entries where the field is absent contribute an empty-string
// value to the count.
package top

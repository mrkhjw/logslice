// Package pivot provides field-based grouping and counting of log entries,
// producing a tabular pivot summary.
//
// # Overview
//
// Apply groups a slice of parser.Entry values by a chosen key field and
// counts occurrences per group. An optional value field collects the
// distinct values of a second field for each group. The result can be
// limited to the top-N most frequent groups.
//
// # Example
//
//	opts := pivot.Options{
//	    KeyField:   "level",
//	    ValueField: "host",
//	    TopN:       5,
//	}
//	res := pivot.Apply(entries, opts)
//	pivot.PrintSummary(os.Stdout, res)
//
// # Output format
//
// PrintSummary renders a fixed-width table with one row per group,
// sorted by count descending.
package pivot

// Package stats computes and displays summary statistics for parsed log entries.
//
// It provides:
//   - Compute: aggregates total counts, per-level breakdowns, time range, and
//     field value frequencies from a slice of parser.Entry values.
//   - Print: renders a human-readable summary table to any io.Writer.
//
// Typical usage:
//
//	summary := stats.Compute(entries)
//	stats.Print(os.Stdout, summary)
package stats

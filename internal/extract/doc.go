// Package extract provides field-extraction for log entries.
//
// It allows callers to select a subset of fields from each [parser.Entry],
// optionally renaming them, and to include or exclude the built-in message,
// level, and timestamp values.
//
// Basic usage:
//
//	opts := extract.NewBuilder().
//		Field("user").
//		FieldAs("svc", "service").
//		WithTimestamp(false).
//		Build()
//
//	result := extract.Apply(entries, opts)
//
// Each output entry's Fields map will contain only the extracted keys.
// The original entries are never mutated.
package extract

// Package cli implements the command-line interface for logslice.
//
// It wires together the parser, filter, and output packages into
// a single pipeline driven by flags:
//
//	-level   filter entries to a specific log level
//	-since   discard entries before this RFC3339 timestamp
//	-until   discard entries after this RFC3339 timestamp
//	-field   filter by a structured field in key=value form
//	-format  output format: text (default), json, or table
//
// Input is read from a file path supplied as the first positional
// argument, or from stdin when no path is given (or "-" is used).
package cli

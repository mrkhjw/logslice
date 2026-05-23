// Package highlight provides ANSI terminal color utilities for
// log level visualization in logslice text output.
//
// It maps log severity levels (ERROR, WARN, INFO, etc.) to
// corresponding ANSI colors, making terminal output easier to
// scan at a glance.
//
// Usage:
//
//	import "github.com/yourorg/logslice/internal/highlight"
//
//	colored := highlight.ForLevel("ERROR")  // red-colored "ERROR"
//	bold := highlight.BoldText("2024-01-01") // bold timestamp
//
// Supported log levels and their colors:
//
//	ERROR, FATAL  -> Red
//	WARN          -> Yellow
//	INFO          -> Green
//	DEBUG         -> Cyan
//	TRACE         -> Blue
//	(unknown)     -> no color applied
//
// Color output is written using ANSI escape codes and is intended
// for use in TTY-compatible terminals. Callers that write to files
// or non-terminal outputs should disable highlighting accordingly.
package highlight

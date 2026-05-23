// Package output provides formatters for rendering parsed log entries
// to various output formats including plain text, JSON (newline-delimited),
// and human-readable table format.
//
// Supported formats:
//
//   - "text"  – plain text, one entry per line
//   - "json"  – newline-delimited JSON (NDJSON), one JSON object per line
//   - "table" – human-readable table with aligned columns
//
// Usage:
//
//	f, err := output.ParseFormat("json")
//	if err != nil {
//		log.Fatal(err)
//	}
//	formatter := output.New(os.Stdout, f)
//	if err := formatter.Write(entries); err != nil {
//		log.Fatal(err)
//	}
package output

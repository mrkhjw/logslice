// Package output provides formatters for rendering parsed log entries
// to various output formats including plain text, JSON (newline-delimited),
// and human-readable table format.
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

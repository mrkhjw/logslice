// Package grep provides regex-based filtering of parsed log entries.
//
// It supports matching against the log message body, individual structured
// fields, or both simultaneously. Results can be inverted (grep -v style)
// to exclude matching entries instead.
//
// Usage:
//
//	c, err := grep.Compile(grep.Options{
//		Pattern:      "timeout",
//		FieldPattern: map[string]string{"service": "payment"},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	filtered := c.Apply(entries)
//
// All patterns are compiled as Go regular expressions (regexp/syntax).
package grep

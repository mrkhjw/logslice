// Package route implements rule-based routing of log entries into named destination buckets.
//
// Each Rule specifies a field to inspect, a regular expression to match, and a
// destination bucket name. Apply processes a slice of entries, testing each
// against the provided rules in order; the first matching rule wins and the
// entry is placed into the corresponding bucket. Entries that match no rule are
// collected in Result.Unmatched.
//
// Supported field selectors:
//
//	"level"   — matches against the parsed log level string
//	"message" — matches against the log message (default when Field is empty)
//	any other string — matches against the named extra field in Entry.Fields
//
// Example:
//
//	rules := []route.Rule{
//	    {Field: "level",   Pattern: regexp.MustCompile(`(?i)error`), Dest: "errors"},
//	    {Field: "service", Pattern: regexp.MustCompile(`auth`),      Dest: "auth"},
//	}
//	res := route.Apply(entries, rules)
//	// res.Buckets["errors"] — error-level entries
//	// res.Buckets["auth"]   — entries from the auth service
//	// res.Unmatched         — everything else
package route

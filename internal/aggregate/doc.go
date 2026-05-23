// Package aggregate groups parsed log entries by a specified field and
// returns sorted frequency counts. It is useful for summarising log data
// by level, hostname, service name, or any arbitrary structured field.
//
// Basic usage:
//
//	results := aggregate.Apply(entries, aggregate.Options{
//		Field: "level",
//		TopN:  5,
//	})
//	for _, r := range results {
//		fmt.Printf("%s=%s  count=%d\n", r.Key, r.Value, r.Count)
//	}
package aggregate

// Package scope partitions log entries into named time-window buckets.
//
// It is useful for grouping high-volume log streams into discrete intervals
// for analysis, reporting, or downstream aggregation.
//
// Basic usage:
//
//	opts := scope.DefaultOptions()          // 1-minute buckets
//	buckets := scope.Apply(entries, opts)
//	for _, b := range buckets {
//		fmt.Printf("%s: %d entries\n", b.Name, len(b.Entries))
//	}
//
// Custom bucket size:
//
//	opts := scope.Options{
//		Duration: 5 * time.Minute,
//		Label:    "2006-01-02T15:04",
//	}
//	buckets := scope.Apply(entries, opts)
//
// Entries with a zero timestamp are placed in a special "unscoped" bucket.
package scope

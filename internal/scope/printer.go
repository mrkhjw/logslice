package scope

import (
	"fmt"
	"io"
)

// PrintSummary writes a human-readable summary of buckets to w.
// Each bucket is shown with its name, time range, and entry count.
func PrintSummary(w io.Writer, buckets []Bucket) {
	if len(buckets) == 0 {
		fmt.Fprintln(w, "no scoped buckets")
		return
	}
	fmt.Fprintf(w, "%-30s  %-25s  %-25s  %s\n", "Bucket", "From", "To", "Count")
	fmt.Fprintf(w, "%-30s  %-25s  %-25s  %s\n",
		"------------------------------",
		"-------------------------",
		"-------------------------",
		"-----")
	for _, b := range buckets {
		fromStr := "-"
		toStr := "-"
		if !b.From.IsZero() {
			fromStr = b.From.UTC().Format("2006-01-02T15:04:05Z")
		}
		if !b.To.IsZero() {
			toStr = b.To.UTC().Format("2006-01-02T15:04:05Z")
		}
		fmt.Fprintf(w, "%-30s  %-25s  %-25s  %d\n",
			b.Name, fromStr, toStr, len(b.Entries))
	}
}

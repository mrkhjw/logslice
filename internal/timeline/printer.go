package timeline

import (
	"fmt"
	"io"
	"strings"
)

const (
	// maxBarWidth is the maximum number of characters in a histogram bar.
	maxBarWidth = 40
)

// PrintSummary writes a compact histogram of bucket counts to w.
func PrintSummary(w io.Writer, buckets []Bucket, layout string) {
	if len(buckets) == 0 {
		fmt.Fprintln(w, "(no entries)")
		return
	}
	if layout == "" {
		layout = "2006-01-02 15:04"
	}

	// Find the maximum count for scaling.
	max := 0
	for _, b := range buckets {
		if len(b.Entries) > max {
			max = len(b.Entries)
		}
	}

	fmt.Fprintf(w, "%-20s  %6s  %s\n", "bucket", "count", "bar")
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 60))
	for _, b := range buckets {
		count := len(b.Entries)
		barLen := 0
		if max > 0 {
			barLen = count * maxBarWidth / max
		}
		fmt.Fprintf(w, "%-20s  %6d  %s\n",
			b.Start.Format(layout),
			count,
			strings.Repeat("█", barLen),
		)
	}
}

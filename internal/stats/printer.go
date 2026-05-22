package stats

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Print writes a human-readable summary to w.
func Print(w io.Writer, s Summary) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, "Total entries:\t%d\n", s.Total)

	if s.FirstTime != nil {
		fmt.Fprintf(tw, "First entry:\t%s\n", s.FirstTime.Format("2006-01-02 15:04:05"))
	}
	if s.LastTime != nil {
		fmt.Fprintf(tw, "Last entry:\t%s\n", s.LastTime.Format("2006-01-02 15:04:05"))
	}
	if d := s.Duration(); d > 0 {
		fmt.Fprintf(tw, "Duration:\t%s\n", d.String())
	}

	fmt.Fprintf(tw, "\nLevel breakdown:\n")
	levels := make([]string, 0, len(s.ByLevel))
	for l := range s.ByLevel {
		levels = append(levels, l)
	}
	sort.Strings(levels)
	for _, l := range levels {
		pct := 0.0
		if s.Total > 0 {
			pct = float64(s.ByLevel[l]) / float64(s.Total) * 100
		}
		fmt.Fprintf(tw, "  %-8s\t%d\t(%.1f%%)\n", l, s.ByLevel[l], pct)
	}

	tw.Flush()
}

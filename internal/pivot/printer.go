package pivot

import (
	"fmt"
	"io"
	"strings"
)

// PrintSummary writes a human-readable pivot table to w.
func PrintSummary(w io.Writer, res Result) {
	if len(res.Rows) == 0 {
		fmt.Fprintln(w, "(no entries)")
		return
	}

	// Compute column widths.
	keyWidth := len(res.KeyField)
	for _, r := range res.Rows {
		if len(r.Key) > keyWidth {
			keyWidth = len(r.Key)
		}
	}

	header := fmt.Sprintf("%-*s  %8s", keyWidth, res.KeyField, "count")
	if res.ValueField != "" {
		header += fmt.Sprintf("  %s", res.ValueField)
	}
	fmt.Fprintln(w, header)
	fmt.Fprintln(w, strings.Repeat("-", len(header)))

	for _, r := range res.Rows {
		line := fmt.Sprintf("%-*s  %8d", keyWidth, r.Key, r.Count)
		if res.ValueField != "" && len(r.Values) > 0 {
			uniq := unique(r.Values)
			line += fmt.Sprintf("  [%s]", strings.Join(uniq, ", "))
		}
		fmt.Fprintln(w, line)
	}
}

func unique(ss []string) []string {
	seen := make(map[string]struct{}, len(ss))
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}

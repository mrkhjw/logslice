// Package scope provides time-window scoping for log entries,
// allowing entries to be grouped into named time buckets.
package scope

import (
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Bucket holds a named group of log entries within a time window.
type Bucket struct {
	Name    string
	From    time.Time
	To      time.Time
	Entries []parser.Entry
}

// Options controls how entries are scoped into buckets.
type Options struct {
	// Duration is the width of each bucket window.
	Duration time.Duration
	// Label is a strftime-style Go time format string used to name buckets.
	Label string
}

// DefaultOptions returns sensible defaults: 1-minute buckets, RFC3339 minute label.
func DefaultOptions() Options {
	return Options{
		Duration: time.Minute,
		Label:    "2006-01-02T15:04",
	}
}

// Apply partitions entries into time buckets according to opts.
// Entries with zero timestamps are placed in a bucket named "unscoped".
// The returned slice is ordered by bucket start time.
func Apply(entries []parser.Entry, opts Options) []Bucket {
	if len(entries) == 0 {
		return nil
	}
	if opts.Duration <= 0 {
		opts.Duration = DefaultOptions().Duration
	}
	if opts.Label == "" {
		opts.Label = DefaultOptions().Label
	}

	index := make(map[string]*Bucket)
	var order []string

	for _, e := range entries {
		if e.Timestamp.IsZero() {
			key := "unscoped"
			if _, ok := index[key]; !ok {
				index[key] = &Bucket{Name: key}
				order = append(order, key)
			}
			index[key].Entries = append(index[key].Entries, e)
			continue
		}

		trunc := e.Timestamp.Truncate(opts.Duration)
		key := trunc.Format(opts.Label)
		if _, ok := index[key]; !ok {
			index[key] = &Bucket{
				Name: key,
				From: trunc,
				To:   trunc.Add(opts.Duration),
			}
			order = append(order, key)
		}
		index[key].Entries = append(index[key].Entries, e)
	}

	buckets := make([]Bucket, 0, len(order))
	for _, k := range order {
		buckets = append(buckets, *index[k])
	}
	return buckets
}
